package kaligo

import (
    "reflect"
    "regexp"
    // "net/http"
)

// Route is a Router's route
type Route struct {
    Regex          *regexp.Regexp
    Methods        map[string]string
    Params         map[int]string
    ControllerType reflect.Type
    // Middlewares returns the list of middlewares in use by the router.
    // Middlewares() Middlewares
}

// StaticRoute is a Router's route
type StaticRoute struct {
    Prefix    string
    StaticDir string
}

// Middlewares type is a slice of standard middleware handlers with methods
// to compose middleware chains and http.Handler's.
// type Middlewares []func(http.Handler) http.Handler
