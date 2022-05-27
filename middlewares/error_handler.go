package middlewares

import (
    "net/http"
    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/logs"
    "github.com/owner888/kaligo/errors"
)

// 错误处理中间件
func ErrorHandler() kaligo.HandlerFunc {
    return func(c *kaligo.Context) {
        c.Next()    // 把下面代码挪到其他中间件 和 Controller 后面跑, 这样才能拿到所有错误

        // Start from here after all middleware and router processing is complete
        // Check for errors in Context.Errors
        for _, err := range c.Errors {
            // log, handle, etc.
            logs.Error("ErrorHandler ==> ", err)

            // Context.Error 转成 error package 里面的错误成功
            if myErr, ok := err.Err.(*errors.Error); ok {
                c.JSON(http.StatusOK, kaligo.H{
                    "code": myErr.Code,
                    "msg":  myErr.Msg,
                    "data": myErr.Data,
                })
            } else {
                // For example, err is set when save session fails
                c.JSON(http.StatusOK, kaligo.H{
                    "code": 500, 
                    "msg": "Server exception", 
                    "data": err.Error(),
                })
            }
            // 不需要了,因为是最后一个了
            // return
        }

        // 不需要了,因为是最后一个了
        // c.Next()
    }
}
