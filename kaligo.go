package kaligo

import (
    "fmt"
    "log"
    "net/http"
    "path"
    "reflect"
    "regexp"

    // "runtime"
    "strings"
    "sync"
    "time"

    "github.com/owner888/kaligo/cache"
    "github.com/owner888/kaligo/contex"
    "github.com/owner888/kaligo/controller"
    "github.com/owner888/kaligo/database"
    "github.com/owner888/kaligo/routes"
    "github.com/owner888/kaligo/util"

    "github.com/astaxie/beego/logs"
)

// 定义当前package中使用的全局变量
var (
    err         error
    storeTimers sync.Map
    Timer       map[string]*time.Ticker
    Tasker      map[string]*time.Timer
    Mutex       sync.Mutex
)

// var db *database.DB
//
// func init() {
//     var err error
//     // db, err = database.Open(sqlite.Open(config.Get[string]("database.sqlite.file")))
//     db, err = database.Open(mysql.Open(config.Get[string]("database.mysql.dsn")))
//     if err != nil {
//         panic(err)
//     }
// }

func init() {
    logs.SetLogFuncCall(true)
    logs.SetLogFuncCallDepth(3)
}

// KaliGo is use for add Route struct and StaticRoute struct
type KaliGo struct {
    Handler      http.Handler // http.ServeMux
    routes       []*routes.Route
    staticRoutes []*routes.StaticRoute
    DB           *database.DB
    cache        cache.Cache // interface 本身是指针，不需要用 *cache.Cache，struct 才需要
    pool         sync.Pool   // Context 复用
}

func New() *KaliGo {
    kali := &KaliGo{}
    cache, err := cache.New()
    if err != nil {
        panic(err)
    }
    kali.cache = cache
    return kali
}

// AddDB is use for add a db struct
func (a *KaliGo) AddDB(db *database.DB) {
    a.DB = db
}

// DelTasker is the function for delete tasker
func DelTasker(name string) bool {
    tasker, ok := storeTimers.Load(name)
    if ok {
        tasker.(*time.Ticker).Stop()
        return true
    }
    return false
}

// DelTimer is the function for delete timer
func DelTimer(name string) bool {
    tasker, ok := storeTimers.Load(name)
    if ok {
        tasker.(*time.Timer).Stop()
        return true
    }
    return false
}

// AddTasker is the function for add tasker
// AddTasker("default", &control.Task{}, "import_database", "2014-10-15 15:33:00")
// func AddTasker(name string, control any, action string, taskTime string) {
func (a *KaliGo) AddTasker(name, taskTime, m string, c controller.Interface) {
    go func() {
        then, _ := time.ParseInLocation("2006-01-02 15:04:05", taskTime, time.Local)
        dura := then.Sub(time.Now())
        //fmt.Println(dura)
        if dura > 0 {
            timeTasker := time.AfterFunc(dura, func() {
                a.controllerMethodCall(reflect.Indirect(reflect.ValueOf(c)).Type(), m, nil, nil, nil)
            })
            storeTimers.Store(name, timeTasker)
        } else {
            logs.Error("定时任务 --- [ " + name + " ] --- 小于当前时间，将不会被执行")
        }
    }()
}

// AddTimer is the function for add timer, The interval is in microseconds
// router.AddTimer("import_database", 3000, "ImportDatabase", &controller.Get{})
func (a *KaliGo) AddTimer(name string, duration time.Duration, m string, c controller.Interface) {
    go func() {
        timeTicker := time.NewTicker(duration * time.Millisecond)
        storeTimers.Store(name, timeTicker)
        for {
            select {
            case <-timeTicker.C:
                a.controllerMethodCall(reflect.Indirect(reflect.ValueOf(c)).Type(), m, nil, nil, nil)
            }
        }
    }()
}

// AddStaticRoute is use for add a static file route
func (a *KaliGo) AddStaticRoute(prefix, staticDir string) {
    route := &routes.StaticRoute{}
    route.Prefix = prefix
    route.StaticDir = staticDir

    a.staticRoutes = append(a.staticRoutes, route)
}

// AddRoute is use for add a http route
// https://expressjs.com/en/5x/api.html
func (a *KaliGo) AddRoute(pattern string, m map[string]string, c controller.Interface) {
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
                expr = part[index:]
                part = part[1:index]
            } else {
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
    route := &routes.Route{}
    route.Regex = regex
    route.Methods = m
    route.Params = params
    route.ControllerType = reflect.Indirect(reflect.ValueOf(c)).Type()

    a.routes = append(a.routes, route)
}

// ServeHTTP is a http Handler
// type Handler interface {
//     ServeHTTP(ResponseWriter, *Request)
// }
func (a *KaliGo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if a.Handler != nil { // 有注入 Handler 是只使用注入的 Handler，不走默认逻辑
        a.Handler.ServeHTTP(w, r)
        return
    }

    matchRouted := false
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

        params := make(map[string]string)

        if len(route.Params) > 0 {
            // add url parameters to the query param map
            values := r.URL.Query()
            // logs.Debug("values", values)

            for i, match := range matches[1:] {
                values.Add(route.Params[i], match)
                params[route.Params[i]] = match
                // fmt.Println(route.Params[i])
                // fmt.Println(match)
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
        }

        a.controllerMethodCall(route.ControllerType, m, w, r, params)
        matchRouted = true
    }

    // if no matches to url, throw a not found exception
    if matchRouted == false {
        http.NotFound(w, r)
    }
}

func (a *KaliGo) controllerMethodCall(controllerType reflect.Type, m string, w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
    // Invoke the request handler
    vc := reflect.New(controllerType)

    // Init callback
    method := vc.MethodByName("Init")
    contex := &contex.Context{
        ResponseWriter: w,
        Request:        r,
        Params:         params,
        DB:             a.DB,
        Cache:          a.cache,
    }
    args := make([]reflect.Value, 2)
    args[0] = reflect.ValueOf(contex)
    args[1] = reflect.ValueOf(controllerType.Name())
    method.Call(args)

    args = make([]reflect.Value, 0)

    // Prepare callback
    method = vc.MethodByName("Prepare")
    method.Call(args)

    // Request callback
    method = vc.MethodByName(m)
    if !method.IsValid() {
        // if is HTTP callback
        if w != nil {
            http.NotFound(w, r)
        }
        return fmt.Errorf("Controller Method not exist")
    }
    method.Call(args)

    // Finish callback
    method = vc.MethodByName("Finish")
    method.Call(args)

    return err
}

// RunTLS is to run a tls server
func RunTLS(kali *KaliGo, port, crtFile, keyFile string) {
    log.Printf("ListenAndServe: %s - %s - %s", port, crtFile, keyFile)

    err := http.ListenAndServeTLS(port, crtFile, keyFile, kali)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

// Run is to run a router
func Run(kali *KaliGo, port string) {
    log.Printf("ListenAndServe: %s", port)

    err := http.ListenAndServe(port, kali)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
