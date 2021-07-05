package kaligo

import (
    "net/http"
    "html/template"
    "io"
    "fmt"
)
type Tpl struct {
    args map[string]interface{} 
    Response http.ResponseWriter
}

func (this *Tpl) Assign(tplVar string, value interface{}) bool {

    if len(this.args) == 0 {
        this.args = make(map[string]interface{})
    }
    this.args[tplVar] = value
    return true
}


func (this *Tpl) Display(tpl string) bool {

    t, err := template.ParseFiles("template/"+tpl)
    if err != nil {
        io.WriteString(this.Response, fmt.Sprintf("%s", err));
        return false
    }
    t.Execute(this.Response, this.args)

    return true
}
