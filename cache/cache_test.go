package cache

import (
    "testing"

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
