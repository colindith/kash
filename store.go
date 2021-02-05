package kash

import (
	"fmt"
	"sync"
	"time"
)

const (
	triggeringEvictionOptNum = 100
)

type store interface {
	set(key string, value interface{}, timeout time.Duration) error
	get(key string) (interface{}, error)
	delete(key string) error
}

// defaultStore implement with build-in map. Most naive implementation
type defaultStore struct {
	m map[string]unit
	operationCount uint    // memo the number of keys mutated since last time eviction
	mu sync.RWMutex        // The whole map share a single lock is inefficient
}

type unit struct {
	data interface{}
	deadline int64    // timestamp nanosecond
}

func (s *defaultStore) set(key string, value interface{}, timeout time.Duration) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = unit{
		data: value,
		deadline: time.Now().Add(timeout).UnixNano(),
	}
	s.operationCount++
	go func() {
		if s.operationCount >= triggeringEvictionOptNum {
			// TODO: Except for this, there are other occurrences that should trigger eviction
			evictMap(s)   // This function would require the lock
			s.operationCount = 0
		}
	}()
	return nil
}

func (s *defaultStore) get(key string) (value interface{}, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.m[key]
	if !ok {
		return nil, fmt.Errorf("error_cache_not_found")
	}
	if time.Now().UnixNano() > v.deadline {
		// The key was timeout. Evict it.
		delete(s.m, key)
		return nil, fmt.Errorf("error_cache_not_found")
	}
	return v.data, nil
}

func (s *defaultStore) delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
	return nil
}

func getDefaultStore() store {
	return &defaultStore{
		m: make(map[string]unit),
	}
}

// evictMap loop though the map and evict the expired key
func evictMap(s *defaultStore) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UnixNano()
	for k, u := range s.m {
		if u.deadline <= now {
			delete(s.m, k)      // delete map key while loop through map. Is it okay?
		}
	}
}

const (
	shardCount = 32
)

// shardedMapStore
type shardedMapStore struct {
	shardedMaps []shardedMap
}

type shardedMap struct {
	m map[string]entry
	opCount uint    // memo the number of keys mutated since last time eviction
	mu sync.RWMutex
}

type entry struct {
	data interface{}
	deadline int64    // timestamp nanosecond
}


func getShardedMapStore() store {
	s := &shardedMapStore{
		shardedMaps: make([]shardedMap, shardCount),
	}
	i := 0
	for i < len(s.shardedMaps) {
		sm := &s.shardedMaps[i]
		sm.mu.Lock()
		sm.m = make(map[string]entry)
		sm.mu.Unlock()
	}
	return s
}

func (s *shardedMapStore) selectSharedMap(key string) *shardedMap {
	return &s.shardedMaps[fnv32(key)%shardCount]
}

func (s *shardedMapStore) set(key string, value interface{}, timeout time.Duration) (err error) {
	sm := s.selectSharedMap(key)
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.m[key] = entry{
		data: value,
		deadline: time.Now().Add(timeout).UnixNano(),
	}
	sm.opCount++
	go func() {
		if sm.opCount >= triggeringEvictionOptNum {
			// TODO: Except for this, there are other occurrences that should trigger eviction
			evictShardedMap(sm)   // This function would require the lock
			sm.opCount = 0
		}
	}()
	return nil
}

func (s *shardedMapStore) get(key string) (value interface{}, err error) {
	sm := s.selectSharedMap(key)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	v, ok := sm.m[key]
	if !ok {
		return nil, fmt.Errorf("error_cache_not_found")
	}
	if time.Now().UnixNano() > v.deadline {
		// The key was timeout. Evict it.
		delete(sm.m, key)
		return nil, fmt.Errorf("error_cache_not_found")
	}
	return v.data, nil
}

func (s *shardedMapStore) delete(key string) error {
	sm := s.selectSharedMap(key)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.m, key)
	return nil
}

// evictShardedMap loop though the sharded map and evict the expired key
func evictShardedMap(sm *shardedMap) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	now := time.Now().UnixNano()
	for k, e := range sm.m {
		if e.deadline <= now {
			delete(sm.m, k)
		}
	}
}