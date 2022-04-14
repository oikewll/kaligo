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

func (r *Redis) String(key string) string {
    reply, found := r.Get(key)
    if !found {
        return ""
    }
    return reply.(string)
}

func (r *Redis) Int(key string) int {
    reply, found := r.Get(key)
    if !found {
        return 0
    }
    return reply.(int)
}

func (r *Redis) Int64(key string) int64 {
    reply, found := r.Get(key)
    if !found {
        return 0
    }
    return reply.(int64)
}

func (r *Redis) Uint(key string) uint {
    reply, found := r.Get(key)
    if !found {
        return 0
    }
    return reply.(uint)
}

func (r *Redis) Uint64(key string) uint64 {
    reply, found := r.Get(key)
    if !found {
        return 0
    }
    return reply.(uint64)
}

func (r *Redis) Float64(key string) float64 {
    reply, found := r.Get(key)
    if !found {
        return 0
    }
    return reply.(float64)
}

func (r *Redis) Bool(key string) bool {
    reply, found := r.Get(key)
    if !found {
        return false
    }
    return reply.(bool)
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

func (r *Redis) Incr(key string, args ...uint64) int64 {
    conn := r.conn.Get()
    defer conn.Close()

    val, err := redis.Int64(conn.Do("INCR", key))
    if err != nil {
        return 0
    }

    return val
}

func (r *Redis) Decr(key string, args ...uint64) int64 {
    conn := r.conn.Get()
    defer conn.Close()

    val, err := redis.Int64(conn.Do("DECR", key))
    if err != nil {
        return 0
    }

    return val
}

// push insert a value to List
func (r *Redis) push(command, key, value string) {
    conn := r.conn.Get()
    defer conn.Close()

    _, _ = conn.Do(command, key, value)
}

// RPop return a value from List
func (r *Redis) pop(command, key string) string {
    conn := r.conn.Get()
    defer conn.Close()

    reply, err := redis.String(conn.Do(command, key))
    if err != nil {
        return ""
    }

    return reply
}

func (r *Redis) LPush(key string, value any) {
    var valueStr string
    switch v := value.(type) {
    case string:
        valueStr = v
    default:
        jsonStr, _ := json.Marshal(v)
        valueStr = string(jsonStr)
    }
    r.push("LPUSH", key, valueStr)
}

func (r *Redis) RPush(key string, value any) {
    var valueStr string
    switch v := value.(type) {
    case string:
        valueStr = v
    default:
        jsonStr, _ := json.Marshal(v)
        valueStr = string(jsonStr)
    }
    r.push("RPUSH", key, valueStr)
}

// user := &User{}
// LPop("user", &user)
func (r *Redis) LPop(key string, value any) (err error) {
    valueStr := r.pop("LPOP", key)
    switch v := value.(type) {
    case *string:
        value = v
    default:
        err = json.Unmarshal([]byte(valueStr), value)
    }
    return
}

func (r *Redis) RPop(key string, value any) (err error) {
    valueStr := r.pop("RPOP", key)
    switch v := value.(type) {
    case *string:
        value = v
    default:
        err = json.Unmarshal([]byte(valueStr), value)
    }
    return
}
