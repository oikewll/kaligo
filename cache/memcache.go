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

func (mem *Memcache) String(key string) string {
    reply, found := mem.Get(key);
    if  !found {
        return ""
    }
    return reply.(string)
}

func (mem *Memcache) Int(key string) int {
    reply, found := mem.Get(key);
    if  !found {
        return 0
    }
    return reply.(int)
}

func (mem *Memcache) Int64(key string) int64 {
    reply, found := mem.Get(key);
    if  !found {
        return 0
    }
    return reply.(int64)
}

func (mem *Memcache) Uint(key string) uint {
    reply, found := mem.Get(key);
    if  !found {
        return 0
    }
    return reply.(uint)
}

func (mem *Memcache) Uint64(key string) uint64 {
    reply, found := mem.Get(key);
    if  !found {
        return 0
    }
    return reply.(uint64)
}

func (mem *Memcache) Float64(key string) float64 {
    reply, found := mem.Get(key);
    if  !found {
        return 0
    }
    return reply.(float64)
}

func (mem *Memcache) Bool(key string) bool {
    reply, found := mem.Get(key);
    if  !found {
        return false
    }
    return reply.(bool)
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

func (mem *Memcache) LPush(key string, value any) {
}

func (mem *Memcache) RPush(key string, value any) {
}

func (mem *Memcache) LPop(key string, value any) (err error) {
    return
}

func (mem *Memcache) RPop(key string, value any) (err error) {
    return
}
