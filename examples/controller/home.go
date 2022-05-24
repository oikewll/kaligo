package controller

import (
    "fmt"
    "net/http"
    "time"

    "github.com/owner888/kaligo"
)

type Home struct {
    kaligo.Controller
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

// @Summary 账户登陆
// @Tags    Home
// @Success 200 {object} map[string]string
// @Router  /home [POST]
func (c *Home) Login() {
    username := c.FormValue("username")
    c.SetCookie("access_token", username, 1000, "/", "", true, true)
}

func (c *Home) Logout() {

}
