package testui

import (
	"os"
	"syscall/js"
	"testing"

	"github.com/fabregas/nina/nn"
)

var (
	root *testRoot
)

func TestMain(m *testing.M) {
	doc := js.Global().Get("document")
	container := doc.Call("createElement", "div")
	container.Set("id", "root")
	doc.Get("body").Call("appendChild", container)

	root = &testRoot{}
	nn.Mount("root", root)

	rcode := m.Run()

	os.Exit(rcode)
}

func rootNode() js.Value {
	return js.Global().Get("document").Call("getElementById", "root")
}

type testRoot struct {
	nn.BaseComponent

	nodes []nn.AsNode
}

func (r *testRoot) Render(nodes ...nn.AsNode) {
	r.nodes = nodes
	nn.Update(nil)
	<-WaitNextFrame()
}

func (r *testRoot) View() nn.Node {
	return nn.Div().Children(r.nodes...)
}

func WaitNextFrame() <-chan struct{} {
	ch := make(chan struct{})

	js.Global().Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) any {
		close(ch)
		return nil
	}))

	return ch
}
