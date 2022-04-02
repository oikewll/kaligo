package cache

import (
    "testing"
    "time"

    "github.com/bradfitz/gomemcache/memcache"
    "github.com/stretchr/testify/assert"
)

var _ Cache = &Memcache{}

func TestMemcache(t *testing.T) {
    mem := NewMemcache("127.0.0.1:11211")
    var err error
    timeoutDuration := 10 * time.Second
    if err = mem.Set("username", "silenceper", timeoutDuration); err != nil {
        t.Error("set Error", err)
    }

    if !mem.Has("username") {
        t.Error("IsExist Error")
    }
    exists := mem.Has("unknown-key")
    assert.Equal(t, false, exists)

    name, ok := mem.Get("username")
    if name != "" {
        if name != "silenceper" {
            t.Error("get Error")
        }
    }
    assert.True(t, ok)
    data, ok := mem.Get("unknown-key")
    assert.Nil(t, data)
    assert.False(t, ok)

    if err = mem.Del("username"); err != nil {
        t.Errorf("delete Error , err=%v", err)
    }

    err = mem.Del("unknown-key")
    assert.Equal(t, memcache.ErrCacheMiss, err)
}
