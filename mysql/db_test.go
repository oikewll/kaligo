// 要运行这个test，记得先cd mysql，然后再 go test -v -count=1 或者 直接用 alias 好的 gotest
package mysql
// go test -v -count=1 mysql/mysql_test.go

import(
    "testing"
    //"fmt"
    _ "github.com/go-sql-driver/mysql"
    //"time"
    //"strconv"
    //"strings"
    //"regexp"
    "reflect"
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

type Player struct {
    playerName string
}

func NewPlayer() *Player {
    return &Player{
        playerName: "playerName",
    }
}

func (p *Player) setPlayerName() {
    p.playerName = "kakakak"
}

type User struct {
    ID int          `db:"id"`
    Name string     `db:"name"`
    Addr string     `db:"addr"`
    *Player
}

type User1 struct {
    Name string     `db:"name"`
    *Player
}

func NewUser() *User {
    return &User{
        Name: "userName",
        Player: &Player{
            playerName: "playerName",
        },
    }
}

func TestDB(t *testing.T) {
    //p := NewPlayer()
    //u := &User{
        //Name: "userName",
        //Player: p,
    //}
    ////t.Logf("userName = %v\n", u.Name)
    //t.Logf("playerName = %v\n", u.playerName)

    //u1 := &User1{
        //Name: "userName",
        //Player: p,
    //}
    //t.Logf("playerName = %v\n", u1.playerName)

    //u.setPlayerName()
    //t.Logf("playerName = %v\n", u.playerName)
    //t.Logf("playerName = %v\n", u1.playerName)

    //u := &User{1, "kaka", "hk"}
    //modelType := reflect.ValueOf(u).Type()
    //t.Logf("modelType === %T = %v", modelType, modelType)                                 //  *mysql.User
    //t.Logf("modelType.Kind() === %T = %v", modelType.Kind(), modelType.Kind())            // 类型：reflect.Slice reflect.Array reflect.Ptr(指针) reflect.Struct
    //t.Logf("modelType.Elem() === %T = %v", modelType.Elem(), modelType.Elem())            // 元素：mysql.User
    //t.Logf("modelType.Name() === %T = %v", modelType.Name(), modelType.Name())            // 类型：string
    //t.Logf("modelType.PkgPath() === %T = %v", modelType.PkgPath(), modelType.PkgPath())   // 类型：string

    
    //var str *[]string
    //str = &[]string{}
    //t.Logf("%T=%v", str, str)

    //t.Logf("%T = %v", [...]int{1, 2, 3,4, 5}, [...]int{})

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
    sqlStr = db.Select("id", "name").From("user").Execute()
    t.Logf("sqlStr = %v", sqlStr)
    //sqlStr = db.Select("id", "name").From("user").Join("player", "LEFT").Compile()
    //t.Logf("sqlStr = %v", sqlStr)

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
