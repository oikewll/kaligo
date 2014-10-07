package epooll

import (
    "github.com/garyburd/redigo/redis"
)

func InitRedisPool() *redis.Pool {
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

