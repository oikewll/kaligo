package controller

import (
    "examples/model"
    "fmt"
    "net/http"
    "time"

    "github.com/owner888/kaligo"
)

type Home struct {
    kaligo.Controller
}

// @Summary 初始化接口
// @Tags    Home
// @Success 200 {object} map[string]string
// @Router  /init [GET]
func (c *Home) Initialization() {
    // 都不要带 api
    purviews := "GET-/user,GET-/user/:id,POST-/user,PUT-/user/:id,DELETE-/user/:id"
    menus, err := model.Menu{}.LoadDefault(purviews, true)
    result(c.Context, menus, err)
}

// @Summary 权限列表: 权限选择、权限展示
// @Tags    Home
// @Success 200 {object} map[string]string
// @Router  /permissions [GET]
func (c *Home) Permissions() {
    menus, err := model.Menu{}.LoadDefault("*", true)
    result(c.Context, menus, err)
}

func (c *Home) Index() {
    tplName := c.ParamValue("tplName")
    if tplName == "" {
        tplName = "index"
    }
    data := kaligo.H{
        "title":   "Hello, Kaligo!",
        "year":    fmt.Sprintf("%v", time.Now().Year()),
        "tplName": tplName,
    }
    c.HTML(http.StatusOK, tplName, data)
}
