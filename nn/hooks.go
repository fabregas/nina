package nn

type Mounter interface {
	OnMount()
}

type Destroyer interface {
	OnDestroy()
}
