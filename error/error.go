package error

// Check is a helper function to abstract error handling
// away and make checking easier.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}
