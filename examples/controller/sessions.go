package controller

import (
    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/sessions"
)

type Sessions struct {
    kaligo.Controller
}

// Swagger 文档地址: https://github.com/swaggo/swag#api-operation

// @Summary Session 添加
// @Description Session 添加简介
// @Success 200 {object} map[string]string
// @Router  /sessions [POST]
func (c *Sessions) Create() {
    session := sessions.Default(c.Context)
    session.Set("hello", "world")
    session.Save()

    c.JSON(200, kaligo.H{"hello": session.Get("hello")})
}

// @Summary Session 信息
// @Success 200 {object} map[string]string
// @Router  /sessions [GET]
func (c *Sessions) Detail() {
    session := sessions.Default(c.Context)

    c.JSON(200, kaligo.H{"hello": session.Get("hello")})
}

// @Summary Session 删除
// @Success 200 {object} map[string]string
// @Router  /sessions [DELETE]
func (c *Sessions) Delete() {
    session := sessions.Default(c.Context)
    session.Delete("hello")
    session.Save()

    c.JSON(200, kaligo.H{"hello": session.Get("hello")})
}

// @Summary Session 销毁
// @Success 200 {object} map[string]string
// @Router  /sessions/destory [DELETE]
func (c *Sessions) Destory() {
    session := sessions.Default(c.Context)
    session.Clear()
    session.Save()

    c.JSON(200, kaligo.H{"hello": session.Get("hello")})
}
