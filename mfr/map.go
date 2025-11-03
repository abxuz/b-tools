package mfr

type MapIndexFunc[T1, T2 any] = func(i int, v T1) (T2, error)

func MapIndex[T1, T2 any](list []T1, mapper MapIndexFunc[T1, T2], dst ...[]T2) ([]T2, error) {
	var ret []T2
	if len(dst) > 0 {
		ret = dst[0]
	} else {
		ret = make([]T2, 0, len(list))
	}
	for i, v1 := range list {
		v2, err := mapper(i, v1)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v2)
	}
	return ret, nil
}

type MustMapIndexFunc[T1, T2 any] = func(i int, v T1) T2

func MustMapIndex[T1, T2 any](list []T1, mapper MustMapIndexFunc[T1, T2], dst ...[]T2) []T2 {
	ret, _ := MapIndex(list, func(i int, v T1) (T2, error) {
		return mapper(i, v), nil
	}, dst...)
	return ret
}

type MapFunc[T1, T2 any] = func(v T1) (T2, error)

func Map[T1, T2 any](list []T1, mapper MapFunc[T1, T2], dst ...[]T2) ([]T2, error) {
	return MapIndex(list, func(_ int, v T1) (T2, error) {
		return mapper(v)
	}, dst...)
}

type MustMapFunc[T1, T2 any] = func(v T1) T2

func MustMap[T1, T2 any](list []T1, mapper MustMapFunc[T1, T2], dst ...[]T2) []T2 {
	ret, _ := MapIndex(list, func(_ int, v T1) (T2, error) {
		return mapper(v), nil
	}, dst...)
	return ret
}

func NewMapFromSlice[T1 any, T2 comparable](list []T1, getKey func(v T1) T2) map[T2]T1 {
	m := make(map[T2]T1, len(list))
	for _, v := range list {
		m[getKey(v)] = v
	}
	return m
}

func MapValues[Map ~map[K]V, K comparable, V any](m Map) []V {
	ret := make([]V, 0, len(m))
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}
