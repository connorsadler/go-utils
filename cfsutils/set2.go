package cfsutils

// Based on this one with some mods: https://github.com/eliben/gogl/blob/main/hashset/hashset.go

// Set2 is a generic set based on a hash table (map).
type Set2[T KeyProducer] struct {
	// map of 'CalcKey' values, to determine if an item is in the set yet
	m map[string]struct{}
	// actual items in the set
	itemsSlice []T
}

type KeyProducer interface {
	CalcKey() string
}

// New creates a new Set.
func NewSet2[T KeyProducer]() *Set2[T] {
	return &Set2[T]{m: make(map[string]struct{})}
}

// InitWith creates a new Set initialized with vals.
func NewSet2With[T KeyProducer](vals ...T) *Set2[T] {
	hs := NewSet2[T]()
	for _, v := range vals {
		hs.Add(v)
	}
	return hs
}

// Add adds a value to the set.
func (hs *Set2[T]) Add(val T) {
	// Only add if not already there
	if _, ok := hs.m[val.CalcKey()]; !ok {
		hs.m[val.CalcKey()] = struct{}{}
		hs.itemsSlice = append(hs.itemsSlice, val)
	}
}

// Contains reports whether the set contains the given value.
func (hs *Set2[T]) Contains(val T) bool {
	_, ok := hs.m[val.CalcKey()]
	return ok
}

// Len returns the size/length of the set - the number of values it contains.
func (hs *Set2[T]) Len() int {
	return len(hs.m)
}

// Delete removes a value from the set; if the value doesn't exist in the
// set, this is a no-op.
func (hs *Set2[T]) Delete(val T) {
	delete(hs.m, val.CalcKey())
}

// returns a slice of the items - the type of each item is as per the Set2's type declaration
func (hs *Set2[T]) AsSlice() []T {
	return hs.itemsSlice
}

// returns a slice of the items - but you can cast each item to another type "I"
// this can be useful for a NewSet2[myintslice]() where we want a slice of []int rather than myintslice - see the test
func AsSliceWithCast[T KeyProducer, I any](hs *Set2[T], cast func(t T) I) []I {
	result := make([]I, 0)
	for _, item := range hs.itemsSlice {
		castedItem := cast(item)
		result = append(result, castedItem)
	}
	return result
}

// func (hs *Set2[T]) AsSliceWithCast(cast func(t T)) []X {
// 	return hs.itemsSlice
// }

// Union returns the set union of hs with other. It creates a new set.
func (hs *Set2[T]) Union(other *Set2[T]) *Set2[T] {
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
func (hs *Set2[T]) Intersection(other *Set2[T]) *Set2[T] {
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
func (hs *Set2[T]) Difference(other *Set2[T]) *Set2[T] {
	// result := New[T]()
	// for v := range hs.m {
	// 	if !other.Contains(v) {
	// 		result.Add(v)
	// 	}
	// }
	// return result
	panic("Not yet supported")
}
