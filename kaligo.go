package kaligo

import (
    "net/http"
    "time"

    "github.com/owner888/kaligo/tpl"
)

// Router 是一套路由接口，Mux 实现了此接口
type Router interface {
    http.Handler

    AddRoute(pattern string, m map[string]string, c Interface)

    AddStaticRoute(prefix, staticDir string)

    SetHTMLTemplate(dir string, ext string, reloadTime time.Duration) (*tpl.Tpl, error)

    Use(middlewares ...HandlerFunc)

    With(middlewares ...HandlerFunc) Router
}

/* vim: set expandtab: */
