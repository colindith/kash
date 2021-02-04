package kash

import (
	"reflect"
	"testing"
	"time"
)

func Test_defaultStoreGetAndSetFlow(t *testing.T) {
	// TODO: This is bad. Change to table-driven

	s := getDefaultStore()
	v, err := s.get("123")
	if err == nil || err.Error() != "error_cache_not_found" {
		t.Errorf("should_be_error_cache_not_found, got: %v", err)
	}
	if v != nil {
		t.Errorf("value_should_be_nil, got: %v", v)
	}

	err = s.set("123", map[string]string{"jack": "box"}, 1 * time.Minute)
	if err != nil {
		t.Errorf("set_cache_error, err: %v", err)
	}
	v, err = s.get("123")
	if err != nil {
		t.Errorf("get_cache_error, err: %v", err)
	}
	if !reflect.DeepEqual(v, map[string]string{"jack": "box"}) {
		t.Errorf("get_cache_value_incorrect, value: %v, want: %v", v, map[string]string{"jack": "box"})
	}
}
