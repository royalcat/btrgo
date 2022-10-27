package btrgo

import "golang.org/x/exp/constraints"

// slice operations without size changes

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

func Sort[T constraints.Ordered](x []T) {
	n := len(x)
	for {
		swapped := false
		for i := 1; i < n; i++ {
			if x[i] < x[i-1] {
				x[i-1], x[i] = x[i], x[i-1]
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}
