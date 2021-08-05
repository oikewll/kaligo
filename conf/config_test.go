package conf

import (
    "fmt"
    "testing"
)

func Test(t *testing.T) {
    InitConfig("../conf/conf.ini")
    user := Get("db", "user")
    fmt.Printf("Test user: %v\n", user)
    Delete("db", "user")
    user = Get("db", "user")
    if len(user) == 0 {
        fmt.Println("username is not exists") //this stdout username is not exists
    }
    Set("db", "user", "widuu")
    user = Get("db", "user")
    fmt.Println(user) //widuu

    //data := conf.ReadList()
    //fmt.Println(data)
}
