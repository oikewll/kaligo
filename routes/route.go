package routes

import (
    "reflect"
    "regexp"

    "github.com/owner888/kaligo/controller"
)

// App is use for add Route struct and StaticRoute struct
type App interface {
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
