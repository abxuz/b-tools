package bslice

import "github.com/abxuz/b-tools/bset"

func Unique[Key comparable, Obj any](objs []Obj, getKey func(Obj) Key) bool {
	set := bset.New[Key]()
	for _, obj := range objs {
		key := getKey(obj)
		if set.Has(key) {
			return false
		}
		set.Set(key)
	}
	return true
}

func FindIndex[Obj any](objs []Obj, filter func(Obj) bool) int {
	for i, obj := range objs {
		if match := filter(obj); match {
			return i
		}
	}
	return -1
}
