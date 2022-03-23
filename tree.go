package kaligo

import "github.com/owner888/kaligo/util"

// Param is a single URL parameter, consisting of a key and a value.
type Param[T any] struct {
    Key   string
    Value T
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param[string]

// ?name=test&age=10
// ?age=10
// ?name=&age=10
// Get("name", "kk")
// Get returns the value of the first Param which key matches the given name and a boolean true.
// If no matching Param is found, an empty string is returned and a boolean false .
func (ps Params) Get(name string, defaultValue ...string) (string, bool) {
    for _, entry := range ps {
        if entry.Key == name {
            return entry.Value, true
        }
    }
    var ret string = ""
    if len(defaultValue) != 0 {
        ret = defaultValue[0]
    }
    return ret, false
}

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string, defaultValue ...string) (va string) {
    va, _ = ps.Get(name, defaultValue...)
    return
}

func (ps Params) GetAnyKeyValue(key string, defaultValue ...any) (v any, ok bool) {
    v, ok = ps.Get(key)
    if !ok {
        if len(defaultValue) != 0 {
            v = defaultValue[0]
        }
    }
    return
}

type AnyKeyValueGetter interface {
    GetAnyKeyValue(key string, defaultValue ...any) (any, bool)
}

func Get[T any](kv AnyKeyValueGetter, key string, defaultValue ...T) T {
    v, _ := GetValue(kv, key, defaultValue...)
    return v
}

func GetValue[T any](kv AnyKeyValueGetter, key string, defaultValue ...T) (T, bool) {
    anyDefaultValue := util.CastArray[T, any](defaultValue)
    v, ok := kv.GetAnyKeyValue(key, anyDefaultValue...)
    return util.To[T](v), ok
}
