package nn

import (
	"reflect"
	"syscall/js"
)

var (
	// caching global objects for performance reasons
	global   = js.Global()
	document = global.Get("document")
)

func createElement(tag string) js.Value {
	return document.Call("createElement", tag)
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

		// 1. set classes and attributes
		if n.classes != "" {
			el.Set("className", n.classes)
		}
		for key, val := range n.attrs {
			el.Call("setAttribute", key, val)
		}

		// 2. add event listeners
		for eventName, _ := range n.listeners {
			// wrapper that calls handler() and schedule DOM update
			cb := js.FuncOf(func(this js.Value, args []js.Value) any {
				handler := n.listeners[eventName]
				if handler == nil {
					return nil
				}
				var goEvent Event
				if len(args) > 0 {
					goEvent = Event{jsEvent: args[0]}
				}
				handler(goEvent)

				nina.scheduleUpdate(nil)
				return nil
			})
			el.Call("addEventListener", eventName, cb)
			n.activeCallbacks[eventName] = cb
		}

		// 3. add children recursively
		for _, child := range n.children {
			childDOM := createDOM(child)
			el.Call("appendChild", childDOM)
		}

		return el

	case *ComponentNode:
		n.lastRender = n.comp.View()

		nina.registerComp(n.comp, n)

		dom := createDOM(n.lastRender)

		if m, ok := n.comp.(Mounter); ok {
			m.OnMount()
		}

		return dom
	default:
		panic("unknown node type")
	}
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
		patchEvents(old.domNode, old, newEl)

		// 1. update attributes
		patchClasses(old.domNode, old.classes, newEl.classes)
		patchAttributes(old.domNode, old.attrs, newEl.attrs)

		// 2. update children
		patchChildren(old.domNode, old.children, newEl.children)

	case *ComponentNode:
		newComp := newNode.(*ComponentNode)
		oldComp := oldNode.(*ComponentNode)

		newComp.parentDOM = parentDOM
		if oldComp.comp != newComp.comp {
			nina.unregisterComp(oldComp.comp)
		}
		nina.registerComp(newComp.comp, newComp)

		if pureComp, ok := newComp.comp.(Pure); ok {
			newHash := pureComp.Hash()
			newComp.hash = newHash

			if oldComp != nil && oldComp.hash == newHash {
				newComp.lastRender = oldComp.lastRender

				// exit from patch(), don't call View and don't go to children
				return
			}
		}

		newComp.lastRender = newComp.comp.View()

		patch(parentDOM, old.lastRender, newComp.lastRender)
	}
}

func patchEvents(domEl js.Value, oldE *Element, newE *Element) {
	if len(oldE.listeners) == 0 && len(newE.listeners) == 0 {
		return
	}

	if oldE.listeners == nil {
		oldE.listeners = make(map[string]func(Event))
		oldE.activeCallbacks = make(map[string]js.Func)
	}

	for eventName, newHandler := range newE.listeners {
		if _, exists := oldE.listeners[eventName]; !exists {
			cb := js.FuncOf(func(this js.Value, args []js.Value) any {
				handler := oldE.listeners[eventName]
				if handler != nil {
					var goEvent Event
					if len(args) > 0 {
						goEvent = Event{jsEvent: args[0]}
					}
					handler(goEvent)
					nina.scheduleUpdate(nil)
				}
				return nil
			})

			domEl.Call("addEventListener", eventName, cb)

			oldE.activeCallbacks[eventName] = cb
		}

		oldE.listeners[eventName] = newHandler
	}

	// remove old listeners
	for eventName, cb := range oldE.activeCallbacks {
		if _, exists := newE.listeners[eventName]; !exists {
			domEl.Call("removeEventListener", eventName, cb)
			cb.Release()

			delete(oldE.activeCallbacks, eventName)
			delete(oldE.listeners, eventName)
		}
	}

	newE.activeCallbacks = oldE.activeCallbacks
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
	case "class":
		domNode.Set("className", val)
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
		for _, cb := range n.activeCallbacks {
			cb.Release()
		}
		n.activeCallbacks = nil

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
	}
}
