package model

import (
    "fmt"
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
    csrf := util.MakeCsrfToken()
    key  := fmt.Sprintf(AuthDefaultKeyFormat, csrf)
    logs.Debug("Key ==> ", key)
    logs.Debug("Csrf ==> ", csrf)
    m.ctx.Set(key, csrf)

    return csrf
}

func (m *Auth) CheckCsrfToken(csrf string) bool {
    key := fmt.Sprintf(AuthDefaultKeyFormat, csrf)
    ret := m.ctx.GetString(key)
    logs.Debug("Csrf ==> ", csrf)
    logs.Debug("Key ==> ", key)
    logs.Debug("Context.Csrf ==> ", ret)
    return true
}
