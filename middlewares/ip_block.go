package middlewares

import (
    "github.com/owner888/kaligo"
    "golang.org/x/exp/slices"
)

func IpWhiteList(ips []string) kaligo.HandlerFunc {
    return func(ctx *kaligo.Context) {
        ip := ctx.ClientIP()
        if len(ip) > 0 && !slices.Contains(ips, ip) {
            ctx.Abort()
            return
        }
        ctx.Next()
    }
}

func IpBlackList(ips []string) kaligo.HandlerFunc {
    return func(ctx *kaligo.Context) {
        ip := ctx.ClientIP()
        if len(ip) > 0 && slices.Contains(ips, ctx.ClientIP()) {
            ctx.Abort()
            return
        }
        ctx.Next()
    }
}
