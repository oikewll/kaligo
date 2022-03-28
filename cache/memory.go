package cache

import (
    "errors"
    "sync"
    "sync/atomic"
    "time"
)

// Memory struct contains *data
type Memory struct {
    data sync.Map
}

type data struct {
    Data    any
    Expired time.Time
}

// NewMemory create new memcache
func NewMemory() *Memory {
    return &Memory{
        data: sync.Map{},
    }
}

// Get return cached value
func (mem *Memory) Get(key string) (reply any, err error) {
    val, ok := mem.data.Load(key)
    if !ok {
        return nil, errors.New("sync.Map load error.")
    }

    ret := val.(*data)
    if ret.Expired.Before(time.Now()) {
        mem.data.Delete(key)
        return nil, errors.New("key expired.")
    }

    reply = ret.Data
    return reply, nil
}

// IsExist check value exists in memcache.
func (mem *Memory) IsExist(key string) bool {
    _, err := mem.Get(key)
    if err != nil {
        return false
    }
    return true
}

// Set cached value with key and expire time.
func (mem *Memory) Set(key string, value any, timeout time.Duration) (err error) {
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

func (mem *Memory) Incr(key string) int64 {
    var incrID int64    
    ret, err := mem.Get(key)
    if err != nil {
        incrID = 0
    }
    incrID = ret.(int64)
    incrID = atomic.AddInt64(*incrID)
    return incrID
}

func (mem *Memory) Decr(key string) int64 {
    return 0
}

func (mem *Memory) GetAnyKeyValue(key string, defaultValue ...any) (v any, ok bool) {
    v, err := mem.Get(key)
    ok = err == nil
    if !ok {
        if len(defaultValue) != 0 {
            v = defaultValue[0]
        }
    }
    return
}
