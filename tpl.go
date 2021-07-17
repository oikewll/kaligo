package kaligo

import (
    "fmt"
    "html/template"
    "io"
    "net/http"
)
type Tpl struct {
    args map[string]interface{} 
    Response http.ResponseWriter
}

func (t *Tpl) Assign(tplVar string, value interface{}) bool {

    if len(t.args) == 0 {
        t.args = make(map[string]interface{})
    }
    t.args[tplVar] = value
    return true
}


func (t *Tpl) Display(tpl string) bool {

    var htp, err = template.ParseFiles("template/" + tpl)
    if err != nil {
        io.WriteString(t.Response, fmt.Sprintf("%s", err))
        return false
    }
    htp.Execute(t.Response, t.args)

    return true
}
