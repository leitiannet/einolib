package einolib

import (
	"sync"
)

// 并发安全的key-value存储
type SyncMap struct {
	mu sync.RWMutex
	m  map[string]interface{}
}

func NewSyncMap() *SyncMap {
	return &SyncMap{
		m: make(map[string]interface{}),
	}
}

func (s *SyncMap) Set(key string, val interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = val
}

func (s *SyncMap) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.m[key]
	return val, ok
}

func (s *SyncMap) Has(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.m[key]
	return ok
}

func (s *SyncMap) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
}

func (s *SyncMap) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}
