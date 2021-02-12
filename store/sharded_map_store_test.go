package store

import (
	"reflect"
	"testing"
	"time"
)

func Test_shardedMapStoreGetAndSetFlow(t *testing.T) {
	s := GetShardedMapStore()
	// Get key from empty store
	v, err := s.Get("123")
	if err == nil || err.Error() != "error_cache_not_found" {
		t.Errorf("should_be_error_cache_not_found, got: %v", err)
	}
	if v != nil {
		t.Errorf("value_should_be_nil, got: %v", v)
	}

	// set the key
	err = s.SetWithTimeout("123", map[string]string{"jack": "box"}, 1 * time.Minute)
	if err != nil {
		t.Errorf("set_cache_error, err: %v", err)
	}

	// Get the key just set
	v, err = s.Get("123")
	if err != nil {
		t.Errorf("Get_cache_error, err: %v", err)
	}
	if !reflect.DeepEqual(v, map[string]string{"jack": "box"}) {
		t.Errorf("get_cache_value_incorrect, value: %v, want: %v", v, map[string]string{"jack": "box"})
	}

	// delete key
	err = s.Delete("123")
	if err != nil {
		t.Errorf("delete_cache_error, err: %v", err)
	}

	// Get the deleted key
	v, err = s.Get("123")
	if err == nil || err.Error() != "error_cache_not_found" {
		t.Errorf("should_be_error_cache_not_found, got: %v", err)
	}
}

func Test_shardedMapStore_eviction(t *testing.T) {
	s := GetShardedMapStore()
	err := s.SetWithTimeout("evicted_key", "box", 100*time.Millisecond)
	if err != nil {
		t.Errorf("set_cache_error, err: %v", err)
	}
	time.Sleep(100*time.Millisecond)

	i := 0
	for i < 101 {
		err := s.SetWithTimeout("123", "box", 100*time.Millisecond)
		if err != nil {
			t.Errorf("set_cache_error, err: %v", err)
		}
		i++
	}
	// Should trigger eviction
	_, err = s.Get("evicted_key")
	if err == nil || err.Error() != "error_cache_not_found" {
		t.Errorf("cache_key_should_expired, err: %v", err)
	}
}

func Test_shardedMapStore_SetDefaultTime(t *testing.T) {
	s := GetShardedMapStore(SetDefaultTimeout(100 * time.Millisecond))

	// set key with default timeout
	err := s.Set("default_timeout", "box")
	if err != nil {
		t.Errorf("set_cache_error, err: %v", err)
	}
	// get key right after set
	v, err := s.Get("default_timeout")
	if err != nil {
		t.Errorf("get_cache_error, err: %v", err)
	}
	if !reflect.DeepEqual(v, "box") {
		t.Errorf("get_cache_value_incorrect, value: %v, want: %v", v, "box")
	}
	time.Sleep(100*time.Millisecond)

	// get again after timeout
	_, err = s.Get("default_timeout")
	if err == nil || err.Error() != "error_cache_not_found" {
		t.Errorf("cache_key_should_expired, err: %v", err)
	}
}