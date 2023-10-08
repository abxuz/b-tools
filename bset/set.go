package bset

import (
	"fmt"
	"strings"
)

type Set[T comparable] struct {
	m map[T]struct{}
}

type SetAny = Set[any]
type SetString = Set[string]

func NewSet[T comparable](vs ...T) *Set[T] {
	return (&Set[T]{}).Init(vs...)
}

func NewSetFromSlice[T comparable](list []T) *Set[T] {
	set := NewSet[T]()
	for _, v := range list {
		set.Set(v)
	}
	return set
}

func NewSetAny(vs ...any) *SetAny {
	return NewSet[any](vs...)
}

func NewSetString(vs ...string) *SetString {
	return NewSet[string](vs...)
}

func (s *Set[T]) Init(vs ...T) *Set[T] {
	s.m = make(map[T]struct{})
	return s.Set(vs...)
}

func (s *Set[T]) Set(vs ...T) *Set[T] {
	for _, v := range vs {
		s.m[v] = struct{}{}
	}
	return s
}

func (s *Set[T]) Remove(v T) *Set[T] {
	delete(s.m, v)
	return s
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

func (s *Set[T]) String() string {
	elements := make([]string, 0)
	for v := range s.m {
		elements = append(elements, fmt.Sprintf("%v", v))
	}
	return "[" + strings.Join(elements, ",") + "]"
}

func (s *Set[T]) ToSlice() []T {
	list := make([]T, 0)
	for v := range s.m {
		list = append(list, v)
	}
	return list
}
