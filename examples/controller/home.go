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
    menus, err := model.Menu{}.LoadDefault()
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
