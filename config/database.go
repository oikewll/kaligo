package config
// Global variable

import (
    "fmt"
)

// TablePrefix is database table prefix
var (
    TablePrefix  string
    MaxOpenConns int
    MaxIdleConns int

    // CryptKey is a mysql AES_ENCRYPT、AES_DECRYPT crypt key
    CryptKey    string
    // CryptFields is mysql AES_ENCRYPT、AES_DECRYPT crypt fields
    CryptFields map[string][]string
)

// DBUser is ...
var (
    DBUser string
    DBPass string
    DBhost string
    DBPort string
    DBName string
    DBDSN  string
    CheckDBPrivilege bool
)

func init() {

    // --------------------------
    // 数据库基础配置
    // --------------------------
    TablePrefix  = ""
    MaxOpenConns = 100
    MaxIdleConns = 16
    CryptKey     = "tPVPnynVnsiqh"
    CryptFields  = map[string][]string {
        "user"  : {"name", "age"},
        "player": {"nickname"},
    }

    // --------------------------
    // 数据库链接配置
    // --------------------------
    DBUser = "root"
    DBPass = "root"
    DBhost = "127.0.0.1"
    DBPort = "3306"
    DBName = "test"
    DBDSN  = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", DBUser, DBPass, DBhost+":"+DBPort, DBName, "utf8mb4")

    // 是否检查数据库权限，show databases 出现 mysql、sys、information_schema、performance_schema 数据库则报错
    CheckDBPrivilege = false
}
