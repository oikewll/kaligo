package kaligo

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
    Key   string
    Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ?name=test&age=10
// ?age=10
// ?name=&age=10
// Get("name", "kk")
// Get returns the value of the first Param which key matches the given name and a boolean true.
// If no matching Param is found, an empty string is returned and a boolean false .
func (ps Params) Get(name string, defalutValue ...string) (string, bool) {
    for _, entry := range ps {
        if entry.Key == name {
            return entry.Value, true
        }
    }
    var ret string = ""
    if len(defalutValue) != 0 {
        ret = defalutValue[0]
    }
    return ret, false
}

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string, defalutValue ...string) (va string) {
    va, _ = ps.Get(name, defalutValue...)
    return
}
