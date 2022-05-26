package main

import (
    "examples/middleware/auth"
    "examples/model"
    "flag"
    "fmt"
    // "os/exec"

    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/config"
    "github.com/owner888/kaligo/database"
    "github.com/owner888/kaligo/database/driver/mysql"
    "github.com/owner888/kaligo/logs"
    "github.com/owner888/kaligo/sessions"
    "github.com/owner888/kaligo/sessions/cookie"
    // "github.com/owner888/kaligo/sessions/redis"
    "github.com/owner888/kaligo/middlewares"
)

// swag 不要集成到项目,直接安装工具包即可
// go install github.com/swaggo/swag/cmd/swag@latest
// 文档: https://github.com/swaggo/swag#api-operation

// @title Kaligo Example API
// @description 
// @version 1.0
// @host localhost:8080
// @BasePath /api
func main() {
    // cmd := exec.Command(GOPATH+"swag init")
    // cmd := exec.Command("/Users/coffee/Documents/golang/bin/swag init")
    // err := cmd.Run()
    // if err != nil {
    //     logs.Fatalf("cmd.Run() failed with %s\n", err)
    // }

    loadConfig()
    setupLog()
    db := setupDatabase()
    run(db)
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

func run(db *database.DB) {
    r := kaligo.NewRouter()

    // 中间件
    r.Use(middlewares.CORS())
    // r.Use(middlewares.CSRF())    // 防止CSRF攻击
    // 创建基于cookie的存储引擎，secret 参数是用于加密的密钥
    store := cookie.NewStore([]byte("secret"))
    // store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
    r.Use(sessions.Sessions("mysession", store))
    r.Use(auth.Auth())

    r.AddDB(db)
    AddRoutes(r)
    kaligo.Run(r, ":"+config.String("app.server.port", "80"))
}
