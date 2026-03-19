package einolib

import (
	"sync"
)

// 并发安全的key-value存储，key必须为可比较类型（如string、int、float等）
type SyncMap struct {
	mu sync.RWMutex
	m  map[interface{}]interface{}
}

func NewSyncMap() *SyncMap {
	return &SyncMap{
		m: make(map[interface{}]interface{}),
	}
}

func (s *SyncMap) Set(key interface{}, val interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = val
}

func (s *SyncMap) Get(key interface{}) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.m[key]
	return val, ok
}

func (s *SyncMap) Has(key interface{}) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.m[key]
	return ok
}

func (s *SyncMap) Delete(key interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
}

func (s *SyncMap) Keys() []interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]interface{}, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}
