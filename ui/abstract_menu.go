package ui

import (
	"github.com/fabregas/nina/nn"
)

func WrapMenu(menuNode nn.AsNode, onClickCb func(value string)) *nn.Element {
	menuContent, ok := menuNode.AsNode().(*nn.Element)
	if !ok {
		panic("wrapMenu can wrap only Element")
	}

	menuRef := nn.NewRef()

	onContentHover := func(e nn.Event) {
		e.PreventUpdate()

		target := e.Target()
		r := e.Renderer()

		item := r.Closest(target, "[role='option']")
		if item == nil {
			return
		}

		container := e.CurrentTarget()

		current := r.QuerySelector(container, "[data-highlighted]")
		if current != nil && !current.Equal(item) {
			r.RemoveAttribute(current, "data-highlighted")
		}

		if !r.HasAttribute(item, "data-highlighted") {
			r.SetAttribute(item, "data-highlighted", "")
		}
	}

	onMouseDown := func(e nn.Event) {
		e.PreventDefault()
		e.PreventUpdate()
	}

	onMouseLeave := func(e nn.Event) {
		e.PreventUpdate()

		container := e.CurrentTarget()
		r := e.Renderer()

		current := r.QuerySelector(container, "[data-highlighted]")
		if current != nil {
			r.RemoveAttribute(current, "data-highlighted")
		}
	}

	onClick := func(e nn.Event) {
		e.PreventDefault()
		e.PreventUpdate()
		container := e.CurrentTarget()
		r := e.Renderer()

		current := r.QuerySelector(container, "[data-highlighted]")
		if current == nil {
			return
		}

		val := r.GetAttribute(current, "data-value")
		if onClickCb != nil {
			onClickCb(val)
		}
	}

	onKeyDown := func(e nn.Event) {
		key := e.Key()

		switch key {
		case "ArrowDown", "ArrowUp":
		case "Enter":
		default:
			return
		}

		e.PreventDefault()
		e.PreventUpdate()

		r := e.Renderer()
		content := menuRef.Current
		if content == nil {
			return
		}

		items := r.QuerySelectorAll(content, "[role='option']")
		length := len(items)
		if length == 0 {
			return
		}

		currentIndex := -1
		for i, li := range items {
			if r.HasAttribute(li, "data-highlighted") {
				currentIndex = i
				break
			}
		}

		if currentIndex >= 0 && key == "Enter" && onClickCb != nil {
			val := r.GetAttribute(items[currentIndex], "data-value")
			onClickCb(val)
			return
		}

		newIndex := currentIndex
		if key == "ArrowDown" {
			newIndex++
			if newIndex >= length {
				newIndex = 0
			}
		} else if key == "ArrowUp" {
			newIndex--
			if newIndex < 0 {
				newIndex = length - 1
			}
		}

		if currentIndex != -1 {
			r.RemoveAttribute(items[currentIndex], "data-highlighted")
		}

		newItem := items[newIndex]
		r.SetAttribute(newItem, "data-highlighted", "")

		r.ScrollIntoView(newItem, map[string]any{"block": "nearest"})
	}

	return menuContent.
		Ref(menuRef).
		On("pointermove", onContentHover).
		On("mousedown", onMouseDown).
		On("mouseleave", onMouseLeave).
		On("click", onClick).
		OnGlobal("keydown", onKeyDown)
}
