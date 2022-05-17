package auth

import (
    "net/http"

    "github.com/owner888/kaligo"
)

func Auth(c *kaligo.Context) {
    if !GetSession(c).IsValid() {
        c.JSON(http.StatusUnauthorized, kaligo.H{"code": -1, "msg": "not authorized"})
        c.Abort()
        return
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
