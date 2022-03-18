package cache

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

var _ Cache = &Redis{}

func TestRedis(t *testing.T) {
    redis := NewRedis(&RedisOpts{Host: "127.0.0.1:6379"})

    key := "any_key"
    value := "any value"
    assert.NoError(t, redis.Set(key, value, time.Second))
    assert.True(t, redis.IsExist(key))
    assert.Equal(t, value, redis.Get(key))
    assert.NoError(t, redis.Delete(key))
    assert.Nil(t, redis.Get(key))

    // 测试 timeout
    assert.NoError(t, redis.Set(key, value, time.Second))
    assert.Equal(t, value, redis.Get(key))
    time.Sleep(time.Second)
    assert.Nil(t, redis.Get(key))

    // 测试不存在的 key
    key = "unknown_key"
    assert.False(t, redis.IsExist(key))
    assert.Nil(t, redis.Get(key))
    assert.NoError(t, redis.Delete(key))
}
