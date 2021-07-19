package kaligo

import (
    "fmt"
    "testing"
    "github.com/owner888/kaligo/conf"
)

func Test(t *testing.T) {
    conf.InitConfig("../conf/conf.ini")
    user := conf.Get("db", "user")
    fmt.Printf("Test user: %v\n", user)
    conf.Delete("db", "user")
    user = conf.Get("db", "user")
    if len(user) == 0 {
        fmt.Println("username is not exists") //this stdout username is not exists
    }
    conf.Set("db", "user", "widuu")
    user = conf.Get("db", "user")
    fmt.Println(user) //widuu

    //data := conf.ReadList()
    //fmt.Println(data)
}
