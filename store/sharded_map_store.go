package store

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

const (
	shardCount = 32

	// The EvictionPolicy is not implemented
	EvictionVolatileRandom EvictionPolicy = 0
	EvictionVolatileLRU EvictionPolicy    = 1
	EvictionAllRandom EvictionPolicy      = 2
	EvictionAllLRU EvictionPolicy         = 3

	maxInt64 = int64(^uint64(0)>>1)
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
	m map[string]*entry
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
		sm.m = make(map[string]*entry)
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
	// TODO: This set method is very naive. It doesn't put any restriction on the input type.
	// Also, when it get a value, it remove all the pointers above it find the real data. Only save the real data.

	// if timeout == 0, the key will never expire
	deadline := time.Now().Add(timeout).UnixNano()
	if timeout == 0 {
		deadline = maxInt64
	}

	sm := s.selectSharedMap(key)
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.m[key] = &entry{
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
	// TODO: This is terrible. If return the data directly, users can edit the data outside the cache store.
	return v.data, nil
}

func (s *shardedMapStore) Delete(key string) error {
	sm := s.selectSharedMap(key)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.m, key)
	return nil
}

// Increase increase the number stored at the key by one. Set the value to 1 if the key is not exist.
// Return an error if the stored value is not a "integer". Support int, uint32, uint64
func (s *shardedMapStore) Increase(key string) error {
	sm := s.selectSharedMap(key)
	sm.mu.Lock()
	defer sm.mu.Unlock()

	v, ok := sm.m[key]
	if !ok {
		sm.m[key] = &entry{
			data: 1,
			deadline: maxInt64,
		}
		return nil
	}
	switch data := v.data.(type) {
	case int:
		sm.m[key].data = data + 1
	case uint32:
		sm.m[key].data = data + 1
	case uint64:
		sm.m[key].data = data + 1
	default:
		return fmt.Errorf("value_is_not_number_type")
	}
	return nil
}

func (s *shardedMapStore) Close() error {
	s.shardedMaps = nil
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

// dumpAllJSON print all the data in cache in json format including the timeout data
func (s *shardedMapStore) DumpAllJSON() (string, error) {
	// TODO: Maybe can support also dump the timeout of each cache key?
	// TODO: Support limiting the output
	totalSize := 0
	for i := 0; i < len(s.shardedMaps); i++ {
		sm := &s.shardedMaps[i]
		sm.mu.RLock()
		totalSize += len(sm.m)
		sm.mu.RUnlock()
	}

	res := make(map[string]interface{}, totalSize)
	for i := 0; i < len(s.shardedMaps); i++ {
		sm := &s.shardedMaps[i]
		sm.mu.RLock()
		for key, entryValue := range sm.m {
			res[key] = entryValue.data
		}
		sm.mu.RUnlock()
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(resBytes), nil
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

// evictVolatileRandomShardedMap loop though the sharded map and evict the expired key
func evictVolatileRandomShardedMap(sm *shardedMap) {
	// TODO: implement this
}