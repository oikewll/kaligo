package util

import (
    "net/url"
    "strings"
)

//
type UrlValues url.Values

// Value returns the keyed url query value if it exists,
// otherwise it returns an empty string `("")`.
// It is shortcut for `c.Request.URL.Value().Get(key)`
//     GET /path?id=1234&name=Manu&value=
//     c.Value("id") == "1234"
//     c.Value("name") == "Manu"
//     c.Value("value") == ""
//     c.Value("wtf") == ""
func (c UrlValues) Value(key string) (value string) {
    value, _ = c.Get(key)
    return
}

// Default returns the keyed url query value if it exists,
// otherwise it returns the specified defaultValue string.
// See: Query() and GetQuery() for further information.
//     GET /?name=Manu&lastname=
//     c.Default("name", "unknown") == "Manu"
//     c.Default("id", "none") == "none"
//     c.Default("lastname", "none") == ""
func (c UrlValues) Default(key, defaultValue string) string {
    if value, ok := c.Get(key); ok {
        return value
    }
    return defaultValue
}

// Get is like Query(), it returns the keyed url query value
// if it exists `(value, true)` (even when the value is an empty string),
// otherwise it returns `("", false)`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
//     GET /?name=Manu&lastname=
//     ("Manu", true) == c.Get("name")
//     ("", false) == c.Get("id")
//     ("", true) == c.Get("lastname")
func (c UrlValues) Get(key string) (string, bool) {
    if values, ok := c.GetArray(key); ok {
        return values[0], ok
    }
    return "", false
}

// Array returns a slice of strings for a given query key.
// The length of the slice depends on the number of params with the given key.
func (c UrlValues) Array(key string) (values []string) {
    values, _ = c.GetArray(key)
    return
}

// GetArray returns a slice of strings for a given query key, plus
// a boolean value whether at least one value exists for the given key.
func (c UrlValues) GetArray(key string) (values []string, ok bool) {
    values, ok = c[key]
    return
}

// Map returns a map for a given query key.
func (c UrlValues) Map(key string) (dicts map[string]string) {
    dicts, _ = c.GetMap(key)
    return
}

// GetMap returns a map for a given query key, plus a boolean value
// whether at least one value exists for the given key.
func (c UrlValues) GetMap(key string) (map[string]string, bool) {
    return c.get(c, key)
}

// get is an internal method and returns a map which satisfy conditions.
func (c UrlValues) get(m map[string][]string, key string) (map[string]string, bool) {
    dicts := make(map[string]string)
    exist := false
    for k, v := range m {
        if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
            if j := strings.IndexByte(k[i+1:], ']'); j >= 1 {
                exist = true
                dicts[k[i+1:][:j]] = v[0]
            }
        }
    }
    return dicts, exist
}
