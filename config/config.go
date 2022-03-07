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
func Env(key string, defaultValue ...interface{}) interface{} {
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

// Add 新增配置项
func Add(key string, value any) {
    configMaps.Store(key, value)
}

// Set 设置配置项
func Set(key string, value any) {
    var keys []string = strings.Split(key, ".")
    lastIndex := len(keys) - 1
    var maps ConfigMap = &configMaps
    for i, k := range keys {
        if i == lastIndex {
            maps.Store(k, value)
        } else {
            if m := getConfigMap(maps, k); m != nil {
                maps = m
            } else {
                newMap := sync.Map{}
                maps.Store(k, newMap)
                maps = &newMap
            }
        }
    }
}

// Get 获取配置项，允许使用点式获取，如：core.name
func Get[T any](key string, defaultValue ...T) (t T) {
    var value any
    anyDefaultValue := cast[T, any](defaultValue)
    switch any(t).(type) {
    case string:
        value = GetString(key, anyDefaultValue...)
    case int:
        value = GetInt(key, anyDefaultValue...)
    case int64:
        value = GetInt64(key, anyDefaultValue...)
    case int32:
        value = GetInt32(key, anyDefaultValue...)
    case uint:
        value = GetUint(key, anyDefaultValue...)
    case bool:
        value = GetBool(key, anyDefaultValue...)
    default:
        value = Env(key, anyDefaultValue...)
    }
    if value != nil {
        return value.(T)
    }
    return
}

// GetString 获取 String 类型的配置信息
func GetString(path string, defaultValue ...interface{}) string {
    return util.ToString(Env(path, defaultValue...))
}

// GetInt 获取 Int 类型的配置信息
func GetInt(path string, defaultValue ...interface{}) int {
    return util.ToInt(Env(path, defaultValue...))
}

// GetInt64 获取 Int64 类型的配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
    return util.ToInt64(Env(path, defaultValue...))
}

// GetInt32 获取 Int64 类型的配置信息
func GetInt32(path string, defaultValue ...interface{}) int32 {
    return util.ToInt32(Env(path, defaultValue...))
}

// GetUint 获取 Uint 类型的配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
    return util.ToUint(Env(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
    return util.ToBool(Env(path, defaultValue...))
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

func cast[T any, B any](in []T) []B {
    var interfaceSlice []B = make([]B, len(in))
    for i, d := range in {
        interfaceSlice[i] = any(d).(B)
    }
    return interfaceSlice
}
