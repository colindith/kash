package store

import (
	"encoding/json"
	"sync"
	"time"
)

const (
	shardCount = 32

	// The EvictionPolicy is not implemented
	EvictionRandom EvictionPolicy = 0
	EvictionLRU    EvictionPolicy = 1

	maxInt64 = int64(^uint64(0)>>1)
)

type EvictionPolicy uint32

// shardedMapStore
type shardedMapStore struct {
	Store
	shardedMaps []shardedMap

	defaultTimeout time.Duration
	maxMemory int64              // unit: bytes
	capacity int                 // The max number of keys can be stored in cache. 0 means no limit. Default value is 0.
	evictionPolicy EvictionPolicy

	// LRU
	lru bool
	head *entry
	tail *entry
	length int                   // current key count
	linkedListMutex sync.Mutex   // TODO: LL operation should lock with this mutex
}

type shardedMap struct {
	m map[string]*entry
	mu sync.RWMutex

	opCount uint    // memo the number of keys mutated since last time eviction
}

type entry struct {
	data interface{}
	deadline int64    // timestamp nanosecond

	// For LRU
	key string        // TODO: This is bad cause it would need too many additional space. Maybe change it to *string?
	prev *entry
	next *entry
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

	if s.capacity != 0 && s.evictionPolicy == EvictionLRU {
		s.lru = true
	}

	return s
}

func (s *shardedMapStore) selectSharedMap(key string) *shardedMap {
	return &s.shardedMaps[fnv32(key)%shardCount]
}

func (s *shardedMapStore) Set(key string, value interface{}) ErrorCode {
	return s.SetWithTimeout(key, value, s.defaultTimeout)
}

func (s *shardedMapStore) SetWithTimeout(key string, value interface{}, timeout time.Duration) ErrorCode {
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

	var e *entry
	var ok bool
	if e, ok = sm.m[key]; ok {
		// Avoid create new entry obj to reduce non-necessary allocation
		e.data = value
		e.deadline = deadline
	} else {
		if s.lru {
			e = &entry{
				data:     value,
				deadline: deadline,
				key:      key,
			}
		} else {
			e = &entry{
				data:     value,
				deadline: deadline,
			}
		}
		s.length++
		sm.m[key] = e
	}
	sm.opCount++


	if s.lru {
		s.linkedListMutex.Lock()
		if s.head == nil {
			// first key case
			s.head = e
			s.tail = s.head
		} else {
			if ok {
				// The key already be in cache. Don't need to create new node in linked list
				s.moveEntryToFront(e)
			} else {
				s.addEntryToFront(e)

				if s.length > s.capacity {
					s.lruEvict()
				}
			}
		}
		s.linkedListMutex.Unlock()
	}

	if sm.opCount >= triggeringEvictionOptNum {
		go func() {
			// TODO: Except for this, there are other occurrences that should trigger eviction
			evictShardedMap(sm) // This function would require the lock
			sm.opCount = 0
		}()
	}
	return Success
}

func (s *shardedMapStore) Get(key string) (value interface{}, code ErrorCode) {
	sm := s.selectSharedMap(key)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	e, ok := sm.m[key]
	if !ok {
		return nil, KeyNotFound
	}
	if time.Now().UnixNano() > e.deadline {
		// The key was timeout. Evict it.
		delete(sm.m, key)
		return nil, KeyNotFound
	}

	if s.lru {
		s.linkedListMutex.Lock()
		// move the entry to the head of linked list
		s.moveEntryToFront(e)
		s.linkedListMutex.Unlock()
	}

	// TODO: This is terrible. If return the data directly, users can edit the data outside the cache store.
	return e.data, Success
}

func (s *shardedMapStore) Delete(key string) ErrorCode {
	e := s.deleteKeyFromMap(key)
	if s.lru {
		s.linkedListMutex.Lock()
		s.evictEntryFromLL(e)
		s.linkedListMutex.Unlock()
		s.length--
	}
	return Success
}

// Increase increase the number stored at the key by one. Set the value to 1 if the key is not exist.
// Return an error if the stored value is not an "integer". Support int, uint32, uint64
func (s *shardedMapStore) Increase(key string) ErrorCode {
	sm := s.selectSharedMap(key)
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var e *entry

	e, ok := sm.m[key]
	if !ok {
		if s.lru {
			e = &entry{
				data:     1,
				deadline: maxInt64,
				key:      key,
			}
		} else {
			e = &entry{
				data:     1,
				deadline: maxInt64,
			}
		}
		sm.m[key] = e
		s.length++
	} else {
		switch data := e.data.(type) {
		case int:
			e.data = data + 1
		case uint32:
			e.data = data + 1
		case uint64:
			e.data = data + 1
		default:
			return ValueNotNumberType
		}
	}

	if s.lru {
		s.linkedListMutex.Lock()
		// TODO: The following codes are duplicated
		if s.head == nil {
			// first key case
			s.head = e
			s.tail = s.head
		} else {
			if ok {
				// The key already be in cache. Don't need to create new node in linked list
				s.moveEntryToFront(e)
			} else {
				s.addEntryToFront(e)

				if s.length > s.capacity {
					s.lruEvict()
				}
			}
		}
		s.linkedListMutex.Unlock()
	}

	return Success
}

func (s *shardedMapStore) GetTTL(key string) (int64, ErrorCode) {
	sm := s.selectSharedMap(key)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	v, ok := sm.m[key]
	if !ok {
		return 0, KeyNotFound
	}
	if time.Now().UnixNano() > v.deadline {
		// The key was timeout. Evict it.
		delete(sm.m, key)
		return 0, KeyNotFound
	}
	return v.deadline, Success
}

func (s *shardedMapStore) Close() ErrorCode {
	s.shardedMaps = nil
	return Success
}

func (s *shardedMapStore) setDefaultTimeout(timeout time.Duration) {
	s.defaultTimeout = timeout
}

func (s *shardedMapStore) setEvictionPolicy(policy EvictionPolicy) {
	// s method can only be called at init stage of cache
	s.evictionPolicy = policy
}

func (s *shardedMapStore) setMaxMemory(size int64) {
	s.maxMemory = size
}

func (s *shardedMapStore) setCapacity(cap int) {
	// s method can only be called at init stage of cache
	s.capacity = cap
}

// dumpAllJSON print all the data in cache in json format including the timeout data
func (s *shardedMapStore) DumpAllJSON() (string, ErrorCode) {
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
		return "", JSONMarshalErr
	}
	return string(resBytes), Success
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

func (s *shardedMapStore) moveEntryToFront(e *entry) {
	// TODO: There are some duplicated codes.
	if e == s.head {
		return
	} else if e == s.tail {
		prev := e.prev
		prev.next = nil
		s.tail = prev

		e.next = s.head
		s.head.prev = e
		s.head = e
	} else {
		prev := e.prev
		prev.next = e.next
		e.next.prev = prev

		e.next = s.head
		s.head.prev = e
		s.head = e
	}
}

func (s *shardedMapStore) addEntryToFront(e *entry) {
	if s.head == nil {
		s.head = e
		s.tail = s.head
		return
	}
	s.head.prev = e
	e.next = s.head
	s.head = e
}

func (s *shardedMapStore) evictTailFromLL() {
	if s.tail == nil {
		// nothing to evict
		return
	}
	prev := s.tail.prev
	if prev == nil {
		// only one element in linked list
		if s.tail != s.head {
			panic("tail of linked list has no prev")
		}
		s.head, s.tail = nil, nil
		return
	}
	s.tail = prev
	prev.next = nil
}

func (s *shardedMapStore) evictHeadFromLL() {
	if s.head == nil {
		// nothing to evict
		return
	}
	next := s.head.next
	if next == nil {
		// only one element in linked list
		if s.tail != s.head {
			panic("head of linked list has no next")
		}
		s.head, s.tail = nil, nil
		return
	}
	s.head = next
	next.prev = nil
}

func (s *shardedMapStore) evictEntryFromLL(e *entry) {
	if e == s.tail {
		s.evictTailFromLL()
	} else if e == s.head {
		s.evictHeadFromLL()
	} else {
		e.prev.next, e.next.prev = e.next, e.prev
	}
}

func (s *shardedMapStore) deleteKeyFromMap(key string) (e *entry) {
	sm := s.selectSharedMap(key)
	sm.mu.Lock()
	defer sm.mu.Unlock()

	e, _ = sm.m[key]
	delete(sm.m, key)

	return e
}

func (s *shardedMapStore) lruEvict() {
	// TODO: should support evict multiple
	s.deleteKeyFromMap(s.tail.key)

	// delete from double linked list
	s.evictTailFromLL()

	s.length--
}