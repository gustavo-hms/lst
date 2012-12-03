/*
   Package lst provides an implementation of lists using shared vectors.

   To create a new empty list, type:

   	l := lst.New()

   You can also create a list filled with predefined elements:

   	l := lst.NewWithElements(1, 2, 3, 4)

   Or, as a shorthand,

   	l := lst.L(1, 2, 3, 4)

   If you want to use the elements of a slice:

   	l := lst.NewFromSlice(aSlice)
*/
package lst

import (
	"fmt"
	"strings"
)

type Elem interface{}

type group struct {
	elements []Elem
	// All lists have underlying vectors storing its elements, which can be 
	// shared across several lists. To avoid that a list overwrites elements of 
	// other lists with shared data, we must be sure we are not writting at 
	// a position already used in this shared vector. Since slices do not see 
	// anything in the underlying vector except for its own data, we must keep 
	// two informations regarding the original vector: first, the index of its 
	// first empty position (stored in the firsEmpty variable); second, the 
	// index at which the slice begins (offset variable).  So, for instance, if 
	// we have a vector v = [1, 2, 3, 4, 0, 0, 0, 0] and a slice s := v[3:5], 
	// the index at which the slice begins is 3, and the index of the first 
	// empty position is 4. To see how this information is used, see the Cons 
	// implementation.
	accumulatedLen int
	offset         int
	firstEmpty     *int
}

func newGroup(accumulatedLen int) group {
	var g group
	g.elements = make([]Elem, 0)
	g.accumulatedLen = accumulatedLen
	g.firstEmpty = new(int)
	return g
}

func newGroupFromReversedSlice(elems []Elem) group {
	var g group
	length := len(elems)
	g.firstEmpty = &length
	g.elements = make([]Elem, length)
	copy(g.elements, elems)
	return g
}

type List struct {
	groups          []group
	firstEmptyGroup *int
	groupsOffset    int
}

func New() *List {
	l := new(List)
	l.groups = []group{newGroup(0)}
	firstEmptyGroup := 1
	l.firstEmptyGroup = &firstEmptyGroup
	return l
}

func NewFromList(original *List) *List {
	dest := new(List)
	dest.groups = original.groups
	dest.firstEmptyGroup = original.firstEmptyGroup

	return dest
}

func newCopyingGroups(original *List, extraCapacity ...int) *List {
	dest := new(List)
	numberOfGroups := len(original.groups)

	extra := 0
	amountToShrink := 0
	if len(extraCapacity) > 0 {
		extra = extraCapacity[0]
		if extra < 0 {
			amountToShrink = extra
		}
	}

	dest.groups = make([]group, numberOfGroups-amountToShrink, numberOfGroups+extra)
	copy(dest.groups, original.groups)

	firstEmpty := len(dest.groups)
	dest.firstEmptyGroup = &firstEmpty

	return dest
}

func newFromReversedSlice(slice []Elem) *List {
	l := new(List)
	l.groups = []group{newGroupFromReversedSlice(slice)}
	firstEmptyGroup := 1
	l.firstEmptyGroup = &firstEmptyGroup

	return l
}

// Constructs a new list from the given slice
//
// Example:
//
// l := NewFromSlice([]Elem{1, 2, 3, 4})
// -> [1, 2, 3, 4]
func NewFromSlice(slice []Elem) (l *List) {
	reversed := make([]Elem, len(slice))
	for k, v := range slice {
		reversed[len(slice)-k-1] = v
	}

	return newFromReversedSlice(reversed)
}

// Constructs a new list with the given elements.
//
// Example:
//
// l := NewWithElements(1, 2, 3, 4)
// -> [1, 2, 3, 4]
func NewWithElements(elems ...Elem) (l *List) {
	return NewFromSlice(elems)
}

var L = NewWithElements // Just for convenience

// The length of the list
func Len(l *List) int {
	lastGroup := l.groups[len(l.groups)-1]
	return lastGroup.accumulatedLen + len(lastGroup.elements)
}

// Gets the element at index i
func Get(l *List, i int) Elem {
	last := Len(l) - 1
	index := last - i
	group := findGroup(l.groups, index)

	return group.elements[index-group.offset]
}

func findGroup(groups []group, i int) *group {
	numberOfGroups := len(groups)
	if numberOfGroups == 1 {
		return &groups[0]
	}

	middle := numberOfGroups / 2
	if i < groups[middle].accumulatedLen {
		return findGroup(groups[:middle], i)
	}
	return findGroup(groups[middle:], i)
}

func set(l *List, i int, x Elem) {
	last := Len(l) - 1
	index := last - i
	group := findGroup(l.groups, index)

	group.elements[index-group.offset] = x
}

func (l *List) String() string {
	elems := make([]string, Len(l))

	for i := 0; i < Len(l); i++ {
		elems[i] = fmt.Sprintf("%v", Get(l, i))
	}

	return "[" + strings.Join(elems, ", ") + "]"
}

// Gets the head of the list (its most recently inserted element)
func Head(l *List) Elem {
	return Get(l, 0)
}

// Gets all but the head of the list
func Tail(l *List) *List {
	lastGroup := len(l.groups) - 1
	switch {
	case Len(l) == 1:
		return New()
	case len(l.groups[lastGroup].elements) > 1:
		return removeLastElement(l)
	}

	return removeLastGroup(l)
}

func removeLastElement(original *List) *List {
	l := newCopyingGroups(original)
	lastGroup := &l.groups[len(l.groups)-1]
	lastGroup.elements = lastGroup.elements[:len(lastGroup.elements)-1]

	return l
}

func removeLastGroup(original *List) *List {
	l := NewFromList(original)
	l.groups = l.groups[:len(l.groups)-1]
	return l
}

// Gets the last element of the list (the first inserted element)
func Last(l *List) Elem {
	return Get(l, Len(l)-1)
}

// Gets all but the last element of the list
func Init(l *List) (init *List) {
	switch {
	case Len(l) == 1:
		return New()
	case len(l.groups[0].elements) > 1:
		return removeFirstElement(l)
	}

	return removeFirstGroup(l)
}

func removeFirstElement(original *List) *List {
	l := newCopyingGroups(original)
	firstGroup := &l.groups[0]
	firstGroup.elements = firstGroup.elements[:len(firstGroup.elements)-1]

	for i, g := range l.groups[1:] {
		g.accumulatedLen--
		l.groups[i+1] = g
	}

	return l
}

func removeFirstGroup(original *List) *List {
	l := NewFromList(original)
	l.groups = l.groups[1:]
	l.groupsOffset++
	return l
}

// The list constructor. It constructs a new list by inserting a new element in
// the front of an old one.
//
// Example:
//
// first := New()
//
// second := Cons(1, first)
//
// third := Cons(2, second)
//
// Cons(3, third)
// -> [3, 2, 1]
func Cons(x Elem, xs *List) *List {
	if canAppendMoreToGroup(xs.groups[len(xs.groups)-1]) {
		return appendToLastGroup(x, xs)
	}

	return createNewGroupWithElement(x, xs)
}

func canAppendMoreToGroup(g group) bool {
	return len(g.elements)+g.offset >= *g.firstEmpty
}

func appendToLastGroup(x Elem, xs *List) *List {
	l := newCopyingGroups(xs)
	expandGroup([]Elem{x}, &l.groups[len(l.groups)-1])

	return l
}

func expandGroup(elements []Elem, g *group) {
	oldSliceAddress := &g.elements
	g.elements = append(g.elements, elements...)

	if oldSliceAddress != &g.elements {
		// append has copied array to a new address
		g.offset = 0
		firstEmpty := *g.firstEmpty
		g.firstEmpty = &firstEmpty
	}

	length := len(elements)
	g.firstEmpty = &length
}

func createNewGroupWithElement(x Elem, xs *List) *List {
	g := newGroupFromReversedSlice([]Elem{x})

	if canAppendMoreGroups(xs) {
		l := NewFromList(xs)
		oldAddress := &l.groups
		l.groups = append(l.groups, g)

		if oldAddress != &l.groups {
			// append has copied array to a new address
			l.groupsOffset = 0
			firstEmptyGroup := *l.firstEmptyGroup
			l.firstEmptyGroup = &firstEmptyGroup
		}

		*l.firstEmptyGroup++

		return l
	}

	// This case, we need to make a copy of the entire groups slice
	l := newCopyingGroups(xs, 1)
	l.groups = append(l.groups, g)
	*l.firstEmptyGroup++

	return l
}

func canAppendMoreGroups(l *List) bool {
	return len(l.groups)+l.groupsOffset >= *l.firstEmptyGroup
}

func concatenate(l1, l2 *List) *List {
	if Empty(l1) {
		return NewFromList(l2)
	}

	if Empty(l2) {
		return NewFromList(l1)
	}

	var con *List
	if canAppendMoreToGroup(l2.groups[len(l2.groups)-1]) {
		con = newCopyingGroups(l2)
	} else {
		con = newCopyingGroups(l2, 1)
		lastGroup := con.groups[len(con.groups)-1]
		accumulatedLen := lastGroup.accumulatedLen + len(lastGroup.elements)
		con.groups = append(con.groups, newGroup(accumulatedLen))
	}

	lastGroup := &con.groups[len(con.groups)-1]
	for _, g := range l1.groups {
		expandGroup(g.elements, lastGroup)
	}

	return con
}

// Concatenates all the lists given as arguments.
//
// Example:
//
// l1 := NewWithElements(1, 2)
// l2 := NewWithElements(3, 4)
// l3 := NewWithElements(5, 6)
// Concatenate(l1, l2, l3)
// -> [1, 2, 3, 4, 5, 6]
func Concatenate(lists ...*List) (con *List) {
	last := len(lists) - 1
	con = lists[last]
	for i := last - 1; i >= 0; i-- {
		con = concatenate(lists[i], con)
	}
	return
}
