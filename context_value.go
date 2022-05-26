package kaligo

/////////////// Param ///////////////

func (c *Context) ParamInt(key string, defaultValue ...int) int {
    return Get(c.Params, key, defaultValue...)
}

func (c *Context) ParamInt64(key string, defaultValue ...int64) int64 {
    return Get(c.Params, key, defaultValue...)
}

func (c *Context) ParamBool(key string, defaultValue ...bool) bool {
    return Get(c.Params, key, defaultValue...)
}

func (c *Context) ParamFloat64(key string, defaultValue ...float64) float64 {
    return Get(c.Params, key, defaultValue...)
}

/////////////// Form ///////////////

func (c *Context) FormInt(key string, defaultValue ...int) int {
    c.initFormCache()
    return Get(c.FormCache, key, defaultValue...)
}

func (c *Context) FormInt64(key string, defaultValue ...int64) int64 {
    c.initFormCache()
    return Get(c.FormCache, key, defaultValue...)
}

func (c *Context) FormBool(key string, defaultValue ...bool) bool {
    c.initFormCache()
    return Get(c.FormCache, key, defaultValue...)
}

func (c *Context) FormFloat64(key string, defaultValue ...float64) float64 {
    c.initFormCache()
    return Get(c.FormCache, key, defaultValue...)
}

/////////////// Query ///////////////

func (c *Context) QueryInt(key string, defaultValue ...int) int {
    c.initQueryCache()
    return Get(c.QueryCache, key, defaultValue...)
}

func (c *Context) QueryInt64(key string, defaultValue ...int64) int64 {
    c.initQueryCache()
    return Get(c.QueryCache, key, defaultValue...)
}

func (c *Context) QueryBool(key string, defaultValue ...bool) bool {
    c.initQueryCache()
    return Get(c.QueryCache, key, defaultValue...)
}

func (c *Context) QueryFloat64(key string, defaultValue ...float64) float64 {
    c.initQueryCache()
    return Get(c.QueryCache, key, defaultValue...)
}
