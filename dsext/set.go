package dsext

import "iter"

// Set is a generic set backed by a map. Due to the nature of Go's maps,
// the order of elements in the resulting Set is undefined.
type Set[T comparable] map[T]interface{}

// NewSet creates a new Set backed by a map. Equivalent to:
//
//	make(Set[T])
func NewSet[T comparable]() Set[T] {
	set := make(Set[T])
	return set
}

// NewSetCap creates a new Set backed by a map with the specified capacity. Equivalent to:
//
//	make(Set[T], capacity)
func NewSetCap[T comparable](capacity uint) Set[T] {
	set := make(Set[T], capacity)
	return set
}

// NewSetSeq creates a new Set backed by a map containing the elements of the given sequence. Due to the nature of sequences,
// the final Set size is undefined before the iteration is complete, leading to multiple potential allocations.
func NewSetSeq[T comparable](seq iter.Seq[T]) Set[T] {
	set := make(Set[T])
	set.AddAllSeq(seq)
	return set
}

// NewSetSlice creates a new Set backed by a map containing the elements of the given slice. The resulting Set will have the
// same size as the slice, which means that the allocation will be done only once by specifying the slice capacity hint.
func NewSetSlice[T comparable](elements []T) Set[T] {
	set := make(Set[T], len(elements))
	set.AddAll(elements)
	return set
}

// Add adds the given element to the Set.
func (s Set[T]) Add(t T) {
	s[t] = nil
}

// AddAll adds all elements from the given slice to the Set.
func (s Set[T]) AddAll(other []T) {
	for i := range other {
		s.Add(other[i])
	}
}

// AddAllSeq adds all elements from the given sequence to the Set. Due to the nature of sequences, infinite iterations are possible,
// leading to undefined behavior. Make sure the sequence is finite.
func (s Set[T]) AddAllSeq(other iter.Seq[T]) {
	for v := range other {
		s.Add(v)
	}
}

// AddAllSet adds all elements from the given Set to the Set.
func (s Set[T]) AddAllSet(other Set[T]) {
	for e := range other {
		s.Add(e)
	}
}

// Contains checks if the given element is present in the Set. It will return true if the element is present, false otherwise.
func (s Set[T]) Contains(t T) bool {
	_, ok := (s)[t]
	return ok
}

// Remove removes the given element from the Set.
func (s Set[T]) Remove(t T) {
	delete(s, t)
}

// RemoveAll removes all elements from the Set that are also present in the given slice.
func (s Set[T]) RemoveAll(other []T) {
	for i := range other {
		delete(s, other[i])
	}
}

// RemoveAllSeq removes all elements from the Set that are also present in the given sequence. Due to the nature of sequences, infinite iterations are possible,
// leading to undefined behavior. Make sure the sequence is finite.
func (s Set[T]) RemoveAllSeq(other iter.Seq[T]) {
	for t := range other {
		delete(s, t)
	}
}

// RemoveAllSet removes all elements from the Set that are also present in the given Set.
func (s Set[T]) RemoveAllSet(other Set[T]) {
	for t := range other {
		delete(s, t)
	}
}

// Intersect returns a new Set containing the elements that are present in both the current Set and the given one. Nil sets are considered empty.
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	if len(s) == 0 || len(other) == 0 {
		return make(Set[T])
	}

	outer, inner := outerAndInnerSet(s, other)

	// The maximum size of the intersected Set is equal to the smaller Set, which will always be the outer one.
	// So it makes sense to set the capacity hint to the length of the outer Set.
	intersect := make(Set[T], len(outer))
	for v := range outer {
		if inner.Contains(v) {
			intersect.Add(v)
		}
	}
	return intersect
}

// Union returns a new Set containing the elements that are present in either the current Set or the given one. Nil sets are considered empty.
func (s Set[T]) Union(other Set[T]) Set[T] {
	union := make(Set[T], len(s)+len(other))
	union.AddAllSet(s)
	union.AddAllSet(other)
	return union
}

// Difference returns a new Set containing the elements that are present in the current Set but not in the given one. Nil sets are considered empty.
func (s Set[T]) Difference(other Set[T]) Set[T] {
	if len(s) == 0 {
		return s
	}
	if len(other) == 0 {
		return other
	}

	outer, inner := outerAndInnerSet(s, other)

	difference := make(Set[T])
	for v := range outer {
		if !inner.Contains(v) {
			difference.Add(v)
		}
	}
	return difference
}

// Clear removes all elements from the Set.
func (s Set[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Len returns the number of elements in the Set. Equivalent to:
//
//	len(s)
func (s Set[T]) Len() int {
	return len(s)
}

// IsEmpty returns true if the Set contains no elements, false otherwise. Equivalent to:
//
//	len(s) == 0
func (s Set[T]) IsEmpty() bool {
	return len(s) == 0
}

// Slice returns a slice containing all elements in the Set. Equivalent to:
//
//	dsext.Keys(s)
func (s Set[T]) Slice() []T {
	return Keys(s)
}

// Seq returns a sequence containing all elements in the Set. Equivalent to:
//
//	dsext.KeysSeq(s)
func (s Set[T]) Seq() iter.Seq[T] {
	return KeysSeq(s)
}

// outerAndInnerSet returns the outer and inner Set of the given two. The outer Set is always the smaller one.
func outerAndInnerSet[T comparable](a, b Set[T]) (Set[T], Set[T]) {
	outer := a
	inner := b
	// The outer set should be the smaller one, as the Contains() method is O(1)
	if len(outer) > len(inner) {
		tmp := outer
		outer = inner
		inner = tmp
	}
	return outer, inner
}
