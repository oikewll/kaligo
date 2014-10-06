package epooll

import (
    "net/http"
    "fmt"
    "io"
    "strings"
    "reflect"
    "runtime"
)

var (
    controlMap map[string]interface{} 
)

func Run() {
    runtime.GOMAXPROCS(runtime.NumCPU());
    http.HandleFunc("/", loadController)
    //http.Handle("/static/", http.FileServer(http.Dir("./")))
    http.HandleFunc("/static/", staticServer)
    http.ListenAndServe(":9527", nil)
}

func Router(ct string, control interface{}) {
    if controlMap == nil {
        controlMap = make(map[string]interface{})
    }
    controlMap[ct] = control   
}

// 处理静态文件
func staticServer(w http.ResponseWriter, r *http.Request) {
    // 千万不要设置到 ./static ，直接设置 ./ 就好
    staticHandler := http.FileServer(http.Dir("./"))
    staticHandler.ServeHTTP(w, r)
    return
}

func loadController(w http.ResponseWriter, r *http.Request) { 

    path := r.URL.Path[1:]
    // 浏览器每次访问都会请求多一次favicon.ico，这里要把它拦截下来,不然后面的代码会执行两次
    // 如果有数据库写入，就会写入两次，后果非常严重
    if path == "favicon.ico"  {
        http.NotFound(w, r)
        return
    }

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

    // 处理控制器和方法
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

// 反射调用函数, 注意被调用的类要带 &
// util.Invoke(&control.Home{}, ct, w, r)
func callMethod(any interface{}, name string, args ...interface{}) {
    inputs := make([]reflect.Value, len(args))
    for i, _ := range args {
        inputs[i] = reflect.ValueOf(args[i])
    }
    method := reflect.ValueOf(any).MethodByName(name)
    if method.IsValid() {
        method.Call(inputs)
    } else {
        panic("Method "+name+" is not exists!")
    }
}

