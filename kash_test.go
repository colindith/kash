package kash

import (
	"testing"
	"time"
)

func Test_cacheGetSet(t *testing.T) {
	c := &Config{
		evictionPolicy: LRU,
		defaultTimeout: 10 * time.Minute,
	}
	k, err := NewKash(c)
	if err != nil {
		t.Errorf("init_kash_error, err: %v", err)
	}
	err = k.Set("123", "test_value", time.Minute)
	if err != nil {
		t.Errorf("set_cache_error, err: %v", err)
	}
	v, err := k.Get("123")
	if err != nil {
		t.Errorf("get_cache_error, err: %v", err)
	}
	if v != "test_value" {
		t.Errorf("get_cache_data_incorrect, want: %v, got: %v", "test_value", v)
	}

}