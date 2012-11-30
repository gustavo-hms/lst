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

type List struct {
	elements       [][]Elem
	accumulatedLen []int
	// All lists have underlying vectors storing its elements, which can be 
	// shared across several lists. To avoid that a list overwrites elements of 
	// other lists with shared data, we must be sure we are not writting at 
	// a position already used in this shared vector. Since slices do not see 
	// anything in the underlying vector except for its own data, we must keep 
	// two informations regarding the original vector: first, the index of its 
	// first empty position (stored in the firsEmpty variable); second, the 
	// index at which the slice begins (offset variable).  So, for instance, 
	// if we have a vector v = [1, 2, 3, 4, 0, 0, 0, 0] and a slice s := 
	// v[3:5], the index at which the slice begins is 3, and the index of the 
	// first empty position is 4. To see how this information is used, see the 
	// Cons implementation.
	offsets    []int
	firstEmpty []*int
}

func New() *List {
	l := new(List)
	l.elements = make([][]Elem, 1)
	l.elements[0] = make([]Elem, 0)
	l.accumulatedLen = make([]int, 1)
	l.firstEmpty = make([]*int, 1)
	l.firstEmpty[0] = new(int)
	l.offsets = make([]int, 1)
	return l
}

func NewFromList(original *List) (dest *List) {
	dest = new(List)
	numberOfGroups := len(original.elements)

	dest.elements = make([][]Elem, len(original.elements))
	copy(dest.elements, original.elements)

	dest.accumulatedLen = make([]int, numberOfGroups)
	copy(dest.accumulatedLen, original.accumulatedLen)

	dest.firstEmpty = make([]*int, numberOfGroups)
	copy(dest.firstEmpty, original.firstEmpty)

	dest.offsets = make([]int, numberOfGroups)
	copy(dest.offsets, original.offsets)

	return
}

func newFromReversedSlice(slice []Elem) (l *List) {
	l = New()
	l.elements[0] = make([]Elem, len(slice))
	copy(l.elements[0], slice)
	lenSlice := len(slice)
	l.firstEmpty[0] = &lenSlice
	return
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
	numberOfGroups := len(l.accumulatedLen)
	return l.accumulatedLen[numberOfGroups-1] + len(l.elements[numberOfGroups-1])
}

// Gets the element at index i
func Get(l *List, i int) Elem {
	last := Len(l) - 1
	index := last - i
	slice, offset := findSlice(l.accumulatedLen, l.elements, index)

	return slice[index-offset]
}

func findSlice(lengths []int, elements [][]Elem, i int) (group []Elem, offset int) {
	numberOfGroups := len(lengths)
	if numberOfGroups == 1 {
		return elements[0], lengths[0]
	}

	middle := numberOfGroups / 2
	if i < lengths[middle] {
		return findSlice(lengths[:middle], elements[:middle], i)
	}
	return findSlice(lengths[middle:], elements[middle:], i)
}

func set(l *List, i int, x Elem) {
	last := Len(l) - 1
	index := last - i
	slice, offset := findSlice(l.accumulatedLen, l.elements, index)

	slice[index-offset] = x
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
func Tail(l *List) (tail *List) {
	if Len(l) == 1 {
		tail = New()
		return
	}

	tail = NewFromList(l)
	numberOfGroups := len(tail.elements)
	lenLastGroup := len(tail.elements[numberOfGroups-1])
	if lenLastGroup > 1 {
		tail.elements[numberOfGroups-1] = tail.elements[numberOfGroups-1][:lenLastGroup-1]
		return
	}

	lastGroupRemoved := tail.elements[:numberOfGroups-1]
	tail.elements = make([][]Elem, numberOfGroups-1)
	copy(tail.elements, lastGroupRemoved)

	lastOffsetRemoved := tail.offsets[:numberOfGroups-1]
	tail.offsets = make([]int, numberOfGroups-1)
	copy(tail.offsets, lastOffsetRemoved)

	lastAccumulatedLenRemoved := tail.accumulatedLen[:numberOfGroups-1]
	tail.accumulatedLen = make([]int, numberOfGroups-1)
	copy(tail.accumulatedLen, lastAccumulatedLenRemoved)

	lastFirstEmptyRemoved := tail.firstEmpty[:numberOfGroups-1]
	tail.firstEmpty = make([]*int, numberOfGroups-1)
	copy(tail.firstEmpty, lastFirstEmptyRemoved)

	return
}

// Gets the last element of the list (the first inserted element)
func Last(l *List) Elem {
	return Get(l, Len(l)-1)
}

// Gets all but the last element of the list
func Init(l *List) (init *List) {
	if Len(l) == 1 {
		init = New()
		return
	}

	init = NewFromList(l)
	lenFirstGroup := len(init.elements[0])
	if lenFirstGroup > 1 {
		init.elements[0] = init.elements[0][1:]
		init.offsets[0] += 1
		for k, v := range init.accumulatedLen[1:] {
			init.accumulatedLen[k+1] = v - 1
		}
		return
	}

	numberOfGroups := len(init.elements)

	firstGroupRemoved := init.elements[1:]
	init.elements = make([][]Elem, numberOfGroups-1)
	copy(init.elements, firstGroupRemoved)

	firsOffsetRemoved := init.offsets[1:]
	init.offsets = make([]int, numberOfGroups-1)
	copy(init.offsets, firsOffsetRemoved)

	firstAccumulatedLen := init.accumulatedLen[1]
	firstAccumulatedLenRemoved := init.accumulatedLen[1:]
	init.accumulatedLen = make([]int, numberOfGroups-1)
	for k, v := range firstAccumulatedLenRemoved {
		init.accumulatedLen[k] = v - firstAccumulatedLen
	}

	firstEmptyRemoved := init.firstEmpty[1:]
	init.firstEmpty = make([]*int, numberOfGroups-1)
	copy(init.firstEmpty, firstEmptyRemoved)
	return
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
func Cons(x Elem, l *List) (newl *List) {
	lastGroup := len(l.elements) - 1
	newl = NewFromList(l)
	if canAppendMoreToGroup(newl, 1, lastGroup) {
		// A new value can be appended to l.elements[lastGroup] without 
		// overwritting any shared data
		newl.elements[lastGroup] = append(newl.elements[lastGroup], x)
		*newl.firstEmpty[lastGroup]++
		return
	}

	// In this case, we can't append to l.elements[lastGroup], because it would 
	// overwrite a shared value. To circumvent this problem, we simply start 
	// a new elements' group
	newl.elements = append(newl.elements, []Elem{x})
	lenNewGroup := 1
	newl.firstEmpty = append(newl.firstEmpty, &lenNewGroup)
	newl.offsets = append(newl.offsets, 0)
	lenLastGroup := len(newl.elements[lastGroup])
	newl.accumulatedLen =
		append(newl.accumulatedLen, newl.accumulatedLen[lastGroup]+lenLastGroup)
	return
}

func canAppendMoreToGroup(l *List, amount, groupNumber int) bool {
	return len(l.elements[groupNumber])+l.offsets[groupNumber]+amount >=
		*l.firstEmpty[groupNumber]+1
}

func concatenate(l1, l2 *List) (con *List) {
	if Empty(l1) {
		con = NewFromList(l2)
		return
	}

	if Empty(l2) {
		con = NewFromList(l1)
		return
	}

	con = NewFromList(l2)
	for k, group := range l1.elements {
		lastGroup := len(con.elements) - 1
		if canAppendMoreToGroup(con, len(group), lastGroup) {
			con.elements[lastGroup] = append(con.elements[lastGroup], group...)
			*con.firstEmpty[lastGroup] += len(group)
			continue
		}

		con.elements = append(con.elements, group)
		con.offsets = append(con.offsets, l1.offsets[k])
		con.firstEmpty = append(con.firstEmpty, l1.firstEmpty[k])

		lastInserted := len(con.accumulatedLen) - 1
		con.accumulatedLen = append(con.accumulatedLen,
			con.accumulatedLen[lastInserted]+len(con.elements[lastInserted]))
	}
	return
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
