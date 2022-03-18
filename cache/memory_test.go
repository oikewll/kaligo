package cache

import (
    "os"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

var _ Cache = &Memory{}

var cache Memory

func TestMain(m *testing.M) {
    cache = *NewMemory()
    os.Exit(m.Run())
}

func TestMemoryCache(t *testing.T) {
    assert.NotNil(t, cache)
    v, err := cache.Get("unknown_key")
    assert.Nil(t, v)
    assert.Error(t, err)

    key := "any_key"
    value := "any value"
    cache.Set(key, value, time.Millisecond)
    v, err = cache.Get(key)
    assert.Equal(t, value, v)
    assert.NoError(t, err)

    cache.Delete(key)

    v, err = cache.Get(key)
    assert.Nil(t, v)
    assert.Error(t, err)

    cache.Set(key, value, time.Millisecond)
    v, err = cache.Get(key)
    assert.Equal(t, value, v)
    assert.NoError(t, err)
    time.Sleep(time.Millisecond * 2)
    v, err = cache.Get(key)
    assert.Nil(t, v)
    assert.Error(t, err)
}

func BenchmarkMemoryGet(b *testing.B) {
    cache.Set("key", time.Now(), time.Minute)
    for i := 0; i < b.N; i++ {
        _, _ = cache.Get("key")
    }
}

func BenchmarkMemorySet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = cache.Set("key", time.Now(), time.Minute)
    }
}
