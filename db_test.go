package epooll

import(
    "testing"
    "fmt"
    //"reflect"
)

func Test_DB(t *testing.T) {
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

}
