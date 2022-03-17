package config

// import (
//     "fmt"
//     "github.com/owner888/kaligo/config"
// )
//
// func init() {
//     config.Add("database", config.StrMap{
//         "mysql": map[string]any{
//             // 数据库连接信息
//             "open"   : true,
//             "host"   : "127.0.0.1",
//             "port"   : "3306",
//             "name"   : "test",
//             "user"   : "root",
//             "pass"   : "root",
//             "charset": "utf8mb4",
//             "loc"    : "Asia/Shanghai",
//             // 安全性配置
//             "table_prefix" : "",
//             "check_privilege" : true,
//             "crypt_key"    : "",
//             "crypt_fields" : map[string][]string{
//                 "user"  : {"name", "age"},
//                 "player": {"nickname"},
//             },
//             // 连接池配置
//             "max_idle_connections": 300,
//             "max_open_connections": 25,
//             "max_life_seconds":     5*60,
//             // 慢查询日志
//             "log_slow_query": true,
//             "log_slow_time" : 1,    // 记录慢查询时间，单位：秒
//         },
//     })
//
//     user := config.Get[string]("database.mysql.user")
//     pass := config.Get[string]("database.mysql.pass")
//     host := config.Get[string]("database.mysql.host")
//     port := config.Get[string]("database.mysql.port")
//     name := config.Get[string]("database.mysql.name")
//     dsn  := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", user, pass, host+":"+port, name, "utf8mb4")
//     config.Set("database.mysql.dsn", dsn)
// }
