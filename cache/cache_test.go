package cache

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestMemoryCache(t *testing.T) {
    memory := NewMemory()
    memory.Set("name", "kaka", 10)
    name := memory.Get("name")
    t.Logf("name: [ %v ]", name)
}

func TestCache(t *testing.T) {
    cache := New("memory")
    cache.Set("name", "kaka", time.Millisecond)
    assert.Equal(t, "kaka", cache.Get("name"))
    time.Sleep(time.Millisecond * 2)
    assert.Nil(t, cache.Get("name"))
}
