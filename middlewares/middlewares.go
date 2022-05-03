package middlewares

import "github.com/owner888/kaligo"

func init() {
    kaligo.DefaultHandlers = append(kaligo.DefaultHandlers, Log)
}
