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

// TakeWhile creates a new list using the elements of the original one. It will 
// keep the original elements while the predicate given as argument is valid.  
// After that, all elements of the original list are discarded.
//
// Example:
// l1 := L(1, 2, 3, 2, 5, 4, 3, 9, 1)
// l2 := TakeWhile(l1, func(x Elem) bool {
// 	return x.(int) < 5
// })
//
// -> l2 = [1, 2, 3, 2]
func TakeWhile(l *List, f func(x Elem) bool) *List {
	if Empty(l) {
		return New()
	}

	if f(Head(l)) {
		return Cons(Head(l), TakeWhile(Tail(l), f))
	}

	return New()
}

// DropWhile, like TakeWhile, creates a new list using the elements of the 
// original one.  Unlike TakeWhile, however, it will drop the original elements 
// while the predicate given as argument is valid. The remaining elements are 
// used to create the new list.
//
// Example:
// l1 := L(1, 2, 3, 2, 5, 4, 3, 9, 1)
// l2 := DropWhile(l1, func(x Elem) bool {
// 	return x.(int) < 5
// })
//
// -> l2 = [5, 4, 3, 9, 1]
func DropWhile(l *List, f func(x Elem) bool) *List {
	if Empty(l) {
		return New()
	}

	if f(Head(l)) {
		return DropWhile(Tail(l), f)
	}

	return l
}

// Span breaks the original list in two when it finds an element for which the 
// predicate doesn't hold: the first one are the ones it keeps from the 
// original list while the predicate holds; the second one are the remaining 
// elements.
//
// Example:
//
// l := L(1, 2, 3, 2, 5, 4, 3, 9, 1)
// l1, l2 := Span(l, func(x Elem) bool {
// 	return x.(int) < 5
// })
//
// -> l1 = [1, 2, 3, 2]
//    l2 = [5, 4, 3, 9, 1]
func Span(l *List, f func(x Elem) bool) (first, rest *List) {
	// TODO This implementation is inefficient. Write a new one.
	return TakeWhile(l, f), DropWhile(l, f)
}

// Flatten takes a list of lists and transforms it in a flat list.
//
// Example:
//
// l := L(L(1,2), L(3, 4))
// -> l = [[1,2], [3,4]]
// flat := Flatten(l)
// -> flat = [1, 2, 3, 4]
func Flatten(l *List) *List {
	return Foldr1(l, func(x Elem, acc interface{}) interface{} {
		return Concatenate(x.(*List), acc.(*List))
	}).(*List)
}

// Synonym for Flatten
var Concat = Flatten

// Takes a list of booleans and returns true only if all elements are true
//
// Example:
//
// l1 := L(true, true, true, true)
// And(l1)
// -> true
//
// l2 := L(true, true, false, true)
// And(l2)
// -> false
func And(l *List) bool {
	return Foldr1(l, func(x Elem, acc interface{}) interface{} {
		return x.(bool) && acc.(bool)
	}).(bool)
}

// Takes a list of booleans and returns true if there are true values in it.
//
// Example:
//
// l1 := L(false, false, false, true)
// -> true
//
// l2 := L(false, false, false, false)
// -> false
func Or(l *List) bool {
	return Foldr1(l, func(x Elem, acc interface{}) interface{} {
		return x.(bool) || acc.(bool)
	}).(bool)
}

// Returns true if all elements satisfy the predicate. Returns false otherwise.
//
// Example:
//
// l := L(1, 1, 3, 2, 5)
// All(l, func(x Elem) bool {
// 	return x.(int) < 7
// })
// -> true
//
// All(l, func(x Elem) bool {
// 	return x.(int) < 5
// })
// -> false
func All(l *List, f func(Elem) bool) bool {
	return Foldr(true, l, func(x Elem, acc interface{}) interface{} {
		return f(x) && acc.(bool)
	}).(bool)
}

// Returns true if any element satisfies the predicate. Returns false 
// otherwise.
//
// Example:
//
// l := L(1, 1, 3, 2, 5)
// Any(l, func(x Elem) bool {
// 	return x.(int) < 3
// })
// -> true
//
// Any(l, func(x Elem) bool {
// 	return x.(int) < 1
// })
// -> false
func Any(l *List, f func(Elem) bool) bool {
	return Foldr(false, l, func(x Elem, acc interface{}) interface{} {
		return f(x) || acc.(bool)
	}).(bool)
}

// Groups consecutive identical elements into sublists.
//
// Example:
//
// l := L(1, 1, 1, 3, 3, 2, 3, 3, 6, 6, 6)
// Group(l)
// -> [[1,1,1], [3,3], [2], [3,3], [6,6,6]]
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

// Returns two lists. The first one contains all the elements of the original 
// list that satisfy the predicate. The second one contains the remaining one 
// (the ones that do not satisfy the predicate).
//
// Example:
//
// l := L(1, 1, 1, 3, 3, 2, 3, 3, 6, 6, 6)
// Partition(l, func(x Elem) bool {
// 	x.(int)%2 == 0
// })
// -> [2, 6, 6, 6] [1, 1, 1, 3, 3, 3, 3]
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

// Unique returns a list with the elements of the original one but without any 
// repetition.
//
// Example:
//
// l := L(1, 1, 1, 3, 3, 2, 3, 3, 6, 6, 6)
// Unique(l)
// -> [1, 3, 2, 6]
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

// Synonym for Unique
var Nub = Unique

// Deletes the first occurrence of an element from a list.
//
// Example:
//
// l := L(1, 1, 1, 3, 3, 2, 3, 3, 6, 6, 6)
// Delete(3, l)
// -> [1, 1, 1, 3, 2, 3, 3, 6, 6, 6]
func Delete(x Elem, l *List) *List {
	without, with := Span(l, func(y Elem) bool {
		return x != y
	})
	return Concatenate(without, Tail(with))
}

// Removes, from the first list, the elements found in the second one.
//
// Example:
//
// l1 := L(1, 2, 2, 3, 4, 4, 5, 6)
// l2 := L(2, 5, 7, 9, 7, 10)
// Difference(l1, l2)
// -> [1, 3, 4, 4, 6]
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

// Makes the union of the two lists. Duplicated elements of the second list are 
// removed, as well as elements also found in the first list.  However, 
// duplicated elements of the first list aren't removed.
//
// Example:
//
// l1 := L(1, 2, 2, 3, 4, 4, 5, 6)
// l2 := L(2, 5, 7, 9, 7, 10)
// Union(l1, l2)
// -> [1, 2, 2, 3, 4, 4, 5, 6, 7, 9, 7, 10]
func Union(l1, l2 *List) *List {
	return Concatenate(l1, Difference(Unique(l2), l1))
}

// Makes the intersection of the two lists. If the first list contains 
// duplicates, so will the result.
//
// Example:
// l1 := L(1, 2, 2, 3, 4, 4, 5, 6)
// l2 := L(2, 5, 7, 9, 7, 10)
// Intersect(l1, l2)
// -> [2, 2, 5]
func Intersect(l1, l2 *List) *List {
	table := make(map[Elem]bool)
	next := MakeIterator(l2)
	for e := next(); e != nil; e = next() {
		table[e] = true
	}
	return Foldr(New(), l1, func(x Elem, acc interface{}) interface{} {
		if _, ok := table[x]; ok {
			return Cons(x, acc.(*List))
		}
		return acc
	}).(*List)
}

// Returns true if the two lists have equal elements.
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

// Applies the function f to each element in the list, from left to right.
//
// Example:
//
// l := L(1,2,3)
// Each(l, func(x Elem) {
// 	fmt.Println(x, " ")
// })
// -> 1 2 3
func Each(l *List, f func(Elem)) {
	Foldl(New(), l, func(acc interface{}, x Elem) interface{} {
		f(x)
		return x
	})
}
