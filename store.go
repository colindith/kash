package kash

import "time"

type store interface {
	set(key string, value interface{}, timeout time.Duration)
	get(key string) (interface{}, error)
}

type defaultStore struct {}

func (s defaultStore) set(key string, value interface{}, timeout time.Duration) {

}

func (s defaultStore) get(key string) (value interface{}, err error) {
	return
}

func getDefaultStore() *store {
	var s store
	s = defaultStore{}
	return &s
}