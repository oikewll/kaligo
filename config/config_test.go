package config

import (
    "testing"
    "github.com/astaxie/beego/logs"
)

func TestConfig(t *testing.T) {
    Add("database", StrMap{
        "mysql": map[string]interface{}{
            // 数据库连接信息
            "host":     Env("DB_HOST", "127.0.0.1"),
            "port":     Env("DB_PORT", "3306"),
            "database": Env("DB_DATABASE", "test"),
            "username": Env("DB_USERNAME", ""),
            "password": Env("DB_PASSWORD", ""),
            "charset":  "utf8mb4",
            "loc":      Env("DB_LOC", "Asia/Shanghai"),
            // 连接池配置
            "max_idle_connections": Env("DB_MAX_IDLE_CONNECTIONS", 300),
            "max_open_connections": Env("DB_MAX_OPEN_CONNECTIONS", 25),
            "max_life_seconds":     Env("DB_MAX_LIFE_SECONDS", 5*60),
        },
    })
    logs.Debug(Env("DB_HOST", "127.0.0.1"))
    logs.Debug(Get("database"))
    logs.Debug(Get("database.mysql.charset"))
}
