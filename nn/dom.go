package nn

import (
	"reflect"
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

	case *ComponentNode:
		if carrier, ok := n.comp.(stateCarrier); ok {
			carrier.setUpdater(func() { nina.scheduleUpdate(n.comp) })
			carrier.importState(nil)
		}

		prevComponent := currentRenderingComponent
		currentRenderingComponent = n.comp

		n.lastRender = n.comp.View()

		nina.registerComp(n.comp, n)

		dom := createDOM(n.lastRender)

		currentRenderingComponent = prevComponent

		if m, ok := n.comp.(Mounter); ok {
			nina.scheduleMount(m.OnMount)
		}

		return dom

	case *PortalNode:
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
	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		handler := listeners.events[e]
		if handler == nil {
			return nil
		}
		var skipUpdate bool
		goEvent := Event{skipUpdate: &skipUpdate}
		if len(args) > 0 {
			goEvent.jsEvent = args[0]
		}
		handler(goEvent)

		if !skipUpdate {
			nina.scheduleUpdate(listeners.parentComponent)
		}
		return nil
	})

	if e.isGlobal {
		if e.name == "scroll" {
			document.Call("addEventListener", e.name, cb, js.ValueOf(true))
		} else {
			document.Call("addEventListener", e.name, cb)
		}
	} else {
		el.Call("addEventListener", e.name, cb)
	}

	listeners.activeCallbacks[e] = cb
}

func getDOMNode(vnode Node) js.Value {
	switch n := vnode.(type) {
	case *TextNode:
		return n.domNode
	case *Element:
		return n.domNode
	case *ComponentNode:
		if isNilNode(n) || n.lastRender == nil {
			return js.Undefined()
		}
		return n.lastRender.domNode
	case *PortalNode:
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
	case *ComponentNode:
		n2Val := n2.(*ComponentNode)
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
		//fmt.Println("create new node - ", newNode)
		parentDOM.Call("appendChild", createDOM(newNode))
		return
	}

	// case#2: removing node
	if !isNilNode(oldNode) && isNilNode(newNode) {
		//fmt.Println("remove node - ", oldNode)
		parentDOM.Call("removeChild", getDOMNode(oldNode))
		destroy(oldNode)
		return
	}

	// case#3: different nodes types or different element tags - full replace
	if isDifferentType(oldNode, newNode) {
		//fmt.Println("different nodes types, replace to - ", newNode)
		newDOM := createDOM(newNode)
		parentDOM.Call("replaceChild", newDOM, getDOMNode(oldNode))
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

	case *ComponentNode:
		newComp := newNode.(*ComponentNode)
		oldComp := oldNode.(*ComponentNode)

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

	case *PortalNode:
		newPortal := newNode.(*PortalNode)
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
			activeCallbacks: make(map[eventInfo]js.Func),
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
	for eventInfo, cb := range oldE.listeners.activeCallbacks {
		if _, exists := newE.listeners.events[eventInfo]; !exists {
			if eventInfo.isGlobal {
				document.Call("removeEventListener", eventInfo.name, cb)
			} else {
				domEl.Call("removeEventListener", eventInfo.name, cb)
			}
			cb.Release()

			delete(oldE.listeners.activeCallbacks, eventInfo)
			delete(oldE.listeners.events, eventInfo)
		}
	}

	newE.listeners = oldE.listeners
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

		if matchedOld != nil {
			// update attrs and text
			patch(parentDOM, matchedOld, newChild)

			oldIndex, hasOldIndex := oldIndices[key]

			if hasOldIndex {
				// like React ... if old index less than last placed that we phisycally move it forward
				if oldIndex < lastPlacedIndex {
					childNodes := parentDOM.Get("childNodes")
					if i < childNodes.Get("length").Int() {
						currentDOM := childNodes.Call("item", i)
						parentDOM.Call("insertBefore", getDOMNode(matchedOld), currentDOM)
					} else {
						parentDOM.Call("appendChild", getDOMNode(matchedOld))
					}
				} else {
					lastPlacedIndex = oldIndex
				}
			}
		} else {
			// new node... just create
			newDOM := createDOM(newChild)
			childNodes := parentDOM.Get("childNodes")
			if i < childNodes.Get("length").Int() {
				currentDOM := childNodes.Call("item", i)
				parentDOM.Call("insertBefore", newDOM, currentDOM)
			} else {
				parentDOM.Call("appendChild", newDOM)
			}
		}
	}

	// remove nodes
	for _, oldChild := range oldKeyed {
		parentDOM.Call("removeChild", getDOMNode(oldChild))
		destroy(oldChild)
	}

	for i := unkeyedIndex; i < len(oldUnkeyed); i++ {
		parentDOM.Call("removeChild", getDOMNode(oldUnkeyed[i]))
		destroy(oldUnkeyed[i])
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
			for i, cb := range n.listeners.activeCallbacks {
				if i.isGlobal {
					if i.name == "scroll" {
						document.Call("removeEventListener", i.name, cb, js.ValueOf(true))
					} else {
						document.Call("removeEventListener", i.name, cb)
					}
					cb.Release()
				}
			}

			n.listeners.activeCallbacks = nil
		}

		for _, child := range n.children {
			destroy(child)
		}

	case *ComponentNode:
		Storage.unwatchAll(n.comp)

		nina.unregisterComp(n.comp)

		if d, ok := n.comp.(Destroyer); ok {
			d.OnDestroy()
		}

		if n.lastRender != nil {
			destroy(n.lastRender)
		}
	case *PortalNode:
		if n.child != nil {
			destroy(n.child)
		}

		if !n.domNode.IsNull() && !n.domNode.IsUndefined() {
			n.domNode.Call("remove")
		}
	}
}
