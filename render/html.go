package render

import (
    "html/template"
    "net/http"

    "github.com/owner888/kaligo/tpl"
)

type HTMLRender interface {
    Instance(string, interface{}) Render
}

type HTML struct {
    Tpl      *tpl.Tpl
    Template *template.Template
    Name     string
    Data     interface{}
}

var htmlContentType = []string{"text/html; charset=utf-8"}

func (r HTML) Instance(name string, data interface{}) Render {
    r.Name = name
    r.Data = data
    return r
}

func (r HTML) Render(w http.ResponseWriter) error {
    r.WriteContentType(w)
    return r.Tpl.Render(w, r.Name, r.Data)
}

func (r HTML) WriteContentType(w http.ResponseWriter) {
    writeContentType(w, htmlContentType)
}
