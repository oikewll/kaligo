// 要运行这个test，记得先cd mysql，然后再 go test -v -count=1 或者 直接用 alias 好的 gotest
package mysql
// go test -v -count=1 mysql/mysql_test.go

import(
    "testing"
    _ "github.com/go-sql-driver/mysql"
    //"github.com/go-gorm/gorm"
    //"gorm.io/gorm"
    //"os"
    //"log"
    "fmt"
    //"runtime"
    //"strconv"
    //"strings"
    //"regexp"
    "reflect"
    //"encoding/json"
    //"database/sql"
    //_ "github.com/go-sql-driver/mysql"
    //"github.com/owner888/kaligo"
    //"github.com/owner888/kaligo/conf"
    //"github.com/owner888/kaligo/util"
    //"github.com/owner888/kaligo/mysql"
    //"github.com/owner888/kaligo/cache"
    //"github.com/ziutek/mymysql/autorc"
    //"github.com/ziutek/mymysql/mysql" // 不能引入自己
    //"github.com/ziutek/mymysql/native" // Native engine
    // _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
)

//var gormSourceDir string

//// FileWithLineNum return the file name and line number of the current file
//func FileWithLineNum() string {
    //_, file, _, _ := runtime.Caller(0)
    //fmt.Printf("%v\n", file)
    //// compatible solution to get gorm source directory with various operating systems
    ////gormSourceDir = regexp.MustCompile(`utils.utils\.go`).ReplaceAllString(file, "")
    //gormSourceDir = regexp.MustCompile(``).ReplaceAllString(file, "")
    //fmt.Printf("%v\n", gormSourceDir)

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
    ID int
    Name string
}

func TestDB(t *testing.T) {
    u := &User{1, "kaka"}
    modelType := reflect.ValueOf(u).Type()
    fmt.Printf("modelType === %T = %v\n", modelType, modelType)                         //  *mysql.User
    fmt.Printf("modelType.Kind() === %T = %v\n", modelType.Kind(), modelType.Kind())    // 类型：reflect.Slice reflect.Array reflect.Ptr(指针) reflect.Struct
    fmt.Printf("modelType.Elem() === %T = %v\n", modelType.Elem(), modelType.Elem())    // 元素：mysql.User
    fmt.Printf("modelType.Name() === %T = %v\n", modelType.Name(), modelType.Name())    // 类型：string
    fmt.Printf("modelType.PkgPath() === %T = %v\n", modelType.PkgPath(), modelType.PkgPath())    // 类型：string

    //var str *[]string
    //str = &[]string{}
    //t.Logf("%T=%v", str, str)

    //fmt.Printf("%T = %v\n", [...]int{1, 2, 3,4, 5}, [...]int{})

    //str := FileWithLineNum()
    //t.Logf("FileWithLineNum=%v\n", str)
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


    db := NewDB("default")
    //db.Debug = false
    //sqlStr := db.Query("SELECT * FROM user").Compile()
    //t.Logf("sqlStr=%v\n", sqlStr)

    //sqlStr := db.Select("id", "name").From("user").Compile()
    sqlStr := db.Select("id", "name").From("user").Join("player", "LEFT").Compile()
    fmt.Printf("%v\n", sqlStr)

    //// Register initialisation commands
	//db.Register("set names utf8")

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
        //fmt.Println(err)
    //} else {
        //fmt.Println(db.AffectedRows())
    //}

    //if ok, err := db.Insert("user", data); err != nil {
        //fmt.Println(err)
    //} else {
        //fmt.Println(db.InsertID())
    //}

    //sql := "Select a.id,a.name,b.date From `test` a Left Join `test_part` b On a.id=b.test_id;"
    //sql := "Select * From user"
    //row, _ := db.GetOne(sql)
    //t.Logf("%v", row)
    //fmt.Println(row["id"], " --- ", row["name"])

    //rows, _ := db.GetAll(sql)
    //fmt.Printf("%v\n", rows)
    ////fmt.Printf("%#v\n", rows)
    ////fmt.Printf("%+v\n", rows)

    //jsonStr, err := json.Marshal(rows)
    //if err != nil {
        //t.Fatal(err)
    //}

    //t.Logf("Map2Json 得到 json 字符串内容:%s", jsonStr)
    //t.Logf("map: %v", rows)
    ////fmt.Println(rows)
    //for _, v := range rows {
        //fmt.Println("id = ", v["id"], " --- name = ", v["name"])
    //}

    //log.Print(db, err)

}
