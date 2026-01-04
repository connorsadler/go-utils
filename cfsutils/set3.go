package cfsutils

// Based on this one with some mods: https://github.com/eliben/gogl/blob/main/hashset/hashset.go

// Set3 is a generic set based on a hash table (map).
type Set3[T any] struct {
	keyProducer func(item T) string
	// map of 'CalcKey' values, to determine if an item is in the set yet
	m map[string]struct{}
	// actual items in the set
	itemsSlice []T
}

//type KeyProducerSet3 func[T any](item T) string

// New creates a new Set.
func NewSet3[T any](kp func(item T) string) *Set3[T] {
	return &Set3[T]{
		keyProducer: kp,
		m:           make(map[string]struct{}),
		itemsSlice:  []T{},
	}
}

// InitWith creates a new Set initialized with vals.
func NewSet3With[T any](kp func(item T) string, vals ...T) *Set3[T] {
	hs := NewSet3[T](kp)
	for _, v := range vals {
		hs.Add(v)
	}
	return hs
}

// Add adds a value to the set.
func (hs *Set3[T]) Add(val T) {
	key := hs.keyProducer(val)
	// Only add if not already there
	if _, ok := hs.m[key]; !ok {
		hs.m[key] = struct{}{}
		hs.itemsSlice = append(hs.itemsSlice, val)
	}
}

// Contains reports whether the set contains the given value.
func (hs *Set3[T]) Contains(val T) bool {
	key := hs.keyProducer(val)
	_, ok := hs.m[key]
	return ok
}

// Len returns the size/length of the set - the number of values it contains.
func (hs *Set3[T]) Len() int {
	return len(hs.m)
}

// Delete removes a value from the set; if the value doesn't exist in the
// set, this is a no-op.
func (hs *Set3[T]) Delete(val T) {
	key := hs.keyProducer(val)
	delete(hs.m, key)
}

// returns a slice of the items - the type of each item is as per the Set3's type declaration
func (hs *Set3[T]) AsSlice() []T {
	return hs.itemsSlice
}

// Union returns the set union of hs with other. It creates a new set.
func (hs *Set3[T]) Union(other *Set3[T]) *Set3[T] {
	// result := New[T]()
	// for v := range hs.m {
	// 	result.Add(v)
	// }
	// for v := range other.m {
	// 	result.Add(v)
	// }
	// return result
	panic("Not yet supported")
}

// Intersection returns the set intersection of hs with other. It creates a
// new set.
func (hs *Set3[T]) Intersection(other *Set3[T]) *Set3[T] {
	// result := New[T]()
	// for v := range hs.m {
	// 	if other.Contains(v) {
	// 		result.Add(v)
	// 	}
	// }
	// return result
	panic("Not yet supported")
}

// Difference returns the set difference hs - other. It creates a new set.
func (hs *Set3[T]) Difference(other *Set3[T]) *Set3[T] {
	// result := New[T]()
	// for v := range hs.m {
	// 	if !other.Contains(v) {
	// 		result.Add(v)
	// 	}
	// }
	// return result
	panic("Not yet supported")
}
