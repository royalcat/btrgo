package btrgo

// Not thread safe
type Set[V comparable] struct {
	m map[V]bool
}

func SetFromSlice[V comparable](slice []V) Set[V] {
	m := make(map[V]bool, len(slice)/2)
	for _, v := range slice {
		m[v] = true
	}

	return Set[V]{m: m}
}

func (s *Set[V]) Len() int {
	if s.m == nil {
		return 0
	}
	return len(s.m)
}

func (s *Set[V]) Add(v V) {
	if s.m == nil {
		s.m = make(map[V]bool)
	}

	s.m[v] = true
}

func (s *Set[V]) Del(v V) {
	if s.m == nil {
		return
	}

	delete(s.m, v)
}

func (s *Set[V]) List() []V {
	if s.m == nil {
		return []V{}
	}

	arr := make([]V, 0, len(s.m))
	for v := range s.m {
		arr = append(arr, v)
	}
	return arr
}

func (s *Set[V]) Range(f func(V) bool) {
	if s.m == nil {
		return
	}

	for v := range s.m {
		if !f(v) {
			return
		}
	}
}
