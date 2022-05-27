package main

import (
    "examples/middleware/auth"
    "examples/model"
    "flag"
    "fmt"

    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/config"
    "github.com/owner888/kaligo/database"
    "github.com/owner888/kaligo/database/driver/mysql"
    "github.com/owner888/kaligo/logs"
    "github.com/owner888/kaligo/sessions"
    "github.com/owner888/kaligo/sessions/cookie"
    "github.com/owner888/kaligo/middlewares"
)

// swag 不要集成到项目,直接安装工具包即可
// go install github.com/swaggo/swag/cmd/swag@latest
// 文档: https://github.com/swaggo/swag#api-operation
// go:generate swag init

// @title Kaligo Example API
// @description 
// @version 1.0
// @host localhost:8080
// @BasePath /api
func main() {
    loadConfig()
    setupLog()
    run()
}

func loadConfig() {
    configFile := flag.String("f", "config.yaml", "Config file")
    flag.Parse()
    if len(*configFile) > 0 {
        logs.Info("Loading configuration file: ", *configFile)
        config.LoadFiles(*configFile)
    }
}

func setupLog() {
    level := config.String("app.log.level")
    if len(level) > 0 {
        logs.LogMode(logs.ParseLevel(level))
    }
}

func setupDatabase() *database.DB {
    var err error
    c := mysql.NewConfig()
    c.Addr      = fmt.Sprintf("%s:%s", config.String("database.mysql.host"), config.String("database.mysql.port"))
    c.User      = config.String("database.mysql.user")
    c.Passwd    = config.String("database.mysql.pass")
    c.DBName    = config.String("database.mysql.name")
    c.ParseTime = true
    db, err := database.Open(mysql.Open(c.FormatDSN()))
    if err != nil {
        logs.Panic("Connect to database fail. dsn = ", c.FormatDSN(), err)
    }
    model.DB = db
    return db
}

func run() {
    r := kaligo.NewRouter()
    r.AddDB(setupDatabase())

    // 中间件
    // r.Use(middlewares.ErrorHandler())    // 错误处理中间件放第一位
    r.Use(middlewares.CORS())
    // r.Use(middlewares.CSRF())
    r.Use(sessions.Sessions("GOSESSIONID", cookie.NewStore([]byte("yoursecret"))))
    setPulicRouter(r)   // 公共路由不需要cookie认证，在权限验证中间件之前注册
    r.Use(auth.Auth())  // r.Use(middleware.Cookie())
    setPrivateRouter(r) // 私有路由
    kaligo.Run(r, ":"+config.String("app.server.port", "80"))
}
