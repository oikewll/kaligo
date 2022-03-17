// Package cache 注意go语言的接口实现与接口的定义是没有依赖关系的
package cache

import (
    "errors"
    "fmt"
    "strings"
    "time"

    "github.com/owner888/kaligo/config"
)

// Cache interface
type Cache interface {
    Get(key string) any
    Set(key string, val any, timeout time.Duration) error
    IsExist(key string) bool
    Delete(key string) error
}

func New(param ...string) (Cache, error) {
    var driver string
    if len(param) != 0 {
        driver = param[0]
    } else {
        driver = config.Get[string]("cache.config.driver")
    }

    if driver == "memory" {
        return NewMemory(), nil
    } else if driver == "redis" {
        return NewRedis(&RedisOpts{
            Host:        fmt.Sprintf("%s:%d", config.Get[string]("host"), config.Get[int]("port")),
            Password:    config.Get[string]("cache.redis.password"),
            Database:    config.Get[int]("cache.redis.database"),
            MaxIdle:     config.Get[int]("cache.redis.max_idle"),
            MaxActive:   config.Get[int]("cache.redis.max_active"),
            IdleTimeout: config.Get[int]("cache.redis.idle_timeout"),
            Wait:        config.Get[bool]("cache.redis.wait"),
        }), nil
    } else if driver == "memcache" {
        return NewMemcache(strings.Join(config.Get[[]string]("host"), ",")), nil
    }
    return nil, errors.New("driver does not exist")
}
