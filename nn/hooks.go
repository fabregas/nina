package nn

type Initer interface {
	OnInit()
}

type Mounter interface {
	OnMount()
}

type Destroyer interface {
	OnDestroy()
}
