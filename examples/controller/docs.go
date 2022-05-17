package controller

import (
    "github.com/owner888/kaligo"
    httpSwagger "github.com/swaggo/http-swagger"
)

type Docs struct {
    kaligo.Controller
}

var swaggerHandler = httpSwagger.Handler()

func (c *Docs) Index() {
    swaggerHandler.ServeHTTP(c.ResponseWriter, c.Request)
}
