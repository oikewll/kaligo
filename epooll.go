package epooll

import (
    "net/http"
    "fmt"
    "io"
    "strings"
    "reflect"
    "runtime"
    "path"
    "github.com/owner888/epooll/util"
)

var (
    controlMapping map[string]interface{} 
    StaticDir map[string]string 
)

func Run() {
    runtime.GOMAXPROCS(runtime.NumCPU());
    http.HandleFunc("/", loadController)
    //http.Handle("/static/", http.FileServer(http.Dir("./")))
    //http.HandleFunc("/static/", staticServer)
    http.ListenAndServe(":9527", nil)
}

func Router(ct string, control interface{}) {
    if controlMapping == nil {
        controlMapping = make(map[string]interface{})
    }
    controlMapping[ct] = control   
}

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

// 反射调用函数, 注意被调用的类要带 &，表示对象
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

