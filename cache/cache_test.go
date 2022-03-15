package cache

import (
    "log"
    "testing"
    "time"
)

func TestMemoryCache(t *testing.T) {
    memory := NewMemory()
    memory.Set("name", "kaka", 10)
    name := memory.Get("name")
    t.Logf("name: [ %v ]", name)
}

func TestCache(t *testing.T) {
    cache := New("memory")
    cache.Set("name", "kaka", 4000)
    name := cache.Get("name")
    log.Printf("name: [ %v ]", name)
    time.Sleep(time.Second * 4)
    name = cache.Get("name")
    log.Printf("name: [ %v ]", name)
}
