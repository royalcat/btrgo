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
