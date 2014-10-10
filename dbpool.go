package epooll

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
)

func newRedisPool() *redis.Pool {
    fmt.Println("初始化 Redis 连接池")
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
//redisDB := epooll.NewRedisPool()

// 其他package的，用下面调用
//redisDB := RedisConn.Get()

func newMysqlPool() *ConnPool {
    fmt.Println("初始化 Mysql 连接池")
    return &ConnPool{
        MaxActive: 3,
        //Dial: func() (*autorc.Conn, error) {
        Dial: func() (interface{}, error) {
            conf := InitConfig()
            host := conf.GetValue("db", "host")
            port := conf.GetValue("db", "port")
            user := conf.GetValue("db", "user")
            pass := conf.GetValue("db", "pass")
            name := conf.GetValue("db", "name")

            //conn := autorc.New("tcp", "", "localhost:3306", "root", "root", "test")
            //conn.Register("set names utf8")
            db, err := InitDB(host+":"+port, user, pass, name)
            return db, err
        },
    } 
}

var RedisConn = newRedisPool()
var MysqlConn = newMysqlPool()

