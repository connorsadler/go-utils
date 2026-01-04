package cfsutils

// Based on this one with some mods: https://github.com/eliben/gogl/blob/main/hashset/hashset.go

// Set is a generic set based on a hash table (map).
type Set[T comparable] struct {
	m map[T]struct{}
}

// New creates a new Set.
func New[T comparable]() *Set[T] {
	return &Set[T]{m: make(map[T]struct{})}
}

// InitWith creates a new Set initialized with vals.
func InitWith[T comparable](vals ...T) *Set[T] {
	hs := New[T]()
	for _, v := range vals {
		hs.Add(v)
	}
	return hs
}

// Add adds a value to the set.
func (hs *Set[T]) Add(val T) {
	hs.m[val] = struct{}{}
}

// Contains reports whether the set contains the given value.
func (hs *Set[T]) Contains(val T) bool {
	_, ok := hs.m[val]
	return ok
}

// Len returns the size/length of the set - the number of values it contains.
func (hs *Set[T]) Len() int {
	return len(hs.m)
}

// Delete removes a value from the set; if the value doesn't exist in the
// set, this is a no-op.
func (hs *Set[T]) Delete(val T) {
	delete(hs.m, val)
}

func (hs *Set[T]) AsSlice() []T {
	keys := make([]T, 0, len(hs.m))
	for k := range hs.m {
		keys = append(keys, k)
	}
	return keys
}

// Union returns the set union of hs with other. It creates a new set.
func (hs *Set[T]) Union(other *Set[T]) *Set[T] {
	result := New[T]()
	for v := range hs.m {
		result.Add(v)
	}
	for v := range other.m {
		result.Add(v)
	}
	return result
}

// Intersection returns the set intersection of hs with other. It creates a
// new set.
func (hs *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := New[T]()
	for v := range hs.m {
		if other.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

// Difference returns the set difference hs - other. It creates a new set.
func (hs *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := New[T]()
	for v := range hs.m {
		if !other.Contains(v) {
			result.Add(v)
		}
	}
	return result
}
