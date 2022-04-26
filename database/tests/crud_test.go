package tests

import (
    "testing"
)

type User struct {
    ID   uint   `db:"id"`
    Name string `db:"name"`
    Age  uint   `db:"age"`
    Sex  uint   `db:"sex"`
}

func TestCRUD(t *testing.T) {
    //var sqlStr string
    //
    // user := User{
    //     ID   : 1,
    //     Name : "test111",
    //     Age  : 25,
    //     Sex  : 1,
    // }
    // user.DB = db
    // user.Save()

    //var ages []int64
    //q := db.Query("SELECT age FROM user").Scan(&ages).Execute()
    //if q.Error != nil {
        //t.Logf("q.Error = %v\n", q.Error)
    //} else {
        //t.Logf("jsonStr = %v\n", database.FormatJSON(ages))
    //}
}
