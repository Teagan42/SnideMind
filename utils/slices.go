package utils

func Intersection[T comparable](a []T, b []T) []T {
	// Create a map to track elements in 'b'
	set := make(map[T]struct{}, len(b))
	for _, item := range b {
		set[item] = struct{}{}
	}

	// Collect elements from 'a' that are also in 'b'
	var intersection []T
	for _, item := range a {
		if _, found := set[item]; found {
			intersection = append(intersection, item)
		}
	}

	return intersection
}
