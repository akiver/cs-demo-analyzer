package slice

func Contains[T comparable](collection []T, element T) bool {
	for _, item := range collection {
		if item == element {
			return true
		}
	}

	return false
}

func AppendIfNotInSlice[T comparable](slice []T, value T) []T {
	if !Contains(slice, value) {
		slice = append(slice, value)
	}

	return slice
}

type stringable interface {
	String() string
}

func ToStrings[T stringable](slice []T) []string {
	stringSlice := make([]string, len(slice))
	for index, value := range slice {
		stringSlice[index] = value.String()
	}

	return stringSlice
}

func Filter[T comparable](s []T, predicate func(item T, index int) bool) []T {
	result := make([]T, 0, len(s))

	for i, item := range s {
		if predicate(item, i) {
			result = append(result, item)
		}
	}

	return result
}
