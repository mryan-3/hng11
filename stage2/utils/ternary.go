package utils

// Check is a function that takes in a boolean and two values
// of any type and returns one of the values based on the boolean condition
//
// Example:
//
// ternary.Check(true, "Yes", "No") // returns "Yes"
// ternary.Check(false, "Yes", "No") // returns "No"
// ternary.Check(1 == 1, "Yes", "No") // returns "Yes"
func Check[T any](conditional bool, a T, b T) T {
	if conditional {
		return a
	} else {
		return b
	}
}
