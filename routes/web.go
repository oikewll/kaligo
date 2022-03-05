package routes

import (
    "reflect"
    "regexp"
    "github.com/owner888/chrome/controller"
)

type app interface {
    AddRoute(pattern string, m map[string]string, c controller.Interface)
    AddStaticRoute(prefix, staticDir string)
}

// Route is a app route
type Route struct {
    Regex          *regexp.Regexp
    Methods        map[string]string
    Params         map[int]string
    ControllerType reflect.Type
}

// StaticRoute is a app route
type StaticRoute struct {
    Prefix    string
    StaticDir string
}

// AddRoutes is use for add Route
// https://expressjs.com/en/5x/api.html
func AddRoutes(a app) {
    a.AddStaticRoute("/static", "static")

    a.AddRoute("/statistics/:channel(.*)/:apkurl(.*)", map[string]string{
        "GET": "Statistics",
    }, &controller.Get{})

    a.AddRoute("/channel/:ip(.*)", map[string]string{
        "GET": "GetChannel",
    }, &controller.Get{})

    a.AddRoute("/", map[string]string{
        "GET": "Index",
    }, &controller.Get{})

    a.AddRoute("/posts/:post_id([0-9]+)", map[string]string{
        "POST": "Show",
    }, &controller.Post{})
}
