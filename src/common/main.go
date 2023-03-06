package common

func RemoveDuplicates[T comparable](slice []T) []T {
	unique := make(map[T]bool)
	for _, val := range slice {
		unique[val] = true
	}
	result := make([]T, 0, len(unique))
	for key := range unique {
		result = append(result, key)
	}
	return result
}
