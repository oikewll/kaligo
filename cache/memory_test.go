package cache

import (
    "os"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

var cache Memory

func TestMain(m *testing.M) {
    cache = *NewMemory()
    os.Exit(m.Run())
}

func TestMemoryCache(t *testing.T) {
    assert.NotNil(t, cache)
    assert.Nil(t, cache.Get("unknown_key"))

    key := "any_key"
    value := "any value"
    cache.Set(key, value, time.Millisecond)
    assert.Equal(t, value, cache.Get(key))

    cache.Delete(key)
    assert.Nil(t, cache.Get(key))

    cache.Set(key, value, time.Millisecond)
    assert.Equal(t, value, cache.Get(key))
    time.Sleep(time.Millisecond * 2)
    assert.Nil(t, cache.Get(key))
}

func BenchmarkMemoryGet(b *testing.B) {
    cache.Set("key", time.Now(), time.Minute)
    for i := 0; i < b.N; i++ {
        _ = cache.Get("key")
    }
}

func BenchmarkMemorySet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = cache.Set("key", time.Now(), time.Minute)
    }
}
