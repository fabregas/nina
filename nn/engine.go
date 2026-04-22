package nn

import (
	"fmt"
	"sync"
	"syscall/js"
)

var nina *engine

type engine struct {
	mu              sync.Mutex
	registry        map[Component]*componentNode
	dirtyComponents map[Component]bool
	mountQueue      []func()
	updateScheduled bool
	renderCallback  js.Func

	rootComponent  Component
	lastGlobalTree *componentNode
	rootContainer  js.Value
}

func init() {
	nina = &engine{
		registry:        make(map[Component]*componentNode),
		dirtyComponents: make(map[Component]bool),
	}

	nina.renderCallback = js.FuncOf(func(this js.Value, args []js.Value) any {
		nina.performUpdates()
		return nil
	})
}

func (e *engine) registerComp(c Component, node *componentNode) {
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
	} else {
		fmt.Printf("unknown component: %T\n", c)
	}

	if !e.updateScheduled {
		e.updateScheduled = true
		global.Call("requestAnimationFrame", e.renderCallback)
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
		// fmt.Println("[RENDER] FULL")
		newTree := Comp(e.rootComponent)
		patch(e.rootContainer, e.lastGlobalTree, newTree)
		e.lastGlobalTree = newTree

		// just exit, everything is re-rendered
	} else {
		for comp := range queue {
			//fmt.Printf("[RENDER] %T\n", comp)

			e.mu.Lock()
			node, exists := e.registry[comp]
			e.mu.Unlock()

			if !exists || node == nil {
				continue
			}

			currentRenderingComponent = node.comp

			// local patch
			newRender := node.comp.View()
			patch(node.parentDOM, node.lastRender, newRender)
			node.lastRender = newRender
		}
	}

	for _, onMountCb := range e.mountQueue {
		onMountCb()
	}
	e.mountQueue = nil
}

// schedule component re-render
func Update(c Component) { nina.scheduleUpdate(c) }

// entry point
func Mount(containerID string, root Component) {
	initHistory()
	initStorageListener()

	nina.rootComponent = root
	nina.rootContainer = document.Call("getElementById", containerID)
	if nina.rootContainer.IsNull() {
		panic(fmt.Sprintf("unknown element id: %s", containerID))
	}

	nina.scheduleUpdate(nil)
}

func WaitNextFrame() <-chan struct{} {
	ch := make(chan struct{})

	var cb js.Func

	cb = js.FuncOf(func(this js.Value, args []js.Value) any {
		defer cb.Release()

		close(ch)
		return nil
	})

	js.Global().Call("requestAnimationFrame", cb)

	return ch
}

func WaitForPaint() <-chan struct{} {
	ch := make(chan struct{})

	var cb1, cb2 js.Func

	cb2 = js.FuncOf(func(this js.Value, args []js.Value) any {
		defer cb2.Release()
		close(ch)
		return nil
	})

	cb1 = js.FuncOf(func(this js.Value, args []js.Value) any {
		defer cb1.Release()
		js.Global().Call("requestAnimationFrame", cb2)
		return nil
	})

	js.Global().Call("requestAnimationFrame", cb1)

	return ch
}
