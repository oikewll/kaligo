package kaligo

//import (
    //"fmt"
    //"strconv"
    //"kaligo/conf"
    //"github.com/gomodule/redigo/redis"
//)

//func newRedisPool() *redis.Pool {
    //poolNum, _ := strconv.Atoi(conf.GetValue("pool", "redis")) 
    //fmt.Printf("初始化 Redis 连接池，连接数：%d \n", poolNum)
    //return &redis.Pool{
        //MaxIdle: 80,
        //MaxActive: poolNum, // max number of connections
        //Dial: func() (redis.Conn, error) {
            //host := conf.GetValue("redis", "host")
            //port := conf.GetValue("redis", "port")

            //c, err := redis.Dial("tcp", host+":"+port)
            ////if err != nil {
                ////panic(err.Error())
            ////}
            //return c, err
        //},
    //} 
//}

// 不是当前package 的，每次都会重新初始化
//redisDB := kaligo.NewRedisPool()

// 其他package的，用下面调用
//redisDB := RedisConn.Get()

//func newMysqlPool() *ConnPool {
    //poolNum, _ := strconv.Atoi(conf.GetValue("pool", "mysql")) 
    //fmt.Printf("初始化 Mysql 连接池，连接数：%d \n", poolNum)
    //return &ConnPool{
        //MaxActive: poolNum,
        ////Dial: func() (*autorc.Conn, error) {
        //Dial: func() (interface{}, error) {
            ////conn := autorc.New("tcp", "", "localhost:3306", "root", "root", "test")
            ////conn.Register("set names utf8")
            //db, err := InitDB()
            //return db, err
        //},
    //} 
//}

// RedisConn is ...
//var RedisConn = newRedisPool()
// MysqlConn is ...
//var MysqlConn = newMysqlPool()

