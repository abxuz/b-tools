package mfr

type EveryIndexFunc[T any] = func(i int, v T) (bool, error)

func EveryIndex[T any](list []T, every EveryIndexFunc[T]) (bool, error) {
	for i, v := range list {
		hit, err := every(i, v)
		if err != nil {
			return false, err
		}
		if !hit {
			return false, nil
		}
	}
	return true, nil
}

type MustEveryIndexFunc[T any] = func(i int, v T) bool

func MustEveryIndex[T any](list []T, every MustEveryIndexFunc[T]) bool {
	ret, _ := EveryIndex(list, func(i int, v T) (bool, error) {
		return every(i, v), nil
	})
	return ret
}

type EveryFunc[T any] = func(v T) (bool, error)

func Every[T any](list []T, every EveryFunc[T]) (bool, error) {
	return EveryIndex(list, func(_ int, v T) (bool, error) {
		return every(v)
	})
}

type MustEveryFunc[T any] = func(v T) bool

func MustEvery[T any](list []T, every MustEveryFunc[T]) bool {
	ret, _ := EveryIndex(list, func(_ int, v T) (bool, error) {
		return every(v), nil
	})
	return ret
}
