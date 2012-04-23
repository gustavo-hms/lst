package lst

func Empty(l *List) bool {
	return Len(l) <= 0
}

var Null = Empty // Sinônimos

func Foldr1(l *List, f func(value, acc interface{}) interface{}) interface{} {
	return Foldr(Last(l), Init(l), f)
}

func Foldl1(l *List, f func(acc, value interface{}) interface{}) interface{} {
	return Foldl(Head(l), Tail(l), f)
}

func Map(l *List, f func(interface{}) interface{}) *List {
	return Foldr(New(), l, func(value, acc interface{}) interface{} {
		return Cons(f(value), acc.(*List))
	}).(*List)
}

func Filter(l *List, f func(interface{}) bool) *List {
	return Foldr(New(), l, func(value, acc interface{}) interface{} {
		if f(value) {
			return Cons(value, acc.(*List))
		}
		return acc
	}).(*List)
}

func IntSum(l *List) int {
	result := Foldr(0, l, func(e, acc interface{}) interface{} {
		return acc.(int) + e.(int)
	})
	return result.(int)
}

var Sum = IntSum

func FloatSum(l *List) float64 {
	result := Foldr(0, l, func(e, acc interface{}) interface{} {
		return acc.(float64) + e.(float64)
	})
	return result.(float64)
}

func IntProd(l *List) int {
	result := Foldr(1, l, func(e, acc interface{}) interface{} {
		return acc.(int) * e.(int)
	})
	return result.(int)
}

var Prod = IntProd

func FloatProd(l *List) float64 {
	result := Foldr(1, l, func(e, acc interface{}) interface{} {
		return acc.(float64) * e.(float64)
	})
	return result.(float64)
}

func Elem(elem interface{}, l *List) bool {
	switch {
	case Empty(l):
		return false
	case elem == Head(l):
		return true
	}
	return Elem(elem, Tail(l))
}

func NotElem(elem interface{}, l *List) bool {
	return !Elem(elem, l)
}

func ElemIndex(elem interface{}, l *List) (int, bool) {
	switch {
	case Empty(l):
		return -1, false
	case elem == Head(l):
		return 0, true
	}

	i, ok := ElemIndex(elem, Tail(l))
	return i + 1, ok
}

func Zip(l1, l2 *List) *List {
	if Empty(l1) || Empty(l2) {
		return New()
	}
	return Cons(NewWithElements(Head(l1), Head(l2)), Zip(Tail(l1), Tail(l2)))
}

func ZipWith(l1, l2 *List, f func(x, y interface{}) interface{}) *List {
	if Empty(l1) || Empty(l2) {
		return New()
	}
	return Cons(f(Head(l1), Head(l2)), ZipWith(Tail(l1), Tail(l2), f))
}

func TakeWhile(l *List, f func(x interface{}) bool) *List {
	if Empty(l) {
		return New()
	}

	if f(Head(l)) {
		return Cons(Head(l), TakeWhile(Tail(l), f))
	}

	return New()
}

func DropWhile(l *List, f func(x interface{}) bool) *List {
	if Empty(l) {
		return New()
	}

	if f(Head(l)) {
		return DropWhile(Tail(l), f)
	}

	return l
}

func Span(l *List, f func(x interface{}) bool) (first, rest *List) {
	return TakeWhile(l, f), DropWhile(l, f)
}

func Flatten(l *List) *List {
	return Foldr1(l, func(value, acc interface{}) interface{} {
		return Concatenate(value.(*List), acc.(*List))
	}).(*List)
}

var Concat = Flatten

func And(l *List) bool {
	return Foldr1(l, func(value, acc interface{}) interface{} {
		return value.(bool) && acc.(bool)
	}).(bool)
}

func Or(l *List) bool {
	return Foldr1(l, func(value, acc interface{}) interface{} {
		return value.(bool) || acc.(bool)
	}).(bool)
}

func Any(l *List, f func(x interface{}) bool) bool {
	return Foldr(false, l, func(value, acc interface{}) interface{} {
		return f(value) || acc.(bool)
	}).(bool)
}

func All(l *List, f func(x interface{}) bool) bool {
	return Foldr(true, l, func(value, acc interface{}) interface{} {
		return f(value) && acc.(bool)
	}).(bool)
}

func Group(l *List) *List {
	vec := [2]*List{New(), New()} // {final list, sublist}
	result := Foldr(vec, l, func(x, acc interface{}) interface{} {
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

func Partition(l *List, f func(x interface{}) bool) (satisfy, doNot *List) {
	vec := [2]*List{New(), New()} // {satisfy, doNot}
	result := Foldr(vec, l, func(x, acc interface{}) interface{} {
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
	table := make(map[interface{}]bool)
	unique := Foldl(New(), l, func(acc, x interface{}) interface{} {
		if _, ok := table[x]; ok {
			return acc
		}

		table[x] = true
		return Cons(x, acc.(*List))
	}).(*List)
	return Reverse(unique)
}

var Nub = Unique

func Delete(x interface{}, l *List) *List {
	without, with := Span(l, func(y interface{}) bool {
		return x != y
	})
	return Concatenate(without, Tail(with))
}

func Difference(base, subtract *List) *List {
	table := make(map[interface{}]bool)
	next := MakeIterator(subtract)
	for e := next(); e != nil; e = next() {
		table[e] = true
	}
	return Foldr(New(), base, func(x, acc interface{}) interface{} {
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
	table := make(map[interface{}]bool)
	next := MakeIterator(l1)
	for e := next(); e != nil; e = next() {
		table[e] = true
	}
	return Foldr(New(), l2, func(x, acc interface{}) interface{} {
		if _, ok := table[x]; ok {
			return Cons(x, acc.(*List))
		}
		return acc
	}).(*List)
}
