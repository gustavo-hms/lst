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
	elements []Elem
	// See Cons function for a better understanding of the following 2 fields
	firstEmpty *int
	firstUsed  int
}

func New() *List {
	l := new(List)
	l.elements = make([]Elem, 0)
	l.firstEmpty = new(int)

	return l
}

func NewFromList(original *List) (dest *List) {
	dest = new(List)
	dest.elements = original.elements[:]
	dest.firstEmpty = original.firstEmpty
	dest.firstUsed = original.firstUsed
	return
}

func newFromReversedSlice(slice []Elem) (l *List) {
	l = new(List)
	l.elements = make([]Elem, len(slice))
	copy(l.elements, slice)
	length := len(slice)
	l.firstEmpty = &length
	return
}

func NewFromSlice(slice []Elem) (l *List) {
	l = new(List)
	l.elements = make([]Elem, len(slice))

	for k, v := range slice {
		set(l, k, v)
	}

	i := len(slice)
	l.firstEmpty = &i
	return
}

func NewWithElements(elems ...Elem) (l *List) {
	return NewFromSlice(elems)
}

var L = NewWithElements // Just for convenience

func Len(l *List) int {
	return len(l.elements)
}

func (l *List) String() string {
	/* 
	 * Os elementos da lista são armazenados no slice elements na ordem inversa 
	 * à que serão exibidos
	 */
	last := Len(l) - 1
	elems := make([]string, last+1)

	for k, v := range l.elements {
		elems[last-k] = fmt.Sprintf("%v", v)
	}

	return "[" + strings.Join(elems, ", ") + "]"
}

func Get(l *List, i int) Elem {
	last := Len(l) - 1
	return l.elements[last-i]
}

func set(l *List, i int, x Elem) {
	last := Len(l) - 1
	l.elements[last-i] = x
}

// MakeIterator creates a function one can use to iterate over the elements of 
// the list. A "nil" value signalises the end of the loop.
//
// Example:
//
// next := MakeIterator(list)
// for element := next(); element != nil; element = next() {
// 	do something
// }
func MakeIterator(l *List) func() Elem {
	index := -1
	return func() Elem {
		index++
		if index > Len(l)-1 {
			return nil
		}
		return Get(l, index)
	}
}

// MakeReverseIterator creates a function one can use to iterate over the 
// elements of the list in the reverse order. A "nil" value signalises the end 
// of the loop.
//
// Example:
//
// previous := MakeIterator(list)
// for element := previous(); element != nil; element = previous() {
// 	do something
// }
func MakeReverseIterator(l *List) func() Elem {
	index := Len(l)
	return func() Elem {
		index--
		if index < 0 {
			return nil
		}
		return Get(l, index)
	}
}

// Gets the head of the list (its most recently inserted element)
func Head(l *List) Elem {
	return Get(l, 0)
}

// Gets all but the head of the list
func Tail(l *List) (tailList *List) {
	tailList = new(List)
	tailList.elements = l.elements[:Len(l)-1]
	tailList.firstEmpty = l.firstEmpty
	tailList.firstUsed = l.firstUsed
	return
}

// Gets the last element of the list (the first inserted element)
func Last(l *List) Elem {
	return Get(l, Len(l)-1)
}

// Gets all but the last element of the list
func Init(l *List) (initList *List) {
	initList = new(List)
	initList.elements = l.elements[1:]
	initList.firstEmpty = l.firstEmpty
	initList.firstUsed = l.firstUsed + 1
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
func Cons(x Elem, l *List) (newl *List) {
	/*
	 * O vetor que efetivamente armazena os elementos da lista pode ser 
	 * compartilhado por diversas listas. Assim, é preciso ter cuidado ao 
	 * inserir um novo elemento, para não sobrescrever um elemento que tenha 
	 * sido inserido por outra lista que compartilhe esse mesmo vetor. Para 
	 * isso é que servem os atributos firstEmpty e firstUsed
	 */
	if *l.firstEmpty > Len(l)+l.firstUsed {
		// Neste caso, a posição desejada do vetor já está sendo ocupada. É
		// necessário fazer uma cópia portanto
		newl = newFromReversedSlice(l.elements)
	} else {
		newl = NewFromList(l)
	}

	if Len(l) == cap(l.elements) {
		// Neste caso, a função append que será usada a seguir vai alocar um 
		// novo vetor, o que obriga a reconfigurar as variáveis firsEmpty e
		// firstUsed
		*newl.firstEmpty = Len(l) + 1
		newl.firstUsed = 0
	} else {
		*newl.firstEmpty++
	}

	newl.elements = append(newl.elements, x)
	return
}

// Foldr makes a fold in the list from right to left. For each element, it 
// applies the function f given in its third argument as f(e, acc), where e is 
// the current element in the list, and acc is the value returned by f in the 
// previous iteration. In its first iteration, it uses the value in "init" as 
// the value for "acc".
//
// Example:
//
// l := NewWithElements(1, 2, 3, 4)
//
// sum := Foldr(0, l, func(x Elem, acc interface{}) interface{} {
// 	return x.(int) + acc.(int)
// })
func Foldr(init interface{}, l *List, f func(Elem, interface{}) interface{}) (accum interface{}) {
	accum = init
	for _, v := range l.elements {
		accum = f(v, accum)
	}
	return
}

// Similar to Foldr, Foldl makes a fold in the list left to right.  For each 
// element, it applies the function f given in its third argument as f(acc, e), 
// where e is the current element in the list, and acc is the value returned by 
// f in the previous iteration. In its first iteration, it uses the value in 
// "init" as the value for "acc".
//
// Example:
//
// l := NewWithElements(1, 2, 3, 4)
//
// sum := Foldr(0, l, func(acc interface{}, x Elem) interface{} {
// 	return x.(int) + acc.(int)
// })
//
// -> sum = 11
func Foldl(init interface{}, l *List, f func(interface{}, Elem) interface{}) (accum interface{}) {
	accum = init
	for i := 0; i < Len(l); i++ {
		accum = f(accum, Get(l, i))
	}
	return
}

func concatenate(l1, l2 *List) (con *List) {
	cons := func(x Elem, accum interface{}) interface{} {
		return Cons(x, accum.(*List))
	}
	return Foldr(l2, l1, cons).(*List)
}

// Concatenates all the lists given as arguments.
//
// Example:
//
// l1 := NewWithElements(1, 2)
// l2 := NewWithElements(3, 4)
// l3 := NewWithElements(5, 6)
// c := Concatenate(l1, l2, l3)
//
// -> c = [1, 2, 3, 4, 5, 6]
func Concatenate(lists ...*List) (con *List) {
	last := len(lists) - 1
	con = lists[last]
	for i := last - 1; i >= 0; i-- {
		con = concatenate(lists[i], con)
	}
	return
}
