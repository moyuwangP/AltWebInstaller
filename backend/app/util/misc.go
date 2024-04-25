package util

func KeyBy[V any, K comparable](array []V, callback func(V) K) map[K]V {
	m := map[K]V{}
	for _, item := range array {
		m[callback(item)] = item
	}
	return m
}
