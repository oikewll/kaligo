package controller

import (
    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/sessions"
)

type Sessions struct {
    kaligo.Controller
}

<<<<<<< HEAD
// Swagger 文档地址: https://github.com/swaggo/swag#api-operation
=======
// @Summary Session 样例 Set
// @tags    session
// @Success 200 {object} map[string]string
// @Router  /sessions [post]
func (c *Sessions) Set() {
    // 初始化session对象
    session := sessions.Default(c.Context)
>>>>>>> 643c8173ad82ecec031245d66cf6777830d0ccd7

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

<<<<<<< HEAD
// @Summary Session 信息
=======
// @Summary Session 样例 Get
// @tags    session
>>>>>>> 643c8173ad82ecec031245d66cf6777830d0ccd7
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
