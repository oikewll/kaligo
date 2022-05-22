package controller

import (
    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/sessions"
) 

type Sessions struct {
    kaligo.Controller
}

func (c *Sessions) Hello() {
    // 初始化session对象
    session := sessions.Default(c)

    if session.Get("hello") != "world" {
        session.Set("hello", "world")
        session.Delete("tizi365")
        session.Save()
        // session.Clear()
    }

    c.JSON(200, kaligo.H{"hello": session.Get("hello")})
}
