package mfr

type MapFilterIndexFunc[T1, T2 any] = func(i int, v T1) (T2, bool, error)

func MapFilterIndex[T1, T2 any](list []T1, fn MapFilterIndexFunc[T1, T2]) ([]T2, error) {
	ret := make([]T2, 0, len(list))
	for i, v1 := range list {
		v2, hit, err := fn(i, v1)
		if err != nil {
			return nil, err
		}
		if hit {
			ret = append(ret, v2)
		}
	}
	return ret, nil
}

type MustMapFilterIndexFunc[T1, T2 any] = func(i int, v T1) (T2, bool)

func MustMapFilterIndex[T1, T2 any](list []T1, fn MustMapFilterIndexFunc[T1, T2]) []T2 {
	ret, _ := MapFilterIndex(list, func(i int, v1 T1) (T2, bool, error) {
		v2, hit := fn(i, v1)
		return v2, hit, nil
	})
	return ret
}

type MapFilterFunc[T1, T2 any] = func(v T1) (T2, bool, error)

func MapFilter[T1, T2 any](list []T1, fn MapFilterFunc[T1, T2]) ([]T2, error) {
	return MapFilterIndex(list, func(_ int, v1 T1) (T2, bool, error) {
		return fn(v1)
	})
}

type MustMapFilterFunc[T1, T2 any] = func(v T1) (T2, bool)

func MustMapFilter[T1, T2 any](list []T1, fn MustMapFilterFunc[T1, T2]) []T2 {
	ret, _ := MapFilterIndex(list, func(_ int, v1 T1) (T2, bool, error) {
		v2, hit := fn(v1)
		return v2, hit, nil
	})
	return ret
}
