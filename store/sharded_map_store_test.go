package store

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_shardedMapStoreGetAndSetFlow(t *testing.T) {
	s := GetShardedMapStore()
	// Get key from empty store
	v, code := s.Get("123")
	if code != KeyNotFound {
		t.Errorf("should_be_error_cache_not_found, got: %v", code)
	}
	if v != nil {
		t.Errorf("value_should_be_nil, got: %v", v)
	}

	// set the key
	code = s.SetWithTimeout("123", map[string]string{"jack": "box"}, 1 * time.Minute)
	if code != Success {
		t.Errorf("set_cache_error, code: %v", code)
	}

	// Get the key just set
	v, code = s.Get("123")
	if code != Success {
		t.Errorf("Get_cache_error, code: %v", code)
	}
	if !reflect.DeepEqual(v, map[string]string{"jack": "box"}) {
		t.Errorf("get_cache_value_incorrect, value: %v, want: %v", v, map[string]string{"jack": "box"})
	}

	// delete key
	code = s.Delete("123")
	if code != Success {
		t.Errorf("delete_cache_error, code: %v", code)
	}

	// Get the deleted key
	v, code = s.Get("123")
	if code != KeyNotFound {
		t.Errorf("should_be_error_cache_not_found, got: %v", code)
	}
}

func Test_shardedMapStore_eviction(t *testing.T) {
	s := GetShardedMapStore()
	code := s.SetWithTimeout("evicted_key", "box", 100*time.Millisecond)
	if code != Success {
		t.Errorf("set_cache_error, code: %v", code)
	}
	time.Sleep(100*time.Millisecond)

	i := 0
	for i < 101 {
		code = s.SetWithTimeout("123", "box", 100*time.Millisecond)
		if code != Success {
			t.Errorf("set_cache_error, code: %v", code)
		}
		i++
	}
	// cache key were expired
	_, code = s.Get("evicted_key")
	if code != KeyNotFound {
		t.Errorf("cache_key_should_expired, code: %v", code)
	}
	// Should trigger eviction
	want := "{\"123\":\"box\"}"
	if jsonStr, _ := s.DumpAllJSON(); jsonStr != want {
		t.Errorf("cache_key_should_be_evicted, got: %v, want: %v", jsonStr, want)
	}
}

func Test_shardedMapStore_SetDefaultTime(t *testing.T) {
	s := GetShardedMapStore(SetDefaultTimeout(100 * time.Millisecond))

	// set key with default timeout
	code := s.Set("default_timeout", "box")
	if code != Success {
		t.Errorf("set_cache_error, code: %v", code)
	}
	// get key right after set
	v, code := s.Get("default_timeout")
	if code != Success {
		t.Errorf("get_cache_error, code: %v", code)
	}
	if !reflect.DeepEqual(v, "box") {
		t.Errorf("get_cache_value_incorrect, value: %v, want: %v", v, "box")
	}
	time.Sleep(100*time.Millisecond)

	// get again after timeout
	_, code = s.Get("default_timeout")
	if code != KeyNotFound {
		t.Errorf("cache_key_should_expired, code: %v", code)
	}
}

// TODO: user modify the data outside the in-memory cache should not affect the data in the cache
// Means need to deep copy the data?

func Test_DumpAllJSON(t *testing.T) {
	s := GetShardedMapStore()

	code := s.Set("Best", "Bites!")
	if code != Success {
		t.Errorf("set_cache_error, code: %v", code)
	}
	code = s.Set("Timbre", "+")
	if code != Success {
		t.Errorf("set_cache_error, code: %v", code)
	}

	jsonStr, code := s.DumpAllJSON()
	if code != Success {
		t.Errorf("dump_all_json_error, code: %v", code)
	}
	want := "{\"Best\":\"Bites!\",\"Timbre\":\"+\"}"
	if jsonStr != want {
		fmt.Println("json str: ", jsonStr)
		t.Errorf("dump_all_json_incorrect, got: %v, want: %v", jsonStr, want)
	}
}

func Test_Increase(t *testing.T) {
	s := GetShardedMapStore()

	code := s.Increase("desert")
	if code != Success {
		t.Errorf("incr_non_existed_key_err, code: %v", code)
	}

	v, code := s.Get("desert")
	if code != Success {
		t.Errorf("get_cache_error, code: %v", code)
	}
	if v.(int) != 1 {
		t.Errorf("incr_value_err")
	}

	code = s.Increase("desert")
	if code != Success {
		t.Errorf("incr_existed_err, code: %v", code)
	}
	v, code = s.Get("desert")
	if code != Success {
		t.Errorf("get_cache_error, code: %v", code)
	}
	if v.(int) != 2 {
		t.Errorf("incr_value_err, expected=%v, got=%v", 2, v)
	}


	code = s.Set("gossip", uint32(1234))
	if code != Success {
		t.Errorf("set_cache_error, code: %v", code)
	}

	code = s.Increase("gossip")
	if code != Success {
		t.Errorf("incr_existed_key_err, code: %v", code)
	}

	v, code = s.Get("gossip")
	if code != Success {
		t.Errorf("get_cache_error, code: %v", code)
	}
	if v.(uint32) != uint32(1235) {
		t.Errorf("incr_value_err, expected=%v, got=%v", 1235, v)
	}
}

type TestStruct struct {
	Field1 int
	Field2 int
}

func Test_modifyCacheDataFromOutside(t *testing.T) {
	s := GetShardedMapStore()

	_ = s.Set("gossip", &TestStruct{5, 50})

	testStructInterface, _ := s.Get("gossip")

	testStruct := testStructInterface.(*TestStruct)
	testStruct.Field1 = 6

	want := "{\"gossip\":{\"Field1\":5,\"Field2\":50}}"

	if jsonStr, _ := s.DumpAllJSON(); jsonStr != want {
		t.Errorf("cached_data_deing_modified_from_outside | want=%v | got=%v", want, jsonStr)
	}

}