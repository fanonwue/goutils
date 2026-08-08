package dsext

import (
	"iter"
	"reflect"
	"slices"
	"strconv"
	"testing"
)

func sequence[T any](values ...T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, value := range values {
			if !yield(value) {
				return
			}
		}
	}
}

func TestSetConstructorsAndBasicOperations(t *testing.T) {
	set := NewSet[string]()
	if !set.IsEmpty() || set.Len() != 0 {
		t.Fatalf("new set should be empty: %#v", set)
	}
	if capped := NewSetCap[int](4); capped == nil {
		t.Fatal("NewSetCap returned a nil set")
	}

	set.Add("a")
	set.AddAll([]string{"b", "b", "c"})
	set.AddAllSeq(sequence("d", "d"))
	if !setEqual(set, NewSetSlice([]string{"a", "b", "c", "d"})) {
		t.Fatalf("unexpected set contents: %#v", set)
	}
	if !set.Contains("c") || set.Contains("missing") {
		t.Fatal("Contains returned an incorrect result")
	}

	set.Remove("b")
	set.RemoveAll([]string{"a", "missing"})
	set.RemoveAllSeq(sequence("c", "missing"))
	if !setEqual(set, NewSetSlice([]string{"d"})) {
		t.Fatalf("unexpected contents after removal: %#v", set)
	}

	set.AddAllSet(NewSetSlice([]string{"e", "f"}))
	set.RemoveAllSet(NewSetSlice([]string{"e", "missing"}))
	if !setEqual(set, NewSetSlice([]string{"d", "f"})) {
		t.Fatalf("unexpected contents after set operations: %#v", set)
	}

	fromSeq := NewSetSeq(sequence(1, 2, 2, 3))
	if !setEqual(fromSeq, NewSetSlice([]int{1, 2, 3})) {
		t.Fatalf("NewSetSeq returned %#v", fromSeq)
	}
	if !setEqual(NewSetSlice(set.Slice()), set) || !setEqual(collectSet(set.Seq()), set) {
		t.Fatal("Slice or Seq did not contain the set elements")
	}

	set.Clear()
	if !set.IsEmpty() {
		t.Fatal("Clear did not empty the set")
	}
}

func TestSetAlgebra(t *testing.T) {
	a := NewSetSlice([]int{1, 2, 3})
	b := NewSetSlice([]int{3, 4, 5, 6})

	if !setEqual(a.Intersect(b), NewSetSlice([]int{3})) {
		t.Fatal("Intersect returned the wrong set")
	}
	if !setEqual(a.Union(b), NewSetSlice([]int{1, 2, 3, 4, 5, 6})) {
		t.Fatal("Union returned the wrong set")
	}
	if !setEqual(a.Difference(b), NewSetSlice([]int{1, 2, 4, 5, 6})) {
		t.Fatal("Difference returned the wrong symmetric difference")
	}

	var nilSet Set[int]
	if !setEqual(nilSet.Intersect(a), NewSet[int]()) || !setEqual(nilSet.Union(a), a) || !setEqual(nilSet.Difference(a), a) {
		t.Fatal("nil-set algebra returned the wrong result")
	}
	result := nilSet.Difference(a)
	result.Add(99)
	if a.Contains(99) {
		t.Fatal("Difference should return an independent set")
	}
}

func TestSliceFunctions(t *testing.T) {
	if got := Map([]int{1, 2, 3}, func(value int) string { return strconv.Itoa(value * 2) }); !reflect.DeepEqual(got, []string{"2", "4", "6"}) {
		t.Fatalf("Map returned %#v", got)
	}
	if got := Filter([]int{1, 2, 3, 4}, func(value int) bool { return value%2 == 0 }); !reflect.DeepEqual(got, []int{2, 4}) {
		t.Fatalf("Filter returned %#v", got)
	}
	if got := Join([]int{1, 2, 3}, ",", strconv.Itoa); got != "1,2,3" {
		t.Fatalf("Join returned %q", got)
	}
}

func TestIteratorFunctions(t *testing.T) {
	values := sequence(1, 2, 3, 4)
	if got := collect(MapSeq(values, strconv.Itoa)); !reflect.DeepEqual(got, []string{"1", "2", "3", "4"}) {
		t.Fatalf("MapSeq returned %#v", got)
	}
	if got := collect(FilterSeq(values, func(value int) bool { return value > 2 })); !reflect.DeepEqual(got, []int{3, 4}) {
		t.Fatalf("FilterSeq returned %#v", got)
	}
	if got := JoinSeq(values, "-", strconv.Itoa); got != "1-2-3-4" {
		t.Fatalf("JoinSeq returned %q", got)
	}

	seen := []int{}
	for value := range MapSeq(values, func(value int) int { return value * 2 }) {
		seen = append(seen, value)
		if len(seen) == 2 {
			break
		}
	}
	if !reflect.DeepEqual(seen, []int{2, 4}) {
		t.Fatalf("sequence did not stop correctly: %#v", seen)
	}
}

func TestMapFunctions(t *testing.T) {
	m := map[string]int{"one": 1, "two": 2, "three": 3}
	if got := Keys(m); !sameElements(got, []string{"one", "two", "three"}) {
		t.Fatalf("Keys returned %#v", got)
	}
	if got := Values(m); !sameElements(got, []int{1, 2, 3}) {
		t.Fatalf("Values returned %#v", got)
	}
	if got := collect(KeysSeq(m)); !sameElements(got, []string{"one", "two", "three"}) {
		t.Fatalf("KeysSeq returned %#v", got)
	}
	if got := collect(ValuesSeq(m)); !sameElements(got, []int{1, 2, 3}) {
		t.Fatalf("ValuesSeq returned %#v", got)
	}
	if got := ReverseMap(m); !reflect.DeepEqual(got, map[int]string{1: "one", 2: "two", 3: "three"}) {
		t.Fatalf("ReverseMap returned %#v", got)
	}
	if got := FilterMap(m, func(_ string, value int) bool { return value%2 == 1 }); !reflect.DeepEqual(got, map[string]int{"one": 1, "three": 3}) {
		t.Fatalf("FilterMap returned %#v", got)
	}
}

func collect[T any](seq iter.Seq[T]) []T {
	return slices.Collect(seq)
}

func collectSet[T comparable](seq iter.Seq[T]) Set[T] {
	return NewSetSeq(seq)
}

func setEqual[T comparable](a, b Set[T]) bool {
	return reflect.DeepEqual(a, b)
}

func sameElements[T comparable](got, want []T) bool {
	if len(got) != len(want) {
		return false
	}
	counts := make(map[T]int, len(got))
	for _, value := range got {
		counts[value]++
	}
	for _, value := range want {
		counts[value]--
		if counts[value] < 0 {
			return false
		}
	}
	return true
}
