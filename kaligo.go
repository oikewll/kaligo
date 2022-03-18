package kaligo

import (
    "net/http"

    "github.com/owner888/kaligo/routes"
)

// Router 是一套路由接口，Mux 实现了此接口
type Router interface {
    http.Handler
    routes.Router

    Use(middlewares ...func(http.Handler) http.Handler)

    With(middlewares ...func(http.Handler) http.Handler) Router
}
