package tests

import (
    "testing"
    // "github.com/owner888/kaligo/logs"
    // "github.com/owner888/kaligo/database"
    // "github.com/stretchr/testify/assert"
)


func TestInsert(t *testing.T) {
    // SET GLOBAL general_log=on;
    // SET GLOBAL general_log_file='/tmp/general.log';

    // q, err := db.Insert("user", []string{"username", "password", "age"}).Values([]string{"test111", "test111passwd", "2"}).Execute()
    // assert.NoError(t, err)
    // logs.Debug("LastInsertId = ", q.LastInsertId)
}

func TestInsertMutil(t *testing.T) {
    // q, err := db.Insert("user", []string{"username", "password", "age"}).Values([][]string{{"test111", "test111passwd", "20"}, {"test222", "test222passwd", "25"}}).Execute()
    // assert.NoError(t, err)
    // logs.Debug("RowsAffected = ", q.RowsAffected)
}

// func TestInsertFromSelect(t *testing.T) {
//     // 全部字段复制
//     q := db.Query("SELECT * FROM `user_history`", database.SELECT)
//     sqlStr := db.Insert("user").SubSelect(q).Compile()
//     logs.Debug("sqlStr = ", sqlStr)
//
//     // 只复制 id、name 两个字段
//     q  = db.Query("SELECT `id`, `name` FROM `user_history`", database.SELECT)
//     sqlStr = db.Insert("user", []string{"id", "name"}).SubSelect(q).Compile()
//     logs.Debug("sqlStr = ", sqlStr)
// }
//
// func TestInsertOnDuplicateKeyUpdate(t *testing.T) {
//     q := db.Insert("keywords").Columns([]string{`word`, `creator`}).Values([]string{"电影网站", "1"}).OnDuplicateKeyUpdate(map[string]string{`creator`: "3"})
//     logs.Info(q)
// }


