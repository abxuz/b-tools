package mfr

type RangeIndexFunc[T any] = func(i int, v T) error

func RangeIndex[T any](list []T, rangeFn RangeIndexFunc[T]) error {
	for i, v := range list {
		if err := rangeFn(i, v); err != nil {
			return err
		}
	}
	return nil
}

type RangeFunc[T any] = func(v T) error

func Range[T any](list []T, rangeFn RangeFunc[T]) error {
	for _, v := range list {
		if err := rangeFn(v); err != nil {
			return err
		}
	}
	return nil
}

type BulkIndexRangeFunc[T any] = func(list []T, batch int) error

func BulkIndexRange[T any](list []T, size int, fn BulkIndexRangeFunc[T]) error {
	var err error
	for batch := 1; ; batch++ {
		l := len(list)
		if l == 0 {
			break
		}

		p := min(l, size)
		err = fn(list[:p], batch)
		if err != nil {
			break
		}

		list = list[p:]
	}
	return err
}

type BulkRangeFunc[T any] = func(list []T) error

func BulkRange[T any](list []T, size int, fn BulkRangeFunc[T]) error {
	var err error
	for {
		l := len(list)
		if l == 0 {
			break
		}

		p := min(l, size)
		err = fn(list[:p])
		if err != nil {
			break
		}

		list = list[p:]
	}
	return err
}

func BulkSplit[T any](list []T, size int) [][]T {
	bulks := make([][]T, 0, len(list)/size+1)
	BulkRange(list, size, func(l []T) error {
		bulks = append(bulks, l)
		return nil
	})
	return bulks
}
