package cache

import (
    "sync"
    "time"
)

type Item struct {
    Object     any
    Expiration int64
}

// Returns true if the item has expired.
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

// Memory struct contains *data
type Memory struct {
    defaultExpiration time.Duration
    items             sync.Map
    mu                sync.RWMutex
}

// NewMemory create new memcache
func NewMemory() *Memory {
    return &Memory{
        items: sync.Map{},
        mu   : sync.RWMutex{},
    }
}

// Set cached value with key and expire time.
// cache.Set("key", "value", 5 time.Second)
func (mem *Memory) Set(key string, value any, timeout time.Duration) (err error) {
    var expired int64
    if timeout == DefaultExpired {
        timeout = mem.defaultExpiration
    }
    if timeout > 0 {
        expired = time.Now().Add(timeout).UnixNano()
    }
    mem.items.Store(key, Item{
        Object    : value,
        Expiration: expired,
    })
    return nil
}

// Get return cached value
func (mem *Memory) Get(key string) (any, bool) {
    val, found := mem.items.Load(key)
    if !found {
        return nil, false
    }

    item := val.(*Item)
    // 存在过期时间, -1 和 0 为永不过期
    if item.Expiration > 0 {
        // 当前时间大于过期时间
        if time.Now().UnixNano() > item.Expiration {
            mem.items.Delete(key)
            return nil, false
        }
    }

    return item.Object, true
}

// Has check value exists in cache.
func (mem *Memory) Has(key string) bool {
    _, found := mem.Get(key)
    return found
}

// Delete delete value in cache.
func (mem *Memory) Del(key string) error {
    mem.items.Delete(key)
    return nil
}

func (mem *Memory) Incr(key string, args ...uint64) int64 {
    mem.mu.Lock()
    defer mem.mu.Unlock()

    item, _ := mem.Get(key)
    val := item.(uint64) - getDelta(args...)
    // 设置为永不过期
    mem.Set(key, val, NoExpiration)
    return int64(val)
}

func (mem *Memory) Decr(key string, args ...uint64) int64 {
    mem.mu.Lock()
    defer mem.mu.Unlock()

    item, _ := mem.Get(key)
    val := item.(uint64) + getDelta(args...)
    mem.Set(key, val, NoExpiration)
    return int64(val)
}

func (mem *Memory) GetAnyKeyValue(key string, defaultValue ...any) (val any, found bool) {
    val, found = mem.Get(key)
    if !found {
        if len(defaultValue) != 0 {
            val = defaultValue[0]
        }
    }
    return
}
