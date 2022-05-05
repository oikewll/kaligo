package tests

import (
    "testing"
    // "github.com/owner888/kaligo/logs"
    // "github.com/stretchr/testify/assert"
)


type User struct {
    ID   uint   `db:"id"`
    Name string `db:"name"`
    Age  uint   `db:"age"`
    Sex  uint   `db:"sex"`
}

func TestSave(t *testing.T) {
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
}


