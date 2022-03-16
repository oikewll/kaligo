// Package cache 注意go语言的接口实现与接口的定义是没有依赖关系的
package cache

import "time"

// Cache interface
type Cache interface {
	Get(key string) interface{}
	Set(key string, val interface{}, timeout time.Duration) error
	IsExist(key string) bool
	Delete(key string) error
}

func New(driver string) Cache {
	if driver == "memcache" {
		return NewMemcache("")
	} else if driver == "redis" {
		return NewRedis(&RedisOpts{
            Host        : "127.0.0.1",
            Password    : "",
            Database    : 0,
            MaxIdle     : 1,
            MaxActive   : 1,
            IdleTimeout : 1,
        })
    } else {
        return NewMemory()
	}
}
