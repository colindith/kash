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

func Test_GetTTL(t *testing.T) {
	s := GetShardedMapStore()
	_ = s.SetWithTimeout("salmon", "meteor", 5 * time.Minute)

	_, code := s.GetTTL("milktea")
	if code != KeyNotFound {
		t.Errorf("should_be_key_not_found, code: %v", code)
	}

	ttl, code := s.GetTTL("salmon")
	if code != Success {
		t.Errorf("should_be_success, code: %v", code)
	}
	fmt.Println("ttl: ", ttl)
}

func callFunc(s Store, funcName string, args []interface{}) (resp interface{}) {
	switch funcName {
	case "Set":
		s.Set(args[0].(string), args[1])
		return nil
	case "SetWithTimeout":
		s.SetWithTimeout(args[0].(string), args[1], args[1].(time.Duration))
		return nil
	case "Get":
		value, code := s.Get(args[0].(string))
		if code == KeyNotFound {
			return -1
		}
		return value
	case "GetTTL":
		ttl, _ := s.GetTTL(args[0].(string))
		return ttl
	case "Delete":
		s.Delete(args[0].(string))
		return nil
	case "Increase":
		s.Increase(args[0].(string))
		return nil
	default:
		panic("not recognized func name")
	}
}

func callFuncs(s Store, funcNames []string, argsSlice [][]interface{}) []interface{} {
	if len(funcNames) != len(argsSlice) {
		panic("len of funcNames not equal to len of argsSlice")
	}

	res := make([]interface{}, 0, len(funcNames))

	for i := range funcNames {
		res = append(res, callFunc(s, funcNames[i], argsSlice[i]))
	}
	return res
}

func Test_LRU_flow(t *testing.T) {
	testCases := []struct{
		s Store
		funcNames []string
		argsSlice [][]interface{}
		expected []interface{}
	}{
		{
			GetShardedMapStore(SetCapacity(2), SetEvictionPolicy(EvictionLRU)),
			[]string{"Set", "Set", "Get", "Set", "Get", "Set", "Get", "Get", "Get"},
			[][]interface{}{{"1", 1}, {"2", "2"}, {"1"}, {"3", 3}, {"2"}, {"4", 4}, {"1"}, {"3"}, {"4"}},
			[]interface{}{nil, nil, 1, nil, -1, nil, -1, 3, 4},
		},
		{
			GetShardedMapStore(SetCapacity(2), SetEvictionPolicy(EvictionLRU)),
			[]string{"Set","Get"},
			[][]interface{}{{"2",1},{"2"}},
			[]interface{}{nil, 1},
		},
		{
			GetShardedMapStore(SetCapacity(2), SetEvictionPolicy(EvictionLRU)),
			[]string{"Set","Set","Set","Set","Get","Get"},
			[][]interface{}{{"2",1},{"1",1},{"2",3},{"4",1},{"1"},{"2"}},
			[]interface{}{nil,nil,nil,nil,-1,3},
		},
	}

	for _, testCase := range testCases {
		res := callFuncs(testCase.s, testCase.funcNames, testCase.argsSlice)
		if !reflect.DeepEqual(res, testCase.expected) {
			t.Errorf("return value not correct, res=%v, expected=%v", res, testCase.expected)
		}
	}
}