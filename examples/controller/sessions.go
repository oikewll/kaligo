package controller

import (
    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/sessions"
) 

type Sessions struct {
    kaligo.Controller
}

// Hello Session 样例
// @Summary Session 样例
// @Param   session_id   query integer false "SESSION ID"
// @Success 200 {object} map[string]string
// @Router  /sessions [get]
func (c *Sessions) Hello() {
    // 初始化session对象
    session := sessions.Default(c.Context)

    if session.Get("hello") != "world" {
        session.Set("hello", "world")
        session.Delete("tizi365")
        session.Save()
        // session.Clear()
    }

    c.JSON(200, kaligo.H{"hello": session.Get("hello")})
}
