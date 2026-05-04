package nn

type Event interface {
	needUpdate() bool

	PreventUpdate()
	PreventDefault()
	StopPropagation()
	TargetValue() string
	TargetChecked() bool
	Key() string

	CurrentTarget() NativeNode
	Target() NativeNode

	Renderer() Renderer
}
