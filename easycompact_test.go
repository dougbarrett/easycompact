package easycompact

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ttl := time.Duration(1 * time.Second)
	compact := New(&ttl, func(key any, data []any) {
		if key != "key1" {
			t.Errorf("key should be key1, but got %v", key)
		}
		fmt.Printf("data: %v", data)
		if len(data) != 4 {
			t.Errorf("data length should be 4, but got %v", len(data))
		}
	})
	compact.Set("key1", "data1")
	compact.Set("key1", "data2")
	compact.Set("key1", "data3")
	compact.Set("key1", "data4")
	compact.Close()
}
