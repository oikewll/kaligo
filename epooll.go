package epooll

import (
    //"github.com/owner888/epooll"
    "github.com/owner888/epooll/util"
    "net/http"
    "fmt"
    "io"
    "strings"
)

var (
    controlMap map[string]interface{} 
)

func Run() {

    http.HandleFunc("/", loadController)
    http.ListenAndServe(":9527", nil)
}

func Router(ct string, control interface{}) {
    if controlMap == nil {
        controlMap = make(map[string]interface{})
    }
    controlMap[ct] = control   
}

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

    // 处理控制器和方法
    var ct, ac string
    if ct = r.FormValue("ct"); ct == "" {
        ct = "index"
    }
    if ac = r.FormValue("ac"); ac == "" {
        ac = "index"
    }
    ac = strings.Title(ac)
    //imap := map[string]interface{}{
        //"index": &control.Index{},
        //"home": &control.Home{},
        //"user": &control.User{},
    //}
    if v, ok := controlMap[ct]; ok {
        util.Invoke(v, ac, w, r)
    } else {
        panic("Control "+ct+" is not exists!")
    }
}

