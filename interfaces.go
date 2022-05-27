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

// type AuthUser interface {
//     Name() string
//     CheckUser(account, password string, remember ...bool) (err error)
//     SaveUserSession(account, password string, remember ...bool) (err error)
//     PasswordHash(password string) (string, error)
//     PasswordVerify(password, hash string) bool
//     checkAccount(account string) (err error)
//     checkLoginError24(account string) (err error)
//     checkUserStatus(status int32) (err error)
// }

/* vim: set expandtab: */
