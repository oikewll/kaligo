package config

// import (
//     "github.com/owner888/kaligo/config"
// )
//
// func init() {
//     config.Add("cache", config.StrMap{
//         "redis": map[string]any{
//             "host"          : "127.0.0.1",
//             "port"          : 6379,
//             "password"      : "",
//             "database"      : 0,
//             "max_idle"      : 10,        // 最大空闲连接数
//             "max_active"    : 200,       // 最大连接数
//             "idle_timeout"  : 180,       // 空闲连接超时时间(单位: 秒)
//             "wait"          : true,      // 超过最大连接数的操作:等待
//         },
//         "memcache": map[string]any{
//             "host": []string{"10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11211"},
//         },
//     })
// }