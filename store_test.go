package kash

import (
	"reflect"
	"testing"
	"time"
)

func Test_defaultStoreGetAndSetFlow(t *testing.T) {
	// TODO: This is bad. Change to table-driven

	s := getDefaultStore()
	// get key from empty store
	v, err := s.get("123")
	if err == nil || err.Error() != "error_cache_not_found" {
		t.Errorf("should_be_error_cache_not_found, got: %v", err)
	}
	if v != nil {
		t.Errorf("value_should_be_nil, got: %v", v)
	}

	// set the key
	err = s.set("123", map[string]string{"jack": "box"}, 1 * time.Minute)
	if err != nil {
		t.Errorf("set_cache_error, err: %v", err)
	}

	// get the key just set
	v, err = s.get("123")
	if err != nil {
		t.Errorf("get_cache_error, err: %v", err)
	}
	if !reflect.DeepEqual(v, map[string]string{"jack": "box"}) {
		t.Errorf("get_cache_value_incorrect, value: %v, want: %v", v, map[string]string{"jack": "box"})
	}

	// delete key
	err = s.delete("123")
	if err != nil {
		t.Errorf("delete_cache_error, err: %v", err)
	}

	// get the deleted key
	v, err = s.get("123")
	if err == nil || err.Error() != "error_cache_not_found" {
		t.Errorf("should_be_error_cache_not_found, got: %v", err)
	}
}

func Test_defaultStore_eviction(t *testing.T) {
	s := getDefaultStore()
	err := s.set("evicted_key", "box", 100*time.Millisecond)
	if err != nil {
		t.Errorf("set_cache_error, err: %v", err)
	}
	time.Sleep(100*time.Millisecond)

	i := 0
	for i < 101 {
		err := s.set("123", "box", 100*time.Millisecond)
		if err != nil {
			t.Errorf("set_cache_error, err: %v", err)
		}
		i++
	}
	// Should trigger eviction
	_, err = s.get("evicted_key")
	if err == nil || err.Error() != "error_cache_not_found" {
		t.Errorf("cache_key_should_expired, err: %v", err)
	}
}