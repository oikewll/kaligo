package tests

// Select -> Where -> Builder -> Query -> Connection
// Update -> Where -> Builder -> Query -> Connection
// Delete -> Where -> Builder -> Query -> Connection
// Insert -> Builder -> Query -> Connection

import (
	"testing"
	//"github.com/goinggo/mapstructure"
)

func TestDB(t *testing.T) {

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

    //var ages []int64
    //q := db.Query("SELECT age FROM user").Scan(&ages).Execute()
    //if q.Error != nil {
        //t.Logf("q.Error = %v\n", q.Error)
    //} else {
        //t.Logf("jsonStr = %v\n", database.FormatJSON(ages))
    //}

    //var user User
    //sqlStr = "SELECT id, name, age, sex FROM user WHERE id = :id"
    //q := db.Query(sqlStr).Bind(":id", "1").Scan(&user).Execute()
    //if q.Error != nil {
        //t.Logf("q.Error = %v\n", q.Error)
    //} else {
        //t.Logf("jsonStr = %v\n", database.FormatJSON(user))
    //}

    ////var users []User
    //users := []User{}
    //q := db.Query("SELECT id, name, age, sex FROM user").Scan(&users).Execute()
    //if q.Error != nil {
        //t.Logf("q.Error = %v\n", q.Error)
    //} else {
        //t.Logf("jsonStr = %v\n", database.FormatJSON(users))
    //}

    //var count int64
    //q := db.Query("SELECT COUNT(*) FROM user").Scan(&count).Execute();
    //if q.Error != nil {
        //t.Logf("q.Error = %v\n", q.Error)
    //} else {
        //t.Logf("jsonStr = %v\n", database.FormatJSON(count))
    //}

    //var name string
    //q := db.Query("SELECT `name` FROM `user` WHERE `id` = '1'").Scan(&name).Execute();
    //if q.Error != nil {
        //t.Logf("q.Error = %v\n", q.Error)
    //} else {
        //t.Logf("jsonStr = %v\n", database.FormatJSON(name))
    //}

    //result := map[string]interface{}{}
    //q := db.Query("SELECT name, age FROM user WHERE id = :id").Bind("id", "2").Scan(&result).Execute()
    //if q.Error != nil {
        //t.Logf("q.Error = %v\n", q.Error)
    //} else {
        //t.Logf("jsonStr = %v\n", database.FormatJSON(result))
    //}

    //results := []map[string]interface{}{}
    //q := db.Select("id", "name", "age").From("user").Scan(&results).Execute()
    //if q.Error != nil {
        //t.Logf("q.Error = %v\n", q.Error)
    //} else {
        //t.Logf("jsonStr = %v\n", database.FormatJSON(results))
    //}

    ////users := []User{}
    //users := []map[string]interface{}{}
    //db.Select("user.id", "user.name").From("user").
    //Join("player", "LEFT").On("user.uid", "=", "player.uid").
    ////Join("userinfo", "LEFT").On("user.uid", "=", "userinfo.uid").
    //Where("player.room_id", "=", "10").
    //Scan(&users).Execute()
    //t.Logf("jsonStr = %v\n", database.FormatJSON(users))

    //q := db.Insert("user", []string{"name", "age"}).Values([]string{"test111", "20"}).Execute()
    //q := db.Insert("user", []string{"name", "age"}).Values([][]string{{"test111", "20"}, {"test222", "25"}}).Execute()
    //if q.Error != nil {
        //t.Fatal(q.Error)
    //}

    //// 全部字段复制
    //query  = db.Query("SELECT * FROM `user_history`", SELECT)
    //sqlStr = db.Insert("user").SubSelect(query).Compile()
    //t.Logf("sqlStr = %v", sqlStr)
    //// 只复制 id、name 两个字段
    //query  = db.Query("SELECT `id`, `name` FROM `user_history`", SELECT)
    //sqlStr = db.Insert("user", []string{"id", "name"}).SubSelect(query).Compile()
    //t.Logf("sqlStr = %v", sqlStr)

    //sets := map[string]string{"name":"demo111"}
    //q := db.Update("user").Set(sets).Where("id", "=", "1").Execute()
    //if q.Error != nil {
        //t.Logf("%v", q.Error)
    //}
    //sqlStr = db.Update("user").Join("player", "LEFT").On("user.uid", "=", "player.uid").Set(sets).Where("player.room_id", "=", "10").Compile()
    //t.Logf("sqlStr = %v", sqlStr)

    //// 暂时不支持DELETE JOIN写法
    ////sqlStr = db.Delete("user").Join("player", "LEFT").On("user.uid", "=", "player.uid").Where("player.id", "=", "test").Compile()
    //sqlStr = db.Delete("user").Where("nickname", "=", "test").Compile()
    //t.Logf("sqlStr = %v", sqlStr)
}
