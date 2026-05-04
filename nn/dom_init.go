//go:build js && wasm

package nn

func init() {
	nina = &engine{
		registry:        make(map[Component]*componentNode),
		dirtyComponents: make(map[Component]bool),
	}

	nina.renderer = newDomRenderer()
	nina.storage = newDomStorage()
	nina.reqDomUpdate, _ = nina.renderer.initRequestAnimationFrame(nina.performUpdates)
}
