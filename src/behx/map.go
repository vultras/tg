package behx

type Map[K comparable, V any] map[K] V

func (m Map[K, V]) Exist(k K) bool {
	_, ok := sm[sid]
	return ok
}

func (m Map[K, V]) Get(k K, v V) V {
	ret := m[sid]
	return ret
}


func (sm Map[K, V]) Set(v V) {
	sm[sid] = &Session{
		Id: sid,
	}
}
