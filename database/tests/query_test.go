package tests

import (
    "fmt"
    "os"
    "testing"

    "github.com/owner888/kaligo/database"
    mysql "github.com/owner888/kaligo/database/driver/mysql"
)

var db *database.DB

func TestMain(m *testing.M) {
    db, _ = database.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/chorme?charset=utf8mb4"))
    code := m.Run()
    os.Exit(code)
}

func TestUpdate(t *testing.T) {
    sql := db.Insert("keywords").Columns([]string{`word`, `creator`}).Values([]string{"电影网站", "1"}).OnDuplicateKeyUpdate(map[string]string{`creator`: "3"}).Compile()
    fmt.Println(sql)
}
