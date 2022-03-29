package cache

import (
    "encoding/json"
    "time"

    "github.com/gomodule/redigo/redis"
)

//Redis redis cache
type Redis struct {
    conn *redis.Pool
}

//RedisOpts redis 连接属性
type RedisOpts struct {
    Host        string `yml:"host" json:"host"`
    Password    string `yml:"password" json:"password"`
    Database    int    `yml:"database" json:"database"`
    MaxIdle     int    `yml:"max_idle" json:"max_idle"`
    MaxActive   int    `yml:"max_active" json:"max_active"`
    IdleTimeout int    `yml:"idle_timeout" json:"idle_timeout"`
    Wait        bool   `yml:"wait" json:"wait"`
}

// NewRedis 实例化
func NewRedis(opts *RedisOpts) *Redis {
    pool := &redis.Pool{
        MaxIdle:     opts.MaxIdle,                                  // 最大空闲连接数
        MaxActive:   opts.MaxActive,                                // 最大连接数
        IdleTimeout: time.Second * time.Duration(opts.IdleTimeout), // 空闲连接超时时间(单位: 秒)
        Wait:        opts.Wait,                                     // 超过最大连接数的操作:等待
        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", opts.Host,
                redis.DialDatabase(opts.Database),
                redis.DialPassword(opts.Password),
            )
            if err != nil {
                return nil, err
            }
            return c, nil
        },
        TestOnBorrow: func(conn redis.Conn, t time.Time) error {
            if time.Since(t) < time.Minute {
                return nil
            }
            _, err := conn.Do("PING")
            return err
        },
    }
    return &Redis{pool}
}

// Get 获取一个值
func (r *Redis) Get(key string) (any, bool) {
    var err error    
    var reply any    
    conn := r.conn.Get() // 从连接池中获取一个链接
    defer conn.Close()

    var data []byte
    if data, err = redis.Bytes(conn.Do("GET", key)); err != nil {
        return nil, false
    }
    if err = json.Unmarshal(data, &reply); err != nil {
        return nil, false
    }

    return reply, true
}

// Set 设置一个值
func (r *Redis) Set(key string, val any, timeout time.Duration) (err error) {
    conn := r.conn.Get()
    defer conn.Close()

    var data []byte
    if data, err = json.Marshal(val); err != nil {
        return
    }

    _, err = conn.Do("SETEX", key, int64(timeout/time.Second), data)

    return
}

// Has 判断key是否存在
func (r *Redis) Has(key string) bool {
    conn := r.conn.Get()
    defer conn.Close()

    a, _ := conn.Do("EXISTS", key)
    i := a.(int64)
    return i > 0
}

// Delete 删除
func (r *Redis) Del(key string) error {
    conn := r.conn.Get()
    defer conn.Close()

    if _, err := conn.Do("DEL", key); err != nil {
        return err
    }

    return nil
}

func (mem *Redis) Incr(key string) int64 {
    return 0
}

func (mem *Redis) Decr(key string) int64 {
    return 0
}

func (r *Redis) GetAnyKeyValue(key string, defaultValue ...any) (val any, found bool) {
    val, found = r.Get(key)
    if !found {
        if len(defaultValue) != 0 {
            val = defaultValue[0]
        }
    }
    return
}
