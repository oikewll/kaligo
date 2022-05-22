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
)

// @title Kaligo Example API
// @version 1.0
// @host localhost/api

func main() {
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
    c.Addr = fmt.Sprintf("%s:%s", config.String("database.mysql.host"), config.String("database.mysql.port"))
    c.User = config.String("database.mysql.user")
    c.Passwd = config.String("database.mysql.pass")
    c.DBName = config.String("database.mysql.name")
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
    // 创建基于cookie的存储引擎，secret11111 参数是用于加密的密钥
	store := cookie.NewStore([]byte("secret11111"))
    r.Use(auth.Auth)
    r.Use(sessions.Sessions("mysession", store))
    r.AddDB(db)
    AddRoutes(r)
    kaligo.Run(r, ":"+config.String("app.server.port", "80"))
}
