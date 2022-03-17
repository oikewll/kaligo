package routes

import (
    "reflect"
    "regexp"
    // "net/http"

    "github.com/owner888/kaligo/controller"
    "github.com/owner888/kaligo/database"
)

// App is use for add Route struct and StaticRoute struct
type App interface {
    AddRoute(pattern string, m map[string]string, c controller.Interface)
    AddStaticRoute(prefix, staticDir string)
    AddDB(db *database.DB)
}

// Route is a app route
type Route struct {
    Regex          *regexp.Regexp
    Methods        map[string]string
    Params         map[int]string
    ControllerType reflect.Type
    // Middlewares returns the list of middlewares in use by the router.
    // Middlewares() Middlewares
}

// StaticRoute is a app route
type StaticRoute struct {
    Prefix    string
    StaticDir string
}

// Middlewares type is a slice of standard middleware handlers with methods
// to compose middleware chains and http.Handler's.
// type Middlewares []func(http.Handler) http.Handler
