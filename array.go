package btrgo

// operations on constant size slices

func InSlice[V comparable](arr []V, val V) bool {

	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func SliceDiff[V comparable](a, b []V) []V {
	m := make(map[V]bool)
	diff := make([]V, 0, len(a))

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}

func CompareSlices[V comparable](a, b []V) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func SliceUnique[V comparable](s []V) []V {
	keys := make(map[V]bool, len(s))
	list := make([]V, 0, len(s))
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return s
}

func RemoveElements[T comparable](slice []T, elem T) []T {
	n := make([]T, 0)
	for _, v := range slice {
		if v != elem {
			n = append(n, v)
		}
	}
	return n
}
