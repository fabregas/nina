package nn

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"sync"
)

type engine struct {
	mu              sync.Mutex
	registry        map[Component]*componentNode
	dirtyComponents map[Component]bool
	mountQueue      []func()
	updateScheduled bool

	rootComponent  Component
	lastGlobalTree *componentNode
	rootContainer  NativeNode
	reqDomUpdate   func()

	currentRenderingComponent Component

	renderer         Renderer
	documentRenderer DocumentRenderer
	storage          LocalStorage
}

func (e *engine) registerComp(c Component, node *componentNode) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.registry[c] = node
}

func (e *engine) unregisterComp(c Component) {
	e.storage.unwatchAll(c)

	e.mu.Lock()
	delete(e.registry, c)
	delete(e.dirtyComponents, c)
	e.mu.Unlock()

	c.destroy()
}

func (e *engine) scheduleUpdate(c Component) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if c == nil {
		e.dirtyComponents[e.rootComponent] = true
	} else if _, exists := e.registry[c]; exists {
		e.dirtyComponents[c] = true
	} else {
		fmt.Printf("unknown component: %T\n", c)
	}

	if !e.updateScheduled {
		e.updateScheduled = true
		e.reqDomUpdate()
	}
}

func (e *engine) scheduleMount(f func()) {
	e.mountQueue = append(e.mountQueue, f)
}

func (e *engine) performUpdates() {
	//t0 := time.Now()
	//defer func() {
	//	fmt.Println("performUpdates time:", time.Since(t0))
	//}()

	e.mu.Lock()
	e.updateScheduled = false
	queue := e.dirtyComponents
	e.dirtyComponents = make(map[Component]bool)
	e.mu.Unlock()

	// full re-render must be in priority
	if queue[e.rootComponent] {
		newTree := Comp(e.rootComponent)
		e.patch(e.rootContainer, e.lastGlobalTree, newTree)
		e.lastGlobalTree = newTree

		// just exit, everything is re-rendered
	} else {
		for comp := range queue {
			//fmt.Printf("[RENDER] %T - %p\n", comp, comp)

			e.mu.Lock()
			node, exists := e.registry[comp]
			e.mu.Unlock()

			if !exists || node == nil {
				continue
			}

			e.currentRenderingComponent = node.comp

			// local patch
			newRender := node.comp.View()
			e.patch(node.parentDOM, node.lastRender, newRender)
			node.lastRender = newRender
		}
	}

	for _, onMountCb := range e.mountQueue {
		onMountCb()
	}
	e.mountQueue = nil
}

func (e *engine) createElement(tag string) NativeNode {
	// maybe switch?
	isSVG := tag == "svg" || tag == "path" || tag == "circle" || tag == "rect" || tag == "g"

	if isSVG {
		return e.renderer.CreateElementNS("http://www.w3.org/2000/svg", tag)
	} else {
		return e.renderer.CreateElement(tag)
	}
}

// createDOM generates real DOM from virtual Node
func (e *engine) createDOM(parentDOM NativeNode, vnode Node) NativeNode {
	if isNilNode(vnode) {
		return nil
	}

	switch n := vnode.(type) {

	case *TextNode:
		n.domNode = e.renderer.CreateTextNode(n.value)
		return n.domNode

	case *Element:
		el := e.createElement(n.tag)
		n.domNode = el
		if len(n.refs) > 0 {
			for _, r := range n.refs {
				r.Current = el
				r.Renderer = e.renderer
			}
		}

		// 1. set classes and attributes
		if n.classes != "" {
			n.compClasses()
			e.renderer.SetAttribute(el, "class", n.classes)
		}
		for key, val := range n.attrs {
			e.renderer.SetAttribute(el, key, val)
		}

		// 2. add event listeners
		if n.listeners != nil {
			n.listeners.parentComponent = e.currentRenderingComponent
			for eventInfo, _ := range n.listeners.events {
				e.addEventListener(el, n.listeners, eventInfo)
			}
		}

		// 3. add children
		if n.rawHTML != "" {
			e.renderer.SetInnerHTML(el, n.rawHTML)
		} else {
			for _, child := range n.children {
				childDOM := e.createDOM(el, child)
				e.appendDOMChild(el, childDOM)
			}
		}

		return el

	case *groupNode:
		frag := e.renderer.CreateDocumentFragment()

		for _, child := range n.children {
			childDOM := e.createDOM(frag, child)
			e.appendDOMChild(frag, childDOM)
		}

		return frag

	case *componentNode:
		if carrier, ok := n.comp.(stateCarrier); ok {
			carrier.importState(nil)
		}
		n.comp.setUpdater(func() { e.scheduleUpdate(n.comp) })

		prevComponent := e.currentRenderingComponent
		e.currentRenderingComponent = n.comp
		n.comp.setParent(prevComponent)
		n.parentDOM = parentDOM

		n.comp.resolveContexts()

		if i, ok := n.comp.(Initer); ok {
			i.OnInit()
		}

		n.lastRender = n.comp.View()

		e.registerComp(n.comp, n)

		dom := e.createDOM(parentDOM, n.lastRender)

		e.currentRenderingComponent = prevComponent

		if m, ok := n.comp.(Mounter); ok {
			e.scheduleMount(m.OnMount)
		}

		return dom

	case *portalNode:
		targetDOM := e.renderer.QuerySelector(e.renderer.RootNode(), n.targetSelector)

		if targetDOM == nil {
			panic("Portal target not found: " + n.targetSelector)
		}

		var childDOM NativeNode
		if n.child != nil {
			childDOM = e.createDOM(targetDOM, n.child)
			e.appendDOMChild(targetDOM, childDOM)
		}
		n.domNode = childDOM
		n.placeholderNode = e.renderer.CreateComment("portal-placeholder")

		return n.placeholderNode
	default:
		panic("unknown node type")
	}
}

func (e *engine) addEventListener(el NativeNode, listeners *listenersInfo, eInfo eventInfo) {
	// wrapper that calls handler() and schedule DOM update
	handlerFunc := func(event Event) {
		handler := listeners.events[eInfo]
		if handler == nil {
			return
		}

		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[Nina] panic in cb at %T [%s]: %+v\n", listeners.parentComponent, eInfo.name, r)
				fmt.Println(string(debug.Stack()))
			}
		}()
		handler(event)

		if event.needUpdate() {
			e.scheduleUpdate(listeners.parentComponent)
		}
	}

	var cleaner func()

	if eInfo.isGlobal {
		switch eInfo.name {
		case "scroll":
			cleaner = e.renderer.AddEventListenerWithCapture(nil, eInfo.name, handlerFunc)
		case "resize-el":
			cleaner = e.renderer.AddResizeObserver(el, handlerFunc)
		default:
			cleaner = e.renderer.AddEventListener(nil, eInfo.name, handlerFunc)
		}
	} else {
		cleaner = e.renderer.AddEventListener(el, eInfo.name, handlerFunc)
	}
	listeners.activeCallbacks[eInfo] = cleaner

}

// patch compare old and new node and makes changes into parentDOM
func (e *engine) patch(parentDOM NativeNode, oldNode, newNode Node) {
	if isNilNode(oldNode) && isNilNode(newNode) {
		return
	}

	// case#1: create new node
	if isNilNode(oldNode) && !isNilNode(newNode) {
		// fmt.Println("create new node - ", newNode)
		e.appendDOMChild(parentDOM, e.createDOM(parentDOM, newNode))
		return
	}

	// case#2: removing node
	if !isNilNode(oldNode) && isNilNode(newNode) {
		// fmt.Println("remove node - ", oldNode)
		e.destroy(oldNode)
		return
	}

	// case#3: different nodes types or different element tags - full replace
	if isDifferentType(oldNode, newNode) {
		//fmt.Printf("different nodes types at %v  [%v] ~ [%v] \n", parentDOM, oldNode, newNode)
		anchorDOM := getDOMNode(oldNode)
		newDOM := e.createDOM(parentDOM, newNode)

		if anchorDOM != nil {
			e.insertDOMBefore(parentDOM, newDOM, anchorDOM)
		} else {
			// Fallback
			e.appendDOMChild(parentDOM, newDOM)
		}

		e.destroy(oldNode)
		return
	}

	// case#4: the same nodes types, check content
	switch old := oldNode.(type) {

	case *TextNode:
		newText := newNode.(*TextNode)
		if old.value != newText.value {
			e.renderer.SetNodeValue(old.domNode, newText.value)
		}
		newText.domNode = old.domNode

	case *Element:
		newEl := newNode.(*Element)
		newEl.domNode = old.domNode
		if len(newEl.refs) > 0 {
			for _, r := range newEl.refs {
				r.Current = old.domNode
				r.Renderer = e.renderer
			}
		}
		e.patchEvents(old.domNode, old, newEl)

		// 1. update attributes

		newEl.compClasses()
		e.patchClasses(old.domNode, old.classes, newEl.classes)
		e.patchAttributes(old.domNode, old.attrs, newEl.attrs)

		// 2. update children
		if old.rawHTML != "" {
			e.renderer.SetInnerHTML(old.domNode, newEl.rawHTML)
		} else {
			e.patchChildren(old.domNode, old.children, newEl.children)
		}

	case *groupNode:
		newGroup := newNode.(*groupNode)
		e.patchChildren(parentDOM, old.children, newGroup.children)

	case *componentNode:
		newComp := newNode.(*componentNode)
		oldComp := oldNode.(*componentNode)

		if newCarrier, ok := newComp.comp.(stateCarrier); ok {
			oldCarrier := oldComp.comp.(stateCarrier)
			// copy old state into new component
			newCarrier.importState(oldCarrier.exportState())
		}

		newComp.parentDOM = parentDOM
		if oldComp.comp != newComp.comp {
			e.unregisterComp(oldComp.comp)

			newComp.comp.resolveContexts()
		}
		e.registerComp(newComp.comp, newComp)

		newUpdater := func() { e.scheduleUpdate(newComp.comp) }
		newComp.comp.setUpdater(newUpdater)
		oldComp.comp.setUpdater(newUpdater)

		if pureComp, ok := newComp.comp.(Pure); ok {
			newHash := pureComp.Hash()
			newComp.hash = newHash

			if oldComp != nil && oldComp.hash == newHash {
				newComp.lastRender = oldComp.lastRender

				// exit from patch(), don't call View and don't go to children
				return
			}
		}

		prevComponent := e.currentRenderingComponent
		e.currentRenderingComponent = newComp.comp
		newComp.comp.setParent(prevComponent)

		prevRender := old.lastRender // important for old == newComp case
		newComp.lastRender = newComp.comp.View()

		e.patch(parentDOM, prevRender, newComp.lastRender)

		e.currentRenderingComponent = prevComponent

	case *portalNode:
		newPortal := newNode.(*portalNode)
		newPortal.domNode = old.domNode
		newPortal.placeholderNode = old.placeholderNode

		e.patchChildren(old.domNode, []Node{old.child}, []Node{newPortal.child})
	}
}

type diffMaps struct {
	keyed   map[string]Node
	indices map[string]int
}

var diffPool = sync.Pool{
	New: func() any {
		return &diffMaps{
			keyed:   make(map[string]Node),
			indices: make(map[string]int),
		}
	},
}

func (e *engine) patchChildren(parentDOM NativeNode, oldChilds, newChilds []Node) {
	oldMaps := diffPool.Get().(*diffMaps)
	defer func() {
		clear(oldMaps.keyed)
		clear(oldMaps.indices)
		diffPool.Put(oldMaps)
	}()

	oldKeyed := oldMaps.keyed
	oldIndices := oldMaps.indices
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

		var anchorDOM NativeNode
		if i > 0 {
			prevVNode := newChilds[i-1]
			lastPrevDOM := getLastRealDOM(prevVNode)
			if lastPrevDOM != nil {
				anchorDOM = e.renderer.NextSibling(lastPrevDOM)
			}
		} else {
			anchorDOM = e.renderer.FirstChild(parentDOM)
		}

		if matchedOld != nil {
			// update attrs and text
			e.patch(parentDOM, matchedOld, newChild)

			oldIndex, hasOldIndex := oldIndices[key]

			if hasOldIndex {
				if oldIndex < lastPlacedIndex {
					e.moveNode(newChild, parentDOM, anchorDOM)
				} else {
					lastPlacedIndex = oldIndex
				}
			}
		} else {
			// new node... just create
			newDOM := e.createDOM(parentDOM, newChild)

			if anchorDOM != nil {
				e.insertDOMBefore(parentDOM, newDOM, anchorDOM)
			} else {
				e.appendDOMChild(parentDOM, newDOM)
			}
		}
	}

	// remove nodes
	for _, oldChild := range oldKeyed {
		e.destroy(oldChild)
	}

	for i := unkeyedIndex; i < len(oldUnkeyed); i++ {
		e.destroy(oldUnkeyed[i])
	}
}

func (e *engine) appendDOMChild(parentDOM, newDOM NativeNode) {
	if newDOM == nil {
		return
	}

	e.renderer.AppendChild(parentDOM, newDOM)
}

func (e *engine) insertDOMBefore(parentDOM, newDOM, anchorDOM NativeNode) {
	if newDOM == nil {
		return
	}

	e.renderer.InsertBefore(parentDOM, newDOM, anchorDOM)
}

func (e *engine) moveNode(vnode Node, parentDOM, anchorDOM NativeNode) {
	if vnode == nil {
		return
	}

	var actualDom NativeNode
	switch n := vnode.(type) {
	case *Element:
		actualDom = n.domNode
	case *TextNode:
		actualDom = n.domNode
	case *componentNode:
		if isNilNode(n) || n.lastRender == nil {
			return
		}
		e.moveNode(n.lastRender, parentDOM, anchorDOM)
	case *portalNode:
		actualDom = n.placeholderNode

	case *groupNode:
		for _, child := range n.children {
			e.moveNode(child, parentDOM, anchorDOM)
		}
		return
	}

	if actualDom != nil {
		if anchorDOM != nil {
			e.insertDOMBefore(parentDOM, actualDom, anchorDOM)
		} else {
			e.appendDOMChild(parentDOM, actualDom)
		}
	}

}

func (e *engine) patchEvents(domEl NativeNode, oldE *Element, newE *Element) {
	if oldE.listeners == nil && newE.listeners == nil {
		return
	}

	if oldE.listeners == nil {
		oldE.listeners = &listenersInfo{
			events:          make(map[eventInfo]func(Event)),
			activeCallbacks: make(map[eventInfo]func()),
		}
	}

	if newE.listeners == nil {
		newE.listeners = &listenersInfo{}
	}

	oldE.listeners.parentComponent = e.currentRenderingComponent

	for eventInfo, newHandler := range newE.listeners.events {
		if _, exists := oldE.listeners.events[eventInfo]; !exists {
			e.addEventListener(domEl, oldE.listeners, eventInfo)
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

func (e *engine) patchClasses(domNode NativeNode, oldClasses, newClasses string) {
	if oldClasses == newClasses {
		return
	}

	e.renderer.SetAttribute(domNode, "class", newClasses)
}

func (e *engine) patchAttributes(domNode NativeNode, oldAttrs, newAttrs map[string]string) {
	for key, newVal := range newAttrs {
		oldVal, exists := oldAttrs[key]

		if !exists || oldVal != newVal {
			e.renderer.SetAttribute(domNode, key, newVal)
		}
	}

	// remove attributes
	for key := range oldAttrs {
		if _, exists := newAttrs[key]; !exists {
			e.renderer.RemoveAttribute(domNode, key)
		}
	}
}

func (e *engine) destroy(node Node) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *Element:
		if n.listeners != nil {
			for ei, _ := range n.listeners.activeCallbacks {
				destroyEventListeners(n, ei)
			}

			n.listeners.activeCallbacks = make(map[eventInfo]func())
		}

		for _, child := range n.children {
			e.destroy(child)
		}

		if n.domNode != nil {
			e.renderer.Remove(n.domNode)
		}

	case *TextNode:
		if n.domNode != nil {
			e.renderer.Remove(n.domNode)
		}

	case *groupNode:
		for _, child := range n.children {
			e.destroy(child)
		}

	case *componentNode:
		e.unregisterComp(n.comp)

		if d, ok := n.comp.(Destroyer); ok {
			d.OnDestroy()
		}

		if n.lastRender != nil {
			e.destroy(n.lastRender)
		}
	case *portalNode:
		if n.child != nil {
			e.destroy(n.child)
		}

		if n.domNode != nil {
			e.renderer.Remove(n.domNode)
		}

		if n.placeholderNode != nil {
			e.renderer.Remove(n.placeholderNode)
		}
	}
}

func destroyEventListeners(el *Element, ei eventInfo) {
	cleaner := el.listeners.activeCallbacks[ei]
	if cleaner != nil {
		cleaner()
	}
}

func getDOMNode(vnode Node) NativeNode {
	switch n := vnode.(type) {
	case *TextNode:
		return n.domNode
	case *Element:
		return n.domNode
	case *groupNode:
		for _, ch := range n.children {
			dom := getDOMNode(ch)
			if dom != nil {
				return dom
			}
		}

		return nil
	case *componentNode:
		if isNilNode(n) || n.lastRender == nil {
			return nil
		}
		return getDOMNode(n.lastRender)
	case *portalNode:
		return n.placeholderNode

	default:
		return nil
	}
}

func getLastRealDOM(vnode Node) NativeNode {
	if vnode == nil {
		return nil
	}

	switch n := vnode.(type) {
	case *TextNode:
		return n.domNode
	case *Element:
		return n.domNode
	case *groupNode:
		for i := len(n.children) - 1; i >= 0; i-- {
			dom := getLastRealDOM(n.children[i])
			if dom != nil {
				return dom
			}
		}

		return nil
	case *componentNode:
		if isNilNode(n) || n.lastRender == nil {
			return nil
		}
		return getLastRealDOM(n.lastRender)
	case *portalNode:
		return n.placeholderNode
	}

	return nil
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

////////////////////
// GLOBAL
////////////////////

var nina *engine

// schedule component re-render
func Update(c Component) { nina.scheduleUpdate(c) }

// entry point
func Mount(containerID string, root Component) {
	fireInitHooks(nina.renderer)

	nina.rootComponent = root
	nina.rootContainer = nina.renderer.GetElementById(containerID)
	if nina.rootContainer == nil {
		panic(fmt.Sprintf("unknown element id: %s", containerID))
	}

	nina.scheduleUpdate(nil)
}

func WaitNextFrame() <-chan struct{} {
	return nina.renderer.waitNextFrame()
}

func Storage() LocalStorage {
	return nina.storage
}

func Doc() DocumentRenderer {
	return nina.documentRenderer
}

// init hooks

var initHooks []func(Renderer)

func RegisterInitHook(hook func(Renderer)) {
	initHooks = append(initHooks, hook)
}

func fireInitHooks(r Renderer) {
	for _, hook := range initHooks {
		hook(r)
	}
}
