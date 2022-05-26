package model

import (
    // "fmt"
    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/logs"
    "github.com/owner888/kaligo/util"
)

const (
    AuthDefaultKeyFormat = "AuthDefaultKey-%s"       // 根据UID获取用户信息
)

type Auth struct {
    ctx  *kaligo.Context
}

func DefaultAuth(c *kaligo.Context) *Auth {
    return &Auth{ctx: c}
}

func (m *Auth) MakeCsrfToken() string {
    csrf := util.GenerateToken()
    // m.ctx.Set(fmt.Sprintf(AuthDefaultKeyFormat, csrf), csrf)     // Context.Set() 不能用于不同请求之间共享数据, 只能用于中间件之间共享
    m.ctx.SetCookie("csrf_token", csrf, 1000, "/", "", true, true)
    logs.Debug("Csrf Cookie Set ==> ", csrf)

    return csrf
}

func (m *Auth) CheckCsrfToken(csrf string) bool {
    // ret := m.ctx.GetString(fmt.Sprintf(AuthDefaultKeyFormat, csrf))
    ret := m.ctx.CookieValue("csrf_token")
    // 清除 cookie
    m.ctx.SetCookie("csrf_token", "", 0, "/", "", true, true)

    // logs.Debug("Csrf ==> ", csrf)
    // logs.Debug("Csrf Cookie Get ==> ", ret)

    return csrf == ret
}
