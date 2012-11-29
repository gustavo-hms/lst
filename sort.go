package lst

import (
	"sort"
)

/*
 * To use package sort, one needs a structure conforming its interface. The Len 
 * and Swap methods are trivial, but Less needs to be given by the user
 */

type sortable struct {
	elements []Elem
	less     func(Elem, Elem) bool
}

func (s *sortable) Len() int {
	return len(s.elements)
}

func (s *sortable) Swap(i, j int) {
	s.elements[i], s.elements[j] = s.elements[j], s.elements[i]
}

func (s *sortable) Less(i, j int) bool {
	// The slice needs to be in reverse order. So, the negation
	return !s.less(s.elements[i], s.elements[j])
}

/*
 * Sorts a list based on a function returning true if the “x” element is lesser 
 * then “y”
 */
func SortBy(l *List, less func(x, y Elem) bool) *List {
	s := new(sortable)
	s.elements = make([]Elem, Len(l))
	for i := 0; i < Len(l); i++ {
		s.elements[i] = Get(l, i)
	}
	s.less = less
	sort.Sort(s)
	return newFromReversedSlice(s.elements)
}
