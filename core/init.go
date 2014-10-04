package core

import (
	"net/http"
    "epooll/control"
    "epooll/util"
    "strings"
    "io"
    "fmt"
)

type Init struct {
}

func (this *Init) InitServer() bool{
    http.HandleFunc("/", this.loadController)
    http.ListenAndServe(":9527", nil)
    return true
}

func (this *Init) loadController(w http.ResponseWriter, r *http.Request) { 

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
    imap := map[string]interface{}{
        "index": &control.Index{},
        "home": &control.Home{},
        "user": &control.User{},
    }
    if v, ok := imap[ct]; ok {
        util.Invoke(v, ac, w, r)
    } else {
        panic("Control "+ct+" is not exists!")
    }
}

