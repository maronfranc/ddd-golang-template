package util

// ArrayContains checks if a value is present in an array
func ArrayContains[T comparable](arr []T, value T) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
