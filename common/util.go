package common

// Test if a slice contains an element
func Contains[T comparable](testElem T, slice []T) bool {
	for _, sliceElem := range slice {
		if testElem == sliceElem {
			return true
		}
	}
	return false
}
