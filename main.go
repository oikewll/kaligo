package main

import (
    //"epooll/core"
    "github.com/widuu/goini"
    "fmt"
)

func main() {

    conf := goini.SetConfig("/data/golang/src/epooll/conf/app.ini")
    username := conf.GetValue("database", "username") //database是你的[section]，username是你要获取值的key名称
    fmt.Println(username) //root
    //init := core.Init{}
    //init.InitServer()
    //或者
    //reflect.ValueOf(&core.Init{}).MethodByName("InitServer").Call(nil)
    //下面是错误的写法
    //core.Init.InitServer()
}

