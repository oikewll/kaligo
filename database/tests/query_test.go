package tests

import (
    "os"
    "testing"

    // "github.com/owner888/kaligo/logs"
    "github.com/owner888/kaligo/database"
    mysql "github.com/owner888/kaligo/database/driver/mysql"
    "github.com/stretchr/testify/assert"
)

var db *database.DB

func TestMain(m *testing.M) {
    // 开始前做初始化工作
    // db, err := database.Open(sqlite.Open("./test.db"))
    db, _ = database.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4"))

    // 执行增删改查
    code := m.Run()

    // 结束前做数据清理工作
    // ...

    os.Exit(code)
}

func TestQuery(t *testing.T) {
    //var sqlStr string
    //
    // user := User{
    //     ID   : 1,
    //     Name : "test111",
    //     Age  : 25,
    //     Sex  : 1,
    // }
    // user.DB = db
    // user.Save()

    var ages []int64
    // sqlStr = SELECT `age` FROM `user` WHERE `id` = ?";
    // _, err := db.Query("SELECT `age` FROM `user` WHERE `id` = :id").Scan(&ages).Bind(":id", "1").Execute()
    // SELECT AES_DECRYPT(`age`, 'aaa') AS `age` FROM `user` WHERE `id` = ?
    _, err := db.Query("SELECT `age` FROM `user` WHERE `id` = :id").
        // SetCryptKey("aaa").
        // SetCryptFields(map[string][]string{
        //     "user": {"name", "age"},
        // }).
        Scan(&ages).Bind(":id", "1 or 1=1").Execute()
    assert.NoError(t, err)
    assert.NotNil(t, ages)
}

func TestSelectDecrypt(t *testing.T) {
    var ages []int64
    _, err := db.Select("age").From("user").Where("id", "=", "1").
        SetCryptKey("aaa").
        SetCryptFields(map[string][]string{
            "user": {"name", "age"},
        }).
        Scan(&ages).Execute()
    assert.NoError(t, err)
    assert.NotNil(t, ages)
}

// func TestQueryCount(t *testing.T) {
//     var count int64
//     _, err := db.Query("SELECT COUNT(*) FROM `user`").Scan(&count).Execute();
//     assert.NoError(t, err)
// }
//
// func TestQueryBind(t *testing.T) {
//     var user User
//     sqlStr := "SELECT `id`, `name`, `age` FROM `user` WHERE `id` = :id"
//     _, err := db.Query(sqlStr).Bind(":id", "1").Scan(&user).Execute()
//     assert.NoError(t, err)
//
//     var users []User
//     _, err = db.Query("SELECT `id`, `name`, `age`  FROM `user`").Scan(&users).Execute()
//     assert.NoError(t, err)
//
//     // logs.Debug(database.FormatJSON(users))
// }

// func TestQueryStringMap(t *testing.T) {
//     result := map[string]interface{}{}
//     _, err := db.Query("SELECT `name`, `age` FROM `user` WHERE `id` = :id").Bind("id", "1").Scan(&result).Execute()
//     assert.NoError(t, err)
//
//     // logs.Debug(database.FormatJSON(result))
// }
//
// func TestQuerySliceStringMap(t *testing.T) {
//     results := []map[string]interface{}{}
//     _, err := db.Select("id", "name", "age").From("user").Scan(&results).Execute()
//     assert.NoError(t, err)
//
//     // logs.Debug(database.FormatJSON(results))
// }

// func TestQueryJoin(t *testing.T) {
//     // users := []map[string]interface{}{}
//     users := []User{}
//     _, err := db.Select("user.id", "user.name").From("user").
//     Join("player", "LEFT").On("user.uid", "=", "player.uid").
//     Where("player.room_id", "=", "10").
//     Scan(&users).Execute()
//     assert.NoError(t, err)
//
//     logs.Debug(database.FormatJSON(users))
// }

func TestConfig(t *testing.T) {
    cfg := mysql.NewConfig()
    cfg.Addr = "localhost:3308"
    cfg.DBName = "test"
    cfg.User = "root"
    cfg.Passwd = "pw"
    dsn := cfg.FormatDSN()
    assert.Equal(t, "root:pw@tcp(localhost:3308)/test", dsn)
}
