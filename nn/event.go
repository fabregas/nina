package nn

import "syscall/js"

type Event struct {
	jsEvent js.Value
}

// PreventDefault stops default browser behavior (for example link forward)
func (e Event) PreventDefault() {
	if !e.jsEvent.IsUndefined() && !e.jsEvent.IsNull() {
		e.jsEvent.Call("preventDefault")
	}
}

// StopPropagation stops event "popup-ing" up on DOM tree
func (e Event) StopPropagation() {
	if !e.jsEvent.IsUndefined() && !e.jsEvent.IsNull() {
		e.jsEvent.Call("stopPropagation")
	}
}

func (e Event) TargetValue() string {
	if e.jsEvent.IsUndefined() || e.jsEvent.IsNull() {
		return ""
	}
	target := e.jsEvent.Get("target")
	if target.IsUndefined() || target.IsNull() {
		return ""
	}
	return target.Get("value").String()
}

func (e Event) TargetChecked() bool {
	if e.jsEvent.IsUndefined() || e.jsEvent.IsNull() {
		return false
	}
	target := e.jsEvent.Get("target")
	if target.IsUndefined() || target.IsNull() {
		return false
	}

	return target.Get("checked").Bool()
}
