package controller

import (
    "fmt"
    "net/http"
    "time"

    "github.com/owner888/kaligo"
)

type Meta struct {
    Name string
    Icon string
}
type Menu struct {
    ID int
    Name string
    Path string
    Show bool

    Meta
}

type Menus struct {
    Current  Menu
    Children Menu
}

type Home struct {
    kaligo.Controller
}

// @Summary 初始化接口
// @Tags    Home
// @Success 200 {object} map[string]string
// @Router  /init [GET]
func (c *Home) Initialization() {
    // {
    //     id: "0-5",
    //     name: "acquisition",
    //     path: "/acquisition",
    //     meta: {
    //         title: "menu.acquisition-settings",
    //         icon: "el-icon-truck",
    //     },
    //     show: true,
    //     childrens: [
    //     ]
    // },
    c.String(200, "Initialization")
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
