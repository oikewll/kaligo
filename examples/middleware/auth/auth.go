package auth

import (
    "net/http"

    "examples/model"
    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/logs"
    "github.com/owner888/kaligo/sessions"
)

// AuthUserKey is the cookie name for user credential in basic auth.
const AuthUserKey = "user"

func Auth(c *kaligo.Context) {
    // if !GetSession(c).IsValid() {
    //     c.JSON(http.StatusUnauthorized, kaligo.H{"code": -1, "msg": "not authorized"})
    //     c.Abort()
    //     return
    // }

    var req *http.Request
    if req = c.Request; req != nil {
        logs.Debug(req.Method, " ", req.URL)
    }

    session := sessions.Default(c)
    c.UID = session.Get(AuthUserKey).(string)

    u := model.DefaultUser(c)
    // 检查权限
    if !u.CheckPurviews(req.Method, req.URL.Path) {
        c.JSON(http.StatusUnauthorized, kaligo.H{"code": -1, "msg": "not authorized"})
        c.Abort()
    }

    c.Next()
}

type Session struct {
    UserID int64
    Role   string
    Token  string
}

func GetSession(c *kaligo.Context) *Session {
    token := c.CookieValue("access_token")
    return &Session{Token: token}
}

func (s *Session) IsValid() bool {
    return true
}
