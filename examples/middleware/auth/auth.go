package auth

import (
    "net/http"

    "examples/model"
    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/logs"
    "github.com/owner888/kaligo/util"
    "github.com/owner888/kaligo/sessions"
)

// AuthUserKey is the cookie name for user credential in basic auth.
const AuthUserKey = "user"

// func Sessions(name string, store Store) kaligo.HandlerFunc {
//     return func(c *kaligo.Context) {
//         s := &session{name, c.Request, store, nil, false, c.ResponseWriter}
//         c.Set(DefaultKey, s)
//         defer context.Clear(c.Request)
//         c.Next()
//     }
// }

func Auth() kaligo.HandlerFunc {
    return func(c *kaligo.Context) {
        var req *http.Request
        if req = c.Request; req != nil {
            logs.Debug(req.Method, " ", req.URL)
        }

        // 因为 session 使用的是 context.Set() 方法, 要检查服务重启是否有效问题
        session := sessions.Default(c)
        c.UID = util.ToString(session.Get(AuthUserKey))

        u := model.DefaultUser(c)
        // 检查权限
        if !u.CheckPurviews(req.Method, req.URL.Path) {
            c.JSON(http.StatusUnauthorized, kaligo.H{"code": -1, "msg": "not authorized"})
            c.Abort()
        }

        c.Next()
    }
}
