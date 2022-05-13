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
        "title": "Hello, Kaligo!",
        "year":  fmt.Sprintf("%v", time.Now().Year()),
    }
    c.HTML(http.StatusOK, tplName, data)
}
