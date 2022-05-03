package kaligo

import (
    "net/http"
)

// Router 是一套路由接口，Mux 实现了此接口
type Router interface {
    http.Handler

    AddRoute(pattern string, m map[string]string, c Interface)

    AddStaticRoute(prefix, staticDir string)

    Use(middlewares ...HandlerFunc)

    With(middlewares ...HandlerFunc) Router
}

/* vim: set expandtab: */
