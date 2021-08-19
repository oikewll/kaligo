// 要运行这个test，记得先cd mysql，然后再 go test -v -count=1 或者 直接用 alias 好的 gotest
package mysql
// go test -v -count=1 mysql/mysql_test.go

// Select -> Where -> Builder -> Query -> Connection
// Update -> Where -> Builder -> Query -> Connection
// Delete -> Where -> Builder -> Query -> Connection
// Insert -> Builder -> Query -> Connection

import (
    //"encoding/json"
    //"fmt"
    "testing"
    _ "github.com/go-sql-driver/mysql"
    //"strconv"
    //"strings"
    //"reflect"
    //"time"

    //"github.com/goinggo/mapstructure"

    //"github.com/owner888/kaligo"
    //"github.com/owner888/kaligo/conf"
    //"github.com/owner888/kaligo/util"
    //"github.com/owner888/kaligo/mysql"
    //"github.com/owner888/kaligo/cache"
)

//var gormSourceDir string

//// FileWithLineNum return the file name and line number of the current file
//func FileWithLineNum() string {
    //_, file, _, _ := runtime.Caller(0)
    //t.Logf("%v", file)
    //// compatible solution to get gorm source directory with various operating systems
    ////gormSourceDir = regexp.MustCompile(`utils.utils\.go`).ReplaceAllString(file, "")
    //gormSourceDir = regexp.MustCompile(``).ReplaceAllString(file, "")
    //t.Logf("%v", gormSourceDir)

	//// the second caller usually from gorm internal, so set i start from 2
	//for i := 2; i < 15; i++ {
		//_, file, line, ok := runtime.Caller(i)
		//if ok && (!strings.HasPrefix(file, gormSourceDir) || strings.HasSuffix(file, "_test.go")) {
			//return file + ":" + strconv.FormatInt(int64(line), 10)
		//}
	//}

	//return ""
//}

type User struct {
    ID   int    `db:"id"`
    Name string `db:"name"`
    Age  int    `db:"age"`
    Sex  int    `db:"sex"`
}

func (u *User) Set(name string) *User {
    u.Name = name
    return u
}

func (u *User) Get() string {
    return u.Name
}

func TestDB(t *testing.T) {

    //modelType := reflect.ValueOf(u).Type()
    //t.Logf("modelType === %T = %v", modelType, modelType)                                 // *mysql.User
    //t.Logf("modelType.Kind() === %T = %v", modelType.Kind(), modelType.Kind())            // 类型：reflect.Slice reflect.Array reflect.Ptr(指针) reflect.Struct
    //t.Logf("modelType.Elem() === %T = %v", modelType.Elem(), modelType.Elem())            // 元素：mysql.User，能找到地址，所以可以 SetString()
    //t.Logf("modelType.Name() === %T = %v", modelType.Name(), modelType.Name())            // 类型：string
    //t.Logf("modelType.PkgPath() === %T = %v", modelType.PkgPath(), modelType.PkgPath())   // 包路径：string

    //str := FileWithLineNum()
    //t.Logf("FileWithLineNum=%v", str)
    //db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    //if err != nil {
        //panic("failed to connect database")
    //}

    //str, _ := os.Getwd()
    //conf.AppPath = str + "/../"
    //log.Printf("TestDB AppPath: [ %v ]", conf.AppPath)
    //define('APPPATH', __DIR__.'/../app');

    //model := &User{1, "test", "addr"}
    //value := reflect.ValueOf(model).Elem()
	//data := make(map[string]interface{})
	//StructToMap(value, data)
    //t.Logf("%v", data)

    db, err := New()
    if err != nil {
        t.Fatal(err)
    }

    //db.Debug = false
    //func main() {
        //defer db.SqlDB.Close()
        //router := initRouter()
        //router.Run(":8000")
    //}

    var sqlStr string

    //sqlStr = db.Query("SELECT * FROM user WHERE username = :name").Param("name", "test").Compile();
    //t.Logf("sqlStr = %v", sqlStr)

    //sqlStr = db.Query("SELECT * FROM `user`").Compile()
    //t.Logf("sqlStr = %v", sqlStr)

    //sqlStr = db.Select("id", "name").From("user").Compile()
    //t.Logf("sqlStr = %v", sqlStr)
    //sqlStr = db.Select("id", "name").From("user").Execute()
    //t.Logf("sqlStr = %v", sqlStr)

    user := User{}
    sqlStr = "SELECT name, age, sex FROM user WHERE id = :id"
    //t.Logf("sqlStr = %v\n", sqlStr)
    //db.Select("name", "age").From("user").Scan(&user).Execute()
    db.Query(sqlStr).Bind(":id", "1").Scan(&user).Execute();
    t.Logf("user name = %v --- age = %v --- Sex = %v\n", user.Name, user.Age, user.Sex)

    //sqlStr = db.Select("user.id", "user.name").From("user").
    //Join("player", "LEFT").On("user.uid", "=", "player.uid").
    ////Join("userinfo", "LEFT").On("user.uid", "=", "userinfo.uid").
    //Where("player.room_id", "=", "10").Compile()
    //t.Logf("sqlStr = %v", sqlStr)

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

