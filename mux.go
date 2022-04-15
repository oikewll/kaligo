package kaligo

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "path"
    "reflect"
    "regexp"
    "strings"
    "sync"
    "sync/atomic"

    "github.com/owner888/kaligo/cache"
    "github.com/owner888/kaligo/database"
    "github.com/owner888/kaligo/util"
)

// 强制要求 Mux 实现 Router
// 假如 Router 是一个第三方库接口，如果第三方库新版动了这个接口，我们发现不了，编译也不报错，只有用的时候才报错，但写了这句，就编译不过了
var _ Router = &Mux{}

// var muxCounter uint32

// Mux is use for add Route struct and StaticRoute struct
type Mux struct {
    // Handler       http.Handler // http.ServeMux
    routes       []*Route
    staticRoutes []*StaticRoute
    DB           *database.DB
    Cache        cache.Cache // interface 本身是指针，不需要用 *cache.Cache，struct 才需要
    pool         sync.Pool   // Context 复用
    Timer        *Timer
    // 统计
    requestCount  uint32
    responseCount uint32
    notFoundCount uint32
}

// 只会调用一次
func NewRouter() *Mux {
    mux := &Mux{}
    cache, err := cache.New()
    if err != nil {
        panic(err)
    }
    mux.Cache = cache
    mux.Timer = NewTimer(mux)
    mux.pool.New = func() any {
        return &Context{DB: mux.DB, Cache: mux.Cache, Timer: mux.Timer}
    }
    return mux
}

// String is the text representation of the collector.
// It contains useful debug information about the collector's internals
func (a *Mux) String() string {
    // return fmt.Sprintf(
    //     "Requests made: %d (%d responses) | Callbacks: OnRequest: %d, OnHTML: %d, OnResponse: %d, OnError: %d",
    //     atomic.LoadUint32(&a.requestCount),
    //     atomic.LoadUint32(&a.responseCount),
    //     // len(c.requestCallbacks),
    //     // len(c.htmlCallbacks),
    //     // len(c.responseCallbacks),
    //     // len(c.errorCallbacks),
    // )

    return fmt.Sprintf(
        "Requests made: %d (%d responses) | NotFound: %d",
        atomic.LoadUint32(&a.requestCount),
        atomic.LoadUint32(&a.responseCount),
        atomic.LoadUint32(&a.notFoundCount),
    )
}

// AddDB is use for add a db struct
func (a *Mux) AddDB(db *database.DB) {
    a.DB = db
}

// AddStaticRoute is use for add a static file route
func (a *Mux) AddStaticRoute(prefix, staticDir string) {
    route := &StaticRoute{}
    route.Prefix = prefix
    route.StaticDir = staticDir

    a.staticRoutes = append(a.staticRoutes, route)
}

// AddRoute is use for add a http route
// https://expressjs.com/en/5x/api.html
func (a *Mux) AddRoute(pattern string, m map[string]string, c Interface) {
    parts := strings.Split(pattern, "/")

    j := 0
    params := make(map[int]string)
    for i, part := range parts {
        if strings.HasPrefix(part, ":") {
            expr := "([^/]+)"

            // a user may choose to override the defult expression
            // https://expressjs.com/en/5x/api.html
            // similar to expressjs: ‘/user/:id([0-9]+)’
            // This will match paths starting with /abc and /xyz: /\/abc|\/xyz/
            if index := strings.Index(part, "("); index != -1 {
                // 正则匹配：/user/:id([0-9]+)
                expr = part[index:]
                part = part[1:index]
            } else {
                // 参数注册：/:param
                part = part[1:]
            }

            params[j] = part
            parts[i] = expr
            j++
        }
    }

    // recreate the url pattern, with parameters replaced
    // by regular expressions. then compile the regex
    pattern = strings.Join(parts, "/")
    regex, regexErr := regexp.Compile(pattern)
    if regexErr != nil {
        // TODO add error handling here to avoid panic
        panic(regexErr)
    }

    // now create the Route
    route := &Route{
        Regex:          regex,
        Methods:        m,
        Params:         params,
        ControllerType: reflect.Indirect(reflect.ValueOf(c)).Type(),
    }

    a.routes = append(a.routes, route)
}

// ServeHTTP is a http Handler
// type Handler interface {
//     ServeHTTP(ResponseWriter, *Request)
// }
func (a *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 不想走这里的时候才需要使用注入
    // if a.Handler != nil { // 有注入 Handler 时只使用注入的 Handler，不走默认逻辑
    //     a.Handler.ServeHTTP(w, r)
    //     return
    // }

    requestPath := r.URL.RawPath
    if requestPath == "" {
        requestPath = r.URL.Path // 会自动 unescape, 先用 RawPath ,不行才用这个
    }
    requestPath = path.Clean(requestPath) // 多个反斜杠变成一个

    // find a matching Route
    for _, staticRoute := range a.staticRoutes {
        // 如果设置的静态文件目录包含在url 路径信息中
        if strings.HasPrefix(requestPath, staticRoute.Prefix) {
            if len(requestPath) > len(staticRoute.Prefix) && requestPath[len(staticRoute.Prefix)] != '/' {
                continue
            }
            file := path.Join(staticRoute.StaticDir, requestPath[len(staticRoute.Prefix):])
            if util.FileExists(file) {
                http.ServeFile(w, r, file)
            } else {
                http.NotFound(w, r)
            }
            return
        }
    }

    // 请求数 +1
    atomic.AddUint32(&a.requestCount, 1)

    if r.Method == "OPTIONS" {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "*")
        w.Header().Set("Access-Control-Allow-Methods", "*")
        w.WriteHeader(http.StatusNoContent)
        return
    }

    var matchRouted bool

    // find a matching Route
    for _, route := range a.routes {
        // check if Route pattern matches url
        if !route.Regex.MatchString(requestPath) {
            continue
        }

        // get submatches (params)
        matches := route.Regex.FindStringSubmatch(requestPath)

        // double check that the Route matches the URL pattern.
        if len(matches[0]) != len(requestPath) {
            continue
        }

        params := make(Params, len(route.Params))

        if len(route.Params) > 0 {
            // add url parameters to the query param map
            values := r.URL.Query()
            for i, match := range matches[1:] {
                values.Add(route.Params[i], match)
                params[i] = Param[string]{route.Params[i], match}
            }

            // reassemble query params and add to RawQuery
            // r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
            // r.URL.RawQuery = url.Values(values).Encode()
        }

        var ok bool
        var m string
        // Request callback
        if m, ok = route.Methods[r.Method]; !ok {
            http.NotFound(w, r)

            // 错误数 +1
            atomic.AddUint32(&a.notFoundCount, 1)
            return
        }

        a.controllerMethodCall(route.ControllerType, m, w, r, params)
        matchRouted = true
    }

    // 相应数 +1
    atomic.AddUint32(&a.responseCount, 1)

    // if no matches to url, throw a not found exception
    if matchRouted == false {
        http.NotFound(w, r)
    }
}

func (a *Mux) CallController(controller Interface, method string, params Params) (ret any, err error) {
    return a.controllerMethodCall(reflect.Indirect(reflect.ValueOf(controller)).Type(), method, nil, nil, params)
}

func (a *Mux) controllerMethodCall(controllerType reflect.Type, m string, w http.ResponseWriter, r *http.Request, params Params) (ret []reflect.Value, err error) {
    ctx := a.pool.Get().(*Context)
    defer a.pool.Put(ctx)
    ctx.Reset()
    ctx.mux = a
    if a.DB != nil {
        ctx.DB = &database.DB{Config: a.DB.Config}
    }
    ctx.ResponseWriter = w
    ctx.Request = r
    ctx.Params = params
    ret, err = runController(controllerType, m, ctx, params)
    if err != nil && w != nil {
        http.NotFound(w, r)
    }
    return ret, err
}

// Router 接口实现

// Use 添加一个中间件
func (a *Mux) Use(middlewares ...func(http.Handler) http.Handler) {

}

// With adds inline middlewares for an endpoint handler.
func (a *Mux) With(middlewares ...func(http.Handler) http.Handler) Router {
    return a
}

type TLSOption struct {
    crtFile, keyFile string
}

// Run is to run a router
func Run(mux *Mux, addr string, options ...TLSOption) {
    // addr := fmt.Sprintf("%s:%d", cfg.host, cfg.port)
    server := &http.Server{Addr: addr, Handler: mux}

    // Creating a waiting group that waits until the graceful shutdown procedure is done
    var wg sync.WaitGroup
    wg.Add(1)

    go func() {
        stop := make(chan os.Signal, 1)
        signal.Notify(stop, os.Interrupt)
        <-stop
        log.Println("Shutting down the HTTP server...")
        err := server.Shutdown(context.Background())
        if err != nil {
            log.Printf("Error during shutdown: %v\n", err)
        }
        wg.Done()
    }()

    log.Printf("ListenAndServe: %s", addr)

    var err error
    // err = http.ListenAndServe(port, mux)
    if options == nil {
        err = server.ListenAndServe()
    } else {
        err = server.ListenAndServeTLS(
            options[0].crtFile, // cfg.certificatePemFilePath,
            options[0].keyFile, // cfg.certificatePemPrivKeyFilePath,
        )
    }

    if err == http.ErrServerClosed { // graceful shutdown
        log.Println("Commencing server shutdown...")
        wg.Wait()
        log.Println("Server was gracefully shut down.")
    } else if err != nil {
        log.Printf("Server error: %v\n", err)
    }
}

// 2021/03/12 13:39:49 listening on port 3333...
// 2021/03/12 13:39:50 user initiated a request
// 2021/03/12 13:39:54 commencing server shutdown...
// 2021/03/12 13:40:00 user request is fulfilled
// 2021/03/12 13:40:01 server was gracefully shut down.

/* vim: set expandtab: */
