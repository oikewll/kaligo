package main

import (
    "examples/controller"
    "net/http"
    "time"

    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/tpl"
)

// AddRoutes Load Routes
func AddRoutes(r kaligo.Router) {
    addSessionRoute(r)
    addApiRoute(r)
    addHomeRoute(r)
    addStaticRoute(r)
    addDocsRoute(r)
    loadHtmlTemplates(r)
}

//go:generate swag init
func addSessionRoute(r kaligo.Router) {
    r.AddRoute("/api/sessions", map[string]string{
        http.MethodGet:     "Detail",
        http.MethodPost:    "Create",
        http.MethodDelete:  "Delete",
    }, &controller.Sessions{})

    r.AddRoute("/api/sessions/destory", map[string]string{
        http.MethodDelete:  "Destory",
    }, &controller.Sessions{})
}

func addApiRoute(r kaligo.Router) {
    // 优先匹配精准的, 否则和 /api/user/:id 冲突了
    r.AddRoute("/api/auth/login", map[string]string{
        http.MethodPost: "Login",
    }, &controller.Auth{})

	r.AddRoute("/api/auth/logout", map[string]string{
		http.MethodDelete: "Logout",
	}, &controller.Auth{})

	r.AddRoute("/api/auth/token", map[string]string{
		http.MethodGet: "Token",
	}, &controller.Auth{})

    r.AddRoute("/api/auth/check_token", map[string]string{
		http.MethodPost: "CheckToken",
	}, &controller.Auth{})

    r.AddRoute("/api/user/:id", map[string]string{
		http.MethodPut:    "Update",
		http.MethodDelete: "Delete",
        http.MethodGet:    "Detail",
    }, &controller.User{})

	r.AddRoute("/api/user/createform", map[string]string{
		http.MethodGet:    "CreateForm",
	}, &controller.User{})

	r.AddRoute("/api/user", map[string]string{
		http.MethodGet:    "List",
		http.MethodPost:   "Create",
	}, &controller.User{})

	r.AddRoute("/api/todo/:id", map[string]string{
		http.MethodPut:    "Update",
		http.MethodDelete: "Delete",
		http.MethodGet:    "Detail",
	}, &controller.Todo{})
    
	r.AddRoute("/api/todo", map[string]string{
		http.MethodGet:    "List",
		http.MethodPost:   "Create",
	}, &controller.Todo{})
}

func addHomeRoute(r kaligo.Router) {
    r.AddRoute("/api/init", map[string]string{http.MethodGet: "Initialization"}, &controller.Home{})

    r.AddRoute("/", map[string]string{http.MethodGet: "Index"}, &controller.Home{})
    r.AddRoute("/home/:tplName", map[string]string{http.MethodGet: "Index"}, &controller.Home{})
}

func addStaticRoute(r kaligo.Router) {
    r.AddStaticRoute("/admin", webRootPath()+"/admin")
    r.AddStaticRoute("/static", webRootPath()+"/static")
    r.AddStaticRoute("/favicon.ico", webRootPath()+"/favicon.ico")
}

func addDocsRoute(r kaligo.Router) {
    r.AddStaticRoute("/docs", "docs")
    r.AddRoute("/swagger/.*", map[string]string{http.MethodGet: "Index"}, &controller.Docs{})
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
