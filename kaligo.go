package kaligo

import (
    "net/http"

    "github.com/owner888/kaligo/controller"
)

// Router 是一套路由接口，Mux 实现了此接口
type Router interface {
    http.Handler

    AddRoute(pattern string, m map[string]string, c controller.Interface)

    AddStaticRoute(prefix, staticDir string)

    Use(middlewares ...func(http.Handler) http.Handler)

    With(middlewares ...func(http.Handler) http.Handler) Router
}
