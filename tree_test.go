package kaligo

import (
    "testing"

    "github.com/owner888/kaligo/cache"
    "github.com/stretchr/testify/assert"
)

var _ AnyKeyValueGetter = Params{}
var _ AnyKeyValueGetter = cache.NewMemcache()
var _ AnyKeyValueGetter = cache.NewMemory()
var _ AnyKeyValueGetter = cache.NewRedis(&cache.RedisOpts{})

func TestParams(t *testing.T) {
    params := Params{
        Param[string]{"key1", "value1"},
        Param[string]{"key2", "value2"},
        Param[string]{"key3", ""},
    }
    assert.Equal(t, "value1", params.ByName("key1"))
    assert.Equal(t, "", params.ByName("key3", "defaultValue"))
    assert.Equal(t, "defaultValue", params.ByName("unknown_key", "defaultValue"))
    assert.Equal(t, "", params.ByName("unknown_key"))
}
