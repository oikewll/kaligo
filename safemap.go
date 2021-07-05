// 用读写锁实现了 安全的map
package kaligo
 
import (
    "sync"
)
 
type EpoollMap struct {
    lock *sync.RWMutex
    bm   map[interface{}]interface{}
}
 
func NewEpoollMap() *EpoollMap {
    return &EpoollMap{
        lock: new(sync.RWMutex),
        bm:   make(map[interface{}]interface{}),
    }
}
 
//Get from maps return the k's value
func (m *EpoollMap) Get(k interface{}) interface{} {
    m.lock.RLock()
    defer m.lock.RUnlock()
    if val, ok := m.bm[k]; ok {
        return val
    }
    return nil
}
 
// Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *EpoollMap) Set(k interface{}, v interface{}) bool {
    m.lock.Lock()
    defer m.lock.Unlock()
    if val, ok := m.bm[k]; !ok {
        m.bm[k] = v
    } else if val != v {
        m.bm[k] = v
    } else {
        return false
    }
    return true
}
 
// Returns true if k is exist in the map.
func (m *EpoollMap) Check(k interface{}) bool {
    m.lock.RLock()
    defer m.lock.RUnlock()
    if _, ok := m.bm[k]; !ok {
        return false
    }
    return true
}
 
func (m *EpoollMap) Delete(k interface{}) {
    m.lock.Lock()
    defer m.lock.Unlock()
    delete(m.bm, k)
}
