package testui

import (
	"os"
	"testing"

	"github.com/fabregas/nina/nn"
)

var (
	root  *testRoot
	rootN nn.NativeNode
)

func init() {
	nn.RegisterInitHook(func(r nn.Renderer) {
		el := r.CreateElement("div")
		r.SetAttribute(el, "id", "root")
		body := r.QuerySelector(r.RootNode(), "body")
		r.AppendChild(body, el)

		rootN = el
	})
}

func TestMain(m *testing.M) {

	/*
		doc := js.Global().Get("document")
		container := doc.Call("createElement", "div")
		container.Set("id", "root")
		doc.Get("body").Call("appendChild", container)
	*/

	root = &testRoot{}
	nn.Mount("root", root)

	rcode := m.Run()

	os.Exit(rcode)
}

type testRoot struct {
	nn.BaseComponent

	nodes []nn.AsNode
}

func (r *testRoot) Render(nodes ...nn.AsNode) {
	r.nodes = nodes
	nn.Update(nil)
	<-nn.WaitNextFrame()
}

func (r *testRoot) View() nn.Node {
	return nn.Div().Children(r.nodes...)
}
