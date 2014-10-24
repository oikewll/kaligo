package epooll

import (
    "net/http"
    "fmt"
    "log"
    "io"
    "strings"
    "reflect"
    "runtime"
    "path"
    "time"
    "github.com/owner888/epooll/util"
)

// 定义当前package中使用的全局变量
var (
    controlMapping map[string]interface{} 
    StaticDir map[string]string 
    Timer map[string]*time.Ticker
    Tasker map[string]*time.Timer
)

// 设置定时任务
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
// 删除定时任务
func DelTasker(name string) bool{
    return Tasker[name].Stop()
}

// 设置闹钟，间隔单位为微秒
func AddTimer(name string, control interface{}, action string, duration time.Duration) {

    llen := len(Timer)
    if llen == 0 {
        Timer = make(map[string]*time.Ticker)
    }
    //fmt.Println(Timer)
    go func() {
        Timer[name] = time.NewTicker(duration * time.Millisecond)
        for {
            select {
            case <-Timer[name].C:
                action = strings.Title(action)
                callMethod(control, action)
            }
        }
    } ()
}
// 删除闹钟
func DelTimer(name string) bool{
    Timer[name].Stop()
    return true
}

// 启动Web服务
func Run() {
    runtime.GOMAXPROCS(runtime.NumCPU());
    http.HandleFunc("/", loadController)
    conf := InitConfig()
    addr := conf.GetValue("http", "addr") 
    port := conf.GetValue("http", "port") 
    if addr == "no value" {
        addr = ""
    }
    if port == "no value" {
        port = "9527"
    }

    log.Printf("[I] Running on %s:%s", addr, port)
    http.ListenAndServe(addr+":"+port, nil)
}

// 配置动态路由
func Router(ct string, control interface{}) {
    if controlMapping == nil {
        controlMapping = make(map[string]interface{})
    }
    controlMapping[ct] = control
}

// 配置静态路由
func SetStaticPath(key string, value string) {
    StaticDir = make(map[string]string)
    StaticDir[key] = value
}

// 处理静态文件
func staticServer(w http.ResponseWriter, r *http.Request) {

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
                return
            } else {
                http.NotFound(w, r)
                return
            }
        }
        // 如果设置的静态文件目录包含在url 路径信息中
        if strings.HasPrefix(requestPath, prefix) {
            if len(requestPath) > len(prefix) && requestPath[len(prefix)] != '/' {
                continue
            }
            file := path.Join(staticDir, requestPath[len(prefix):])
            fmt.Println(file)
            if util.FileExists(file) {
                http.ServeFile(w, r, file)
                return
            } else {
                http.NotFound(w, r)
                return
            }
        }
    }
}

// 载入控制器
func loadController(w http.ResponseWriter, r *http.Request) { 

    // 拦截静态文件
    staticServer(w, r)

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
    if v, ok := controlMapping[ct]; ok {
        callMethod(v, ac, w, r)
    } else {
        panic("Control "+ct+" is not exists!")
    }
}

// 反射调用函数, 注意被调用的object要带 &，表示指针传递
func callMethod(object interface{}, methodName string, args ...interface{}) {
    params := make([]reflect.Value, len(args))
    for i, _ := range args {
        params[i] = reflect.ValueOf(args[i])
    }
    method := reflect.ValueOf(object).MethodByName(methodName)
    if method.IsValid() {
        method.Call(params)
    } else {
        panic("Method "+methodName+" is not exists!")
    }
}

