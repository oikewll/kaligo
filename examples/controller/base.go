package controller

import (
    "net/http"

    "github.com/owner888/kaligo"
)

func result(ctx *kaligo.Context, data any, err error) {
    if err != nil {
        ctx.JSON(http.StatusBadRequest, kaligo.H{
            "code":    -1,
            "message": err.Error(),
        })
    } else if data != nil {
        ctx.JSON(http.StatusOK, kaligo.H{
            "code": 0,
            "data": data,
        })
    }
}
