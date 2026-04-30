package testui

import (
	"testing"

	"github.com/fabregas/nina/nn"
)

func TestGroupBase(t *testing.T) {
	// render just one div
	root.Render(nn.Div().Text("test"))

	Expect(t).
		ChildrenCount(1).
		Child(0, func(el *DOMAssert) {
			el.Tag("div").Text("test")
		})

	// re-render with group
	root.Render(
		nn.Div().Attr("role", "begin"),
		nn.Group(
			nn.Div().Attr("role", "g1"),
			nn.Div().Attr("role", "g2"),
		),
		nn.Div().Attr("role", "end"),
	)

	Expect(t).
		ChildrenCount(4).
		Child(0, func(el *DOMAssert) {
			el.Tag("div").Text("").Attr("role", "begin")
		}).
		Child(1, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "g1")
		}).
		Child(2, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "g2")
		}).
		Child(3, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "end")
		})

	// add span between divs
	root.Render(
		nn.Div().Attr("role", "begin"),
		nn.Span().Attr("role", "second"),
		nn.Group(
			nn.Div().Attr("role", "g1"),
			nn.Div().Attr("role", "g2"),
		),
		nn.Div().Attr("role", "end"),
	)

	Expect(t).
		ChildrenCount(5).
		Child(0, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "begin")
		}).
		Child(1, func(el *DOMAssert) {
			el.Tag("span").Attr("role", "second")
		}).
		Child(2, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "g1")
		}).
		Child(3, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "g2")
		}).
		Child(4, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "end")
		})

	// add span into group
	root.Render(
		nn.Div().Attr("role", "begin"),
		nn.Group(
			nn.Div().Attr("role", "g1"),
			nn.Span().Attr("role", "ss"),
			nn.Div().Attr("role", "g2"),
		),
		nn.Div().Attr("role", "end"),
	)

	Expect(t).
		ChildrenCount(5).
		Child(0, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "begin")
		}).
		Child(1, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "g1")
		}).
		Child(2, func(el *DOMAssert) {
			el.Tag("span").Attr("role", "ss")
		}).
		Child(3, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "g2")
		}).
		Child(4, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "end")
		})

	// reorder group
	root.Render(
		nn.Div().Attr("role", "begin"),
		nn.Group(
			nn.Div().Attr("role", "g1.0"),
			nn.Div().Attr("role", "g1"),
			nn.Div().Attr("role", "g1.1"),
			nn.Div().Attr("role", "g2"),
			nn.Div().Attr("role", "g2.1"),
		),
		nn.Div().Attr("role", "end"),
	)

	Expect(t).
		ChildrenCount(7).
		Child(0, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "begin")
		}).
		Child(1, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "g1.0")
		}).
		Child(2, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "g1")
		}).
		Child(3, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "g1.1")
		}).
		Child(4, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "g2")
		}).
		Child(5, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "g2.1")
		}).
		Child(6, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "end")
		})

}

func TestRerenderKeyed(t *testing.T) {
	root.Render(
		nn.Div().Attr("role", "list").ID("list").Children(
			nn.Div().Key("g1").Attr("role", "g1").Class("list-item"),
			nn.Div().Key("g2").Attr("role", "g2").Class("list-item"),
		),
		nn.Div().Attr("role", "end"),
	)

	Expect(t).
		ChildrenCount(2).
		Child(0, func(el *DOMAssert) {
			el.Tag("div").
				Attr("role", "list").
				ChildrenCount(2).
				Child(0, func(el *DOMAssert) {
					el.Tag("div").Attr("role", "g1").HasClass("list-item")
				}).
				Child(1, func(el *DOMAssert) {
					el.Tag("div").Attr("role", "g2").HasClass("list-item")
				})
		}).
		Child(1, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "end")
		})

	listNode := rootNode().Call("querySelector", "#list")

	originalDOM_g1 := listNode.Get("children").Index(0)
	originalDOM_g2 := listNode.Get("children").Index(1)

	// add div as first item in group
	root.Render(
		nn.Div().Attr("role", "list").ID("list").Children(
			nn.Div().Key("g0").Attr("role", "g0").Class("list-item"),
			nn.Div().Key("g1").Attr("role", "g1").Class("list-item"),
			nn.Div().Key("g2").Attr("role", "g2").Class("list-item"),
		),
		nn.Div().Attr("role", "end"),
	)

	Expect(t).
		ChildrenCount(2).
		Child(0, func(el *DOMAssert) {
			el.Tag("div").
				Attr("role", "list").
				ChildrenCount(3).
				Child(0, func(el *DOMAssert) {
					el.Tag("div").Attr("role", "g0").HasClass("list-item")
				}).
				Child(1, func(el *DOMAssert) {
					el.Tag("div").Attr("role", "g1").HasClass("list-item")
				}).
				Child(2, func(el *DOMAssert) {
					el.Tag("div").Attr("role", "g2").HasClass("list-item")
				})
		}).
		Child(1, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "end")
		})

	listNode = rootNode().Call("querySelector", "#list")
	newDOM_g1 := listNode.Get("children").Index(1)
	newDOM_g2 := listNode.Get("children").Index(2)

	isSameG1 := newDOM_g1.Call("isSameNode", originalDOM_g1).Bool()
	if !isSameG1 {
		t.Fatal("❌ Key error: Node 'g1' was recreated or mutated")
	}
	isSameG2 := newDOM_g2.Call("isSameNode", originalDOM_g2).Bool()
	if !isSameG2 {
		t.Fatal("❌ Key error: Node 'g2' was recreated or mutated")
	}

}

func TestGroupKeyed(t *testing.T) {
	root.Render(
		nn.Div().Key("lst").Attr("role", "list").ID("list").Children(
			nn.Group(
				nn.Div().Key("g1").Attr("role", "g1").Class("list-item"),
				nn.Div().Key("g2").Attr("role", "g2").Class("list-item"),
			),
		),
		nn.Div().Attr("role", "end"),
	)

	Expect(t).
		ChildrenCount(2).
		Child(0, func(el *DOMAssert) {
			el.Tag("div").
				Attr("role", "list").
				ChildrenCount(2).
				Child(0, func(el *DOMAssert) {
					el.Tag("div").Attr("role", "g1").HasClass("list-item")
				}).
				Child(1, func(el *DOMAssert) {
					el.Tag("div").Attr("role", "g2").HasClass("list-item")
				})
		}).
		Child(1, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "end")
		})

	listNode := rootNode().Call("querySelector", "#list")

	originalDOM_g1 := listNode.Get("children").Index(0)
	originalDOM_g2 := listNode.Get("children").Index(1)

	// add div as first item in group
	root.Render(
		nn.Div().Key("lst").Attr("role", "list").ID("list").Children(
			nn.Group(
				nn.Div().Key("g0").Attr("role", "g0").Class("list-item"),
				nn.Div().Key("g1").Attr("role", "g1").Class("list-item"),
				nn.Div().Key("g2").Attr("role", "g2").Class("list-item"),
			),
		),
		nn.Div().Attr("role", "end"),
	)

	Expect(t).
		ChildrenCount(2).
		Child(0, func(el *DOMAssert) {
			el.Tag("div").
				Attr("role", "list").
				ChildrenCount(3).
				Child(0, func(el *DOMAssert) {
					el.Tag("div").Attr("role", "g0").HasClass("list-item")
				}).
				Child(1, func(el *DOMAssert) {
					el.Tag("div").Attr("role", "g1").HasClass("list-item")
				}).
				Child(2, func(el *DOMAssert) {
					el.Tag("div").Attr("role", "g2").HasClass("list-item")
				})
		}).
		Child(1, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "end")
		})

	listNode = rootNode().Call("querySelector", "#list")
	newDOM_g1 := listNode.Get("children").Index(1)
	newDOM_g2 := listNode.Get("children").Index(2)

	isSameG1 := newDOM_g1.Call("isSameNode", originalDOM_g1).Bool()
	if !isSameG1 {
		t.Fatal("❌ Key error: Node 'g1' was recreated or mutated")
	}
	isSameG2 := newDOM_g2.Call("isSameNode", originalDOM_g2).Bool()
	if !isSameG2 {
		t.Fatal("❌ Key error: Node 'g2' was recreated or mutated")
	}

	// add span before group
	root.Render(
		nn.Span().Key("s"),
		nn.Div().Key("lst").Attr("role", "list").ID("list").Children(
			nn.Group(
				nn.Div().Key("g0").Attr("role", "g0").Class("list-item"),
				nn.Div().Key("g1").Attr("role", "g1").Class("list-item"),
				nn.Div().Key("g2").Attr("role", "g2").Class("list-item"),
			),
		),
		nn.Div().Key("0").Attr("role", "end"),
	)

	Expect(t).
		ChildrenCount(3).
		Child(0, func(el *DOMAssert) {
			el.Tag("span")
		}).
		Child(1, func(el *DOMAssert) {
			el.Tag("div").
				Attr("role", "list").
				ChildrenCount(3).
				Child(0, func(el *DOMAssert) {
					el.Tag("div").Attr("role", "g0").HasClass("list-item")
				}).
				Child(1, func(el *DOMAssert) {
					el.Tag("div").Attr("role", "g1").HasClass("list-item")
				}).
				Child(2, func(el *DOMAssert) {
					el.Tag("div").Attr("role", "g2").HasClass("list-item")
				})
		}).
		Child(2, func(el *DOMAssert) {
			el.Tag("div").Attr("role", "end")
		})

	listNode = rootNode().Call("querySelector", "#list")
	newDOM_g1 = listNode.Get("children").Index(1)
	newDOM_g2 = listNode.Get("children").Index(2)

	isSameG1 = newDOM_g1.Call("isSameNode", originalDOM_g1).Bool()
	if !isSameG1 {
		t.Fatal("❌ Key error: Node 'g1' was recreated or mutated")
	}
	isSameG2 = newDOM_g2.Call("isSameNode", originalDOM_g2).Bool()
	if !isSameG2 {
		t.Fatal("❌ Key error: Node 'g2' was recreated or mutated")
	}
}
