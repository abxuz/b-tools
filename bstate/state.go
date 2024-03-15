package bstate

type State[T any] struct {
	v        T
	watchers map[string]func(prev, cur T)
}

func NewState[T any]() *State[T] {
	return &State[T]{
		watchers: make(map[string]func(prev T, cur T)),
	}
}

func (s *State[T]) Get() T {
	return s.v
}

func (s *State[T]) Set(v T) {
	prev := s.v
	s.v = v
	for _, f := range s.watchers {
		f(prev, v)
	}
}

func (s *State[T]) Watch(name string, f func(prev, cur T)) {
	s.watchers[name] = f
}
