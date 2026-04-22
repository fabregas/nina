package nn

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"syscall/js"
)

var (
	// caching global objects for performance reasons
	global   = js.Global()
	document = global.Get("document")

	currentRenderingComponent Component
)

func createElement(tag string) js.Value {
	// maybe switch?
	isSVG := tag == "svg" || tag == "path" || tag == "circle" || tag == "rect" || tag == "g"

	if isSVG {
		return document.Call("createElementNS", "http://www.w3.org/2000/svg", tag)
	} else {
		return document.Call("createElement", tag)
	}
}

func createTextNode(text string) js.Value {
	return document.Call("createTextNode", text)
}

// createDOM generates real DOM from virtual Node
func createDOM(vnode Node) js.Value {
	switch n := vnode.(type) {

	case *TextNode:
		n.domNode = createTextNode(n.value)
		return n.domNode

	case *Element:
		el := createElement(n.tag)
		n.domNode = el
		if n.ref != nil {
			n.ref.Current = el
		}

		// 1. set classes and attributes
		if n.classes != "" {
			n.compClasses()
			el.Call("setAttribute", "class", n.classes)
		}
		for key, val := range n.attrs {
			el.Call("setAttribute", key, val)
		}

		// 2. add event listeners
		if n.listeners != nil {
			n.listeners.parentComponent = currentRenderingComponent
			for eventInfo, _ := range n.listeners.events {
				addEventListener(el, n.listeners, eventInfo)
			}
		}

		// 3. add children
		if n.rawHTML != "" {
			el.Set("innerHTML", n.rawHTML)
		} else {
			for _, child := range n.children {
				childDOM := createDOM(child)
				el.Call("appendChild", childDOM)
			}
		}

		return el

	case *groupNode:
		frag := document.Call("createDocumentFragment")

		for _, child := range n.children {
			childDOM := createDOM(child)
			frag.Call("appendChild", childDOM)
		}

		return frag

	case *componentNode:
		if carrier, ok := n.comp.(stateCarrier); ok {
			carrier.setUpdater(func() { nina.scheduleUpdate(n.comp) })
			carrier.importState(nil)
		}

		prevComponent := currentRenderingComponent
		currentRenderingComponent = n.comp

		if i, ok := n.comp.(Initer); ok {
			i.OnInit()
		}

		n.lastRender = n.comp.View()

		nina.registerComp(n.comp, n)

		dom := createDOM(n.lastRender)

		currentRenderingComponent = prevComponent

		if m, ok := n.comp.(Mounter); ok {
			nina.scheduleMount(m.OnMount)
		}

		return dom

	case *portalNode:
		targetDOM := document.Call("querySelector", n.targetSelector)

		if targetDOM.IsNull() || targetDOM.IsUndefined() {
			panic("Portal target not found: " + n.targetSelector)
		}

		var childDOM js.Value
		if n.child != nil {
			childDOM = createDOM(n.child)
			targetDOM.Call("appendChild", childDOM)
		}
		n.domNode = childDOM
		n.placeholderNode = document.Call("createComment", "portal-placeholder")

		return n.placeholderNode
	default:
		panic("unknown node type")
	}
}

func addEventListener(el js.Value, listeners *listenersInfo, e eventInfo) {
	// wrapper that calls handler() and schedule DOM update
	cbFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		handler := listeners.events[e]
		if handler == nil {
			return nil
		}
		var skipUpdate bool
		goEvent := Event{skipUpdate: &skipUpdate}
		if len(args) > 0 {
			goEvent.jsEvent = args[0]
		}
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[Nina] panic in cb at %T [%s]: %+v\n", listeners.parentComponent, e.name, r)
				fmt.Println(string(debug.Stack()))
			}
		}()
		handler(goEvent)

		if !skipUpdate {
			nina.scheduleUpdate(listeners.parentComponent)
		}
		return nil
	})

	cb := cbInfo{fn: cbFunc}
	if e.isGlobal {
		switch e.name {
		case "scroll":
			document.Call("addEventListener", e.name, cbFunc, js.ValueOf(true))
		case "resize-el":
			// TODO global observer singleton
			observer := global.Get("ResizeObserver").New(cbFunc)

			observer.Call("observe", el)
			cb.closeFn = func() {
				observer.Call("disconnect")
			}
		default:
			document.Call("addEventListener", e.name, cbFunc)
		}
	} else {
		el.Call("addEventListener", e.name, cbFunc)
	}

	listeners.activeCallbacks[e] = cb
}

func getDOMNode(vnode Node) js.Value {
	switch n := vnode.(type) {
	case *TextNode:
		return n.domNode
	case *Element:
		return n.domNode
	case *groupNode:
		for _, ch := range n.children {
			dom := getDOMNode(ch)
			if !dom.IsNull() && !dom.IsUndefined() {
				return dom
			}
		}

		return js.Null()
	case *componentNode:
		if isNilNode(n) || n.lastRender == nil {
			return js.Undefined()
		}
		return n.lastRender.domNode
	case *portalNode:
		return n.placeholderNode

	default:
		return js.Undefined()
	}
}

func isDifferentType(n1, n2 Node) bool {
	n1IsNil := isNilNode(n1)
	n2IsNil := isNilNode(n2)

	if n1IsNil && n2IsNil {
		return false
	}

	if n1IsNil || n2IsNil {
		return true
	}

	if reflect.TypeOf(n1) != reflect.TypeOf(n2) {
		return true
	}

	switch n1Val := n1.(type) {
	case *TextNode:
		return false
	case *Element:
		n2Val := n2.(*Element)

		return n1Val.tag != n2Val.tag
	case *groupNode:
		return false
	case *componentNode:
		n2Val := n2.(*componentNode)
		oldType := reflect.TypeOf(n1Val.comp)
		newType := reflect.TypeOf(n2Val.comp)

		return oldType != newType
	}

	return false
}

// patch compare old and new node and makes changes into parentDOM
func patch(parentDOM js.Value, oldNode, newNode Node) {
	if isNilNode(oldNode) && isNilNode(newNode) {
		return
	}

	// case#1: create new node
	if isNilNode(oldNode) && !isNilNode(newNode) {
		// fmt.Println("create new node - ", newNode)
		parentDOM.Call("appendChild", createDOM(newNode))
		return
	}

	// case#2: removing node
	if !isNilNode(oldNode) && isNilNode(newNode) {
		// fmt.Println("remove node - ", oldNode)
		destroy(oldNode)
		return
	}

	// case#3: different nodes types or different element tags - full replace
	if isDifferentType(oldNode, newNode) {
		// fmt.Printf("different nodes types  [%v] ~ [%v] \n", oldNode, newNode)
		anchorDOM := getDOMNode(oldNode)
		newDOM := createDOM(newNode)

		if !anchorDOM.IsNull() && !anchorDOM.IsUndefined() {
			parentDOM.Call("insertBefore", newDOM, anchorDOM)
		} else {
			// Fallback
			parentDOM.Call("appendChild", newDOM)
		}

		destroy(oldNode)
		return
	}

	// case#4: the same nodes types, check content
	switch old := oldNode.(type) {

	case *TextNode:
		newText := newNode.(*TextNode)
		if old.value != newText.value {
			//fmt.Println("update TextNode ", newText.value)
			old.domNode.Set("nodeValue", newText.value)
		}
		newText.domNode = old.domNode

	case *Element:
		newEl := newNode.(*Element)
		newEl.domNode = old.domNode
		if newEl.ref != nil {
			newEl.ref.Current = old.domNode
		}
		patchEvents(old.domNode, old, newEl)

		// 1. update attributes

		newEl.compClasses()
		patchClasses(old.domNode, old.classes, newEl.classes)
		patchAttributes(old.domNode, old.attrs, newEl.attrs)

		// 2. update children
		if old.rawHTML != "" {
			old.domNode.Set("innerHTML", newEl.rawHTML)
		} else {
			patchChildren(old.domNode, old.children, newEl.children)
		}

	case *groupNode:
		newGroup := newNode.(*groupNode)
		patchChildren(parentDOM, old.children, newGroup.children)

	case *componentNode:
		newComp := newNode.(*componentNode)
		oldComp := oldNode.(*componentNode)

		newComp.parentDOM = parentDOM
		if oldComp.comp != newComp.comp {
			nina.unregisterComp(oldComp.comp)
		}
		nina.registerComp(newComp.comp, newComp)

		if newCarrier, ok := newComp.comp.(stateCarrier); ok {
			oldCarrier := oldComp.comp.(stateCarrier)

			newCarrier.setUpdater(func() { nina.scheduleUpdate(newComp.comp) })
			// copy old state into new component
			newCarrier.importState(oldCarrier.exportState())

		}

		if pureComp, ok := newComp.comp.(Pure); ok {
			newHash := pureComp.Hash()
			newComp.hash = newHash

			if oldComp != nil && oldComp.hash == newHash {
				newComp.lastRender = oldComp.lastRender

				// exit from patch(), don't call View and don't go to children
				return
			}
		}

		prevComponent := currentRenderingComponent
		currentRenderingComponent = newComp.comp

		newComp.lastRender = newComp.comp.View()

		patch(parentDOM, old.lastRender, newComp.lastRender)

		currentRenderingComponent = prevComponent

	case *portalNode:
		newPortal := newNode.(*portalNode)
		newPortal.domNode = old.domNode
		newPortal.placeholderNode = old.placeholderNode

		patchChildren(old.domNode, []Node{old.child}, []Node{newPortal.child})
	}
}

func patchEvents(domEl js.Value, oldE *Element, newE *Element) {
	if oldE.listeners == nil && newE.listeners == nil {
		return
	}

	if oldE.listeners == nil {
		oldE.listeners = &listenersInfo{
			events:          make(map[eventInfo]func(Event)),
			activeCallbacks: make(map[eventInfo]cbInfo),
		}
	}

	if newE.listeners == nil {
		newE.listeners = &listenersInfo{}
	}

	oldE.listeners.parentComponent = currentRenderingComponent

	for eventInfo, newHandler := range newE.listeners.events {
		if _, exists := oldE.listeners.events[eventInfo]; !exists {
			addEventListener(domEl, oldE.listeners, eventInfo)
		}

		oldE.listeners.events[eventInfo] = newHandler
	}

	// remove old listeners
	for eventInfo, _ := range oldE.listeners.activeCallbacks {
		if _, exists := newE.listeners.events[eventInfo]; !exists {
			destroyEventListeners(oldE, eventInfo)

			delete(oldE.listeners.activeCallbacks, eventInfo)
			delete(oldE.listeners.events, eventInfo)
		}
	}

	newE.listeners = oldE.listeners
}

func destroyEventListeners(el *Element, ei eventInfo) {
	cb := el.listeners.activeCallbacks[ei]
	if ei.isGlobal {
		switch ei.name {
		case "scroll":
			document.Call("removeEventListener", ei.name, cb.fn, js.ValueOf(true))
		case "resize-el":
		default:
			document.Call("removeEventListener", ei.name, cb.fn)
		}
	} else {
		el.domNode.Call("removeEventListener", ei.name, cb.fn)
	}

	if cb.closeFn != nil {
		cb.closeFn()
	}

	cb.fn.Release()
}

func patchChildren(parentDOM js.Value, oldChilds, newChilds []Node) {
	oldKeyed := make(map[string]Node, len(oldChilds))
	oldIndices := make(map[string]int, len(oldChilds))
	var oldUnkeyed []Node

	for i, oldChild := range oldChilds {
		if oldChild == nil {
			continue
		}
		key := oldChild.getKey()
		if key != "" {
			oldKeyed[key] = oldChild
			oldIndices[key] = i
		} else {
			oldUnkeyed = append(oldUnkeyed, oldChild)
		}
	}

	unkeyedIndex := 0
	lastPlacedIndex := 0

	for i, newChild := range newChilds {
		if newChild == nil {
			continue
		}

		key := newChild.getKey()
		var matchedOld Node

		if key != "" {
			if old, exists := oldKeyed[key]; exists {
				matchedOld = old
				delete(oldKeyed, key)
			}
		} else {
			if unkeyedIndex < len(oldUnkeyed) {
				matchedOld = oldUnkeyed[unkeyedIndex]
				unkeyedIndex++
			}
		}

		var anchorDOM js.Value = js.Null()
		if i > 0 {
			prevVNode := newChilds[i-1]
			lastPrevDOM := getLastRealDOM(prevVNode)
			if !lastPrevDOM.IsNull() && !lastPrevDOM.IsUndefined() {
				anchorDOM = lastPrevDOM.Get("nextSibling")
			}
		} else {
			anchorDOM = parentDOM.Get("firstChild")
		}

		if matchedOld != nil {
			// update attrs and text
			patch(parentDOM, matchedOld, newChild)

			oldIndex, hasOldIndex := oldIndices[key]

			if hasOldIndex {
				if oldIndex < lastPlacedIndex {
					moveNode(newChild, parentDOM, anchorDOM)
				} else {
					lastPlacedIndex = oldIndex
				}
			}
		} else {
			// new node... just create
			newDOM := createDOM(newChild)

			if !anchorDOM.IsNull() && !anchorDOM.IsUndefined() {
				parentDOM.Call("insertBefore", newDOM, anchorDOM)
			} else {
				parentDOM.Call("appendChild", newDOM)
			}
		}
	}

	// remove nodes
	for _, oldChild := range oldKeyed {
		destroy(oldChild)
	}

	for i := unkeyedIndex; i < len(oldUnkeyed); i++ {
		destroy(oldUnkeyed[i])
	}
}

func getLastRealDOM(vnode Node) js.Value {
	if vnode == nil {
		return js.Null()
	}

	switch n := vnode.(type) {
	case *TextNode:
		return n.domNode
	case *Element:
		return n.domNode
	case *groupNode:
		for i := len(n.children) - 1; i >= 0; i-- {
			dom := getLastRealDOM(n.children[i])
			if !dom.IsNull() && !dom.IsUndefined() {
				return dom
			}
		}

		return js.Null()
	case *componentNode:
		if isNilNode(n) || n.lastRender == nil {
			return js.Undefined()
		}
		return n.lastRender.domNode
	case *portalNode:
		return n.placeholderNode
	}

	return js.Null()
}

func moveNode(vnode Node, parentDOM js.Value, anchorDOM js.Value) {
	if vnode == nil {
		return
	}

	var actualDom js.Value
	switch n := vnode.(type) {
	case *Element:
		actualDom = n.domNode
	case *TextNode:
		actualDom = n.domNode
	case *componentNode:
		if isNilNode(n) || n.lastRender == nil {
			return
		}
		actualDom = n.lastRender.domNode
	case *portalNode:
		actualDom = n.placeholderNode

	case *groupNode:
		for _, child := range n.children {
			moveNode(child, parentDOM, anchorDOM)
		}
		return
	}

	if !actualDom.IsNull() && !actualDom.IsUndefined() {
		if !anchorDOM.IsNull() && !anchorDOM.IsUndefined() {
			parentDOM.Call("insertBefore", actualDom, anchorDOM)
		} else {
			parentDOM.Call("appendChild", actualDom)
		}
	}

}

func patchClasses(domNode js.Value, oldClasses, newClasses string) {
	if oldClasses == newClasses {
		return
	}

	setAttribute(domNode, "class", newClasses)
}

func patchAttributes(domNode js.Value, oldAttrs, newAttrs map[string]string) {
	for key, newVal := range newAttrs {
		oldVal, exists := oldAttrs[key]

		if !exists || oldVal != newVal {
			setAttribute(domNode, key, newVal)
		}
	}

	// remove attributes
	for key := range oldAttrs {
		if _, exists := newAttrs[key]; !exists {
			removeAttribute(domNode, key)
		}
	}
}

func setAttribute(domNode js.Value, key, val string) {
	//fmt.Println("set attribute ", key, val)
	switch key {
	case "className":
		domNode.Call("setAttribute", "class", val)
	case "value":
		// for <input> and <textarea> we chould change DOM attribute value,
		// bcs setAttribute changes only initial value
		domNode.Set("value", val)
	case "checked", "disabled":
		if val == "true" || val == "" {
			domNode.Set(key, true)
		} else {
			domNode.Set(key, false)
		}
	default:
		// for all other HTML attributes (id, src, href, style...)
		domNode.Call("setAttribute", key, val)
	}
}

func removeAttribute(domNode js.Value, key string) {
	switch key {
	case "class":
		domNode.Set("className", "")
	case "value":
		domNode.Set("value", "")
	case "checked", "disabled":
		domNode.Set(key, false)
	default:
		domNode.Call("removeAttribute", key)
	}
}

func destroy(node Node) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *Element:
		if n.listeners != nil {
			for ei, _ := range n.listeners.activeCallbacks {
				destroyEventListeners(n, ei)
			}

			n.listeners.activeCallbacks = make(map[eventInfo]cbInfo)
		}

		for _, child := range n.children {
			destroy(child)
		}

		if !n.domNode.IsNull() && !n.domNode.IsUndefined() {
			n.domNode.Call("remove")
		}

	case *TextNode:
		if !n.domNode.IsNull() && !n.domNode.IsUndefined() {
			n.domNode.Call("remove")
		}

	case *groupNode:
		for _, child := range n.children {
			destroy(child)
		}

	case *componentNode:
		Storage.unwatchAll(n.comp)

		nina.unregisterComp(n.comp)

		if d, ok := n.comp.(Destroyer); ok {
			d.OnDestroy()
		}

		if n.lastRender != nil {
			destroy(n.lastRender)
		}
	case *portalNode:
		if n.child != nil {
			destroy(n.child)
		}

		if !n.domNode.IsNull() && !n.domNode.IsUndefined() {
			n.domNode.Call("remove")
		}
	}
}
