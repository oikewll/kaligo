package tests

import (
    "os"
    "testing"

    "github.com/owner888/kaligo/logs"
    "github.com/owner888/kaligo/database"
    mysql "github.com/owner888/kaligo/database/driver/mysql"
    "github.com/stretchr/testify/assert"
)

var db *database.DB

func TestMain(m *testing.M) {
    logs.LogMode(logs.LevelDebug)
    db, _ = database.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/chrome?charset=utf8mb4"))
    code := m.Run()
    os.Exit(code)
}

func TestMigratorCurrentDatabase(t *testing.T) {
    databases := db.Migrator().CurrentDatabase()
    assert.Equal(t, databases, "chrome")
}

func TestMigratorListDatabases(t *testing.T) {
    databases, _ := db.Migrator().ListDatabases("chrome")
    t.Logf("jsonStr = %v\n", database.FormatJSON(databases))
}

func TestUpdate(t *testing.T) {
    q := db.Insert("keywords").Columns([]string{`word`, `creator`}).Values([]string{"电影网站", "1"}).OnDuplicateKeyUpdate(map[string]string{`creator`: "3"})
    logs.Info(q)
}

func TestConfig(t *testing.T) {
    cfg := mysql.NewConfig()
    cfg.Addr    = "localhost:3308"
    cfg.DBName  = "test"
    cfg.User    = "root"
    cfg.Passwd  = "pw"
    dsn := cfg.FormatDSN()
    assert.Equal(t, "root:pw@tcp(localhost:3308)/test", dsn)
}
