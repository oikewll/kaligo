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
            "enable":   true,
            "loc":      Env("DB_LOC", "Asia/Shanghai"),
            // 连接池配置
            "max_idle_connections": Env("DB_MAX_IDLE_CONNECTIONS", 300),
            "max_open_connections": Env("DB_MAX_OPEN_CONNECTIONS", 25),
            "max_life_seconds":     Env("DB_MAX_LIFE_SECONDS", 5*60),
        },
        "sqlite": "wrong",
    })
    logs.Debug(GetString("database.mysql.charset"))
    assertEqual(t, Env("DB_HOST", "127.0.0.1").(string), "127.0.0.1")
    assertEqual(t, Get[string]("database.mysql.charset"), "utf8mb4")
    assertEqual(t, Get[bool]("database.mysql.enable"), true)
    assertNotEqual(t, Get[bool]("database.mysql.enable"), false)
    assertEqual(t, Get("database.sqlite", "default value"), "wrong")
    assertEqual(t, Get[string]("database.sqlite"), "wrong")
    assertEqual(t, Get("database.sqlite.host", "localhost"), "localhost")
    assert(t, Get[any]("database.sqlite.host") == nil)
}

func assertEqual[T comparable](t *testing.T, expected, actual T, messages ...any) {
    assert(t, expected == actual, append(messages, actual, "is not expected", expected)...)
}

func assertNotEqual[T comparable](t *testing.T, expected, actual T, messages ...any) {
    assert(t, expected != actual, append(messages, actual, "is not expected")...)
}

func assert(t *testing.T, result bool, messages ...any) {
    if !result {
        t.Fatal(append([]any{"Fail!"}, messages...)...)
    }
}
