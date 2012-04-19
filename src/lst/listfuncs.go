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
