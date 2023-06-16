package bmap

type Map[KeyT comparable, ValueT any] map[KeyT]ValueT

func NewMap[KeyT comparable, ValueT any]() Map[KeyT, ValueT] {
	return make(Map[KeyT, ValueT])
}

func NewMapFromSlice[KeyT comparable, ValueT any](list []ValueT, getKey func(ValueT) KeyT) Map[KeyT, ValueT] {
	m := make(Map[KeyT, ValueT])
	for _, v := range list {
		k := getKey(v)
		m[k] = v
	}
	return m
}

func (m Map[KeyT, ValueT]) Keys() []KeyT {
	list := make([]KeyT, len(m))
	i := 0
	for k := range m {
		list[i] = k
		i++
	}
	return list
}

func (m Map[KeyT, ValueT]) Values() []ValueT {
	list := make([]ValueT, len(m))
	i := 0
	for _, v := range m {
		list[i] = v
		i++
	}
	return list
}
