package btrgo

func KeysOfMap[K comparable, V any](m map[K]V) []K {
	keys := make([]K, len(m))

	i := 0
	for key := range m {
		keys[i] = key
		i++
	}

	return keys
}

func ValuesOfMap[K comparable, V any](m map[K]V) []V {
	values := make([]V, len(m))

	i := 0
	for _, v := range m {
		values[i] = v
		i++
	}

	return values
}

func FilterMap[K comparable, V any](keys []K, data map[K]V) map[K]V {
	newMap := make(map[K]V, len(keys))

	for _, key := range keys {
		newMap[key] = data[key]
	}

	return newMap
}
