//go:build !js

package nn

import (
	"fmt"
	"strconv"
	"testing"
)

func BenchmarkPatchChildren_Reorder(b *testing.B) {
	parentDOM := newMockRenderer().CreateElement("div")

	const numItems = 100
	stateA := make([]Node, numItems)
	stateB := make([]Node, numItems)

	for i := 0; i < numItems; i++ {
		key := fmt.Sprintf("item-%d", i)
		nodeA := Div().Key(key).Children(Text(fmt.Sprintf("Text %d", i)))

		stateA[i] = nodeA
		stateB[numItems-1-i] = nodeA
	}

	nina.patchChildren(parentDOM, nil, stateA)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			nina.patchChildren(parentDOM, stateA, stateB)
		} else {
			nina.patchChildren(parentDOM, stateB, stateA)
		}
	}
}

func BenchmarkVDOMCreation_Flat(b *testing.B) {
	var pregenKeys []string

	for i := 0; i < 1000; i++ {
		pregenKeys = append(pregenKeys, strconv.Itoa(i))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		children := make([]AsNode, 1000)
		for j := 0; j < 1000; j++ {
			children[j] = Div().Key(pregenKeys[j]).Children(Text("List Item"))
		}

		_ = Div().Class("container").Children(children...)
	}
}

func BenchmarkVDOMCreation_Deep(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var node Node = Text("Deep core")
		for j := 0; j < 100; j++ {
			node = Div().Class("wrapper").Children(node)
		}
		_ = node
	}
}

func BenchmarkSignal_ReadWrite(b *testing.B) {
	sig := NewSignal[int](0)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sig.Set(i)

		_ = sig.Get(nil)
	}
}

type testComp struct {
	BaseComponent
	n Node
}

func (c *testComp) View() Node {
	return c.n
}

func BenchmarkSignal_Notify100(b *testing.B) {
	sig := NewSignal[int](0)

	for i := 0; i < 100; i++ {
		mockComponent := &testComp{n: Div()}

		_ = sig.Get(mockComponent)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sig.Set(i)
	}
}
