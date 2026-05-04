package nn

type Ref struct {
	Current NativeNode

	Renderer Renderer
}

func NewRef() *Ref {
	return &Ref{}
}
