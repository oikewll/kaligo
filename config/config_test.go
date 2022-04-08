package config

import (
    "math/rand"
    "strings"
    "testing"
    "time"

    "github.com/owner888/kaligo/logs"
    "github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
    Set("database", StrMap{
        "mysql": StrMap{
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
        "custom": []string{"1", "2"},
    })
    logs.Debug(String("database.mysql.charset"))
    assert.Equal(t, Env("DB_HOST", "127.0.0.1").(string), "127.0.0.1")
    assert.Equal(t, Get[string]("database.mysql.charset"), "utf8mb4")
    assert.Equal(t, Get[bool]("database.mysql.enable"), true)
    assert.NotEqual(t, Get[bool]("database.mysql.enable"), false)
    assert.Equal(t, Get("database.sqlite", "default value"), "wrong")
    assert.Equal(t, Get[string]("database.sqlite"), "wrong")
    assert.Equal(t, Get("database.sqlite.host", "localhost"), "localhost")
    assert.Nil(t, Get[any]("database.sqlite.host"))
    assert.Equal(t, Get[[]string]("database.custom"), []string{"1", "2"})
}

func TestSet(t *testing.T) {
    Set("database.mysql.port", "1234")
    assert.Equal(t, Get[string]("database.mysql.port"), "1234")
    Set("database.mysql.port", "4321")
    assert.Equal(t, Get[string]("database.mysql.port"), "4321")
}

func TestGoroutine(t *testing.T) {
    rwConfig := func() {
        key := "test.test2"
        value := rand.Int()
        Set(key, value)
        _ = Get[int](key)
    }
    for i := 0; i < 100; i++ {
        go rwConfig()
    }
    time.Sleep(time.Second * 1)
}

func TestLoadFiles(t *testing.T) {
    Set("database", StrMap{
        "mysql": StrMap{
            "host": "127.0.0.1",
            "port": "3308",
        },
        "sqlite": "unknown",
    })
    yaml := `
    database:
      mysql:
        host: 192.168.0.1
      sqlite:
        host: 192.168.0.2
    custom: 1
    `
    LoadConfig(strings.NewReader(yaml), "yaml")
    assert.Equal(t, Get("database.mysql.host", "localhost"), "192.168.0.1")
    assert.Equal(t, Get[string]("database.mysql.port"), "3308")
    assert.Equal(t, Get[string]("database.sqlite.host"), "192.168.0.2")
    assert.Equal(t, Get[int]("custom"), 1)
}
