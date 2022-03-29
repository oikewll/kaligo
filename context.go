package kaligo

import (
    // "log"

    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "net/url"
    "os"
    "strings"
    "sync"
    "time"

    "github.com/owner888/kaligo/cache"
    "github.com/owner888/kaligo/database"
    "github.com/owner888/kaligo/render"
    "github.com/owner888/kaligo/util"
)

type SuccessJSON struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
    Data any    `json:"data"`
}

type ErrorJSON struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
}

// Context is use for ServeHTTP goroutine
type Context struct {
    ResponseWriter http.ResponseWriter
    Request        *http.Request

    Params   Params
    fullPath string

    DB *database.DB

    // Cache is a key/value pair exclusively for the context of all request.
    Cache cache.Cache

    // Timer 定时任务
    Timer *Timer

    // Keys is a key/value pair exclusively for the context of each request.
    // Keys map[string]any
    Keys sync.Map

    // queryCache caches the query result from c.Request.URL.Query().
    queryCache util.UrlValues

    // formCache caches c.Request.PostForm, which contains the parsed form data from POST, PATCH,
    // or PUT body parameters.
    formCache util.UrlValues

    // SameSite allows a server to define a cookie attribute making it impossible for
    // the browser to send this cookie along with cross-site requests.
    sameSite http.SameSite
}

var MaxMultipartMemory int64 = 32 << 20 // 32 MiB

func (c *Context) Reset() {
    c.Params = c.Params[:0]
    c.fullPath = ""
    c.Keys = sync.Map{}
    c.queryCache = nil
    c.formCache = nil
    c.sameSite = 0
    fmt.Print()
}

// FullPath returns a matched route full path. For not found routes
// returns an empty string.
//     router.GET("/user/:id", func(c *gin.Context) {
//         c.FullPath() == "/user/:id" // true
//     })
func (c *Context) FullPath() string {
    return c.fullPath
}

// ClientIP returns a client ip. returns 127.0.0.1 if request from local machine
func (c *Context) ClientIP() (ip string) {
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
    c.Render(-1, render.Redirect{
        Code:     code,
        Location: location,
        Request:  c.Request,
    })
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(code int, obj any) {
    c.Render(code, render.JSON{Data: obj})
}

// String writes the given string into the response body.
func (c *Context) String(code int, format string, values ...any) {
    c.Render(code, render.String{Format: format, Data: values})
}

// Data writes some data into the body stream and updates the HTTP code.
func (c *Context) Data(code int, contentType string, data []byte) {
    c.Render(code, render.Data{
        ContentType: contentType,
        Data:        data,
    })
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

/************************************/
/************ INPUT DATA ************/
/************************************/

// Param returns the value of the URL param.
// It is a shortcut for c.Params.ByName(key)
//     router.GET("/user/:id", func(c *gin.Context) {
//         // a GET request to /user/john
//         id := c.Param("id") // id == "john"
//     })
func (c *Context) Param(key string, defaultValue ...string) string {
	return c.Params.ByName(key, defaultValue...)
}

// AddParam adds param to context and
// replaces path param key with given value for e2e testing purposes
// Example Route: "/user/:id"
// AddParam("id", 1)
// Result: "/user/1"
// func (c *Context) AddParam(key, value string) {
//     c.Params = append(c.Params, Params{Key: key, Value: value})
// }

// GetHeader returns value from request headers.
func (c *Context) GetHeader(key string) string {
    return c.requestHeader(key)
}

// GetRawData returns stream data.
func (c *Context) GetRawData() ([]byte, error) {
    return ioutil.ReadAll(c.Request.Body)
}

func (c *Context) initQueryCache() {
    if c.queryCache == nil {
        if c.Request != nil {
            c.queryCache = util.UrlValues(c.Request.URL.Query())
        } else {
            c.queryCache = util.UrlValues{}
        }
    }
}

func (c *Context) getDefaultValue(defaultValue ...string) string {
    if len(defaultValue) > 0 {
        return defaultValue[0]
    }
    return ""
}

// QueryValue returns the keyed url query value if it exists,
// It is shortcut for `c.Request.URL.Query().Get(key)`
//     GET /path?id=1234&name=Manu&value=
// 	   c.QueryValue("id") == "1234"
// 	   c.QueryValue("name") == "Manu"
// 	   c.QueryValue("value") == ""
// 	   c.QueryValue("wtf") == ""
func (c *Context) QueryValue(key string, defaultValue ...string) string {
    c.initQueryCache()
    if value, ok := c.queryCache.Get(key); ok {
        return value
    }
    return c.getDefaultValue(defaultValue...)
}

func (c *Context) initFormCache() {
    if c.formCache == nil {
        c.formCache = util.DefaultMIMEParsers.ParseValues(c.Request)
    }
}

// FormValue returns the specified key from a POST urlencoded form or multipart form
func (c *Context) FormValue(key string, defaultValue ...string) string {
    c.initFormCache()
    if value, ok := c.formCache.Get(key); ok {
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

func (c *Context) CookieValue(key string, defaultValue ...string) string {
    ret, err := c.Cookie(key)
    if err != nil {
        if len(defaultValue) > 0 {
            return defaultValue[0]
        }
        return ""
    }
    return ret
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
        if val := c.Get(keyAsString); val != nil {
            return val
        }
    }
    if c.Request == nil || c.Request.Context() == nil {
        return nil
    }
    return c.Request.Context().Value(key)
}

/* vim: set expandtab: */
