package main

import (
    "examples/controller"
    "net/http"
    "time"

    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/tpl"
)

//go:generate swag init
// AddRoutes Load Routes
func AddRoutes(r kaligo.Router) {
    addApiRoute(r)
    addHomeRoute(r)
    addStaticRoute(r)
    loadHtmlTemplates(r)
}

func addApiRoute(r kaligo.Router) {
    r.AddRoute("/user/login", map[string]string{
        http.MethodPost: "Login",
    }, &controller.User{})

    r.AddRoute("/api/user/logout", map[string]string{
        http.MethodPost: "Logout",
    }, &controller.User{})

    r.AddRoute("/api/todo", map[string]string{
        http.MethodGet:    "List",
        http.MethodPost:   "Create",
        http.MethodPut:    "Update",
        http.MethodDelete: "Delete",
    }, &controller.Todo{})

    r.AddRoute("api/todo/:id", map[string]string{
        http.MethodGet: "Detail",
    }, &controller.Todo{})
}

func addHomeRoute(r kaligo.Router) {
    r.AddRoute("/", map[string]string{http.MethodGet: "Index"}, &controller.Home{})
    r.AddRoute("/home/:tplName", map[string]string{http.MethodGet: "Index"}, &controller.Home{})
}

func addStaticRoute(r kaligo.Router) {
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
