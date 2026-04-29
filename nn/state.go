package nn

type stateCarrier interface {
	exportState() any
	importState(any)
	setUpdater(trigger func())
}

// state structure that customer add to component
type State[T any] struct {
	Data *T

	updater func()
}

func (s *State[T]) setUpdater(trigger func()) {
	s.updater = trigger
}

func (s *State[T]) Update() {
	if s.updater != nil {
		s.updater()
	}
}

func (s *State[T]) exportState() any {
	return s
}

func (s *State[T]) importState(oldState any) {
	if oldState == nil {
		// first render: allocate memory for new state
		if s.Data == nil {
			s.Data = new(T)
		}
	} else {
		// next renders: just copy pointer from old tree
		old := oldState.(*State[T])
		s.Data = old.Data
		old.updater = s.updater // remap new updater func to old one (for anync calls)
	}
}
