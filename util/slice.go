package util

// GroupBy groups slice by parameter returned by provided func `getKey`.
// Returns map, where key - unique value, received from `getKey`, and value - group of elements.
func GroupBy[K comparable, V any](items []V, getKey func(V) K) map[K][]V {
	out := make(map[K][]V)

	for _, it := range items {
		key := getKey(it)

		out[key] = append(out[key], it)
	}

	return out
}

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
