package mfr

type AnyIndexFunc[T any] = func(i int, v T) (bool, error)

func AnyIndex[T any](list []T, anyFn AnyIndexFunc[T]) (bool, error) {
	for i, v := range list {
		hit, err := anyFn(i, v)
		if err != nil {
			return false, err
		}
		if hit {
			return true, nil
		}
	}
	return false, nil
}

type MustAnyIndexFunc[T any] = func(i int, v T) bool

func MustAnyIndex[T any](list []T, anyFn MustAnyIndexFunc[T]) bool {
	ret, _ := AnyIndex(list, func(i int, v T) (bool, error) {
		return anyFn(i, v), nil
	})
	return ret
}

type AnyFunc[T any] = func(v T) (bool, error)

func Any[T any](list []T, anyFn AnyFunc[T]) (bool, error) {
	return AnyIndex(list, func(_ int, v T) (bool, error) {
		return anyFn(v)
	})
}

type MustAnyFunc[T any] = func(v T) bool

func MustAny[T any](list []T, anyFn MustAnyFunc[T]) bool {
	ret, _ := AnyIndex(list, func(_ int, v T) (bool, error) {
		return anyFn(v), nil
	})
	return ret
}
