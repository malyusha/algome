package util

// Filter returns only those elements, that are matched by given filter func.
func Filter[T any](items []T, filter func(item T) bool) []T {
	n := 0
	for _, item := range items {
		if filter(item) {
			items[n] = item
			n++
		}
	}

	return items[:n]
}
