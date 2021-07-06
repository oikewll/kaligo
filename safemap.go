/**
 * Realize a safe map with a read-write lock
 *
 * @copyright   (C) 2014  seatle
 * @lastmodify  2021-07-06
 *
 */

package kaligo
 
import (
    "sync"
)

// SafeMap is the struct for Safemap
type SafeMap struct {
    lock *sync.RWMutex
    bm   map[interface{}]interface{}
}

// NewSafeMap is the function for create a new Safemap struct
func NewSafeMap() *SafeMap {
    return &SafeMap{
        lock: new(sync.RWMutex),
        bm:   make(map[interface{}]interface{}),
    }
}

// Get from maps return the k's value
func (m *SafeMap) Get(k interface{}) interface{} {
    m.lock.RLock()
    defer m.lock.RUnlock()
    if val, ok := m.bm[k]; ok {
        return val
    }
    return nil
}

// Set is the function for maps the given key and value, if the key is already in the map and changes nothing.
func (m *SafeMap) Set(k interface{}, v interface{}) bool {
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

// Check is the function for returns true if k is exist in the map.
func (m *SafeMap) Check(k interface{}) bool {
    m.lock.RLock()
    defer m.lock.RUnlock()
    if _, ok := m.bm[k]; !ok {
        return false
    }
    return true
}

// Delete is the function for delete the corresponding key values
func (m *SafeMap) Delete(k interface{}) {
    m.lock.Lock()
    defer m.lock.Unlock()
    delete(m.bm, k)
}
