package contex

import (
    "io/ioutil"
    "net/http"
    "net/url"
    "sync"
    "time"

    "github.com/owner888/kaligo/cache"
    "github.com/owner888/kaligo/database"
    "github.com/owner888/kaligo/render"
    "github.com/owner888/kaligo/timer"
)

type SuccJSON struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
    Data any    `json:"data"`
}

type FailJSON struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
}

// Context is use for ServeHTTP goroutine
type Context struct {
    ResponseWriter http.ResponseWriter
    Request        *http.Request

    Params   map[string]string
    fullPath string

    DB *database.DB

    // Cache is a key/value pair exclusively for the context of all request.
    Cache cache.Cache

    // Timer 定时任务
    Timer *timer.Timer

    // Keys is a key/value pair exclusively for the context of each request.
    // Keys map[string]any
    Keys sync.Map

    // SameSite allows a server to define a cookie attribute making it impossible for
    // the browser to send this cookie along with cross-site requests.
    sameSite http.SameSite
}

// FullPath returns a matched route full path. For not found routes
// returns an empty string.
//     router.GET("/user/:id", func(c *gin.Context) {
//         c.FullPath() == "/user/:id" // true
//     })
func (c *Context) FullPath() string {
    return c.fullPath
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value any) {
    c.Keys.Store(key, value)
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exist it returns (nil, false)
// func (c *Context) Get(key string) (value any, exists bool) {
func (c *Context) Get(key string) any {
    val, ok := c.Keys.Load(key)
    if ok {
        // return (any)(val).(T)
        return val
    }
    return nil
}

// Del is used to delete value from store with a key.
func (c *Context) Del(key string) {
    c.Keys.Delete(key)
}

// Clear is used to Clear map.
func (c *Context) Clear(key string) {
    c.Keys = sync.Map{}
}

// Redirect returns an HTTP redirect to the specific location.
func (c *Context) Redirect(code int, location string) {
    http.Redirect(c.ResponseWriter, c.Request, location, code)
}

func (c *Context) ApiJSON(code int, msg string, param ... any) {
    // c.Header("Access-Control-Allow-Origin", "*")             //允许访问所有域
    // c.Header("Access-Control-Allow-Headers", "Content-Type") //header的类型
    // c.Header("content-type", "application/json")             //返回数据格式是json
    if len(param) == 0 {
        obj := &FailJSON{
            Code: code,
            Msg:  msg,
        }
        c.Render(http.StatusOK, render.JSON{Data: obj})
    } else {
        data := param[0]
        obj := &SuccJSON{
            Code: code,
            Msg:  msg,
            Data: data,
        }
        c.Render(http.StatusOK, render.JSON{Data: obj})
    }
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(code int, obj any) {
    c.Render(code, render.JSON{Data: obj})
}

// String writes the given string into the response body.
func (c *Context) String(code int, format string, values ...any) {
    // c.Render(code, render.String{Format: format, Data: values})
}

// Data writes some data into the body stream and updates the HTTP code.
func (c *Context) Data(code int, contentType string, data []byte) {
    // c.Render(code, render.Data{
    //     ContentType: contentType,
    //     Data:        data,
    // })
}

func (c *Context) requestHeader(key string) string {
    return c.Request.Header.Get(key)
}

/************************************/
/******** RESPONSE RENDERING ********/
/************************************/

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
    switch {
    case status >= 100 && status <= 199:
        return false
    case status == http.StatusNoContent:
        return false
    case status == http.StatusNotModified:
        return false
    }
    return true
}

// Status sets the HTTP response code.
func (c *Context) Status(code int) {
    c.ResponseWriter.WriteHeader(code)
}

// Header is an intelligent shortcut for c.Writer.Header().Set(key, value).
// It writes a header in the response.
// If value == "", this method removes the header `c.Writer.Header().Del(key)`
func (c *Context) Header(key, value string) {
    if value == "" {
        c.ResponseWriter.Header().Del(key)
        return
    }
    c.ResponseWriter.Header().Set(key, value)
}

// GetHeader returns value from request headers.
func (c *Context) GetHeader(key string) string {
    return c.requestHeader(key)
}

// GetRawData returns stream data.
func (c *Context) GetRawData() ([]byte, error) {
    return ioutil.ReadAll(c.Request.Body)
}

// SetSameSite with cookie
func (c *Context) SetSameSite(samesite http.SameSite) {
    c.sameSite = samesite
}

// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
    if path == "" {
        path = "/"
    }
    http.SetCookie(c.ResponseWriter, &http.Cookie{
        Name:     name,
        Value:    url.QueryEscape(value),
        MaxAge:   maxAge,
        Path:     path,
        Domain:   domain,
        SameSite: c.sameSite,
        Secure:   secure,
        HttpOnly: httpOnly,
    })
}

// Cookie returns the named cookie provided in the request or
// ErrNoCookie if not found. And return the named cookie is unescaped.
// If multiple cookies match the given name, only one cookie will
// be returned.
func (c *Context) Cookie(name string) (string, error) {
    cookie, err := c.Request.Cookie(name)
    if err != nil {
        return "", err
    }
    val, _ := url.QueryUnescape(cookie.Value)
    return val, nil
}

// Render writes the response headers and calls render.Render to render data.
func (c *Context) Render(code int, r render.Render) {
    c.Status(code)

    if !bodyAllowedForStatus(code) {
        r.WriteContentType(c.ResponseWriter)
        // c.ResponseWriter.WriteHeader()
        // c.ResponseWriter.WriteHeaderNow()
        return
    }

    if err := r.Render(c.ResponseWriter); err != nil {
        panic(err)
    }
}

/************************************/
/***** GOLANG.ORG/X/NET/CONTEXT *****/
/************************************/

// Deadline returns that there is no deadline (ok==false) when c.Request has no Context.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
    if c.Request == nil || c.Request.Context() == nil {
        return
    }
    return c.Request.Context().Deadline()
}

// Done returns nil (chan which will wait forever) when c.Request has no Context.
func (c *Context) Done() <-chan struct{} {
    if c.Request == nil || c.Request.Context() == nil {
        return nil
    }
    return c.Request.Context().Done()
}

// Err returns nil when c.Request has no Context.
func (c *Context) Err() error {
    if c.Request == nil || c.Request.Context() == nil {
        return nil
    }
    return c.Request.Context().Err()
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *Context) Value(key interface{}) interface{} {
    if key == 0 {
        return c.Request
    }
    if keyAsString, ok := key.(string); ok {
        if val := c.Get(keyAsString); val != nil {
            return val
        }
    }
    if c.Request == nil || c.Request.Context() == nil {
        return nil
    }
    return c.Request.Context().Value(key)
}
