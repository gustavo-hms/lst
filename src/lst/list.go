package lst

/*
 * Este arquivo contém a definição da estrutura List, bem como todas as funções 
 * que dependem de algum conhecimento mais interno dessa estrutura
 */

import (
	"fmt"
	"strings"
)

type List struct {
	elements   []interface{}
	firstEmpty *int // O primeiro índice vazio no vetor a que o slice faz referência
	firstUsed  int  // O primeiro índice que o slice enxerga no vetor
}

func New() *List {
	l := new(List)
	var vec [64]interface{}
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

func NewFromSlice(slice []interface{}) (l *List) {
	l = new(List)
	l.elements = make([]interface{}, len(slice))
	copy(l.elements, slice) // Copiando para evitar modificações inadvertidas
	var i int
	l.firstEmpty = &i
	return
}

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

func get(l *List, i int) interface{} {
	last := Len(l) - 1
	return l.elements[last-i]
}

func set(l *List, i int, value interface{}) {
	last := Len(l) - 1
	l.elements[last-i] = value
}

func Head(l *List) interface{} {
	return get(l, 0)
}

func Tail(l *List) (tailList *List) {
	tailList = new(List)
	tailList.elements = l.elements[:Len(l)-1]
	tailList.firstEmpty = l.firstEmpty
	tailList.firstUsed = l.firstUsed
	return
}

func Last(l *List) interface{} {
	return get(l, Len(l)-1)
}

func Init(l *List) (initList *List) {
	initList = new(List)
	initList.elements = l.elements[1:]
	initList.firstEmpty = l.firstEmpty
	initList.firstUsed = l.firstUsed + 1
	return
}

func Cons(value interface{}, l *List) (newl *List) {
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
		newl = NewFromSlice(l.elements)
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

func concat(l1, l2 *List) (con *List) {
	if Len(l1) == 0 {
		con = NewFromList(l2)
		return
	}

	con = Cons(Head(l1), concat(Tail(l1), l2))
	return
}

func Concat(lists ...*List) (con *List) {
	last := len(lists) - 1
	con = lists[last]
	for i := last - 1; i >= 0; i-- {
		con = concat(lists[i], con)
	}
	return
}

func Foldr(init interface{}, l *List, f func(value, acc interface{}) interface{}) (accum interface{}) {
	accum = init
	for _, v := range l.elements {
		accum = f(v, accum)
	}
	return
}

func Foldl(init interface{}, l *List, f func(acc, value interface{}) interface{}) (accum interface{}) {
	accum = init
	for i := 0; i < Len(l); i++ {
		accum = f(accum, get(l, i))
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
		rev = Cons(get(l, i), rev)
	}
	return rev
}
