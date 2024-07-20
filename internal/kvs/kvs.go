package kvs

import "cmp"

type FilterFn[K cmp.Ordered, T cmp.Ordered] func(Pair[K, T], int) bool

type Pairs[K cmp.Ordered, T cmp.Ordered] interface {
	Append(Pair[K, T])
	Keys() []K
	Values() []T
	Filter(FilterFn[K, T]) Pairs[K, T]
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

type Pair[K cmp.Ordered, T cmp.Ordered] struct {
	Key K
	Val T
}

func KeySortedPairs[K cmp.Ordered, T cmp.Ordered](kvMap map[K]T) Pairs[K, T] {
	pairs := make(keySortedPairs[K, T], len(kvMap))

	i := 0
	for k, v := range kvMap {
		pairs[i] = Pair[K, T]{Key: k, Val: v}
		i++
	}

	return pairs
}

func ValueSortedPairs[K cmp.Ordered, T cmp.Ordered](kvMap map[K]T) Pairs[K, T] {
	pairs := make(valueSortedPairs[K, T], len(kvMap))

	i := 0
	for k, v := range kvMap {
		pairs[i] = Pair[K, T]{Key: k, Val: v}
		i++
	}

	return pairs
}

type keySortedPairs[K cmp.Ordered, T cmp.Ordered] []Pair[K, T]

func (p keySortedPairs[K, T]) Append(pair Pair[K, T]) {
	p = append(p, pair)
}

func (p keySortedPairs[K, T]) Keys() []K {
	keys := make([]K, p.Len())
	for i, pair := range p {
		keys[i] = pair.Key
	}

	return keys
}
func (p keySortedPairs[K, T]) Values() []T {
	vals := make([]T, p.Len())
	for i, pair := range p {
		vals[i] = pair.Val
	}

	return vals
}
func (p keySortedPairs[K, T]) Filter(filter FilterFn[K, T]) Pairs[K, T] {
	res := make(keySortedPairs[K, T], 0, p.Len())

	for i, p := range p {
		if !filter(p, i) {
			res = append(res, p)
		}
	}

	return res
}
func (p keySortedPairs[K, T]) Len() int           { return len(p) }
func (p keySortedPairs[K, T]) Less(i, j int) bool { return p[i].Key < p[j].Key }
func (p keySortedPairs[K, T]) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type valueSortedPairs[K cmp.Ordered, T cmp.Ordered] []Pair[K, T]

func (p valueSortedPairs[K, T]) Append(pair Pair[K, T]) {
	p = append(p, pair)
}

func (p valueSortedPairs[K, T]) Keys() []K {
	keys := make([]K, p.Len())
	for i, pair := range p {
		keys[i] = pair.Key
	}

	return keys
}
func (p valueSortedPairs[K, T]) Values() []T {
	vals := make([]T, p.Len())
	for i, pair := range p {
		vals[i] = pair.Val
	}

	return vals
}
func (p valueSortedPairs[K, T]) Filter(filter FilterFn[K, T]) Pairs[K, T] {
	res := make(valueSortedPairs[K, T], 0, p.Len())

	for i, p := range p {
		if !filter(p, i) {
			res = append(res, p)
		}
	}

	return res
}
func (p valueSortedPairs[K, T]) Len() int           { return len(p) }
func (p valueSortedPairs[K, T]) Less(i, j int) bool { return p[i].Val < p[j].Val }
func (p valueSortedPairs[K, T]) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
