package epooll

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
    "github.com/ziutek/mymysql/autorc" 
)

func newRedisPool() *redis.Pool {
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
//redisDB := epooll.NewRedisPool()

// 其他package的，用下面调用
//redisDB := RedisConn.Get()

func newMysqlPool() *ConnPool {
    fmt.Println("newMysqlPool")
    return &ConnPool{
        MaxActive: 2,
        //Dial: func() (*autorc.Conn, error) {
        Dial: func() (interface{}, error) {
            conn := autorc.New("tcp", "", "localhost:3306", "root", "root", "test")
            // Register initialisation commands
            conn.Register("set names utf8")
            return conn, nil
        },
    } 
}

var RedisConn = newRedisPool()
var MysqlConn = newMysqlPool()

