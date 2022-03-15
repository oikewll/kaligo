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

func New(c Cache) *Cache {
    
}

