package lst

import (
	"testing"
)

func TestReverse(t *testing.T) {
	l := NewFromSlice(elements[:])
	r := Reverse(l)

	if Len(r) != Len(l) {
		t.Error("Reversed list doesn't have the same number of arguments as its original list")
	}

	it := MakeReverseIterator(r)
	index := 0
	for elem := it(); elem != nil; elem = it() {
		if elem != elements[index] {
			t.Errorf("Mismatched elements at index %d", index)
		}

		index++
	}
}

func TestEmpty(t *testing.T) {
	l := New()

	if !Empty(l) {
		t.Error("Function Empty reporting a non-empty list")
	}

	l2 := Cons(1, l)
	if Empty(l2) {
		t.Error("Function Empty reporting an empty list")
	}
}

func TestFoldr1(t *testing.T) {
	l := NewFromSlice(elements[:])

	listSum := Foldr1(l, func(e Elem, accum interface{}) interface{} {
		return e.(int) + accum.(int)/2
	}).(int)

	sliceSum := 0
	for i := N - 1; i >= 0; i-- {
		sliceSum = elements[i].(int) + sliceSum/2
	}

	if listSum != sliceSum {
		t.Errorf("Got different sums: %d and %d", listSum, sliceSum)
	}
}

func TestFoldl1(t *testing.T) {
	l := NewFromSlice(elements[:])

	listSum := Foldl1(l, func(accum interface{}, e Elem) interface{} {
		return e.(int) + accum.(int)/2
	}).(int)

	sliceSum := 0
	for _, v := range elements[:] {
		sliceSum = v.(int) + sliceSum/2
	}

	if listSum != sliceSum {
		t.Errorf("Got different sums: %d and %d", listSum, sliceSum)
	}
}

func TestMap(t *testing.T) {
	l := NewFromSlice(elements[:])
	isEvenList := Map(l, func(x Elem) Elem {
		return x.(int)%2 == 0
	})

	for k, v := range elements[:] {
		isEven := v.(int)%2 == 0
		if isEven != Get(isEvenList, k) {
			t.Errorf("Mismatched elements at index %d", k)
		}
	}
}

func TestFilter(t *testing.T) {
	l := NewFromSlice(elements[:])
	evenList := Filter(l, func(x Elem) bool {
		return x.(int)%2 == 0
	})

	evenSlice := make([]Elem, 0, N)
	for _, v := range elements[:] {
		if v.(int)%2 == 0 {
			evenSlice = append(evenSlice, v)
		}
	}

	for k, v := range evenSlice {
		if v != Get(evenList, k) {
			t.Errorf("Mismatched elements at index %d", k)
		}
	}
}

func TestElement(t *testing.T) {
	l := New()
	for i := 0; i < N; i++ {
		l = Cons(1, l)
	}

	if Element(2, l) {
		t.Error("Function Element reporting number 2 is in the list")
	}

	set(l, N-3, 2)

	if !Element(2, l) {
		t.Error("Didn't find number 2 in the list")
	}
}

func TestNotElement(t *testing.T) {
	l := New()
	for i := 0; i < N; i++ {
		l = Cons(1, l)
	}

	if !NotElement(2, l) {
		t.Error("Function Element reporting number 2 is in the list")
	}

	set(l, N-3, 2)

	if NotElement(2, l) {
		t.Error("Didn't find number 2 in the list")
	}
}

func TestElemIndex(t *testing.T) {
	l := New()
	for i := 0; i < N; i++ {
		l = Cons(1, l)
	}

	ind, found := ElemIndex(2, l)
	if found {
		t.Error("Wrongly found element 2 in list")
	}

	set(l, N-3, 2)

	ind, found = ElemIndex(2, l)
	if !found {
		t.Error("Didn't find element 2 in list")
	}
	if ind != N-3 {
		t.Error("Element 2 found in wrong position: %d instead of %d", ind, N-3)
	}
}

func TestElemIndices(t *testing.T) {
	l := New()
	for i := 0; i < N; i++ {
		l = Cons(1, l)
	}

	set(l, N-10, 2)
	set(l, N-7, 2)
	set(l, N-3, 2)

	ind := ElemIndices(2, l)

	if Len(ind) != 3 {
		t.Error("Found %d times intead of 3", Len(ind))
	}

	if Head(ind) != N-10 || Get(ind, 1) != N-7 || Get(ind, 2) != N-3 {
		t.Error("Found at indices %d, %d and %d instead of %d, %d and %d", Head(ind), Get(ind, 1), Get(ind, 2), N-10, N-7, N-3)
	}
}

func TestZip(t *testing.T) {
	l1 := NewFromSlice(elements[:N/2])
	l2 := NewFromSlice(elements[N/2:])

	zip := Zip(l1, l2)

	var min int
	if Len(l1) < Len(l2) {
		min = Len(l1)
	} else {
		min = Len(l2)
	}

	if Len(zip) != min {
		t.Errorf("Wrong list's length: %d instead of %d", Len(zip), min)
	}

	for i := 0; i < min; i++ {
		elem := Get(zip, i).(*List)
		if Head(elem) != Get(l1, i) || Last(elem) != Get(l2, i) {
			t.Errorf("Mismatched elements at index %d", i)
		}
	}
}

func TestZipWith(t *testing.T) {
	l1 := NewFromSlice(elements[:N/2])
	l2 := NewFromSlice(elements[N/2:])

	zip := ZipWith(l1, l2, func(x, y Elem) Elem {
		return x.(int) + y.(int)
	})

	var min int
	if Len(l1) < Len(l2) {
		min = Len(l1)
	} else {
		min = Len(l2)
	}

	if Len(zip) != min {
		t.Errorf("Wrong list's length: %d instead of %d", Len(zip), min)
	}

	for i := 0; i < min; i++ {
		if Get(zip, i) != Get(l1, i).(int)+Get(l2, i).(int) {
			t.Errorf("Mismatched elements at index %d", i)
		}
	}
}

func TestTakeWhile(t *testing.T) {
	l := New()
	for i := 0; i < N; i++ {
		l = Cons(1, l)
	}

	set(l, N-10, 2)
	set(l, N-7, 2)
	set(l, N-3, 2)

	l2 := TakeWhile(l, func(x Elem) bool {
		return x != 2
	})

	if Len(l2) != N-10 {
		t.Errorf("New list has %d elements instead of %d", Len(l2), N-10)
	}
}

func TestDropWhile(t *testing.T) {
	l := New()
	for i := 0; i < N; i++ {
		l = Cons(1, l)
	}

	set(l, N-10, 2)
	set(l, N-7, 2)
	set(l, N-3, 2)

	l2 := DropWhile(l, func(x Elem) bool {
		return x != 2
	})

	if Len(l2) != 10 {
		t.Errorf("New list has %d elements instead of %d", Len(l2), 10)
	}
}

func TestSpan(t *testing.T) {
	l := New()
	for i := 0; i < N; i++ {
		l = Cons(1, l)
	}

	set(l, N-10, 2)
	set(l, N-7, 2)
	set(l, N-3, 2)

	l1, l2 := Span(l, func(x Elem) bool {
		return x != 2
	})

	if Len(l1) != N-10 {
		t.Errorf("First list has %d elements instead of %d", Len(l1), N-10)
	}

	if Len(l2) != 10 {
		t.Errorf("Second list has %d elements instead of %d", Len(l2), 10)
	}
}

func TestFlatten(t *testing.T) {
	l1 := NewFromSlice(elements[:N/2])
	l2 := NewFromSlice(elements[N/2:])

	zip := Zip(l1, l2)
	flat := Flatten(zip)

	var min int
	if Len(l1) < Len(l2) {
		min = Len(l1)
	} else {
		min = Len(l2)
	}
	if Len(flat) != 2*min {
		t.Errorf("List has %d elements instead of %d", Len(flat), 2*min)
	}

	for i := 0; i < min; i++ {
		if Get(flat, 2*i) != Get(l1, i) {
			t.Errorf("Mismatched elements at index %d", 2*i)
			break
		}
		if Get(flat, 2*i+1) != Get(l2, i) {
			t.Errorf("Mismatched elements at index %d", 2*i+1)
			break
		}
	}
}

func TestAnd(t *testing.T) {
	l := New()
	for i := 0; i < N; i++ {
		l = Cons(true, l)
	}

	if !And(l) {
		t.Error("Returning false for an all true list")
	}

	set(l, N-2, false)

	if And(l) {
		t.Error("Returning true for a list with a false element")
	}
}

func TestOr(t *testing.T) {
	l := New()
	for i := 0; i < N; i++ {
		l = Cons(false, l)
	}

	if Or(l) {
		t.Error("Returning true for an all false list")
	}

	set(l, N-2, true)

	if !Or(l) {
		t.Error("Returning false for a list with a true element")
	}
}

func TestAll(t *testing.T) {
	l := New()
	for i := 0; i < N; i++ {
		l = Cons(1, l)
	}

	allOne := All(l, func(x Elem) bool {
		return x == 1
	})

	if !allOne {
		t.Error("Returning false for a true statement")
	}

	set(l, N-2, 0)

	allOne = All(l, func(x Elem) bool {
		return x == 1
	})

	if allOne {
		t.Error("Returning true for a false statement")
	}
}

func TestAny(t *testing.T) {
	l := New()
	for i := 0; i < N; i++ {
		l = Cons(1, l)
	}

	anyZero := Any(l, func(x Elem) bool {
		return x == 0
	})

	if anyZero {
		t.Error("Returning true for a false statement")
	}

	set(l, N-2, 0)

	anyZero = Any(l, func(x Elem) bool {
		return x == 0
	})

	if !anyZero {
		t.Error("Returning false for a true statement")
	}
}
