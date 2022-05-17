package controller

import (
    "net/http"
    "time"

    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/tpl"
)

func AddRoutes(r kaligo.Router) {
    r.AddRoute("/user/login", map[string]string{
        http.MethodPost: "Login",
    }, &User{})

    r.AddRoute("/api/user/logout", map[string]string{
        http.MethodPost: "Logout",
    }, &User{})

    r.AddRoute("/api/todo", map[string]string{
        http.MethodGet:    "List",
        http.MethodPost:   "Create",
        http.MethodPut:    "Update",
        http.MethodDelete: "Delete",
    }, &Todo{})

    r.AddRoute("api/todo/:id", map[string]string{
        http.MethodGet: "Detail",
    }, &Todo{})

    r.AddRoute("/", map[string]string{http.MethodGet: "Index"}, &Home{})
    r.AddRoute("/home/:tplName", map[string]string{http.MethodGet: "Index"}, &Home{})

    loadHtmlTemplates(r)
    AddStaticRoute(r)
}

func AddStaticRoute(r kaligo.Router) {
    r.AddStaticRoute("/static", webRootPath()+"/static")
    r.AddStaticRoute("/favicon.ico", webRootPath()+"/favicon.ico")
}

var Tpls *tpl.Tpl

func loadHtmlTemplates(r kaligo.Router) {
    dir := templatePath() + "/default"
    t, err := r.SetHTMLTemplate(dir, ".html", time.Second*30)
    if err != nil {
        panic("load html template error:" + err.Error())
    }
    Tpls = t
}

func webRootPath() string { return "wwwroot" }

func templatePath() string { return webRootPath() + "/template" }