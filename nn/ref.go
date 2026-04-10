package nn

import "syscall/js"

type Ref struct {
	Current js.Value
}

func NewRef() *Ref {
	return &Ref{Current: js.Undefined()}
}
