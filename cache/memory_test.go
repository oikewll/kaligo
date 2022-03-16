package cache

import (
	"testing"
	"time"
)

func BenchmarkMemoryGet(b *testing.B) {
	cache := NewMemory()
	cache.Set("key", time.Now(), time.Minute)
	for i := 0; i < b.N; i++ {
		_ = cache.Get("key")
	}
}

func BenchmarkMemorySet(b *testing.B) {
	cache := NewMemory()
	for i := 0; i < b.N; i++ {
		_ = cache.Set("key", time.Now(), time.Minute)
	}
}
