package epooll

import(
    "testing"
    "fmt"
    //"reflect"
    "github.com/ziutek/mymysql/autorc" 
    //"github.com/ziutek/mymysql/mysql" 
    //"github.com/ziutek/mymysql/native" // Native engine 
    // _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine 
)

//func initConn() (*autorc.Conn, error) {
func initConn() (interface{}, error) {
    conn := autorc.New("tcp", "", "localhost:3306", "root", "root", "test")
    // Register initialisation commands
	conn.Register("set names utf8")
    fmt.Println(" --- autorc --- ", conn.Raw, " --- ")
    //rows, res, err := db.Query("Select * From `test` Limit 1")
    //fmt.Println(rows)
    //fmt.Println(res)
    //fmt.Println(err)
    return conn, nil
}

var dbPool = new(ConnPool)
var err = dbPool.InitPool(3, initConn)

func Do() {

    db := dbPool.Get().(*autorc.Conn)
    //fmt.Println(db)
    defer dbPool.Release(db)
    _, res, _ := db.Query("Select * From `test` Limit 1")
    fmt.Println(res)
}

func Test_Pool(t *testing.T) {
    Do()
    //for i := 0; i <= 1000; i++ {
        //Do()
    //}

    //count := len(dbPool.conn)
    //for i := 0; i < count; i++ {
        //fmt.Println(<-dbPool.conn)
    //}
}
/*func Test_DB(t *testing.T) {
    //db := &DB{} 
    db := InitDB()
    
    data := map[string]string {
        "name": "nam'e111",
        "pass": "pas\"s111",
        //"sex":"111",
    }
    //if ok, err := db.Update("test", data, "Where `id`=1"); !ok {
        //fmt.Println(err)
    //} else {
        //fmt.Println(db.AffectedRows())
    //}

    if ok, err := db.Insert("test", data); !ok {
        fmt.Println(err)
    } else {
        fmt.Println(db.InsertId())
    }

    //sql := "Select a.id,a.name,b.date From `test` a Left Join `test_part` b On a.id=b.test_id; "
    //row, _ := db.GetOne(sql)
    //fmt.Println(row["id"], " --- ", row["name"])

    //for _, v := range results {
        //fmt.Println("id = ", v["id"], " --- name = ", v["name"])
    //}
    //fmt.Println(results)
    //fmt.Println(rows)

    //log.Print(db, err)

}*/
