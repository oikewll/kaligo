package tests

import (
    "testing"
    // "github.com/owner888/kaligo/logs"
    // "github.com/owner888/kaligo/database"
    "github.com/stretchr/testify/assert"
)


func TestInsert(t *testing.T) {
    q, err := db.Insert("demo_user", []string{"username", "password", "gender"}).
    Values([]any{"test1", "test1passwd", 2}).
    Execute()
    assert.NoError(t, err)
    // logs.Debug("LastInsertId = ", q.LastInsertId)

    _, err = db.Insert("demo_player", []string{"uid", "room_id"}).
    Values([]any{q.LastInsertId, 1}).
    Execute()
    assert.NoError(t, err)
}

func TestInsertMutil(t *testing.T) {
    _, err := db.Insert("demo_user", []string{"username", "password", "gender"}).Values([][]any{{"test2", "test2passwd", 20}, {"test3", "test3passwd", 25}}).Execute()
    assert.NoError(t, err)
    // logs.Debug("RowsAffected = ", q.RowsAffected)
}

func TestInsertCryptData(t *testing.T) {
    _, err := db.Insert("demo_user", []string{"realname", "gender"}).
    SetCryptKey("aaa").
    SetCryptFields(map[string][]string{
        "demo_user": {"realname"},
    }).
    Values([]any{"demoname", 1}).
    Execute()

    assert.NoError(t, err)
    // logs.Debug("LastInsertId = ", q.LastInsertId)
}

func TestInsertFromSelect(t *testing.T) {
    // 全部字段复制
    // q := db.Query("SELECT * FROM `demo_user`", database.SELECT)
    // _, err := db.Insert("demo_user_tmp").SubSelect(q).Execute()
    // assert.NoError(t, err)

    // // 只复制 id、username 两个字段
    // q  = db.Query("SELECT `id`, `username` FROM `user`", database.SELECT)
    // sqlStr = db.Insert("demo_user_tmp", []string{"id", "username"}).SubSelect(q).Compile()
    // logs.Debug("sqlStr = ", sqlStr)
}

// TO DO
// func TestInsertOnDuplicateKeyUpdate(t *testing.T) {
//     q := db.Insert("demo_user").Columns([]string{`username`, `creator`}).Values([]any{"demoname", 1}).OnDuplicateKeyUpdate(map[string]any{`creator`: 3})
//     logs.Info(q)
// }


