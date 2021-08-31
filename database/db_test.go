package database

// Select -> Where -> Builder -> Query -> Connection
// Update -> Where -> Builder -> Query -> Connection
// Delete -> Where -> Builder -> Query -> Connection
// Insert -> Builder -> Query -> Connection

import (
	//"encoding/json"
    "fmt"
	"testing"

     //"database/sql"
    _ "github.com/go-sql-driver/mysql"
    //_ "github.com/mattn/go-sqlite3" 
	//"strconv"
    //"strings"
    //"regexp"
	//"reflect"
	//"time"
	//"github.com/goinggo/mapstructure"
	//"github.com/owner888/kaligo"
	//"github.com/owner888/kaligo/conf"
	//"github.com/owner888/kaligo/util"
	//"github.com/owner888/kaligo/mysql"
	//"github.com/owner888/kaligo/cache"
)

type User struct {
    ID   int    `db:"id"`
    Name string `db:"name"`
    Age  uint   `db:"age"`
    Sex  int    `db:"sex"`
}

func TestDB(t *testing.T) {

    //var sqlStr string

    //str, _ := os.Getwd()
    //conf.AppPath = str + "/../"
    //log.Printf("TestDB AppPath: [ %v ]", conf.AppPath)
    //define('APPPATH', __DIR__.'/../app');

    dbuser := "root"
    dbpass := "root"
    dbhost := "127.0.0.1"
    dbport := "3306"
    dbname := "test"
    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", dbuser, dbpass, dbhost+":"+dbport, dbname, "utf8mb4")
    db, err := Open("mysql", dsn)
    //db, err := Open("sqlite3", "./test.db")
    if err != nil {
        t.Fatal(err)
    }

    tables := db.Schema().ListTables("user")
    t.Logf("jsonStr = %v\n", FormatJSON(tables))

    //columns := db.Schema().ListColumns("user")
    //t.Logf("jsonStr = %v\n", FormatJSON(columns))

    //db.Schema().CreateDatabase("demo2")
    //db.Schema().DropDatabase("demo2")
    
    //fields := map[string] map[string]interface{}{
        //"id": {
            //"type": "int",
            //"constraint": 11,
            //"auto_increment": true,
        //},
        //"name": {
            //"type": "varchar",
            //"constraint": 50,
        //},
        //"title": {
            //"type": "varchar",
            //"constraint": 50,
            //"default": "mr.",
        //},
    //}
    //db.Schema().CreateTable("users2", fields, []string{"id"})
    //db.Schema().RenameTable("users2", "users3")
    //db.Schema().TruncateTable("users2")
    //db.Schema().TableExists("users3")

    //fieldExists := db.Schema().FieldExists("users3", "name")
    //fieldExists := db.Schema().FieldExists("users3", []string{"name"})
    //t.Logf("FieldExists = %v", fieldExists)

    //dropFields := db.Schema().DropFields("users3", "title")
    //dropFields := db.Schema().DropFields("users3", []string{"title"})
    //t.Logf("dropFields = %v", dropFields)

    //fields := map[string] map[string]interface{}{
        //"name": {
            //"type": "varchar",
            //"constraint": 50,
        //},
        //"title": {
            //"type": "varchar",
            //"constraint": 50,
            //"default": "mr.",
        //},
    //}
    //modifyFields := db.Schema().ModifyFields("users3", fields)
    //t.Logf("modifyFields = %v", modifyFields)

    //fields := map[string] map[string]interface{}{
        //"name5": {
            //"type": "varchar",
            //"constraint": 50,
        //},
        //"title5": {
            //"type": "varchar",
            //"constraint": 50,
            //"default": "mr.",
        //},
    //}
    //addFields := db.Schema().AddFields("users3", fields)
    //t.Logf("addFields = %v", addFields)

    //checkTable := db.Schema().OptimizeTable("users3")
    //checkTable := db.Schema().CheckTable("users3")
    //t.Logf("checkTable = %v", checkTable)

    //createIndex := db.Schema().CreateIndex("users3", "name", "name", "UNIQUE")
    //createIndex := db.Schema().CreateIndex("users3", []string{"name", "name2"}, "name333", "UNIQUE")
    //t.Logf("checkTable = %v", createIndex)
    if err = db.Schema().DropIndex("users3", "name"); err != nil {
        t.Logf("dropIndex Err = %v", err)
    }

    //foreignKey := []map[string]interface{}{
        //{
            //"constraint": "fk_uid",
            //"key": "uid",
            //"reference": map[string]string {
                //"table" : "user",   // 要关联的表
                //"column": "uid",    // 要关联的表的字段
            //},
            //"on_update": "CASCADE",
            //"on_delete": "RESTRICT",
        //},
    //}
    //addForeignKey := db.Schema().AddForeignKey("player", foreignKey)
    //t.Logf("addForeignKey = %v", addForeignKey)
    //db.Schema().DropForeignKey("player", "fk_uid")

    //db.Transaction(func(tx *DB) error {
        //// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
        //sqlStr = "insert into user(name, age, sex) values('test222', '30', '1')"
        ////_, err = db.Exec(sqlStr)
        //db.Query(sqlStr).Execute()
        //if err != nil {
            //t.Logf("%q: %s\n", err, sqlStr)
            //// 返回任何错误都会回滚事务
            //return err
        //}

        //t.Logf("RowsAffected = %d: %d\n", tx.RowsAffected, tx.LastInsertId)
        //// 返回 nil 提交事务
        //return nil
    //})

    //db.Begin()
    ////defer db.Rollback()

    ////sqlStr = "drop table user;"
    ////sqlStr = "create table user (id integer not null primary key, name text, age integer, sex integer);"
    ////sqlStr = "delete from user;"
    //sqlStr = "insert into user(name, age, sex) values('test222', '30', '1')"
    //_, err = db.Exec(sqlStr)
    //if err != nil {
        //t.Logf("%q: %s\n", err, sqlStr)
        //return
    //}

    ////db.Rollback()
    //db.Commit()



    //db.Debug = false
    //func main() {
    //defer db.SqlDB.Close()
    //router := initRouter()
    //router.Run(":8000")
    //}

    //var ages []int64
    //db.Query("SELECT age FROM user").Scan(&ages).Execute()
    //for _, v := range ages {
        //t.Logf("ages %T = %v\n", v, v)
    //}
    //t.Logf("jsonStr = %v\n", FormatJSON(ages))

    //var user User
    //sqlStr = "SELECT id, name, age, sex FROM user WHERE id = :id"
    //q := db.Query(sqlStr).Bind(":id", "1").Scan(&user).Execute()
    //if q.Error != nil {
        //t.Logf("q.Error = %v\n", q.Error)
    //}
    //t.Logf("jsonStr = %v\n", FormatJSON(user))

    ////var users []User
    //users := []User{}
    //q := db.Query("SELECT id, name, age, sex FROM user").Scan(&users).Execute()
    //if q.Error != nil {
        //t.Logf("q.Error = %v\n", q.Error)
    //}
    //t.Logf("jsonStr = %v\n", FormatJSON(users))

    //var count int64
    //q = db.Query("SELECT COUNT(*) FROM user").Scan(&count).Execute();
    //if q.Error != nil {
        //t.Logf("q.Error = %v\n", q.Error)
    //}
    //t.Logf("jsonStr = %v\n", FormatJSON(count))

    //var name string
    //q = db.Query("SELECT name FROM user WHERE id = 1").Bind(":id", "1").Scan(&name).Execute();
    //if q.Error != nil {
        //t.Logf("q.Error = %v\n", q.Error)
    //}
    //t.Logf("jsonStr = %v\n", FormatJSON(name))

    //result := map[string]interface{}{}
    //q = db.Query("SELECT name, age FROM user WHERE id = 1").Bind(":id", "1").Scan(&result).Execute()
    //if q.Error != nil {
        //t.Logf("q.Error = %v\n", q.Error)
    //}
    //t.Logf("jsonStr = %v\n", FormatJSON(result))

    //results := []map[string]interface{}{}
    //q = db.Query("SELECT name, age FROM user").Scan(&results).Execute()
    //if q.Error != nil {
        //t.Logf("q.Error = %v\n", q.Error)
    //}
    //t.Logf("jsonStr = %v\n", FormatJSON(results))

    ////users := []User{}
    //users := []map[string]interface{}{}
    //db.Select("user.id", "user.name").From("user").
    //Join("player", "LEFT").On("user.uid", "=", "player.uid").
    ////Join("userinfo", "LEFT").On("user.uid", "=", "userinfo.uid").
    //Where("player.room_id", "=", "10").
    //Scan(&users).Execute()
    //t.Logf("jsonStr = %v\n", FormatJSON(users))

    //sqlStr = db.Insert("user", []string{"id", "name"}).Values([]string{"10", "test"}).Compile()
    //sqlStr = db.Insert("user", []string{"id", "name"}).Values([][]string{{"10", "test"}, {"20", "demo"}}).Compile()
    //var query *Query
    //// 全部字段复制
    //query  = db.Query("SELECT * FROM `user_history`", SELECT)
    //sqlStr = db.Insert("user").SubSelect(query).Compile()
    //t.Logf("sqlStr = %v", sqlStr)
    //// 只复制 id、name 两个字段
    //query  = db.Query("SELECT `id`, `name` FROM `user_history`", SELECT)
    //sqlStr = db.Insert("user", []string{"id", "name"}).SubSelect(query).Compile()
    //t.Logf("sqlStr = %v", sqlStr)

    //sets := map[string]string{"id": "10", "name":"demo"}
    //sqlStr = db.Update("user").Join("player", "LEFT").On("user.uid", "=", "player.uid").Set(sets).Where("player.room_id", "=", "10").Compile()
    //t.Logf("sqlStr = %v", sqlStr)

    //// 暂时不支持DELETE JOIN写法
    ////sqlStr = db.Delete("user").Join("player", "LEFT").On("user.uid", "=", "player.uid").Where("player.id", "=", "test").Compile()
    //sqlStr = db.Delete("user").Where("nickname", "=", "test").Compile()
    //t.Logf("sqlStr = %v", sqlStr)

    //data := map[string]string {
    //"name": "nam'e111",
    //"pass": "pas\"s111",
    //"sex":"1",
    //}
    //t.Fatal(data)   // 输出并中断程序

    //if ok, err := db.Update("user", data, "`id`=1"); err != nil {
    //t.Log(err)
    //} else {
    //t.Log(db.AffectedRows())
    //}

    //if ok, err := db.Insert("user", data); err != nil {
    //t.Log(err)
    //} else {
    //t.Log(db.InsertID())
    //}

    //sql := "Select a.id,a.name,b.date From `test` a Left Join `test_part` b On a.id=b.test_id;"
    //sql := "Select * From user"
    //row, _ := db.GetOne(sql)
    //t.Logf("%v", row)
    //t.Log(row["id"], " --- ", row["name"])

    //rows, _ := db.GetAll(sql)
    //t.Logf("%v", rows)
    ////t.Logf("%#v", rows)

    //jsonStr, err := json.Marshal(rows)
    //if err != nil {
    //t.Fatal(err)
    //}

    //t.Logf("Map2Json 得到 json 字符串内容:%s", jsonStr)
    //t.Logf("map: %v", rows)
    ////t.Log(rows)
    //for _, v := range rows {
    //t.Log("id = ", v["id"], " --- name = ", v["name"])
    //}

    //log.Print(db, err)

}
