package bset

type Set[T comparable] map[T]struct{}

func New[T comparable](list ...T) Set[T] {
	s := make(Set[T], len(list))
	return s.Set(list...)
}

func (s Set[T]) Set(list ...T) Set[T] {
	for _, v := range list {
		s[v] = struct{}{}
	}
	return s
}

func (s Set[T]) Has(v T) (has bool) {
	_, has = s[v]
	return
}

func (s Set[T]) Del(v T) Set[T] {
	delete(s, v)
	return s
}

func (s Set[T]) Values() []T {
	list := make([]T, 0, len(s))
	for v := range s {
		list = append(list, v)
	}
	return list
}
