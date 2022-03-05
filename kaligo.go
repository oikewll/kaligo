package kaligo

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "path"
    "reflect"
    "regexp"
    // "runtime"
    "strings"
    "sync"
    "time"

    mysql "github.com/owner888/kaligo/database/driver/mysql"
    "github.com/owner888/kaligo/database"
    "github.com/owner888/kaligo/contex"
    "github.com/owner888/kaligo/controller"
    "github.com/owner888/kaligo/config"
    // "github.com/owner888/kaligo/util"

    "github.com/owner888/kaligo/routes"
    "github.com/astaxie/beego/logs"
)

// 定义当前package中使用的全局变量
var (
    db *database.DB // global variable to share it between all HTTP handler
    controlMap map[string]func()
    StaticDir  map[string]string 
    Timer      map[string]*time.Ticker
    Tasker     map[string]*time.Timer
    Mutex      sync.Mutex
)

func init() {
    logs.SetLogFuncCall(true)
    logs.SetLogFuncCallDepth(3)
    if config.Get[bool]("database.mysql.open") {
        var err error    
        // db, err := database.Open(sqlite.Open("./test.db"))
        db, err = database.Open(mysql.Open(config.Get[string]("database.mysql.dsn")))
        if err != nil {
            panic(err)
        }
    }
}

// App is a app
type App struct {
    http.Handler // http.ServeMux
	routes       []*routes.Route
	staticRoutes []*routes.StaticRoute
}

// AddStaticRoute is use for add a static file route
func (a *App) AddStaticRoute(prefix, staticDir string) {
    route := &routes.StaticRoute{}
	route.Prefix    = prefix
	route.StaticDir = staticDir

	a.staticRoutes = append(a.staticRoutes, route)
}


// AddRoute is use for add a http route
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

    // logs.Debug("params", params)
    
	// now create the Route
	t := reflect.Indirect(reflect.ValueOf(c)).Type()
	route := &routes.Route{}
	route.Regex          = regex
	route.Methods        = m
	route.Params         = params
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
        requestPath = r.URL.Path    // 会自动 unescape, 先用 RawPath ,不行才用这个
    }
    requestPath = path.Clean(requestPath)   // 多个反斜杠变成一个

    // find a matching Route
    // for _, staticRoute := range a.staticRoutes {
    //
    //     // logs.Debug(requestPath, staticRoute.Prefix, staticRoute.StaticDir)
    //     // 如果设置的静态文件目录包含在url 路径信息中
    //     if strings.HasPrefix(requestPath, staticRoute.Prefix) {
    //         if len(requestPath) > len(staticRoute.Prefix) && requestPath[len(staticRoute.Prefix)] != '/' {
    //             continue
    //         }
    //         file := path.Join(staticRoute.StaticDir, requestPath[len(staticRoute.Prefix):])
    //         if util.FileExists(file) {
    //             http.ServeFile(w, r, file)
    //         } else {
    //             http.NotFound(w, r)
    //         }
    //         return
    //     }
    // }

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

		// Invoke the request handler
		vc := reflect.New(route.ControllerType)

		// Init callback
		method := vc.MethodByName("Init")
		contex := &contex.Context{
			ResponseWriter: w,
			Request:        r,
			Params:         params,
            DB:             db,
		}
		args := make([]reflect.Value, 2)
		args[0] = reflect.ValueOf(contex)
		args[1] = reflect.ValueOf(route.ControllerType.Name())
		method.Call(args)

		args = make([]reflect.Value, 0)

		// Prepare callback
		method = vc.MethodByName("Prepare")
		method.Call(args)

		// Request callback
		if _, ok := route.Methods[r.Method]; !ok {
			http.NotFound(w, r)
		}
		method = vc.MethodByName(route.Methods[r.Method])
		if !method.IsValid() {
			http.NotFound(w, r)
		}
		method.Call(args)

		// Finish callback
		method = vc.MethodByName("Finish")
		method.Call(args)
		matchRouted = true
	}

	// if no matches to url, throw a not found exception
	if matchRouted == false {
        http.NotFound(w, r)
	}
}

// Run is to run a app
func Run() {
    app := &App{}

    routes.AddRoutes(app)

    // http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    err := http.ListenAndServe(":9090", app) 
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

// Run is the function for start the web service
// func Run() {
    // runtime.GOMAXPROCS(runtime.NumCPU())
    // http.HandleFunc("/", loadController)
    //
    // addr := config.Get("http", "addr")
    // port := config.Get("http", "port")
    //
    // str := util.Colorize(fmt.Sprintf("[I] Running on %s:%s", addr, port), "note")
    // log.Printf(str)
    //log.Printf("[I] Running on %s:%s", addr, port)
    // err := http.ListenAndServe(addr+":"+port, nil)
    // if err != nil {
	//     log.Fatal("ListenAndServe: ", err)
	// }
// }

// AddTasker is the function for add tasker
func AddTasker(name string, control interface{}, action string, taskTime string) {

    llen := len(Tasker)
    if llen == 0 {
        Tasker = make(map[string]*time.Timer)
    }
    then, _ := time.ParseInLocation("2006-01-02 15:04:05", taskTime, time.Local)
    dura := then.Sub(time.Now())
    //fmt.Println(dura)
    if dura > 0 {
        Tasker[name] = time.AfterFunc(dura , func () {
            action = strings.Title(action)
            callMethod(control, action)
        })
    } else {
        fmt.Println("定时任务 --- [ "+name+" ] --- 小于当前时间，将不会被执行")
    }
}

// DelTasker is the function for delete tasker
func DelTasker(name string) bool{
    return Tasker[name].Stop()
}

// AddTimer is the function for add timer, The interval is in microseconds
func AddTimer(name string, control func(), action string, duration time.Duration) {

    llen := len(Timer)
    if llen == 0 {
        Timer = make(map[string]*time.Ticker)
    }
    // 这里如果不用协程，main.go的主进程就会被堵住
    // 用协程，要注意map要用读写锁，不然协程间会抢占导致空指针异常
    go func() {
        Mutex.Lock()
        // 下面这个 defer 是永远不会执行的了，因为这个协程一直都不会退出
        //defer Mutex.Unlock()
        // Timer是一个全局的 map，多个协程写的时候会冲突，所以这里写之前要先Lock，写完Unlock
        Timer[name] = time.NewTicker(duration * time.Millisecond)
        Mutex.Unlock()
        for {
            select {
            case <-Timer[name].C:
            //case <-Timer.Get(name).(*time.Ticker).C:
                action = strings.Title(action)
                callMethod(control, action)
            }
        }
    } ()
}

// DelTimer is the function for delete timer
func DelTimer(name string) bool{
    if timer, ok := Timer[name]; ok {
        timer.Stop()
    } else {
        fmt.Println("定时任务 --- [ "+name+" ] --- 不存在，考虑是否在协程里面生成而且尚未生成，不能执行Stop()")
    }
    //Timer[name].Stop()
    return true
}

// Router is the function for configure dynamic routing
//func Router(ct string, control interface{}) {
func Router(ct string, control func()) {
    if controlMap == nil {
        controlMap = make(map[string]func())
    }
    controlMap[ct] = control
}

// SetStaticPath is the function forsconfigure static file path
func SetStaticPath(key string, value string) {
    if len(StaticDir) == 0 {
        StaticDir = make(map[string]string)
    }
    StaticDir[key] = value
}

// Processing static files
// func staticServer(w http.ResponseWriter, r *http.Request) bool {
//
//     // 如果没有设置静态路径
//     if len(StaticDir) == 0 {
//         return false
//     }
//     // 获取请求路径
//     requestPath := path.Clean(r.URL.Path)
//     for prefix, staticDir := range StaticDir {
//         if len(prefix) == 0 {
//             continue
//         }
//         // 浏览器每次访问都会请求多一次favicon.ico，这里要把它拦截下来,不然后面的代码会执行两次
//         // 如果有数据库写入，就会写入两次，后果非常严重
//         if requestPath == "/favicon.ico" || requestPath == "/robots.txt" {
//             file := path.Join(staticDir, requestPath)
//             if util.FileExists(file) {
//                 http.ServeFile(w, r, file)
//             } else {
//                 http.NotFound(w, r)
//             }
//             return true
//         }
//         // 如果设置的静态文件目录包含在url 路径信息中
//         if strings.HasPrefix(requestPath, prefix) {
//             if len(requestPath) > len(prefix) && requestPath[len(prefix)] != '/' {
//                 continue
//             }
//             file := path.Join(staticDir, requestPath[len(prefix):])
//             if util.FileExists(file) {
//                 http.ServeFile(w, r, file)
//             } else {
//                 http.NotFound(w, r)
//             }
//             return true
//         }
//     }
//     return false
// }

// Load controller
func loadController(w http.ResponseWriter, r *http.Request) { 

    //r.ParseForm()
    //forms := make(map[string]string)
    //for k, v := range r.Form {
        //forms[k] = v[0]
    //}

    defer func() {
        // 捕获异常并处理
        if r := recover(); r != nil {
            str := fmt.Sprintf("%s", r)
            io.WriteString(w, str);
        }
    }()

    // 拦截静态文件
    // if !staticServer(w, r) {
    //     // the default control、action
    //     var ct, ac string
    //     if ct = r.FormValue("ct"); ct == "" {
    //         ct = "index"
    //     }
    //     if ac = r.FormValue("ac"); ac == "" {
    //         ac = "index"
    //     }
    //     // 首字母转大写
    //     ac = strings.Title(ac)
    //     if v, ok := controlMap[ct]; ok {
    //         callMethod(v, ac, w, r)
    //     } else {
    //         panic("Control "+ct+" is not exists!")
    //     }
    // }
}

// 反射调用函数, 注意被调用的object要带 &，表示指针传递
func callMethod(object interface{}, methodName string, args ...interface{}) {
    params := make([]reflect.Value, len(args))
    for i := range args {
        params[i] = reflect.ValueOf(args[i])
    }
    method := reflect.ValueOf(object).MethodByName(methodName)
    if method.IsValid() {
        method.Call(params)
    } else {
        panic("Method "+methodName+" is not exists!")
    }
}

