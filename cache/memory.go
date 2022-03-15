package cache

import (
    "sync"
    "time"
)

// Memory struct contains *data
type Memory struct {
    data sync.Map
}

type data struct {
    Data    interface{}
    Expired time.Time
}

// NewMemory create new memcache
func NewMemory() *Memory {
    return &Memory{
        data: sync.Map{},
    }
}

// Get return cached value
func (mem *Memory) Get(key string) interface{} {
    val, ok := mem.data.Load(key)
    if !ok {
        return nil
    }

    ret := val.(*data)
    if ret.Expired.Before(time.Now()) {
        mem.data.Delete(key)
        return nil
    }

    return ret.Data
}

// IsExist check value exists in memcache.
func (mem *Memory) IsExist(key string) bool {
    return mem.Get(key) == nil
}

// Set cached value with key and expire time.
func (mem *Memory) Set(key string, value interface{}, timeout time.Duration) (err error) {
    mem.data.Store(key, &data{
        Data:    value,
        Expired: time.Now().Add(timeout),
    })
    return nil
}

// Delete delete value in memcache.
func (mem *Memory) Delete(key string) error {
    mem.data.Delete(key)
    return nil
}
