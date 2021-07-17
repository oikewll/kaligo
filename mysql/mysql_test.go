package mysql
// go test -v -count=1 mysql/mysql_test.go

import(
    "testing"
    "os"
    "log"
    //"strconv"
    //"reflect"
    //"encoding/json"
    //"database/sql"
    //_ "github.com/go-sql-driver/mysql"
    //"github.com/owner888/kaligo"
    "github.com/owner888/kaligo/conf"
    //"github.com/owner888/kaligo/util"
    "github.com/owner888/kaligo/mysql"
    //"github.com/owner888/kaligo/cache"
    //"github.com/ziutek/mymysql/autorc" 
    //"github.com/ziutek/mymysql/mysql" 
    //"github.com/ziutek/mymysql/native" // Native engine 
    // _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine 
)

func TestDB(t *testing.T) {

    //memory := cache.NewMemory()
    //memory.Set("name", "kaka", 10)
    //name := memory.Get("name")
    //t.Logf("name: [ %v ]", name)
    str, _ := os.Getwd()
    conf.AppPath = str + "/../"
    log.Printf("TestDB AppPath: [ %v ]", conf.AppPath)
    
    //define('APPPATH', __DIR__.'/../app');

    //t.Logf("%s", kaligo.Int64ToStr(100))

    db := mysql.NewDB()
    //db := NewDB()
    //db.Debug = false

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
    sql := "Select * From user"

    row, _ := db.GetOne(sql)
    t.Logf("%v", row)
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
