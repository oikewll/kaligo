// 要运行这个test，记得先cd mysql，然后再 go test -v -count=1 或者 直接用 alias 好的 gotest
package mysql
// go test -v -count=1 mysql/mysql_test.go

// Select -> Where -> Builder -> Query -> Connection
// Update -> Where -> Builder -> Query -> Connection
// Delete -> Where -> Builder -> Query -> Connection
// Insert -> Builder -> Query -> Connection

import(
    "testing"
    //"fmt"
    _ "github.com/go-sql-driver/mysql"
    //"time"
    //"strconv"
    //"strings"
    "reflect"
    //"regexp"
    //"encoding/json"
    //"database/sql"

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

func TestDB(t *testing.T) {

    //u := &User{1, "kaka", "hk"}
    //modelType := reflect.ValueOf(u).Type()
    //t.Logf("modelType === %T = %v", modelType, modelType)                                 //  *mysql.User
    //t.Logf("modelType.Kind() === %T = %v", modelType.Kind(), modelType.Kind())            // 类型：reflect.Slice reflect.Array reflect.Ptr(指针) reflect.Struct
    //t.Logf("modelType.Elem() === %T = %v", modelType.Elem(), modelType.Elem())            // 元素：mysql.User
    //t.Logf("modelType.Name() === %T = %v", modelType.Name(), modelType.Name())            // 类型：string
    //t.Logf("modelType.PkgPath() === %T = %v", modelType.PkgPath(), modelType.PkgPath())   // 类型：string

    //str := FileWithLineNum()
    //t.Logf("FileWithLineNum=%v", str)
    //db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    //if err != nil {
        //panic("failed to connect database")
    //}

    //memory := cache.NewMemory()
    //memory.Set("name", "kaka", 10)
    //name := memory.Get("name")
    //t.Logf("name: [ %v ]", name)


    //str, _ := os.Getwd()
    //conf.AppPath = str + "/../"
    //log.Printf("TestDB AppPath: [ %v ]", conf.AppPath)
    //define('APPPATH', __DIR__.'/../app');

    //model := &User{1, "test", "addr"}
    //value := reflect.ValueOf(model).Elem()
	//data := make(map[string]interface{})
	//mapStructToMap(value, data)
    //t.Logf("%v", data)

    db := New()
    //db.Debug = false

    //// 测试 once.Do()，确实所有协程都只能拿到一样的对象
    //for i := 0; i < 5; i++ {
        //go func(i int) {
            //name := "db" + strconv.Itoa(i)
            //db := New(name)
            //t.Logf("%T = %p", db, db)
            //t.Log(name)
        //} (i)
    //}
    //time.Sleep(2 * time.Second)

    var sqlStr string    
    //sqlStr = db.Query("SELECT * FROM `user`").Compile()
    //t.Logf("sqlStr = %v", sqlStr)

    //sqlStr = db.Select("id", "name").From("user").Compile()
    //t.Logf("sqlStr = %v", sqlStr)
    //sqlStr = db.Select("id", "name").From("user").Execute()
    //t.Logf("sqlStr = %v", sqlStr)
    sqlStr = db.Select("user.id", "user.name").From("user").
    Join("player", "LEFT").On("user.uid", "=", "player.uid").
    //Join("userinfo", "LEFT").On("user.uid", "=", "userinfo.uid").
    Where("player.room_id", "=", "10").Compile()
    t.Logf("sqlStr = %v", sqlStr)

    //sqlStr = db.Insert("user", []string{"id", "name"}).Values([]string{"10", "test"}).Compile()
    //sqlStr = db.Insert("user", []string{"id", "name"}).Values([][]string{{"10", "test"}, {"20", "demo"}}).Compile()
    var query *Query    
    // 全部字段复制
    query  = db.Query("SELECT * FROM `user_history`", SELECT)
    sqlStr = db.Insert("user").SubSelect(query).Compile()
    t.Logf("sqlStr = %v", sqlStr)
    // 只复制 id、name 两个字段
    query  = db.Query("SELECT `id`, `name` FROM `user_history`", SELECT)
    sqlStr = db.Insert("user", []string{"id", "name"}).SubSelect(query).Compile()
    t.Logf("sqlStr = %v", sqlStr)

    sets := map[string]string{"id": "10", "name":"demo"}
    sqlStr = db.Update("user").Join("player", "LEFT").On("user.uid", "=", "player.uid").Set(sets).Where("player.room_id", "=", "10").Compile()
    t.Logf("sqlStr = %v", sqlStr)

    // 暂时不支持DELETE JOIN写法
    //sqlStr = db.Delete("user").Join("player", "LEFT").On("user.uid", "=", "player.uid").Where("player.id", "=", "test").Compile()
    sqlStr = db.Delete("user").Where("nickname", "=", "test").Compile()
    t.Logf("sqlStr = %v", sqlStr)

    //// my is in unconnected state
	////mysql.checkErr(t, c.Use(dbname), nil)

    ////t.Logf("%s", conn)

    //defer db.Close()

    //if err != nil {
        //t.Fatal(err)
    //}

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
    ////t.Logf("%+v", rows)

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

// mapStructToMap 将一个结构体所有字段(包括通过组合得来的字段)到一个map中
// value:结构体的反射值
// data:存储字段数据的map
func mapStructToMap(value reflect.Value, data map[string]interface{}) {
    if value.Kind() != reflect.Struct {
        return
    }

    for i := 0; i < value.NumField(); i++ {
        var fieldValue = value.Field(i)
        if fieldValue.CanInterface() {
            var fieldType = value.Type().Field(i)
            if fieldType.Anonymous {
                // 匿名组合字段,进行递归解析
                mapStructToMap(fieldValue, data)
            } else {
                // 非匿名字段
                var fieldName = fieldType.Tag.Get("db")
                if fieldName == "-" {
                    continue
                }
                if fieldName == "" {
                    fieldName = transFieldName(fieldType.Name)
                }
                data[fieldName] = fieldValue.Interface()
                //t.Log(fieldName + ":" + fieldValue.Interface().(string))
            }
        }
    }
}
