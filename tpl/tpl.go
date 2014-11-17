package tpl

import (
    "net/http"
    "html/template"
)

var (
    args map[string]interface{} 
    tpls map[string]string 
)

func Assign(key string, val interface{}) bool {

    if len(args) == 0 {
        args = make(map[string]interface{})
    }
    args[key] = val
    return true
}

func Layout(key string, val string) bool {

    if len(tpls) == 0 {
        tpls = make(map[string]string)
    }
    tpls[key] = val
    return true
}


func Display(w http.ResponseWriter, tpl string) error {

    t, err := template.ParseFiles("template/"+tpl)
    if err != nil {
        return err
    }

    for _, v := range tpls {
        t, err = t.ParseGlob("template/"+v)
        if err != nil {
            return err
        }
    }

    //t.ExecuteTemplate(w, tpl, args)
    t.Execute(w, args)
    return err
}
