package nn

type stateCarrier interface {
	exportState() any
	importState(any)
}

// state structure that customer add to component
type State[T any] struct {
	S *T
}

func (s *State[T]) exportState() any {
	return s.S
}

func (s *State[T]) importState(oldState any) {
	if oldState == nil {
		// first render: allocate memory for new state
		s.S = new(T)
	} else {
		// next renders: just copy pointer from old tree
		s.S = oldState.(*T)
	}
}
