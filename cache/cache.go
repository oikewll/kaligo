// Package cache 注意go语言的接口实现与接口的定义是没有依赖关系的
package cache

import (
    "errors"
    "fmt"
    "log"
    "strings"
    "time"

    "github.com/owner888/kaligo/config"
    "golang.org/x/exp/constraints"
)

var (
    // defaultCache 默认的全局缓存，全局的 Get Set 使用此缓存
    defaultCache Cache
)

// 支持存取到 Cache 的类型
type CacheValue interface {
    constraints.Integer | constraints.Float | ~bool | ~string
}

// Cache interface
type Cache interface {
    // kaligo.AnyKeyValueGetter // TODO: import cycle
    Get(key string) (any, error)
    Set(key string, val any, timeout time.Duration) error
    IsExist(key string) bool
    Delete(key string) error
    Incr(key string) int64
    Decr(key string) int64
}

func New(param ...string) (Cache, error) {
    var driver string
    if len(param) != 0 {
        driver = param[0]
    } else {
        driver = config.Get[string]("cache.config.driver")
    }

    log.Printf("%v", config.Get[string]("cache.config.driver"))
    if driver == "memory" {
        return NewMemory(), nil
    } else if driver == "redis" {
        return NewRedis(&RedisOpts{
            Host:        fmt.Sprintf("%s:%d", config.Get[string]("cache.redis.host"), config.Get[int]("cache.redis.port")),
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

// SetDefaultCache 设置默认的全局缓存，全局的 Get Set 使用此缓存
func SetDefaultCache(cache Cache) {
    defaultCache = cache
}

// Get 从默认的 Cache 获取 T 类型 value
func Get[T CacheValue](key string) T {
    return GetCache[T](defaultCache, key)
}

// GetCache 从 cache 获取 T 类型 value
func GetCache[T CacheValue](cache Cache, key string) (value T) {
    v, err := cache.Get(key)
    if err != nil {
        return
    }
    if v, ok := v.(T); ok {
        return v
    }
    return
}

// Set 设置 value 到默认 Cache
func Set[T CacheValue](key string, value T, timeout time.Duration) {
    SetCache(defaultCache, key, value, timeout)
}

// Set 设置 value 到 cache
func SetCache[T CacheValue](cache Cache, key string, value T, timeout time.Duration) {
    cache.Set(key, value, timeout)
}
