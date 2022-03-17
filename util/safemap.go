package util

import (
    "sync"
)

// SafeMap is the struct for Safemap
type SafeMap struct {
    sync.Mutex // 内嵌结构，无法当作继承使用，不支持重载（无法访问子类的方法、数据），不支持LSP
    data       map[string]any
}

// NewSafeMap 实例化
func NewSafeMap() *SafeMap {
    return &SafeMap{
        data: make(map[string]any),
    }
}

// Get from maps return the k's value
func (m *SafeMap) Get(key string) any {
    if val, ok := m.data[key]; ok {
        return val
    }

    return nil
}

// Set is the function for maps the given key and value, if the key is already in the map and changes nothing.
func (m *SafeMap) Set(key string, val any) bool {
    m.Lock()
    defer m.Unlock()

    // key 对应的值不存在
    if val, ok := m.data[key]; !ok {
        m.data[key] = val
    } else if val != val { // 存在值但是不同，替换
        m.data[key] = val
    } else {
        return false
    }

    return true
}

// IsExist is the function for returns true if k is exist in the map.
func (m *SafeMap) IsExist(key string) bool {
    m.Lock()
    defer m.Unlock()

    if _, ok := m.data[key]; !ok {
        return false
    }

    return true
}

// Delete is the function for delete the corresponding key values
func (m *SafeMap) Delete(key string) {
    m.Lock()
    defer m.Unlock()

    delete(m.data, key)
}
