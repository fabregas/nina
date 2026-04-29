package nn

type Signal[T any] struct {
	value       T
	subscribers map[Component]struct{}
}

func NewSignal[T any](initial T) *Signal[T] {
	return &Signal[T]{
		value:       initial,
		subscribers: make(map[Component]struct{}),
	}
}

func (s *Signal[T]) Get(subscriber Component) T {
	if subscriber != nil {
		if _, exists := s.subscribers[subscriber]; !exists {
			s.subscribers[subscriber] = struct{}{}

			subscriber.AddCleanup(func() {
				delete(s.subscribers, subscriber)
			})
		}
	}

	return s.value
}

func (s *Signal[T]) Set(newValue T) {
	s.value = newValue

	for sub := range s.subscribers {
		Update(sub)
	}
}
