package bset

type Set[T comparable] struct {
	m map[T]struct{}
}

type SetAny = Set[any]
type SetString = Set[string]

func NewSet[T comparable]() *Set[T] {
	return (&Set[T]{}).Init()
}

func NewSetFromSlice[T comparable](list []T) *Set[T] {
	set := NewSet[T]()
	for _, v := range list {
		set.Set(v)
	}
	return set
}

func NewSetAny() *SetAny {
	return NewSet[any]()
}

func NewSetString() *SetString {
	return NewSet[string]()
}

func (s *Set[T]) Init() *Set[T] {
	s.m = make(map[T]struct{})
	return s
}

func (s *Set[T]) Set(v T) {
	s.m[v] = struct{}{}
}

func (s *Set[T]) Remove(v T) {
	delete(s.m, v)
}

func (s *Set[T]) Has(v T) bool {
	_, exists := s.m[v]
	return exists
}

func (s *Set[T]) Size() int {
	return len(s.m)
}

func (s *Set[T]) Range(f func(T) bool) {
	for v := range s.m {
		if ok := f(v); !ok {
			break
		}
	}
}
