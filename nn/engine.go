package nn

import (
	"sync"
	"syscall/js"
)

var nina *engine

type engine struct {
	mu              sync.Mutex
	registry        map[Component]*ComponentNode
	dirtyComponents map[Component]bool
	updateScheduled bool
	renderCallback  js.Func

	rootComponent  Component
	lastGlobalTree *ComponentNode
	rootContainer  js.Value
}

func init() {
	nina = &engine{
		registry:        make(map[Component]*ComponentNode),
		dirtyComponents: make(map[Component]bool),
	}

	nina.renderCallback = js.FuncOf(func(this js.Value, args []js.Value) any {
		nina.performUpdates()
		return nil
	})
}

func (e *engine) registerComp(c Component, node *ComponentNode) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.registry[c] = node
}

func (e *engine) unregisterComp(c Component) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.registry, c)
	delete(e.dirtyComponents, c)
}

func (e *engine) scheduleUpdate(c Component) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if c == nil {
		e.dirtyComponents[e.rootComponent] = true
	} else if _, exists := e.registry[c]; exists {
		e.dirtyComponents[c] = true
	}

	if !e.updateScheduled {
		e.updateScheduled = true
		global.Call("requestAnimationFrame", e.renderCallback)
	}
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
		newTree := C(e.rootComponent)
		patch(e.rootContainer, e.lastGlobalTree, newTree)
		e.lastGlobalTree = newTree

		// just exit, everything is re-rendered
		return
	}

	for comp := range queue {
		e.mu.Lock()
		node, exists := e.registry[comp]
		e.mu.Unlock()

		if !exists || node == nil {
			continue
		}

		// local patch
		newRender := node.comp.View()
		patch(node.parentDOM, node.lastRender, newRender)
		node.lastRender = newRender
	}
}

// schedule component re-render
func Update(c Component) { nina.scheduleUpdate(c) }

// entry point
func Mount(containerID string, root Component) {
	nina.rootComponent = root
	nina.rootContainer = document.Call("getElementById", containerID)

	nina.scheduleUpdate(nil)
}
