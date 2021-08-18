package cache

import (
	"testing"
)

func TestCache(t *testing.T) {
    memory := NewMemory()
    memory.Set("name", "kaka", 10)
    name := memory.Get("name")
    t.Logf("name: [ %v ]", name)
}
