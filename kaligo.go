package kaligo

import (
    "log"
    "net/http"
    "path"
    "reflect"
    "regexp"

    // "runtime"
    "strings"
    "sync"
    "time"

    "github.com/owner888/kaligo/contex"
    "github.com/owner888/kaligo/controller"
    "github.com/owner888/kaligo/database"
    "github.com/owner888/kaligo/routes"
    "github.com/owner888/kaligo/util"

    "github.com/astaxie/beego/logs"
)

// 定义当前package中使用的全局变量
var (
    storeTimers sync.Map
    Timer      map[string]*time.Ticker
    Tasker     map[string]*time.Timer
    Mutex      sync.Mutex
)

func init() {
    logs.SetLogFuncCall(true)
    logs.SetLogFuncCallDepth(3)
}

// App is use for add Route struct and StaticRoute struct
type App struct {
    http.Handler // http.ServeMux
    routes       []*routes.Route
    staticRoutes []*routes.StaticRoute
    db           *database.DB
}

// AddDB is use for add a db
func (a *App) AddDB(db *database.DB) {
    a.db = db
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
// func AddTasker(name string, control interface{}, action string, taskTime string) {
func (a *App) AddTasker(name, taskTime, m string, c controller.Interface) {
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
// app.AddTimer("import_database", 3000, "ImportDatabase", &controller.Get{})
func (a *App) AddTimer(name string, duration time.Duration, m string, c controller.Interface) {
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
func (a *App) AddStaticRoute(prefix, staticDir string) {
    route := &routes.StaticRoute{}
    route.Prefix = prefix
    route.StaticDir = staticDir

    a.staticRoutes = append(a.staticRoutes, route)
}

// AddRoute is use for add a http route
// https://expressjs.com/en/5x/api.html
func (a *App) AddRoute(pattern string, m map[string]string, c controller.Interface) {
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
    t := reflect.Indirect(reflect.ValueOf(c)).Type()
    route := &routes.Route{}
    route.Regex = regex
    route.Methods = m
    route.Params = params
    route.ControllerType = t

    a.routes = append(a.routes, route)
}

// ServeHTTP is a http Handler
// type Handler interface {
//     ServeHTTP(ResponseWriter, *Request)
// }
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    matchRouted := false
    requestPath := r.URL.RawPath
    if requestPath == "" {
        requestPath = r.URL.Path // 会自动 unescape, 先用 RawPath ,不行才用这个
    }
    requestPath = path.Clean(requestPath) // 多个反斜杠变成一个

    // find a matching Route
    for _, staticRoute := range a.staticRoutes {

        // logs.Debug(requestPath, staticRoute.Prefix, staticRoute.StaticDir)
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

            logs.Debug("values", values)

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
        var m  string    
        // Request callback
        if m, ok = route.Methods[r.Method]; !ok {
            http.NotFound(w, r)
        }

        c := route.ControllerType
        a.controllerMethodCall(reflect.Indirect(reflect.ValueOf(c)).Type(), m, w, r, params)
        matchRouted = true
    }

    // if no matches to url, throw a not found exception
    if matchRouted == false {
        http.NotFound(w, r)
    }
}

// Run is to run a app
func Run(app *App) {
    err := http.ListenAndServe(":9090", app)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

// func (a *App) controllerMethodCall(c controller.Interface, m string, w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
func (a *App) controllerMethodCall(t reflect.Type, m string, w http.ResponseWriter, r *http.Request, params map[string]string) (err error) {
    // controllerType := reflect.Indirect(reflect.ValueOf(c)).Type()
    // Invoke the request handler
    vc := reflect.New(t)

    // Init callback
    method := vc.MethodByName("Init")
    contex := &contex.Context{
        ResponseWriter: nil,
        Request:        nil,
        Params:         params,
        DB:             a.db,
    }
    args := make([]reflect.Value, 2)
    args[0] = reflect.ValueOf(contex)
    args[1] = reflect.ValueOf(t.Name())
    method.Call(args)

    args = make([]reflect.Value, 0)

    // Prepare callback
    method = vc.MethodByName("Prepare")
    method.Call(args)

    // Request callback
    method = vc.MethodByName(m)
    if !method.IsValid() {
        // http.NotFound(w, r)
    }
    method.Call(args)

    // Finish callback
    method = vc.MethodByName("Finish")
    method.Call(args)

    return err
}

