package middlewares

import (
    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/logs"
)

func Log(ctx *kaligo.Context) {
    if req := ctx.Request; req != nil {
        logs.Debug(req.Method, " ", req.URL)
    }
    ctx.Next()
}
