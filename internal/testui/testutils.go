package testui

import (
	"strings"
	"syscall/js"
	"testing"
)

func GetByRole(role string) js.Value {
	return rootNode().Call("querySelector", `[role="`+role+`"]`)
}

func FireClick(element js.Value) {
	event := js.Global().Get("MouseEvent").New("click", map[string]any{"bubbles": true})
	element.Call("dispatchEvent", event)
}

type DOMAssert struct {
	t  *testing.T
	el js.Value
}

func Expect(t *testing.T) *DOMAssert {
	el := rootNode().Get("firstElementChild")
	if el.IsNull() || el.IsUndefined() {
		t.Fatal("❌ expecte DOM-element as container, but actual null/undefined")
	}

	return &DOMAssert{t: t, el: el}
}

func (a *DOMAssert) Tag(expected string) *DOMAssert {
	actual := strings.ToLower(a.el.Get("tagName").String())
	expected = strings.ToLower(expected)
	if actual != expected {
		a.t.Errorf("Expected tag <%s>, but actual <%s>", expected, actual)
	}
	return a
}

func (a *DOMAssert) Attr(name, expected string) *DOMAssert {
	actual := a.el.Call("getAttribute", name)
	if actual.IsNull() {
		a.t.Errorf("Attr '%s' not found on element <%s>", name, a.el.Get("tagName").String())
		return a
	}
	if actual.String() != expected {
		a.t.Errorf("Attr '%s': expected '%s', actual '%s'", name, expected, actual.String())
	}
	return a
}

func (a *DOMAssert) HasClass(className string) *DOMAssert {
	if !a.el.Get("classList").Call("contains", className).Bool() {
		a.t.Errorf("Element does not have class '%s'. Actual classes: '%s'",
			className, a.el.Get("className").String())
	}
	return a
}

func (a *DOMAssert) Text(expected string) *DOMAssert {
	actual := a.el.Get("textContent").String()
	if actual != expected {
		a.t.Errorf("Expected text '%s', but actual '%s'", expected, actual)
	}
	return a
}

func (a *DOMAssert) ContainsText(substring string) *DOMAssert {
	actual := a.el.Get("textContent").String()
	if !strings.Contains(actual, substring) {
		a.t.Errorf("Element <%s> does not contains text '%s'. Actual text: '%s'",
			a.el.Get("tagName").String(), substring, actual)
	}
	return a
}

func (a *DOMAssert) ChildrenCount(expected int) *DOMAssert {
	actual := a.el.Get("children").Get("length").Int()
	if actual != expected {
		a.t.Errorf("Expected %d children, but actual %d", expected, actual)
	}
	return a
}

func (a *DOMAssert) Child(index int, assertions func(*DOMAssert)) *DOMAssert {
	children := a.el.Get("children")
	if index >= children.Get("length").Int() {
		a.t.Fatalf("Access attempt to child by index %d, but actual has %d",
			index, children.Get("length").Int())
	}

	childEl := children.Index(index)
	assertions(&DOMAssert{t: a.t, el: childEl})

	return a
}
