package mfr

type ReduceIndexFunc[T1, T2 any] = func(i int, v T1, r T2) (T2, error)

func ReduceIndex[T1, T2 any](list []T1, reducer ReduceIndexFunc[T1, T2], initial T2) (T2, error) {
	var err error
	for i, v := range list {
		initial, err = reducer(i, v, initial)
		if err != nil {
			return initial, err
		}
	}
	return initial, nil
}

type MustReduceIndexFunc[T1, T2 any] = func(i int, v T1, r T2) T2

func MustReduceIndex[T1, T2 any](list []T1, reducer MustReduceIndexFunc[T1, T2], initial T2) T2 {
	ret, _ := ReduceIndex(list, func(i int, v T1, r T2) (T2, error) {
		return reducer(i, v, r), nil
	}, initial)
	return ret
}

type ReduceFunc[T1, T2 any] = func(v T1, r T2) (T2, error)

func Reduce[T1, T2 any](list []T1, reducer ReduceFunc[T1, T2], initial T2) (T2, error) {
	return ReduceIndex(list, func(_ int, v T1, r T2) (T2, error) {
		return reducer(v, r)
	}, initial)
}

type MustReduceFunc[T1, T2 any] = func(v T1, r T2) T2

func MustReduce[T1, T2 any](list []T1, reducer MustReduceFunc[T1, T2], initial T2) T2 {
	ret, _ := ReduceIndex(list, func(_ int, v T1, r T2) (T2, error) {
		return reducer(v, r), nil
	}, initial)
	return ret
}
