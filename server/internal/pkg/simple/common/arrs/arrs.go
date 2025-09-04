package arrs

func Contains[T comparable](arr []T, target T) bool {
	for _, val := range arr {
		if val == target {
			return true
		}
	}
	return false
}
