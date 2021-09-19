package kaligo

import (
    "fmt"
    "github.com/owner888/kaligo/config"
    "github.com/owner888/kaligo/util"
    "io"
    "log"
    "net/http"
    "path"
    "reflect"
    "runtime"
    "strings"
    "sync"
    "time"
)

// 定义当前package中使用的全局变量
var (
    controlMap map[string]func()
    StaticDir  map[string]string 
    Timer      map[string]*time.Ticker
    Tasker     map[string]*time.Timer
    Mutex      sync.Mutex
)

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

// Run is the function for start the web service
func Run() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    http.HandleFunc("/", loadController)

    addr := config.Get("http", "addr") 
    port := config.Get("http", "port") 

    str := util.Colorize(fmt.Sprintf("[I] Running on %s:%s", addr, port), "note")
    log.Printf(str)
    //log.Printf("[I] Running on %s:%s", addr, port)
    err := http.ListenAndServe(addr+":"+port, nil)
    if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
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
func staticServer(w http.ResponseWriter, r *http.Request) bool {

    // 如果没有设置静态路径
    if len(StaticDir) == 0 {
        return false
    }
    // 获取请求路径
    requestPath := path.Clean(r.URL.Path)
    for prefix, staticDir := range StaticDir {
		if len(prefix) == 0 {
			continue
		}
        // 浏览器每次访问都会请求多一次favicon.ico，这里要把它拦截下来,不然后面的代码会执行两次
        // 如果有数据库写入，就会写入两次，后果非常严重
        if requestPath == "/favicon.ico" || requestPath == "/robots.txt" {
            file := path.Join(staticDir, requestPath)
            if util.FileExists(file) {
                http.ServeFile(w, r, file)
            } else {
                http.NotFound(w, r)
            }
            return true
        }
        // 如果设置的静态文件目录包含在url 路径信息中
        if strings.HasPrefix(requestPath, prefix) {
            if len(requestPath) > len(prefix) && requestPath[len(prefix)] != '/' {
                continue
            }
            file := path.Join(staticDir, requestPath[len(prefix):])
            if util.FileExists(file) {
                http.ServeFile(w, r, file)
            } else {
                http.NotFound(w, r)
            }
            return true
        }
    }
    return false
}

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
    if !staticServer(w, r) {
        // the default control、action
        var ct, ac string
        if ct = r.FormValue("ct"); ct == "" {
            ct = "index"
        }
        if ac = r.FormValue("ac"); ac == "" {
            ac = "index"
        }
        // 首字母转大写
        ac = strings.Title(ac)
        if v, ok := controlMap[ct]; ok {
            callMethod(v, ac, w, r)
        } else {
            panic("Control "+ct+" is not exists!")
        }
    }
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

