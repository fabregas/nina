package testui

import (
	"testing"

	"github.com/fabregas/nina/nn"
)

type testCompWithCtx struct {
	nn.BaseComponent

	child nn.Node
}

func (t *testCompWithCtx) View() nn.Node {
	return nn.Div().Children(t.child)
}

func BenchmarkContextResolution(b *testing.B) {
	ctxComp := &testCompWithCtx{}
	nn.ProvideContext(ctxComp, "test")

	var current *testCompWithCtx

	current = ctxComp
	ctxComp.child = nn.Comp(current)
	for i := 0; i < 100; i++ {
		child := &testCompWithCtx{}
		current.child = nn.Comp(child)
		current = child
	}

	leafNode := current

	root.Render(nn.Comp(ctxComp))

	val := nn.GetContext[string](leafNode)
	if val != "test" {
		b.Errorf("invalid context: %s", val)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		val := nn.GetContext[string](leafNode)
		_ = val
	}
}
