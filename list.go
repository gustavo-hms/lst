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
	elements    [][]Elem
	accumLength []int
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
	firstEmpty []int
	offsets    []int
}

func New() *List {
	l := new(List)
	l.elements = make([][]Elem, 1)
	l.elements[0] = make([]Elem, 0)
	l.accumLength = make([]int, 1)
	l.firstEmpty = make([]int, 1)
	l.offsets = make([]int, 1)
	return l
}

func NewFromList(original *List) (dest *List) {
	dest = new(List)
	dest.elements = original.elements
	dest.accumLength = original.accumLength
	dest.firstEmpty = original.firstEmpty
	dest.offsets = original.offsets
	return
}

func newFromReversedSlice(slice []Elem) (l *List) {
	l = New()
	l.elements[0] = make([]Elem, len(slice))
	copy(l.elements[0], slice)
	l.firstEmpty[0] = len(slice)
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
	numberOfGroups := len(l.accumLength)
	return l.accumLength[numberOfGroups-1] + len(l.elements[numberOfGroups-1])
}

// Gets the element at index i
func Get(l *List, i int) Elem {
	last := Len(l) - 1
	index := last - i
	slice, offset := findSlice(l.accumLength, l.elements, index)

	return slice[index-offset]
}

func findSlice(lengths []int, elements [][]Elem, i int) (group []Elem, offset int) {
	numberOfGroups := len(lengths)
	if numberOfGroups == 1 {
		return elements[0], lengths[0]
	}

	middle := numberOfGroups/2 + 1
	if i < lengths[middle] {
		return findSlice(lengths[:middle], elements[:middle], i)
	}
	return findSlice(lengths[middle:], elements[middle:], i)
}

func set(l *List, i int, x Elem) {
	last := Len(l) - 1
	index := last - i
	slice, offset := findSlice(l.accumLength, l.elements, index)

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
	tail = NewFromList(l)
	numberOfGroups := len(tail.elements)
	lenLastGroup := len(tail.elements[numberOfGroups-1])
	if lenLastGroup > 1 {
		tail.elements[numberOfGroups-1] = tail.elements[numberOfGroups-1][:lenLastGroup-1]
		return
	}

	tail.elements = tail.elements[:numberOfGroups-1]
	tail.offsets = tail.offsets[:numberOfGroups-1]
	return
}

// Gets the last element of the list (the first inserted element)
func Last(l *List) Elem {
	return Get(l, Len(l)-1)
}

// Gets all but the last element of the list
func Init(l *List) (init *List) {
	init = NewFromList(l)
	lenFirstGroup := len(init.elements[0])
	if lenFirstGroup > 1 {
		init.elements[0] = init.elements[0][1:]
		init.offsets[0] += 1
		for k, v := range init.accumLength[1:] {
			init.accumLength[k] = v - 1
		}
		return
	}

	init.elements = init.elements[1:]
	init.offsets = init.offsets[1:]
	for k, v := range init.accumLength[1:] {
		init.accumLength[k] = v - init.accumLength[1]
	}
	init.accumLength = init.accumLength[1:]
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
	if newl.firstEmpty[lastGroup] <= len(newl.elements[lastGroup])+newl.offsets[lastGroup] {
		// A new value can be appended to l.elements[lastGroup] without 
		// overwritting any shared data
		newl.elements[lastGroup] = append(newl.elements[lastGroup], x)
		newl.firstEmpty[lastGroup]++
		return
	}

	// In this case, we can't append to l.elements[lastGroup], because it would 
	// overwrite a shared value. To circumvent this problem, we simply start 
	// a new elements' group
	newl.elements = append(newl.elements, []Elem{x})
	newl.firstEmpty = append(newl.firstEmpty, 1)
	newl.offsets = append(newl.offsets, 0)
	lenLastGroup := len(newl.elements[lastGroup])
	newl.accumLength = append(newl.accumLength, newl.accumLength[lastGroup]+lenLastGroup)
	return
}

func concatenate(l1, l2 *List) (con *List) {
	con = NewFromList(l2)
	for k, group := range l1.elements {
		con.elements = append(con.elements, group)
		con.offsets = append(con.offsets, l1.offsets[k])
		con.firstEmpty = append(con.firstEmpty, l1.firstEmpty[k])

		lastInserted := len(con.accumLength) - 1
		con.accumLength = append(con.accumLength,
			con.accumLength[lastInserted]+len(con.elements[lastInserted]))
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
