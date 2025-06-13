package pkg

func SliceToMap[T any](list []T, getKey func(T) string) map[string]T {
	result := make(map[string]T)
	for _, item := range list {
		key := getKey(item)
		result[key] = item
	}
	return result
}

func SliceToAnyMap[T any](list []T, keyFn func(T) string) map[string]any {
	result := make(map[string]any, len(list))
	for _, item := range list {
		result[keyFn(item)] = item
	}
	return result
}
