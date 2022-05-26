package middlewares

import (
    "github.com/owner888/kaligo"
)

// 支持 CORS 跨域访问
func CORS() kaligo.HandlerFunc {
    return func(c *kaligo.Context) {
        c.ResponseWriter.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
        c.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET,PUT, DELETE, UPDATE")
        c.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        c.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(200)
            return
        }
        c.Next()
    }
}
