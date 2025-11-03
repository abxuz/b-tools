package mfr

type FilterIndexFunc[T any] = func(i int, v T) (bool, error)

func FilterIndex[T any](list []T, filter FilterIndexFunc[T], dst ...[]T) ([]T, error) {
	var ret []T
	if len(dst) > 0 {
		ret = dst[0]
	} else {
		ret = make([]T, 0, len(list))
	}
	for i, v := range list {
		hit, err := filter(i, v)
		if err != nil {
			return nil, err
		}
		if hit {
			ret = append(ret, v)
		}
	}
	return ret, nil
}

type MustFilterIndexFunc[T any] = func(i int, v T) bool

func MustFilterIndex[T any](list []T, filter MustFilterIndexFunc[T], dst ...[]T) []T {
	ret, _ := FilterIndex(list, func(i int, v T) (bool, error) {
		return filter(i, v), nil
	}, dst...)
	return ret
}

type FilterFunc[T any] = func(v T) (bool, error)

func Filter[T any](list []T, filter FilterFunc[T], dst ...[]T) ([]T, error) {
	return FilterIndex(list, func(_ int, v T) (bool, error) {
		return filter(v)
	}, dst...)
}

type MustFilterFunc[T any] = func(v T) bool

func MustFilter[T any](list []T, filter MustFilterFunc[T], dst ...[]T) []T {
	ret, _ := FilterIndex(list, func(_ int, v T) (bool, error) {
		return filter(v), nil
	}, dst...)
	return ret
}
