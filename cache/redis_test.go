package cache

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

var _ Cache = &Redis{}

func TestRedis(t *testing.T) {
    redis := NewRedis(&RedisOpts{Host: "127.0.0.1:6379"})
    assert.NotNil(t, redis)

    key := "any_key"
    value := "any value"
    assert.NoError(t, redis.Set(key, value, time.Second))
    assert.True(t, redis.Has(key))
    v, err := redis.Get(key)
    assert.True(t, err)
    assert.Equal(t, value, v)
    assert.NoError(t, redis.Del(key))
    v, err = redis.Get(key)
    assert.Nil(t, v)
    assert.False(t, err)

    // 测试 timeout
    assert.NoError(t, redis.Set(key, value, time.Second))
    v, err = redis.Get(key)
    assert.True(t, err)
    assert.Equal(t, value, v)
    time.Sleep(time.Second)
    v, err = redis.Get(key)
    assert.Nil(t, v)
    assert.False(t, err)

    // 测试不存在的 key
    key = "unknown_key"
    assert.False(t, redis.Has(key))
    v, err = redis.Get(key)
    assert.Nil(t, v)
    assert.False(t, err)
    assert.NoError(t, redis.Del(key))
}
