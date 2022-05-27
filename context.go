package kaligo

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "math"
    "mime/multipart"
    "net/http"
    "net/url"
    "os"
    "strings"

    // "strconv"
    "sync"
    "time"

    "github.com/owner888/kaligo/cache"
    "github.com/owner888/kaligo/database"
    "github.com/owner888/kaligo/render"
    "github.com/owner888/kaligo/util"
)

// abortIndex represents a typical value used in abort functions.
const abortIndex int8 = math.MaxInt8 >> 1

// Context is use for ServeHTTP goroutine
type Context struct {
    UID string
    mux *Mux

    ResponseWriter http.ResponseWriter
    Request        *http.Request

    Params   Params
    handlers []HandlerFunc
    index    int8
    fullPath string

    // Keys is a key/value pair exclusively for the context of each request.
    // 用于中间件之间共享数据, 不同请求之间是不共享的, 需要使用 session 或者 redis
    // Keys map[string]any
    Keys sync.Map

	// Errors is a list of errors attached to all the handlers/middlewares who used this context.
	Errors errorMsgs

	// Accepted defines a list of manually accepted formats for content negotiation.
	Accepted []string

    // QueryCache caches the query result from c.Request.URL.Query().
    QueryCache util.UrlValues

    // FormCache caches c.Request.PostForm, which contains the parsed form data from POST, PATCH,
    // or PUT body parameters.
    FormCache util.UrlValues

    // SameSite allows a server to define a cookie attribute making it impossible for
    // the browser to send this cookie along with cross-site requests.
    sameSite http.SameSite

    // Cache is a key/value pair exclusively for the context of all request.
    Cache cache.Cache

    DB *database.DB

    Timer *Timer
}

var MaxMultipartMemory int64 = 32 << 20 // 32 MiB

func (c *Context) Reset() {
    c.Params = c.Params[:0]
    c.handlers = nil
    c.index = -1

    c.fullPath = ""
    c.Keys = sync.Map{}
    c.Errors = c.Errors[0:0]
    c.Accepted = nil
    c.QueryCache = nil
    c.FormCache = nil
    c.sameSite = 0
}

// FullPath returns a matched route full path. For not found routes
// returns an empty string.
//     router.GET("/user/:id", func(c *gin.Context) {
//         c.FullPath() == "/user/:id" // true
//     })
func (c *Context) FullPath() string {
    return c.fullPath
}

/************************************/
/*********** FLOW CONTROL ***********/
/************************************/
func (c *Context) Next() {
    c.index++
    for c.index < int8(len(c.handlers)) {
        c.handlers[c.index](c)
        c.index++
    }
}

func (c *Context) IsAborted() bool {
    return c.index >= abortIndex
}

func (c *Context) Abort() {
    c.index = abortIndex
}

// AbortWithStatus calls `Abort()` and writes the headers with the specified status code.
// For example, a failed attempt to authenticate a request could use: context.AbortWithStatus(401).
func (c *Context) AbortWithStatus(code int) {
    c.Status(code)
    c.ResponseWriter.WriteHeader(code)
    c.Abort()
}

// AbortWithStatusJSON calls `Abort()` and then `JSON` internally.
// This method stops the chain, writes the status code and return a JSON body.
// It also sets the Content-Type as "application/json".
func (c *Context) AbortWithStatusJSON(code int, jsonObj any) {
    c.Abort()
    c.JSON(code, jsonObj)
}

// AbortWithError calls `AbortWithStatus()` and `Error()` internally.
// This method stops the chain, writes the status code and pushes the specified error to `c.Errors`.
// See Context.Error() for more details.
func (c *Context) AbortWithError(code int, err error) *Error {
    c.AbortWithStatus(code)
    return c.Error(err)
}

/************************************/
/********* ERROR MANAGEMENT *********/
/************************************/

// Error attaches an error to the current context. The error is pushed to a list of errors.
// It's a good idea to call Error for each error that occurred during the resolution of a request.
// A middleware can be used to collect all the errors and push them to a database together,
// print a log, or append it in the HTTP response.
// Error will panic if err is nil.
func (c *Context) Error(err error) *Error {
    if err == nil {
        panic("err is nil")
    }

    parsedError, ok := err.(*Error)
    if !ok {
        parsedError = &Error{
            Err:  err,
            Type: ErrorTypePrivate,
        }
    }

    c.Errors = append(c.Errors, parsedError)
    return parsedError
}

// CallController 调用其他接口（非异步）
func (a *Context) CallController(controller Interface, method string, params Params) (ret any, err error) {
    if a.mux != nil {
        return a.mux.CallController(controller, method, params)
    }
    return nil, errors.New("service not initialized.")
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value any) {
    c.Keys.Store(key, value)
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exist it returns (nil, false)
// func (c *Context) Get(key string) (value any, exists bool) {
func (c *Context) Get(key string) (value any, exists bool) {
    val, ok := c.Keys.Load(key)
    return val, ok
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) interface{} {
    if value, exists := c.Get(key); exists {
        return value
    }
    panic("Key \"" + key + "\" does not exist")
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) {
    if val, ok := c.Get(key); ok && val != nil {
        s, _ = val.(string)
    }
    return
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) {
    if val, ok := c.Get(key); ok && val != nil {
        b, _ = val.(bool)
    }
    return
}

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) (i int) {
    if val, ok := c.Get(key); ok && val != nil {
        i, _ = val.(int)
    }
    return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Context) GetInt64(key string) (i64 int64) {
    if val, ok := c.Get(key); ok && val != nil {
        i64, _ = val.(int64)
    }
    return
}

// GetUint returns the value associated with the key as an unsigned integer.
func (c *Context) GetUint(key string) (ui uint) {
    if val, ok := c.Get(key); ok && val != nil {
        ui, _ = val.(uint)
    }
    return
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func (c *Context) GetUint64(key string) (ui64 uint64) {
    if val, ok := c.Get(key); ok && val != nil {
        ui64, _ = val.(uint64)
    }
    return
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) (f64 float64) {
    if val, ok := c.Get(key); ok && val != nil {
        f64, _ = val.(float64)
    }
    return
}

// GetTime returns the value associated with the key as time.
func (c *Context) GetTime(key string) (t time.Time) {
    if val, ok := c.Get(key); ok && val != nil {
        t, _ = val.(time.Time)
    }
    return
}

// GetDuration returns the value associated with the key as a duration.
func (c *Context) GetDuration(key string) (d time.Duration) {
    if val, ok := c.Get(key); ok && val != nil {
        d, _ = val.(time.Duration)
    }
    return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Context) GetStringSlice(key string) (ss []string) {
    if val, ok := c.Get(key); ok && val != nil {
        ss, _ = val.([]string)
    }
    return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {
    if val, ok := c.Get(key); ok && val != nil {
        sm, _ = val.(map[string]interface{})
    }
    return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Context) GetStringMapString(key string) (sms map[string]string) {
    if val, ok := c.Get(key); ok && val != nil {
        sms, _ = val.(map[string]string)
    }
    return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
    if val, ok := c.Get(key); ok && val != nil {
        smss, _ = val.(map[string][]string)
    }
    return
}

// Del is used to delete value from store with a key.
func (c *Context) Del(key string) {
    c.Keys.Delete(key)
}

// Clear is used to Clear map.
func (c *Context) Clear(key string) {
    c.Keys = sync.Map{}
}

// """"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
// => INPUT DATA
// """"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""

// HeaderValue returns value from request headers.
func (c *Context) HeaderValue(key string) string {
    return c.requestHeader(key)
}

// RawDataValue returns stream data.
func (c *Context) RawDataValue() ([]byte, error) {
    return ioutil.ReadAll(c.Request.Body)
}

// key 找不到 value 时返回的默认值，用于
func (c *Context) getDefaultValue(defaultValue ...string) string {
    if len(defaultValue) > 0 {
        return defaultValue[0]
    }
    return ""
}

// Param returns the value of the URL param.
// It is a shortcut for c.Params.ByName(key)
//    r.AddRoute("/user/:id([0-9]+)", map[string]string{
//        "GET"   : "GetUser",
//    }, &controller.UserController{})
//
//    // a GET request to /user/10
//    id := c.RouterValue("id") // id == "10"
func (c *Context) ParamValue(key string, defaultValue ...string) string {
    return c.Params.ByName(key, defaultValue...)
}

// func (c *Context) ParamInt(key string, defaultValue ...any) int {
//     param := c.Params.ByName(key, defaultValue...)
//     defaultVal := strconv.Itoa(param)
//     intVar, err := strconv.Atoi(param)
//     if err != nil {
//         return 0
//     }
//
//     return intVar
// }

// AddParam adds param to context and
// replaces path param key with given value for e2e testing purposes
// Example Route: "/user/:id"
// AddParam("id", 1)
// Result: "/user/1"
// func (c *Context) AddParam(key, value string) {
//     c.Params = append(c.Params, Params{Key: key, Value: value})
// }

func (c *Context) initQueryCache() {
    if c.QueryCache == nil {
        if c.Request != nil {
            c.QueryCache = util.UrlValues(c.Request.URL.Query())
        } else {
            c.QueryCache = util.UrlValues{}
        }
    }
}

// QueryValue returns the keyed url query value if it exists,
// It is shortcut for `c.Request.URL.Query().Get(key)`
//     GET /path?id=1234&name=Manu&value=
//     c.QueryValue("id") == "1234"
//     c.QueryValue("name") == "Manu"
//     c.QueryValue("value") == ""
//     c.QueryValue("wtf") == ""
func (c *Context) QueryValue(key string, defaultValue ...string) string {
    c.initQueryCache()
    if value, ok := c.QueryCache.Get(key); ok {
        return value
    }
    return c.getDefaultValue(defaultValue...)
}

// QueryArray returns a slice of strings for a given query key.
// The length of the slice depends on the number of params with the given key.
func (c *Context) QueryArray(key string) (values []string) {
    c.initQueryCache()
    values, _ = c.QueryCache[key]
    return
}

// QueryMap returns a map for a given query key.
func (c *Context) QueryMap(key string) (dicts map[string]string) {
    c.initQueryCache()
    dicts, _ = c.get(c.QueryCache, key)
    return
}

// get is an internal method and returns a map which satisfy conditions.
func (c *Context) get(m map[string][]string, key string) (map[string]string, bool) {
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

func (c *Context) initFormCache() {
    if c.FormCache == nil {
        c.FormCache = util.DefaultMIMEParsers.ParseValues(c.Request)
    }
}

// FormValue returns the specified key from a POST urlencoded form or multipart form
func (c *Context) FormValue(key string, defaultValue ...string) string {
    c.initFormCache()
    if value, ok := c.FormCache.Get(key); ok {
        return value
    }
    return c.getDefaultValue(defaultValue...)
}

// JsonBodyValue 解析 application/json、XML、ymal 等数据到 struct
//    var user User
//    c.JsonBodyValue(&user)
//    var users []User
//    c.JsonBodyValue(&users)
func (c *Context) JsonBodyValue(obj any) error {
    defer c.Request.Body.Close()
    return json.NewDecoder(c.Request.Body).Decode(obj)
}

// FormFile returns the first file for the provided form key.
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
    if c.Request.MultipartForm == nil {
        if err := c.Request.ParseMultipartForm(MaxMultipartMemory); err != nil {
            return nil, err
        }
    }
    f, fh, err := c.Request.FormFile(name)
    if err != nil {
        return nil, err
    }
    f.Close()
    return fh, err
}

// MultipartForm is the parsed multipart form, including file uploads.
func (c *Context) MultipartForm() (*multipart.Form, error) {
    err := c.Request.ParseMultipartForm(MaxMultipartMemory)
    return c.Request.MultipartForm, err
}

// SaveUploadedFile uploads the form file to specific dst.
func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, src)
    return err
}

// ClientIP returns a client ip. returns 127.0.0.1 if request from local machine
// 需要重构
func (c *Context) ClientIP() (ip string) {
    if c.Request == nil {
        return ""
    }
    remoteAddr := c.requestHeader("Remote_addr")
    if remoteAddr == "" {
        remoteAddr = c.Request.RemoteAddr
    }

    if remoteAddr != "" {
        remoteAddrArr := strings.Split(remoteAddr, ":")
        ip = remoteAddrArr[0]
    }
    if ip == "" || ip == "[" {
        ip = "127.0.0.1"
    }
    return
}

func (c *Context) requestHeader(key string) string {
    if c.Request != nil {
        return c.Request.Header.Get(key)
    }
    return ""
}

// """"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
// => RESPONSE RENDERING
// """"""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""

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

func (c *Context) CookieValue(name string, defaultValue ...string) string {
    cookie, err := c.Cookie(name)
    if err != nil {
        if len(defaultValue) > 0 {
            return defaultValue[0]
        }
        return ""
    }
    return cookie
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
        c.Status(code)
        // c.ResponseWriter.WriteHeaderNow()
        return
    }

    if err := r.Render(c.ResponseWriter); err != nil {
        panic(err)
    }
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(code int, obj any) {
    c.Render(code, render.JSON{Data: obj})
}

func (c *Context) HTML(code int, name string, obj any) {
    instance := c.mux.HTMLRender.Instance(name, obj)
    c.Render(code, instance)
}

func (c *Context) View(name string, obj interface{}) {
    c.HTML(http.StatusOK, name, obj)
}

// String writes the given string into the response body.
func (c *Context) String(code int, format string, values ...any) {
    c.Render(code, render.String{Format: format, Data: values})
}

// Redirect returns an HTTP redirect to the specific location.
func (c *Context) Redirect(code int, location string) {
    c.Render(-1, render.Redirect{
        Code:     code,
        Location: location,
        Request:  c.Request,
    })
}

// Data writes some data into the body stream and updates the HTTP code.
func (c *Context) Data(code int, contentType string, data []byte) {
    c.Render(code, render.Data{
        ContentType: contentType,
        Data:        data,
    })
}

// File writes the specified file into the body stream in an efficient way.
func (c *Context) File(filepath string) {
    http.ServeFile(c.ResponseWriter, c.Request, filepath)
}

// FileFromFS writes the specified file from http.FileSystem into the body stream in an efficient way.
func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
    defer func(old string) {
        c.Request.URL.Path = old
    }(c.Request.URL.Path)

    c.Request.URL.Path = filepath

    http.FileServer(fs).ServeHTTP(c.ResponseWriter, c.Request)
}

// FileAttachment writes the specified file into the body stream in an efficient way
// On the client side, the file will typically be downloaded with the given filename
func (c *Context) FileAttachment(filepath, filename string) {
    c.ResponseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
    http.ServeFile(c.ResponseWriter, c.Request, filepath)
}

/***** 当前时间 *****/
// NowFunc 当前时间获取函数类型
type NowFunc func() time.Time

// 当前时间获取函数，默认获取当前时区
//  kaligo.Now = func() time.Time {
//      return time.Now().UTC()
//  }
var Now NowFunc = func() time.Time {
    return time.Now().Local()
}

// Now 获取当前时间（默认当前时区）
func (c *Context) Now() time.Time {
    return Now()
}

// NowTimestamp 获取当前时间戳，统一时间戳格式（默认当前时区秒数）
func (c *Context) NowTimestamp() int64 {
    return Now().Unix()
}

// SetAccepted sets Accept header data.
func (c *Context) SetAccepted(formats ...string) {
    c.Accepted = formats
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
func (c *Context) Value(key any) any {
    if key == 0 {
        return c.Request
    }
    if keyAsString, ok := key.(string); ok {
        val, _ := c.Get(keyAsString)
        return val
    }
    if c.Request == nil || c.Request.Context() == nil {
        return nil
    }
    return c.Request.Context().Value(key)
}

/* vim: set expandtab: */
