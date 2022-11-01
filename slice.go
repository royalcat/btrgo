package btrgo

import "golang.org/x/exp/constraints"

func Paginate[V any](x []V, skip int, size int) []V {
	if skip > len(x) {
		skip = len(x)
	}

	end := skip + size
	if end > len(x) {
		end = len(x)
	}

	return x[skip:end]
}

func PaginateInvert[V any](x []V, skip int, size int) []V {
	if size < 0 || skip < 0 {
		return []V{}
	}

	if skip > len(x) {
		return []V{}
	}

	end := len(x) - skip - size
	if end < 0 {
		end = 0
	}

	out := make([]V, 0, size)
	for i := len(x) - 1 - skip; i >= end; i-- {
		out = append(out, x[i])
	}
	return out
}

// Split slice in chunks of constant size
// based on https://stackoverflow.com/a/67011816
func Chunks[V any](xs []V, chunkSize int) [][]V {
	if chunkSize < 1 {
		return [][]V{xs}
	}
	if len(xs) == 0 {
		return nil
	}
	divided := make([][]V, (len(xs)+chunkSize-1)/chunkSize)
	prev := 0
	i := 0
	till := len(xs) - chunkSize
	for prev < till {
		next := prev + chunkSize
		divided[i] = xs[prev:next]
		prev = next
		i++
	}
	divided[i] = xs[prev:]
	return divided
}

func InSlice[V comparable](arr []V, val V) bool {

	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
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

func SliceDiffSplited[V comparable](a, b []V) (adds, dels []V) {
	aMap := map[V]bool{}
	for _, v := range a {
		aMap[v] = true
	}

	bMap := map[V]bool{}
	for _, v := range b {
		bMap[v] = true
		delete(aMap, v)
	}

	for _, v := range a {
		delete(bMap, v)
	}

	return KeysOfMap(bMap), KeysOfMap(aMap)

}
