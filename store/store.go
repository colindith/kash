package store

import (
	"github.com/docker/go-units"
	"log"
	"time"
)

const (
	triggeringEvictionOptNum = 100
)

// ErrorCode
type ErrorCode uint32

const (
	Success            = 1
	KeyNotFound        = 9001
	ValueNotNumberType = 9002
	JSONMarshalErr     = 9003
)

type Store interface {
	Set(key string, value interface{}) ErrorCode
	SetWithTimeout(key string, value interface{}, timeout time.Duration) ErrorCode
	Get(key string) (interface{}, ErrorCode)
	Delete(key string) ErrorCode
	Increase(key string) ErrorCode

	GetTTL(key string) (int64, ErrorCode)

	setDefaultTimeout(timeout time.Duration)
	setEvictionPolicy(policy EvictionPolicy)
	setMaxMemory(size int64)
	setCapacity(cap int)

	DumpAllJSON() (string, ErrorCode)

	Close() ErrorCode
}

type Option func(s Store)

// SetDefaultTimeout generate a Option for setting the default time for Store concrete type
func SetDefaultTimeout(timeout time.Duration) Option {
	return func(s Store) {
		s.setDefaultTimeout(timeout)
	}
}

// SetEvictionPolicy set the policy when the memory usage exceed the threshold
func SetEvictionPolicy(policy EvictionPolicy) Option {
	return func(s Store) {
		s.setEvictionPolicy(policy)
	}
}

// SetMaxMemory generate an Option for setting the max memory used by the data
// Note it only limit the mem usage of the values, not the mem used by the whole process
// When mem usage exceed this threshold, the stored data would be evicted according to the eviction policy
func SetMaxMemory(sizeHuman string) Option {
	return func(s Store) {
		size, err := units.FromHumanSize(sizeHuman)
		if err != nil {
			log.Fatal("invalid_max_memory_option, err: ", err)
			return
		}
		s.setMaxMemory(size)
	}
}

// SetCapacity generate an Option for setting the max keys can be stored in the cache
// When mem usage exceed this threshold, the stored data would be evicted according to the eviction policy
func SetCapacity(cap int) Option {
	return func(s Store) {
		s.setCapacity(cap)
	}
}

// defaultStore implement with build-in map. Most naive implementation
//type defaultStore struct {
//	Store
//	m map[string]unit
//	mu sync.RWMutex        // The whole map share a single lock is inefficient
//
//	operationCount uint    // memo the number of keys mutated since last time eviction
//
//	defaultTimeout time.Duration
//	maxMemory int64              // unit: bytes
//	evictionPolicy EvictionPolicy
//}
//
//type unit struct {
//	data interface{}
//	deadline int64    // timestamp nanosecond
//}
//
//
//func (s *defaultStore) Set(key string, value interface{}) (err error) {
//	return s.SetWithTimeout(key, value, s.defaultTimeout)
//}
//
//func (s *defaultStore) SetWithTimeout(key string, value interface{}, timeout time.Duration) (err error) {
//	deadline := time.Now().Add(timeout).UnixNano()
//	if timeout == 0 {
//		deadline = maxInt64
//	}
//
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	s.m[key] = unit{
//		data: value,
//		deadline: deadline,
//	}
//	s.operationCount++
//	go func() {
//		if s.operationCount >= triggeringEvictionOptNum {
//			// TODO: Except for this, there are other occurrences that should trigger eviction
//			evictMap(s)   // This function would require the lock
//			s.operationCount = 0
//		}
//	}()
//	return nil
//}
//
//func (s *defaultStore) Get(key string) (value interface{}, err error) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	v, ok := s.m[key]
//	if !ok {
//		return nil, fmt.Errorf("error_cache_not_found")
//	}
//	if time.Now().UnixNano() > v.deadline {
//		// The key was timeout. Evict it.
//		delete(s.m, key)
//		return nil, fmt.Errorf("error_cache_not_found")
//	}
//	return v.data, nil
//}
//
//func (s *defaultStore) Delete(key string) error {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	delete(s.m, key)
//	return nil
//}
//
//
//func (s *defaultStore) setDefaultTimeout(timeout time.Duration) {
//	s.defaultTimeout = timeout
//}
//
//func (s *defaultStore) setEvictionPolicy(policy EvictionPolicy) {
//	s.evictionPolicy = policy
//}
//
//func (s *defaultStore) setMaxMemory(size int64) {
//	s.maxMemory = size
//}
//
//
//func GetDefaultStore(opts... Option) Store {
//	s := &defaultStore{
//		m: make(map[string]unit),
//	}
//	for _, opt := range opts {
//		opt(s)
//	}
//	return s
//}
//
//// evictMap loop though the map and evict the expired key
//func evictMap(s *defaultStore) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	now := time.Now().UnixNano()
//	for k, u := range s.m {
//		if u.deadline <= now {
//			delete(s.m, k)      // delete map key while loop through map. Is it okay?
//		}
//	}
//}
