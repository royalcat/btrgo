package btrgo

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

// return unique elements from given slice
// doesn't sort elements, expect output elements in random order
func SliceUnique[V comparable](s []V) []V {
	keys := make(map[V]bool, len(s))
	list := make([]V, 0, len(s))
	for _, entry := range s {
		if _, ok := keys[entry]; !ok {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
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

func SliceDiffSplited[V comparable](a, b []V) (dels, adds []V) {
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

	return KeysOfMap(aMap), KeysOfMap(bMap)
}
