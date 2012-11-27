package lst

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

const N = 100

var elements [N]Elem

func init() {
	for i := 0; i < N; i++ {
		elements[i] = rand.Int()
	}
}

func TestNewFromSlice(t *testing.T) {
	l := NewFromSlice(elements[:])
	for k, v := range elements {
		if v != l.elements[N-k-1] {
			t.Error("Some list's elements differ from the elements in slice")
		}
	}
}

func TestNewFromReversedSlice(t *testing.T) {
	l := newFromReversedSlice(elements[:])
	for k, v := range elements {
		if v != l.elements[k] {
			t.Error("Some list's elements differ from the elements in slice")
		}
	}
}

func TestGet(t *testing.T) {
	l := NewFromSlice(elements[:])
	for k, v := range elements {
		if v != Get(l, k) {
			t.Error("Some list's elements differ from the elements in slice")
		}
	}
}

func TestSet(t *testing.T) {
	l := NewFromSlice(elements[:])
	cp := make([]Elem, N)
	copy(cp, elements[:])
	for i := 0; i < N; i += 10 {
		randomNumber := rand.Int()
		cp[i] = randomNumber
		set(l, i, randomNumber)
	}

	for k, v := range cp {
		if v != Get(l, k) {
			t.Error("Some list's elements differ from the elements in slice")
		}
	}
}

func TestLen(t *testing.T) {
	l := New()
	if Len(l) != 0 {
		t.Errorf("Function Len reported a length of %d instead of 0", Len(l))
	}

	l = NewFromSlice(elements[:])
	if Len(l) != N {
		t.Errorf("Function Len reported a length of %d instead of %d", Len(l), N)
	}
}

func TestNewFromList(t *testing.T) {
	l1 := NewFromSlice(elements[:])
	l2 := NewFromList(l1)

	if Len(l1) != Len(l2) {
		t.Errorf("The two lists have different lengths: %d and %d", Len(l1), Len(l2))
	}

	for i := 0; i < N; i++ {
		if Get(l1, i) != Get(l2, i) {
			t.Errorf("Elements in position %d differ between the two lists", i)
		}
	}
}

func TestNewWithElements(t *testing.T) {
	l := NewWithElements(elements[:]...)
	for k, v := range elements[:] {
		if v != Get(l, k) {
			t.Errorf("Element in position %d differ from the original array", k)
		}
	}
}

func TestString(t *testing.T) {
	l := NewFromSlice(elements[:])
	var elementsAsString [N]string
	for k, v := range elements[:] {
		elementsAsString[k] = fmt.Sprintf("%v", v)
	}

	desired := "[" + strings.Join(elementsAsString[:], ", ") + "]"

	if desired != l.String() {
		t.Error("Wrong string representation")
	}
}

func TestMakeIterator(t *testing.T) {
	l := NewFromSlice(elements[:])
	it := MakeIterator(l)

	for k, v := range elements[:] {
		if v != it() {
			t.Errorf("Mismatched elements at index %d", k)
		}
	}

	if it() != nil {
		t.Error("Iterator didn't finish with nil")
	}
}

func TestMakeReverseIterator(t *testing.T) {
	l := NewFromSlice(elements[:])
	it := MakeReverseIterator(l)

	for i := N - 1; i >= 0; i-- {
		if elements[i] != it() {
			t.Errorf("Mismatched elements at index %d", i)
		}
	}

	if it() != nil {
		t.Error("Iterator didn't finish with nil")
	}
}

func TestHead(t *testing.T) {
	l := NewFromSlice(elements[:])
	if Head(l) != elements[0] {
		t.Errorf("Head returned a value of %d, but first element of array is %d", Head(l), elements[0])
	}
}

func TestTail(t *testing.T) {
	l := NewFromSlice(elements[:])
	tail := Tail(l)
	for k, v := range elements[1:] {
		if v != Get(tail, k) {
			t.Errorf("Mismatched elements at index %d", k)
		}
	}
}

func TestLast(t *testing.T) {
	l := NewFromSlice(elements[:])
	if Last(l) != elements[N-1] {
		t.Errorf("Head returned a value of %d, but first element of array is %d", Head(l), elements[N-1])
	}
}

func TestInit(t *testing.T) {
	l := NewFromSlice(elements[:])
	tail := Init(l)
	for k, v := range elements[:N-2] {
		if v != Get(tail, k) {
			t.Errorf("Mismatched elements at index %d", k)
		}
	}
}

func TestCons(t *testing.T) {
	l := New()
	for _, v := range elements {
		l = Cons(v, l)
	}

	if Len(l) != N {
		t.Errorf("Number of elements (%d) differ from the %d elements expected", Len(l), N)
	}

	it := MakeReverseIterator(l)
	i := 0
	for elem := it(); elem != nil; elem = it() {
		if elem != elements[i] {
			t.Errorf("Mismatched elements at index %d", N-i-1)
		}
		i++
	}

	l2 := Cons(Head(l).(int)-1, Tail(l))
	if Head(l) == Head(l2) {
		t.Error("Cons is overwritting elements in lists")
	}
}

func TestFoldr(t *testing.T) {
	l := NewFromSlice(elements[:])

	listSum := Foldr(0, l, func(e Elem, accum interface{}) interface{} {
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

func TestFoldl(t *testing.T) {
	l := NewFromSlice(elements[:])

	listSum := Foldl(0, l, func(accum interface{}, e Elem) interface{} {
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

func TestConcatenate(t *testing.T) {
	l1 := NewFromSlice(elements[0 : N/3])
	l2 := NewFromSlice(elements[N/3 : (2*N)/3])
	l3 := NewFromSlice(elements[(2*N)/3 : N])

	conc := Concatenate(l1, l2, l3)

	if Len(conc) != N {
		t.Errorf("Concatenated list has a length of %d, but %d was expected. Lengths: %d, %d, %d", Len(conc), N, Len(l1), Len(l2), Len(l3))
	}

	for k, v := range elements[:] {
		if v != Get(conc, k) {
			t.Errorf("Mismatched elements at index %d", k)
		}
	}
}
