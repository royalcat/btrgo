package btrstruct

import "sync"

// REWORK
type Tree[K comparable, V any] struct {
	m    sync.RWMutex
	root treeBranch[K, V]
}

type treeBranch[K comparable, V any] struct {
	branches map[K]treeBranch[K, V]
	value    V
	hasValue bool
}

func (t *Tree[K, V]) Get(branch []K) (V, bool) {
	t.m.RLock()
	defer t.m.RUnlock()

	cur := t.root
	var ok bool
	for _, key := range branch {
		cur, ok = cur.branches[key]
		if !ok {
			return cur.value, false
		}
	}

	if cur.hasValue {
		return cur.value, true
	}

	return cur.value, false
}

func (t *Tree[K, V]) Set(branch []K, value V) {
	t.m.Lock()
	defer t.m.Unlock()

	cur := t.root
	var ok bool
	for _, key := range branch {
		cur, ok = cur.branches[key]
		if !ok {
			return
		}
	}

	cur.value = value
	cur.hasValue = true
}

func (t *Tree[K, V]) Has(branch []K) bool {

	cur := t.root
	var ok bool
	for _, key := range branch {
		cur, ok = cur.branches[key]
		if !ok {
			return false
		}
	}

	return cur.hasValue
}
