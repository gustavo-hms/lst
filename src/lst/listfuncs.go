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

func Map(l *List, f func(interface{}) interface{}) interface{} {
	return Foldr(New(), l, func(value, acc interface{}) interface{} {
		return Cons(f(value), acc.(*List))
	})
}

func Filter(l *List, f func(interface{}) bool) interface{} {
	return Foldr(New(), l, func(value, acc interface{}) interface{} {
		if f(value) {
			return Cons(value, acc.(*List))
		}
		return acc
	})
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

/*
type sortable interface {
	Less(otherValue sortable) bool
}

func Max(l *List) interface{} {
	if _, ok := Head(l).(sortable); ok {
		return sortableMax(l)
	}
	return operatorMax(l)
}

func sortableMax(l *List) sortable {
	result := Foldr1(l, func(e, acc interface{}) interface{} {
		if e.(sortable).Less(acc.(sortable)) {
			return acc
		}
		return e
	})
	return result.(sortable)
}

func operatorMax(l *List) interface{} {
	t := reflect.TypeOf(Head(l))
	return Foldr1(l, func(e, acc interface{}) interface{} {
		if e.(t) < acc.(t) {
			return acc
		}
		return e
	})
}*/
