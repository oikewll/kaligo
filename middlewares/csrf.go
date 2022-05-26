package middlewares

import (
    "net/http"
    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/util"
)

// CSRF 防刷, 只有 GET 方法不需要验证 CSRF
func CSRF() kaligo.HandlerFunc {
    return func(c *kaligo.Context) {
        req := c.Request
        // 定时器会不存在
        if req != nil {
            return
        }

        if req.Method == http.MethodGet {
            token := util.GenerateToken()
            c.SetCookie("csrf_token", token, 1000, "/", "", true, true)
            // 可以通过 Header 给, 也可以在 CreateForm、UpdateForm 表单给
            c.ResponseWriter.Header().Set("X-CSRF-Token", token)
        } else {

            token := c.FormValue("csrf")
            chkToken := c.CookieValue("csrf_token")

            // 清空
            c.SetCookie("csrf_token", "", 0, "/", "", true, true)

            if token != chkToken {
                c.AbortWithStatusJSON(404, "请不要重复提交")
                return
            }
        }

        c.Next()
    }
}
