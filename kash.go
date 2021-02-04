package kash

import "time"

type Kash struct {
	config Config
	store store
	close chan struct{}    // Need to consider when should the cache close
}


func (k *Kash) Set(key string, value interface{}, timeout time.Duration) error {
	// TODO: this timeout should make as an option
	return k.store.set(key, value, timeout)
}
func (k *Kash) Get(key string) (interface{}, error) {
	return k.store.get(key)
}

func (k *Kash) setConfig(c *Config) (err error) {
	// valid config
	k.config = *c
	return
}

func (k *Kash) setStore(s *store) (err error) {
	k.store = *s
	return
}

func NewKash(c *Config) (k *Kash, err error) {
	k = &Kash{}
	err = k.setConfig(c)
	if err != nil {
		return nil, err
	}

	s := getDefaultStore()
	err = k.setStore(&s)
	if err != nil {
		return nil, err
	}

	return k, nil
}