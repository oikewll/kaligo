/**
 * @Description: Connection
 * @Version: 1.0.0
 * @Author: Seatle
 * @Date: 2014-07-06 21:39
 * @LastEditors: Seatle
 * @LastEditTime: 2021-07-06 21:39
 * @Other: 实现一个Mysql 的 Profilling: https://www.cnblogs.com/sxdcgaq8080/p/11844079.html
 */

package mysql

import (
    "database/sql"
    //"fmt"
    //"log"
    "time"
    //"regexp"
    //"strings"
    //"sort"
    //"container/list"
    //"sync"
    //"reflect"
    //"github.com/owner888/kaligo/conf"
    //"github.com/owner888/kaligo/util"
    //_ "github.com/go-sql-driver/mysql"  // 空白导入必须在main.go、testing，否则就必须在这里写注释
    //"github.com/owner888/mymysql/autorc"
    //"github.com/owner888/mymysql/mysql"
    ////_ "github.com/ziutek/mymysql/native"    // 普通模式
    //_ "github.com/ziutek/mymysql/thrsafe" // 用了连接池之后连接都是重复利用的，没必要用线程安全模式
)

type serverInfo struct {
    protVer byte
    servVer []byte
    thrID   uint32
    scramble [20]byte
    caps     uint32
    lang     byte
    plugin   []byte
}

// Connection is the struct for MySQL connection handler
type Connection struct {
    //readonlys map[string]*Connection 
    //instances map[string]*Connection 

    active string       // instance name
    dbDsn  string       // Database Dsn URL
    dbIdle int          // Database Max Idle Conns

    info  serverInfo    // MySQL server information
    seq   byte          // MySQL sequence number

    timeout time.Duration   // Timeout for connect SetConnMaxLifetime(timeout)
    lastUse time.Time       // The last use time

    //db    *autorc.Conn    // MySQL connection
    db          *sql.DB     // MySQL connection
    tx          *sql.Tx     // MySQL connection for Transaction
	autoCommit   bool
    initCmds     []string   // MySQL commands/queries executed after connect

    schema      *Schema
    Debug        bool       // Debug logging. You may change it at any time.
    logSlowQuery bool
    logSlowTime  int
}

// NewConnection 实例化数据库连接
// (读+写)连接数据库+选择数据库
func NewConnection(name string, dbDsn string, dbIdle int, writable bool) *Connection {
    c := &Connection{
        dbDsn:  dbDsn,
        dbIdle: dbIdle,
    }
    c.Connect()
    return c
}

// Connect is Establishes a connection with MySQL server version 4.1 or later.
func (c *Connection) Connect() (err error) {
    defer catchError(&err)

    // 当前函数推出时不能关闭链接，否则无法调用 Query()、Execute() ，应该在调用 Execute()后close，释放链接到链接池去
    //defer c.db.Close()

    // 数据库的抽象（*sql.DB），并不是真的数据库连接
    //c.db := autorc.New("tcp", "", address, user, pass, name)
    c.db, err = sql.Open("mysql", c.dbDsn)
    if err != nil {
        c.db = nil
        return ErrAuthentication
    }

    // See "Important settings" section.
    //c.db.SetConnMaxLifetime(0)      // 设置为0的话意味着没有最大生命周期，连接总是可重用(默认行为)。
    //c.db.SetMaxOpenConns(0)         // 数据库的最大连接数，超过请求就只能等待了，所以不要设置，直接用mysql默认100个就好了
    c.db.SetMaxIdleConns(c.dbIdle)    // 连接池中的保持连接的最大连接数

    // 初始化一个数据库连接，sql.Open 的时候实际上是返回一个数据库的抽象而已，并没有真的和mysql链接上
    err = c.db.Ping()
    if err != nil {
        c.db = nil
        return ErrAuthentication
    }

    c.db.Query("SET NAMES utf8");

    return
}

func (c *Connection) close() (err error) {
    defer catchError(&err)

    // Always close and invalidate connection
    defer func() {
        err = c.db.Close()
        c.db = nil // Mark that we disconnect
    }()

    return
}

// Close connection to the server
func (c *Connection) Close() (err error) {
    if c.db == nil {
        return ErrNotConn
    }

    err = c.close()
    return err
}

// Caching Per connection cache controller setter/getter
func (c *Connection) Caching() bool {
    //return \Arr::get($this->_config, 'enable_cache', true);
    return false
}

// QueryRow is the function for query one row
func (c *Connection) QueryRow(query string) *sql.Row {
    row := c.db.QueryRow(query, 1) // 查询一条
    return row
}

// Query is the function for query multi rows
func (c *Connection) Query(query string) (*sql.Rows, error) {
    rows, err := c.db.Query(query) // 查询多条
    return rows, err
}

// Exec is the function for Insert、Update、Delete
func (c *Connection) Exec(query string) (sql.Result, error) {
    res, err := c.db.Exec(query)
    return res, err
}

// ListTables If a table name is given it will return the table name with the configured
// prefix. If not, then just the prefix is returnd
func (c *Connection) ListTables(like string) []string {
    var sqlStr string    
    if  like != "" {
        sqlStr += "SHOW TABLES LIKE " + quote(like)
    } else {
        sqlStr += "SHOW TABLES"
    }

    var tables []string    
    tables = append(tables, "111")
    return tables
}

// ListColumns Lists all of the columns in a table. Optionally, a LIKE string can be
// used to search for specific fields.
func (c *Connection) ListColumns(table string, like string) map[string] map[string]string {
    table = quoteTable(table)

    var sqlStr string    
    if  like != "" {
        sqlStr += "SHOW FULL COLUMNS FROM " + table + " LIKE " + quote(like)
    } else {
        sqlStr += "SHOW FULL COLUMNS FROM " + table
    }

    var column map[string]string
    column["name"] = "Field"
    var columns map[string] map[string]string
    columns["Field"] = column

    return columns
}

// ListIndexes Lists all of the idexes in a table. Optionally, a LIKE string can be
// used to search for specific indexes by name.
func (c *Connection) ListIndexes(table string, like string) []map[string]string {
    table = quoteTable(table)

    var sqlStr string    
    if  like != "" {
        sqlStr += "SHOW INDEX FROM " + table + " WHERE " + quoteIdentifier("Key_name") + " LIKE " + quote(like)
    } else {
        sqlStr += "SHOW INDEX FROM " + table 
    }

    var indexes []map[string]string    
    mapName := map[string]string {
        "name"      : "Key_name",
        "column"    : "Column_name",
        "order"     : "Seq_in_index",
        "type"      : "Index_type",
        "primary"   : "true",
        "unique"    : "Non_unique",
        "null"      : "YES",
        "ascending" : "Collation",
    }
    indexes = append(indexes, mapName)

    return indexes
}

