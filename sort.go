package lst

import (
	"sort"
)

/*
 * Para poder usar o pacote sort, é necessária uma estrutura em conformidade 
 * com a interface definida lá. Os métodos Len e Swap podem ter uma implantação 
 * default sem problemas, mas o Less precisa ser escrito pelo usuário.
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
	// O slice precisa estar ordenado em ordem inversa. Por isso a negação
	return !s.less(s.elements[i], s.elements[j])
}

/*
 * Sorts a list based on a function returning true if the “x” element is lesser 
 * then “y”
 */
func SortWith(l *List, less func(x, y Elem) bool) *List {
	s := new(sortable)
	s.elements = make([]Elem, Len(l))
	copy(s.elements, l.elements)
	s.less = less
	sort.Sort(s)
	return newFromReversedSlice(s.elements)
}
