lst
===

A simple implementation of lists in Go using vectors

What's the difference between this package and the standard `container/list`?  
There are some. First of all, the latter is implemented using linked lists, and 
the former uses Go vectors as the underlying data. While `container/list` 
resembles C++ list implementation in some way, `lst` approach tries to be 
closer to the Haskell one, even in the function's names provided by the 
package.  Finally, to allow the lists to be efficiently used in recursions, it 
uses a copy on write strategy, meaning creating a new list based on another one 
(using functions like `lst.NewFromList` and `lst.Tail`) only copies the 
reference to the underlying data.  A copy of the hole data is only processed 
when appending new elements would overwrite elements in the shared data. This 
is handled automatically by the package.

Usage
-----

### Creating a new list

To create a new empty list, type:

	l := lst.New()

You can also create a list filled with predefined elements:

	l := lst.NewWithElements(1, 2, 3, 4)

Or, as a shorthand,

	l := lst.L(1, 2, 3, 4)

If you want to use the elements of a slice:

	l := lst.NewFromSlice(aSlice)

### Adding elements

To add elements, use the Cons function:

	l := lst.Cons(element, l)

### Manipulation

Note that a list is never modified in place. Instead, a new list is created 
each time an operation is executed. Since the package uses copy on write 
optimisation, this does not allocates extra memory.

To add 17 to all elements in a list:

	added17 := lst.Map(l, func(elem lst.Elem) lst.Elem {
		return elem.(int) + 17
	})

For Haskell programmers, it's worthy noting that here the function to be 
applied is the last argument. It's a pattern used in this package for all 
higher order functions, just to be easier to write lambdas.

To take only the even elements in a list:

	even := lst.Filter(l, func(elem lst.Elem) bool {
		return elem.(int)%2 == 0
	})

Other functions from Haskell library are also implemented. You can take a look 
at the files inside the package to see them.
