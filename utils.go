package lst

/*
 * Utility functions not directly related with lists
 */

// Flips function's arguments
func Flip(f func(x, y interface{}) interface{}) func(interface{}, interface{}) interface{} {
	return func(x, y interface{}) interface{} {
		return f(y, x)
	}
}
