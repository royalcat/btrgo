package btrgo

type Set[V comparable] struct {
	m map[V]bool
}

func (s *Set[V]) Len() int { return len(s.m) }

func (s *Set[V]) Add(v V) {
	s.m[v] = true
}

func (s *Set[V]) Del(v V) {
	delete(s.m, v)
}

func (s *Set[V]) List() []V {
	arr := make([]V, 0, len(s.m))
	for v := range s.m {
		arr = append(arr, v)
	}
	return arr
}

func (s *Set[V]) Range(f func(V) bool) {
	for v := range s.m {
		if !f(v) {
			return
		}
	}
}
