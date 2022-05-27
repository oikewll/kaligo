package controller

import (
    "fmt"
    "time"
    "github.com/owner888/kaligo"
    httpSwagger "github.com/swaggo/http-swagger"
)

type Docs struct {
    kaligo.Controller
}

func (c *Docs) Index() {
    var swaggerHandler = httpSwagger.Handler(
        httpSwagger.URL(fmt.Sprintf("/docs/swagger.json?ts=%v", time.Now().Unix())),
    )
    swaggerHandler.ServeHTTP(c.ResponseWriter, c.Request)
}
