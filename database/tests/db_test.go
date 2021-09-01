package tests

// Select -> Where -> Builder -> Query -> Connection
// Update -> Where -> Builder -> Query -> Connection
// Delete -> Where -> Builder -> Query -> Connection
// Insert -> Builder -> Query -> Connection

import (
	//"encoding/json"
    "fmt"
	"testing"

    //"database/sql"
    //_ "github.com/go-sql-driver/mysql"
    //_ "github.com/mattn/go-sqlite3" 
    "github.com/owner888/kaligo/database"
    //sqlite "github.com/owner888/kaligo/database/driver/sqlite"
    mysql "github.com/owner888/kaligo/database/driver/mysql"
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
    t.Logf("%v", dsn)
    db, err := database.Open(mysql.Open(dsn))
    //db, err := database.Open(sqlite.Open("./test.db"))
    if err != nil {
        t.Fatal(err)
    }

    //databases := db.Schema().CurrentDatabase()
    //t.Logf("jsonStr = %v\n", database.FormatJSON(databases))

    //databases := db.Schema().ListDatabases("demo")
    //t.Logf("jsonStr = %v\n", database.FormatJSON(databases))

    //tables := db.Schema().ListTables("user")
    //t.Logf("jsonStr = %v\n", database.FormatJSON(tables))

    //columns := db.Schema().ListColumns("user")
    //t.Logf("jsonStr = %v\n", database.FormatJSON(columns))

    //indexes := db.Schema().ListIndexes("user")
    //t.Logf("jsonStr = %v\n", database.FormatJSON(indexes))

    //db.Schema().CreateDatabase("demo2")
    //db.Schema().DropDatabase("demo2")
    
    //fields := []map[string]interface{}{
        //{
            //"name": "id",
            //"type": "int",
            //"constraint": 11,
            //"auto_increment": true,
        //},
        //{
            //"name": "name",
            //"type": "varchar",
            //"constraint": 50,
        //},
        //{
            //"name": "title",
            //"type": "varchar",
            //"constraint": 50,
            //"default": "mr.",
        //},
    //}
    //err = db.Schema().CreateTable("user", fields, []string{"id"})
    //err = db.Schema().RenameTable("user", "user222")
    //err = db.Schema().TruncateTable("user")
    //if err != nil {
        //t.Logf("Operation Table Err = %v", err)
    //}
    //ok := db.Schema().TableExists("user")
    //t.Logf("TableExists ok = %v", ok)

    //ok := db.Schema().FieldExists("user", "name")
    //ok := db.Schema().FieldExists("user", []string{"name"})
    //t.Logf("FieldExists ok = %v", ok)

    //err = db.Schema().DropFields("user", "title")
    //err = db.Schema().DropFields("user", []string{"name5", "title5"})

    //fields := []map[string]interface{}{
        //{
            //"oldname" : "name",
            //"name": "name111",
            //"type": "varchar",
            //"constraint": 50,
        //},
        //{
            //"name": "title",
            //"type": "varchar",
            //"constraint": 50,
            //"default": "mr.",
        //},
    //}
    //err = db.Schema().ModifyFields("user", fields)

    //fields := []map[string]interface{}{
        //{
            //"name": "name6",
            //"type": "varchar",
            //"constraint": 50,
        //},
    //}
    //err = db.Schema().AddFields("user", fields)
    //if err != nil {
        //t.Logf("Operation Fields Err = %v", err)
    //}

    //ok := db.Schema().OptimizeTable("user")
    //ok := db.Schema().CheckTable("user")
    //t.Logf("Operation Table ok = %v", ok)

    //err = db.Schema().CreateIndex("user", "name", "name_idx6", "UNIQUE")
    //err = db.Schema().CreateIndex("user", []string{"age", "sex"}, "name_idx5", "UNIQUE")
    //err = db.Schema().RenameIndex("user", "ageindex2", "ageindex3")
    //err = db.Schema().DropIndex("user", "name")
    //if err != nil {
        //t.Logf("Operation Index Err = %v", err)
    //}

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
    //err = db.Schema().AddForeignKey("player", foreignKey)
    ////err = db.Schema().DropForeignKey("player", "fk_uid")
    //if err != nil {
        //t.Logf("Operation Foreign Key Err = %v", err)
    //}

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


    ////db.Rollback()
    //db.Commit()

    db.Begin()
    db.Insert("user", []string{"name", "age"}).Values([]string{"test111", "20"}).Execute()

    db.SavePoint("sp1")
    db.Insert("user", []string{"name", "age"}).Values([]string{"test222", "23"}).Execute()
    db.RollbackTo("sp1")    // Rollback the user name is test222

    db.Commit()  // Commit the user name is test111 


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
