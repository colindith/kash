package kash

import (
	"fmt"
	"time"
)

type store interface {
	set(key string, value interface{}, timeout time.Duration) error
	get(key string) (interface{}, error)
}

type defaultStore struct {
	m map[string]unit     // TODO: This map is not thread safe
}

type unit struct {
	data interface{}
	deadline int64    // timestamp nanosecond
}

func (s defaultStore) set(key string, value interface{}, timeout time.Duration) (err error) {
	s.m[key] = unit{
		data: value,
		deadline: time.Now().Add(timeout).UnixNano(),
	}
	return nil
}

func (s defaultStore) get(key string) (value interface{}, err error) {
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

func getDefaultStore() store {
	return defaultStore{
		m: make(map[string]unit),
	}
}