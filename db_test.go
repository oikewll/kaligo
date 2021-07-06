package kaligo

import(
    "testing"
    "fmt"
    //"strconv"
    //"reflect"
    "database/sql"
    //"kaligo/conf"
    //"kaligo/util"
    //_ "github.com/go-sql-driver/mysql"
    //"github.com/ziutek/mymysql/autorc" 
    //"github.com/ziutek/mymysql/mysql" 
    //"github.com/ziutek/mymysql/native" // Native engine 
    // _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine 
)

//func initConn() (*autorc.Conn, error) {
func initConn() (interface{}, error) {
    //conn := autorc.New("tcp", "", "localhost:3306", "root", "root", "test")
    //conn.Register("set names utf8")

    conn, err := sql.Open("mysql", "root:root@localhost:3306/test")
    checkErr(err)
    //fmt.Println(" --- autorc --- ", conn.Raw, " --- ")
    //rows, res, err := db.Query("Select * From `test` Limit 1")
    //fmt.Println(rows, res, err)
    return conn, nil
}


func Test_DB(t *testing.T) {

    db, err := New()

    defer db.Close()

    if err != nil {
        t.Fatal(err)
    }

    //data := map[string]string {
        //"name": "nam'e111",
        //"pass": "pas\"s111",
        //"sex":"1",
    //}
    //t.Fatal(data)   // 输出并中断程序

    //if ok, err := db.Update("user", data, "`id`=1"); !ok {
        //fmt.Println(err)
    //} else {
        //fmt.Println(db.AffectedRows())
    //}

    //if ok, err := db.Insert("user", data); !ok {
        //fmt.Println(err)
    //} else {
        //fmt.Println(db.InsertID())
    //}

    //sql := "Select a.id,a.name,b.date From `test` a Left Join `test_part` b On a.id=b.test_id;"
    sql := "Select * From user"

    //row, _ := db.GetOne(sql)
    //fmt.Println(row)
    //fmt.Println(row["id"], " --- ", row["name"])

    rows, _ := db.GetAll(sql)
    //fmt.Println(rows)
    for _, v := range rows {
        fmt.Println("id = ", v["id"], " --- name = ", v["name"])
    }

    //log.Print(db, err)

}
