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
    "io"
    "os"
    "strings"
    "sync"

    // "github.com/astaxie/beego/logs"

    "github.com/owner888/kaligo/util"
    "gopkg.in/yaml.v3"
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

// ConfigValue 是 config 中存储的动态值类型
type ConfigValue interface {
    // 获取值，ok == false 则使用 config.Get 的默认值
    Load() (value any, ok bool)
}

// 函数类型的 ConfigValue
type ConfigValueFunc func() (value any, ok bool)

func (f ConfigValueFunc) Load() (value any, ok bool) {
    return f()
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
                if val, ok := val.(ConfigValue); ok {
                    if val, ok := val.Load(); ok {
                        return val
                    }
                    break
                }
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

// 新增配置项(合并)，所有子级必须是 StrMap
func Add(value StrMap) {
    merge(&configMaps, value)
}

func merge(from ConfigMap, to StrMap) {
    for k, v := range to {
        switch t := v.(type) {
        case StrMap:
            subMap := getConfigMap(from, k)
            if subMap == nil {
                subMap = &StrMap{}
                from.Store(k, subMap)
            }
            merge(subMap, t)
        default:
            from.Store(k, v)
        }
    }
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
    anyDefaultValue := util.CastSlice[T, any](defaultValue)
    switch any(t).(type) {
    case string:
        value = String(key, anyDefaultValue...)
    case int:
        value = Int(key, anyDefaultValue...)
    case int64:
        value = Int64(key, anyDefaultValue...)
    case uint:
        value = Uint(key, anyDefaultValue...)
    case uint64:
        value = Uint64(key, anyDefaultValue...)
    case float64:
        value = Float64(key, anyDefaultValue...)
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

// Uint 获取 Uint 类型的配置信息
func Uint(key string, defaultValue ...any) uint {
    return util.ToUint(Env(key, defaultValue...))
}

// Uint64 获取 Uint64 类型的配置信息
func Uint64(key string, defaultValue ...any) uint64 {
    return util.ToUint64(Env(key, defaultValue...))
}

// Float64 获取 Float64 类型的配置信息
func Float64(key string, defaultValue ...any) float64 {
    return util.ToFloat64(Env(key, defaultValue...))
}

// Bool 获取 Bool 类型的配置信息
func Bool(key string, defaultValue ...any) bool {
    return util.ToBool(Env(key, defaultValue...))
}

// StringMapStringSlice 获取 map[string][]string 类型的配置信息
func StringMapStringSlice(key string, defaultValue ...any) map[string][]string {
    return util.ToStringMapStringSlice(Env(key, defaultValue...))
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

// LoadFiles 加载配置文件，支持 yaml
func LoadFiles(files ...string) {
    for _, v := range files {
        file, _ := os.Open(v)
        defer file.Close()
        LoadConfig(file, "yaml")
    }
}

// LoadConfig 加载配置 fileType = yaml
func LoadConfig(reader io.Reader, fileType string) {
    var value StrMap
    if fileType == "yaml" {
        yaml.NewDecoder(reader).Decode(&value)
        Add(value)
    }
}
