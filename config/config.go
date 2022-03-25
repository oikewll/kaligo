/**
 * Read the configuration file
 *
 * @copyright   (C) 2014  seatle
 * @lastmodify  2021-07-06
 *
 * Point-to-point acquisition
 * config.Get[string]("database.mysql.charset")
 *
 */

package config

import (
    "strings"
    "sync"

    // "github.com/astaxie/beego/logs"
    "github.com/owner888/kaligo/util"
)

var (
    configMaps sync.Map
)

// ConfigMap 是用于保存多层级数据的 map，已实现此接口的类型: sync.Map, StrMap
type ConfigMap interface {
    // Load 读取值
    Load(key any) (value any, ok bool)
    // Store 保存值
    Store(key, value any)
}

// StrMap is use for string -> map
type StrMap map[string]any

// Load is used to read value by key
func (m StrMap) Load(key any) (value any, ok bool) {
    value = m[key.(string)]
    ok = value != nil
    return
}

// Store is used to write value by key
func (m StrMap) Store(key, value any) {
    m[key.(string)] = value
}

// Env 读取环境变量(configMaps存入的值)，支持默认值
func Env(key string, defaultValue ...any) any {
    var keys []string = strings.Split(key, ".")
    lastIndex := len(keys) - 1
    var maps ConfigMap = &configMaps
    for i, k := range keys {
        if i == lastIndex {
            if val, ok := maps.Load(k); ok {
                return val
            }
        } else if m := getConfigMap(maps, k); m != nil {
            maps = m
        }
    }
    if len(defaultValue) > 0 {
        return defaultValue[0]
    }
    return nil
}

// Deprecated: Use config.Set instead.Add 新增配置项
func Add(key string, value any) {
    Set(key, value)
}

// Set 设置配置项
func Set(key string, value any) {
    var keys []string = strings.Split(key, ".")
    lastIndex := len(keys) - 1
    var maps ConfigMap = &configMaps
    for i, k := range keys {
        if i == lastIndex {
            if v, ok := value.(StrMap); ok {
                value = v.toSyncMap()
            }
            maps.Store(k, value)
        } else {
            if m := getConfigMap(maps, k); m != nil {
                maps = m
            } else {
                newMap := &sync.Map{}
                maps.Store(k, newMap)
                maps = newMap
            }
        }
    }
}

// Get 获取配置项，允许使用点式获取，如：core.name
func Get[T any](key string, defaultValue ...T) (t T) {
    var value any
    anyDefaultValue := util.CastArray[T, any](defaultValue)
    switch any(t).(type) {
    case string:
        value = String(key, anyDefaultValue...)
    case int:
        value = Int(key, anyDefaultValue...)
    case int64:
        value = Int64(key, anyDefaultValue...)
    case int32:
        value = Int32(key, anyDefaultValue...)
    case uint:
        value = Uint(key, anyDefaultValue...)
    case bool:
        value = Bool(key, anyDefaultValue...)
    default:
        value = Env(key, anyDefaultValue...)
    }
    if value != nil {
        return value.(T)
    }
    return
}

// String 获取 String 类型的配置信息
func String(key string, defaultValue ...any) string {
    return util.ToString(Env(key, defaultValue...))
}

// Int 获取 Int 类型的配置信息
func Int(key string, defaultValue ...any) int {
    return util.ToInt(Env(key, defaultValue...))
}

// Int64 获取 Int64 类型的配置信息
func Int64(key string, defaultValue ...any) int64 {
    return util.ToInt64(Env(key, defaultValue...))
}

// Int32 获取 Int64 类型的配置信息
func Int32(key string, defaultValue ...any) int32 {
    return util.ToInt32(Env(key, defaultValue...))
}

// Uint 获取 Uint 类型的配置信息
func Uint(key string, defaultValue ...any) uint {
    return util.ToUint(Env(key, defaultValue...))
}

// Bool 获取 Bool 类型的配置信息
func Bool(key string, defaultValue ...any) bool {
    return util.ToBool(Env(key, defaultValue...))
}

// 获取下一级设置的 map，key 不存在或 value 不是 ConfigMap 则返回 nil
func getConfigMap(config ConfigMap, key string) ConfigMap {
    if val, ok := config.Load(key); ok {
        if m, ok := val.(ConfigMap); ok {
            return m
        } else if m, ok := val.(map[string]any); ok {
            return StrMap(m)
        }
    }
    return nil
}

// toSyncMap 转换为 sync.Map 以支持并发安全
func (m StrMap) toSyncMap() *sync.Map {
    syncMap := sync.Map{}
    for k, v := range m {
        if strMap, ok := v.(StrMap); ok {
            v = strMap.toSyncMap()
        }
        syncMap.Store(k, v)
    }
    return &syncMap
}
