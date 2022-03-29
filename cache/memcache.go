package cache

import (
    "encoding/json"
    "time"

    "github.com/bradfitz/gomemcache/memcache"
)

// Memcache struct contains *memcache.Client
type Memcache struct {
    conn *memcache.Client
}

// NewMemcache create new memcache
func NewMemcache(server ...string) *Memcache {
    mc := memcache.New(server...)
    return &Memcache{mc}
}

// Set cached value with key and expire time.
func (mem *Memcache) Set(key string, val any, timeout time.Duration) (err error) {
    var data []byte
    if data, err = json.Marshal(val); err != nil {
        return err
    }

    item := &memcache.Item{Key: key, Value: data, Expiration: int32(timeout / time.Second)}
    return mem.conn.Set(item)
}

// Get return cached value
func (mem *Memcache) Get(key string) (any, bool) {
    var item *memcache.Item
    var err error    
    var reply any    

    if item, err = mem.conn.Get(key); err != nil {
        return nil, false
    }
    if err = json.Unmarshal(item.Value, &reply); err != nil {
        return nil, false
    }

    return reply, true
}

// Has check value exists in memcache.
func (mem *Memcache) Has(key string) bool {
    if _, err := mem.conn.Get(key); err != nil {
        return false
    }
    return true
}

// Delete delete value in memcache.
func (mem *Memcache) Del(key string) error {
    return mem.conn.Delete(key)
}

func (mem *Memcache) Incr(key string, args ...uint64) int64 {
    num, err := mem.conn.Increment(key, getDelta(args...))
    if err != nil {
        return 0
    }
    return int64(num)
}

func (mem *Memcache) Decr(key string, args ...uint64) int64 {
    num, err := mem.conn.Decrement(key, getDelta(args...))
    if err != nil {
        return 0
    }
    return int64(num)
}

func (mem *Memcache) GetAnyKeyValue(key string, defaultValue ...any) (val any, found bool) {
    val, found = mem.Get(key)
    if !found {
        if len(defaultValue) != 0 {
            val = defaultValue[0]
        }
    }
    return
}
