package config

// import (
//     "github.com/owner888/kaligo/config"
// )
//
// func init() {
//     config.Set("database", config.StrMap{
//         "mysql": map[string]any{
//             // 数据库连接信息
//             "host":    "127.0.0.1",
//             "port":    "3306",
//             "name":    "chrome",
//             "user":    "root",
//             "pass":    "root",
//             "charset": "utf8mb4",
//             "loc":     "Asia/Shanghai",
//             // 安全性配置
//             "table_prefix":    "",
//             "check_privilege": true,
//             "crypt_key":       "",
//             "crypt_fields": map[string][]string{
//                 "user":   {"name", "age"},
//                 "player": {"nickname"},
//             },
//             // 连接池配置
//             "max_idle_connections": 300,
//             "max_open_connections": 25,
//             "max_life_seconds":     5 * 60,
//             // 慢查询日志
//             "log_slow_query": true,
//             "log_slow_time":  1, // 记录慢查询时间，单位：秒
//             "init_cmds": []string{
//                 "SET NAMES utf8mb4",
//             },
//         },
//         "sqlite": map[string]any{
//             "file": "./test.db",
//         },
//     })
// }
