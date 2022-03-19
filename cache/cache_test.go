package cache

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
    cache, err := New("redis")
    assert.IsType(t, &Redis{}, cache)
    cache, _ = New("memcache")
    assert.IsType(t, &Memcache{}, cache)
    cache, _ = New("memory")
    assert.IsType(t, &Memory{}, cache)
    cache, err = New("")
    assert.Error(t, err)
    assert.Nil(t, cache)
    cache, err = New()
    assert.Error(t, err)
    assert.Nil(t, cache)
}

func TestDefaultCache(t *testing.T) {
    assert.Nil(t, defaultCache)
    cache, _ := New("memory")
    SetDefaultCache(cache)
    assert.NotNil(t, defaultCache)

    key, value := "any_key", "any value"
    Set(key, value, time.Millisecond)
    assert.Equal(t, Get[string](key), value)
    assert.Equal(t, Get[int](key), 0)
}
