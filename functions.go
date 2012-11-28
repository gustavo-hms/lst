package lst

// Gives the reverse of some list.
func Reverse(l *List) *List {
	return Foldl(New(), l, func(xs interface{}, x Elem) interface{} {
		return Cons(x, xs.(*List))
	}).(*List)
}

// Tells if a list is empty
func Empty(l *List) bool {
	return Len(l) <= 0
}

// Synonym for Empty
var Null = Empty

// Same as Foldr, but uses the last element of the list as the initial value
func Foldr1(l *List, f func(x Elem, acc interface{}) interface{}) interface{} {
	return Foldr(Last(l), Init(l), f)
}

// Same as Foldl, but uses the first element of the list as the initial value
func Foldl1(l *List, f func(acc interface{}, x Elem) interface{}) interface{} {
	return Foldl(Head(l), Tail(l), f)
}

// Map creates a new list with the same size as the original one, whose 
// elements are obtained applying the function f to each element of the 
// original list.
//
// Example:
//
// l := NewWithElements(1, 2, 3, 4)
// plus17 := Map(l, func(x Elem) Elem {
// 	return x.(int) + 17
// })
//
// -> plus17 = [18, 19, 20, 21]
func Map(l *List, f func(Elem) Elem) *List {
	return Foldr(New(), l, func(x Elem, acc interface{}) interface{} {
		return Cons(f(x), acc.(*List))
	}).(*List)
}

// Filter creates a new list using only the elements in the original one that 
// satisfy the given predicate.
//
// Example:
//
// l := NewWithElements(1, 2, 3, 4)
// even := Filter(l, func(x Elem) bool {
// 	return x.(int)%2 == 0
// })
//
// -> even = [2, 4]
func Filter(l *List, f func(Elem) bool) *List {
	return Foldr(New(), l, func(x Elem, acc interface{}) interface{} {
		if f(x) {
			return Cons(x, acc.(*List))
		}
		return acc
	}).(*List)
}

// Sums all elements of a list of integers
func IntSum(l *List) int {
	result := Foldr(0, l, func(x Elem, acc interface{}) interface{} {
		return acc.(int) + x.(int)
	})
	return result.(int)
}

// Synonym for IntSum
var Sum = IntSum

// Sums all elements of a list of float64
func FloatSum(l *List) float64 {
	result := Foldr(0, l, func(x Elem, acc interface{}) interface{} {
		return acc.(float64) + x.(float64)
	})
	return result.(float64)
}

// Gives the accumulated product of all elements of a list of integers
func IntProd(l *List) int {
	result := Foldr(1, l, func(x Elem, acc interface{}) interface{} {
		return acc.(int) * x.(int)
	})
	return result.(int)
}

// Synonym for IntProd
var Prod = IntProd

// Gives the accumulated product of all elements of a list of integers
func FloatProd(l *List) float64 {
	result := Foldr(1, l, func(x Elem, acc interface{}) interface{} {
		return acc.(float64) * x.(float64)
	})
	return result.(float64)
}

// Tells if some element belongs to the given list
func Element(x Elem, l *List) bool {
	switch {
	case Empty(l):
		return false
	case x == Head(l):
		return true
	}
	return Element(x, Tail(l))
}

// Tells if some element does not belongs to the given list
func NotElement(x Elem, l *List) bool {
	return !Element(x, l)
}

// ElemIndex returns 2 items. The first one is the index of the first 
// occurrence of the element x in the list l if such an element belongs to the 
// list. The second item it returns is true if the element could be found in 
// the list, or false otherwise.
func ElemIndex(x Elem, l *List) (int, bool) {
	switch {
	case Empty(l):
		return -1, false
	case x == Head(l):
		return 0, true
	}

	i, ok := ElemIndex(x, Tail(l))
	return i + 1, ok
}

// ElemIndices returns a list with the indices of all occurrences of the 
// element x in the list xs
func ElemIndices(x Elem, xs *List) *List {
	count := -1
	return indices(x, xs, count)
}

func indices(y Elem, ys *List, count int) *List {
	if Empty(ys) {
		return ys
	}

	count++
	if Head(ys) == y {
		return Cons(count, indices(y, Tail(ys), count))
	}
	return indices(y, Tail(ys), count)
}

// Zip merges two lists together, creating a new list where each element is 
// a list containing two elements: one from each of the original lists. The 
// length of the new list is equal to the length of the smallest one.
//
// Example: 
//
// l1 := L(1, 2, 3, 4)
// l2 := L(5, 6, 7)
// zipped := Zip(l1, l2)
//
// -> zipped = [[1,5], [2,6], [3, 7]]
func Zip(l1, l2 *List) *List {
	if Empty(l1) || Empty(l2) {
		return New()
	}
	return Cons(L(Head(l1), Head(l2)), Zip(Tail(l1), Tail(l2)))
}

// ZipWith is simillar to Zip, but instead of automatically combining the 
// elements of the two original lists creating sublists, it uses the function 
// provided as an argument to mix such elements.
//
// Example:
// l1 := L(1, 2, 3, 4)
// l2 := L(5, 6, 7)
// zipped := ZipWith(l1, l2, func(x, y Elem) Elem {
// 	return x.(int) * y.(int)
// })
//
// -> zipped = [5, 12, 21]
func ZipWith(l1, l2 *List, f func(x, y Elem) Elem) *List {
	if Empty(l1) || Empty(l2) {
		return New()
	}
	return Cons(f(Head(l1), Head(l2)), ZipWith(Tail(l1), Tail(l2), f))
}

func TakeWhile(l *List, f func(x Elem) bool) *List {
	if Empty(l) {
		return New()
	}

	if f(Head(l)) {
		return Cons(Head(l), TakeWhile(Tail(l), f))
	}

	return New()
}

func DropWhile(l *List, f func(x Elem) bool) *List {
	if Empty(l) {
		return New()
	}

	if f(Head(l)) {
		return DropWhile(Tail(l), f)
	}

	return l
}

func Span(l *List, f func(x Elem) bool) (first, rest *List) {
	return TakeWhile(l, f), DropWhile(l, f)
}

func Flatten(l *List) *List {
	return Foldr1(l, func(x Elem, acc interface{}) interface{} {
		return Concatenate(x.(*List), acc.(*List))
	}).(*List)
}

var Concat = Flatten

func And(l *List) bool {
	return Foldr1(l, func(x Elem, acc interface{}) interface{} {
		return x.(bool) && acc.(bool)
	}).(bool)
}

func Or(l *List) bool {
	return Foldr1(l, func(x Elem, acc interface{}) interface{} {
		return x.(bool) || acc.(bool)
	}).(bool)
}

func All(l *List, f func(Elem) bool) bool {
	return Foldr(true, l, func(x Elem, acc interface{}) interface{} {
		return f(x) && acc.(bool)
	}).(bool)
}

func Any(l *List, f func(Elem) bool) bool {
	return Foldr(false, l, func(x Elem, acc interface{}) interface{} {
		return f(x) || acc.(bool)
	}).(bool)
}

func Group(l *List) *List {
	vec := [2]*List{New(), New()} // {final list, sublist}
	result := Foldr(vec, l, func(x Elem, acc interface{}) interface{} {
		v := acc.([2]*List)
		switch {
		case Empty(v[1]):
			v[1] = Cons(x, v[1])
		case x == Head(v[1]):
			v[1] = Cons(x, v[1])
		default:
			v[0] = Cons(v[1], v[0])
			v[1] = NewWithElements(x)
		}
		return v
	})

	r := result.([2]*List)
	return Cons(r[1], r[0])
}

func Partition(l *List, f func(Elem) bool) (satisfy, doNot *List) {
	vec := [2]*List{New(), New()} // {satisfy, doNot}
	result := Foldr(vec, l, func(x Elem, acc interface{}) interface{} {
		v := acc.([2]*List)
		if f(x) {
			v[0], v[1] = Cons(x, v[0]), v[1]
		} else {
			v[0], v[1] = v[0], Cons(x, v[1])
		}
		return v
	})

	r := result.([2]*List)
	return r[0], r[1]
}

func Unique(l *List) *List {
	table := make(map[Elem]bool)
	return unique(l, table)
}

func unique(l *List, table map[Elem]bool) *List {
	if Empty(l) {
		return l
	}

	head := Head(l)
	if _, ok := table[head]; ok {
		return unique(Tail(l), table)
	}

	table[head] = true
	return Cons(head, unique(Tail(l), table))
}

var Nub = Unique

func Delete(x Elem, l *List) *List {
	without, with := Span(l, func(y Elem) bool {
		return x != y
	})
	return Concatenate(without, Tail(with))
}

func Difference(base, subtract *List) *List {
	table := make(map[Elem]bool)
	next := MakeIterator(subtract)
	for e := next(); e != nil; e = next() {
		table[e] = true
	}
	return Foldr(New(), base, func(x Elem, acc interface{}) interface{} {
		if _, ok := table[x]; ok {
			return acc
		}
		return Cons(x, acc.(*List))
	}).(*List)
}

func Union(l1, l2 *List) *List {
	return Concatenate(l1, Difference(l2, l1))
}

func Intersect(l1, l2 *List) *List {
	table := make(map[Elem]bool)
	next := MakeIterator(l1)
	for e := next(); e != nil; e = next() {
		table[e] = true
	}
	return Foldr(New(), l2, func(x Elem, acc interface{}) interface{} {
		if _, ok := table[x]; ok {
			return Cons(x, acc.(*List))
		}
		return acc
	}).(*List)
}

func Equal(l1, l2 *List) bool {
	lenL1 := Len(l1)
	lenL2 := Len(l2)
	if lenL1 != lenL2 {
		return false
	}

	for i := 0; i < lenL1; i++ {
		if Get(l1, i) != Get(l2, i) {
			return false
		}
	}
	return true
}

func Each(l *List, f func(Elem)) {
	Foldl(New(), l, func(acc interface{}, x Elem) interface{} {
		f(x)
		return x
	})
}
