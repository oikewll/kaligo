package main

import (
	"examples/controller"
	"flag"

	"github.com/owner888/kaligo"
	"github.com/owner888/kaligo/config"
	"github.com/owner888/kaligo/database"
	"github.com/owner888/kaligo/logs"
)

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
		logs.Info("加载配置文件", *configFile)
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
	return nil
}

func run(db *database.DB) {
	r := kaligo.NewRouter()
	r.AddDB(db)
	controller.AddRoutes(r)
	kaligo.Run(r, ":"+config.String("app.server.port", "80"))
}
