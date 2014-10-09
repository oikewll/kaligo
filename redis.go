package epooll

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
)

func NewRedisPool() *redis.Pool {
    fmt.Println("newRedisPool")
    return &redis.Pool{
        MaxIdle: 80,
        MaxActive: 12000, // max number of connections
        Dial: func() (redis.Conn, error) {
            conf := InitConfig()
            host := conf.GetValue("redis", "host")
            port := conf.GetValue("redis", "port")

            c, err := redis.Dial("tcp", host+":"+port)
            //if err != nil {
                //panic(err.Error())
            //}
            return c, err
        },
    } 
}

// 不是当前package 的，每次都会重新初始化
//redisDB := epooll.NewRedisPool().Get()
var RedisConn = NewRedisPool().Get()
// 其他package的，用下面调用
//redisDB := RedisConn

