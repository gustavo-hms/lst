package lst

/*
 * Este arquivo contém a definição da estrutura List, bem como todas as funções 
 * que dependem de algum conhecimento mais interno dessa estrutura
 */

import (
	"fmt"
	"strings"
)

type Elem interface{}

type List struct {
	elements   []Elem
	firstEmpty *int // O primeiro índice vazio no vetor a que o slice faz referência
	firstUsed  int  // O primeiro índice que o slice enxerga no vetor
}

func New() *List {
	l := new(List)
	var vec [64]Elem
	l.elements = vec[0:0]
	var i int
	l.firstEmpty = &i

	return l
}

func NewFromList(original *List) (dest *List) {
	dest = new(List)
	dest.elements = original.elements[:]
	dest.firstEmpty = original.firstEmpty
	dest.firstUsed = original.firstUsed
	return
}

func newFromAlreadyReversedSlice(slice []Elem) (l *List) {
	l = new(List)
	l.elements = make([]Elem, len(slice))
	copy(l.elements, slice) // Copiando para evitar modificações inadvertidas
	i := len(slice)
	l.firstEmpty = &i
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

var L = NewWithElements // Apenas por conveniência

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

	return "[" + strings.Join(elems, " ") + "]"
}

func Get(l *List, i int) Elem {
	last := Len(l) - 1
	return l.elements[last-i]
}

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

func set(l *List, i int, value Elem) {
	last := Len(l) - 1
	l.elements[last-i] = value
}

func Head(l *List) Elem {
	return Get(l, 0)
}

func Tail(l *List) (tailList *List) {
	tailList = new(List)
	tailList.elements = l.elements[:Len(l)-1]
	tailList.firstEmpty = l.firstEmpty
	tailList.firstUsed = l.firstUsed
	return
}

func Last(l *List) Elem {
	return Get(l, Len(l)-1)
}

func Init(l *List) (initList *List) {
	initList = new(List)
	initList.elements = l.elements[1:]
	initList.firstEmpty = l.firstEmpty
	initList.firstUsed = l.firstUsed + 1
	return
}

func Cons(value Elem, l *List) (newl *List) {
	newl = new(List)
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
		newl = newFromAlreadyReversedSlice(l.elements)
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

	newl.elements = append(newl.elements, value)
	return
}

func concatenate(l1, l2 *List) (con *List) {
	if Len(l1) == 0 {
		con = NewFromList(l2)
		return
	}

	con = Cons(Head(l1), concatenate(Tail(l1), l2))
	return
}

func Concatenate(lists ...*List) (con *List) {
	last := len(lists) - 1
	con = lists[last]
	for i := last - 1; i >= 0; i-- {
		con = concatenate(lists[i], con)
	}
	return
}

func Foldr(init interface{}, l *List, f func(value Elem, acc interface{}) interface{}) (accum interface{}) {
	accum = init
	for _, v := range l.elements {
		accum = f(v, accum)
	}
	return
}

func Foldl(init interface{}, l *List, f func(acc interface{}, value Elem) interface{}) (accum interface{}) {
	accum = init
	for i := 0; i < Len(l); i++ {
		accum = f(accum, Get(l, i))
	}
	return
}

/*
 * A função Reverse foi colocada neste arquivo por eficiência: implantá-la como 
 * uma função recursiva usando concatenação é muito oneroso
 */
func Reverse(l *List) (rev *List) {
	rev = New()
	for i := 0; i < Len(l); i++ {
		rev = Cons(Get(l, i), rev)
	}
	return
}
