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
	return &Memcache{mc} // bradfitz/gomemcache/memcache 的方法全部合并到当前的 type Memcache struct 去
}

// Get return cached value
func (mem *Memcache) Get(key string) any {
	var err error
	var item *memcache.Item
	if item, err = mem.conn.Get(key); err != nil {
		return nil
	}
	var result any
	if err = json.Unmarshal(item.Value, &result); err != nil {
		return nil
	}
	return result
}

// IsExist check value exists in memcache.
func (mem *Memcache) IsExist(key string) bool {
	if _, err := mem.conn.Get(key); err != nil {
		return false
	}
	return true
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

// Delete delete value in memcache.
func (mem *Memcache) Delete(key string) error {
	return mem.conn.Delete(key)
}
