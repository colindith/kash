package store

import (
	"fmt"
	"sync"
	"time"
)

const (
	shardCount = 32

	EvictionRandom EvictionPolicy = 0
)

type EvictionPolicy uint32

// shardedMapStore
type shardedMapStore struct {
	Store
	shardedMaps []shardedMap

	defaultTimeout time.Duration
	maxMemory int64              // unit: bytes
	evictionPolicy EvictionPolicy
}

type shardedMap struct {
	m map[string]entry
	mu sync.RWMutex

	opCount uint    // memo the number of keys mutated since last time eviction
}

type entry struct {
	data interface{}
	deadline int64    // timestamp nanosecond
}

func GetShardedMapStore(opts... Option) Store {
	s := &shardedMapStore{
		shardedMaps: make([]shardedMap, shardCount),
	}
	i := 0
	for i < len(s.shardedMaps) {
		sm := &s.shardedMaps[i]
		sm.mu.Lock()
		sm.m = make(map[string]entry)
		sm.mu.Unlock()
		i++
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *shardedMapStore) selectSharedMap(key string) *shardedMap {
	return &s.shardedMaps[fnv32(key)%shardCount]
}

func (s *shardedMapStore) Set(key string, value interface{}) (err error) {
	return s.SetWithTimeout(key, value, s.defaultTimeout)
}

func (s *shardedMapStore) SetWithTimeout(key string, value interface{}, timeout time.Duration) (err error) {
	// if timeout == 0, the key will never expire
	deadline := time.Now().Add(timeout).UnixNano()
	if timeout == 0 {
		deadline = maxTime.UnixNano()
	}

	sm := s.selectSharedMap(key)
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.m[key] = entry{
		data: value,
		deadline: deadline,
	}
	sm.opCount++
	if sm.opCount >= triggeringEvictionOptNum {
		go func() {
			// TODO: Except for this, there are other occurrences that should trigger eviction
			evictShardedMap(sm) // This function would require the lock
			sm.opCount = 0
		}()
	}
	return nil
}

func (s *shardedMapStore) Get(key string) (value interface{}, err error) {
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

func (s *shardedMapStore) Delete(key string) error {
	sm := s.selectSharedMap(key)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.m, key)
	return nil
}

func (s *shardedMapStore) setDefaultTimeout(timeout time.Duration) {
	s.defaultTimeout = timeout
}

func (s *shardedMapStore) setEvictionPolicy(policy EvictionPolicy) {
	s.evictionPolicy = policy
}

func (s *shardedMapStore) setMaxMemory(size int64) {
	s.maxMemory = size
}

// evictShardedMap loop though the sharded map and evict the expired key
func evictShardedMap(sm *shardedMap) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	now := time.Now().UnixNano()
	for k, e := range sm.m {    // TODO: This is incorrect. This cannot loop through the map
		if e.deadline <= now {
			delete(sm.m, k)
		}
	}
}