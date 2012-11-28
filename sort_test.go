package lst

import (
	"sort"
	"testing"
)

func TestSortBy(t *testing.T) {
	l := NewFromSlice(elements[:])
	sortedList := SortBy(l, func(x, y Elem) bool {
		return x.(int) < y.(int)
	})

	slice := make([]int, N)
	for k, v := range elements[:] {
		slice[k] = v.(int)
	}
	sort.Ints(slice)

	for k, v := range slice {
		if v != Get(sortedList, k) {
			t.Errorf("Wrong order at index %d", k)
			break
		}
	}
}
