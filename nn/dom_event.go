//go:build js && wasm

package nn

import "syscall/js"

type domEvent struct {
	jsEvent js.Value

	renderer *domRenderer

	skipUpdate bool
}

func newDOMEvent(jsEvent js.Value, renderer *domRenderer) *domEvent {
	return &domEvent{
		jsEvent:    jsEvent,
		skipUpdate: false,
		renderer:   renderer,
	}
}
func (e *domEvent) needUpdate() bool {
	return !e.skipUpdate
}
func (e *domEvent) PreventUpdate() {
	e.skipUpdate = true
}

// PreventDefault stops default browser behavior (for example link forward)
func (e *domEvent) PreventDefault() {
	if !e.jsEvent.IsUndefined() && !e.jsEvent.IsNull() {
		e.jsEvent.Call("preventDefault")
	}
}

// StopPropagation stops event "popup-ing" up on DOM tree
func (e *domEvent) StopPropagation() {
	if !e.jsEvent.IsUndefined() && !e.jsEvent.IsNull() {
		e.jsEvent.Call("stopPropagation")
	}
}

func (e *domEvent) CurrentTarget() NativeNode {
	if e.jsEvent.IsUndefined() || e.jsEvent.IsNull() {
		return nil
	}

	ret := e.jsEvent.Get("currentTarget")

	return domNode{ret}
}

func (e *domEvent) Target() NativeNode {
	if e.jsEvent.IsUndefined() || e.jsEvent.IsNull() {
		return nil
	}

	ret := e.jsEvent.Get("target")

	return domNode{ret}
}

func (e *domEvent) TargetValue() string {
	if e.jsEvent.IsUndefined() || e.jsEvent.IsNull() {
		return ""
	}
	target := e.jsEvent.Get("target")
	if target.IsUndefined() || target.IsNull() {
		return ""
	}
	return target.Get("value").String()
}

//func (e *domEvent) Raw() js.Value {
//	return e.jsEvent
//}

func (e *domEvent) TargetChecked() bool {
	if e.jsEvent.IsUndefined() || e.jsEvent.IsNull() {
		return false
	}
	target := e.jsEvent.Get("target")
	if target.IsUndefined() || target.IsNull() {
		return false
	}

	return target.Get("checked").Bool()
}

func (e *domEvent) Key() string {
	var ret js.Value

	if e.jsEvent.IsUndefined() || e.jsEvent.IsNull() {
		return ""
	}

	ret = e.jsEvent.Get("key")
	if ret.IsUndefined() || ret.IsNull() {
		return ""
	}

	return ret.String()
}

func (e *domEvent) Renderer() Renderer {
	return e.renderer
}
