package ui

import (
	"github.com/fabregas/nina/nn"
)

func WrapMenu(menuContent *nn.Element, onClickCb func(value string)) *nn.Element {
	onContentHover := func(e nn.Event) {
		e.PreventUpdate()

		target := e.Target()

		item := target.Call("closest", "[role='option']")
		if item.IsNull() || item.IsUndefined() {
			return
		}

		container := e.CurrentTarget()

		current := container.Call("querySelector", "[data-highlighted]")
		if !current.IsNull() && !current.Equal(item) {
			current.Call("removeAttribute", "data-highlighted")
		}

		if !item.Call("hasAttribute", "data-highlighted").Bool() {
			item.Call("setAttribute", "data-highlighted", "")
		}
	}

	onMouseDown := func(e nn.Event) {
		e.PreventDefault()
		e.PreventUpdate()
	}

	onMouseLeave := func(e nn.Event) {
		e.PreventUpdate()

		container := e.CurrentTarget()

		current := container.Call("querySelector", "[data-highlighted]")
		if !current.IsNull() {
			current.Call("removeAttribute", "data-highlighted")
		}
	}

	onClick := func(e nn.Event) {
		e.PreventDefault()
		e.PreventUpdate()
		container := e.CurrentTarget()
		current := container.Call("querySelector", "[data-highlighted]")
		if current.IsNull() {
			return
		}

		val := current.Call("getAttribute", "value").String()
		if onClickCb != nil {
			onClickCb(val)
		}
	}

	menuRef := nn.NewRef()
	menuContent.Ref(menuRef)

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

		content := menuRef.Current
		if content.IsUndefined() || content.IsNull() {
			return
		}

		items := content.Call("querySelectorAll", "[role='option']")
		length := items.Get("length").Int()
		if length == 0 {
			return
		}

		currentIndex := -1
		for i := 0; i < length; i++ {
			if items.Index(i).Call("hasAttribute", "data-highlighted").Bool() {
				currentIndex = i
				break
			}
		}

		if currentIndex >= 0 && key == "Enter" {
			val := items.Index(currentIndex).Call("getAttribute", "value").String()
			if onClickCb != nil {
				onClickCb(val)
			}
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
			items.Index(currentIndex).Call("removeAttribute", "data-highlighted")
		}

		newItem := items.Index(newIndex)
		newItem.Call("setAttribute", "data-highlighted", "")

		scrollOptions := map[string]any{"block": "nearest"}
		newItem.Call("scrollIntoView", scrollOptions)
	}

	return menuContent.
		On("pointermove", onContentHover).
		On("mousedown", onMouseDown).
		On("mouseleave", onMouseLeave).
		On("click", onClick).
		OnGlobal("keydown", onKeyDown)
}
