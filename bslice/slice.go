package bslice

import "github.com/xbugio/b-tools/bset"

func Unique(size int, getKey func(i int) any) bool {
	set := bset.NewSetAny()
	for i := 0; i < size; i++ {
		key := getKey(i)
		if set.Has(key) {
			return false
		}
		set.Set(key)
	}
	return true
}

func FindIndex(size int, filter func(int) bool) int {
	for i := 0; i < size; i++ {
		if match := filter(i); match {
			return i
		}
	}
	return -1
}
