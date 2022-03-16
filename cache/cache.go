// Package cache 注意go语言的接口实现与接口的定义是没有依赖关系的
package cache

import (
	"time"

	"github.com/owner888/kaligo/config"
)

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
            Host        : config.Get[string]("host"),
            Password    : config.Get[string]("password"),
            Database    : config.Get[int]("database"),
            MaxIdle     : config.Get[int]("max_idle"),
            MaxActive   : config.Get[int]("max_active"),
            IdleTimeout : config.Get[int]("idle_timeout"),
            Wait        : config.Get[bool]("wait"),
        })
    } else {
        return NewMemory()
	}
}
