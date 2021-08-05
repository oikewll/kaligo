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
    "fmt"
    "time"
    "regexp"
    //"reflect"
    "strings"
    "sync"
    //"github.com/owner888/mymysql/autorc"
    //"github.com/owner888/mymysql/mysql"
    ////_ "github.com/ziutek/mymysql/native"    // 普通模式
    //_ "github.com/ziutek/mymysql/thrsafe" // 用了连接池之后连接都是重复利用的，没必要用线程安全模式
)

// Connection is the struct for MySQL connection handler
type Connection struct {
    //readonlys map[string]*Connection
    //instances map[string]*Connection

    active string       // instance name
    dbDsn  string       // Database Dsn URL

    queryCount int
    lastQuery  string

    timeout time.Duration   // Timeout for connect SetConnMaxLifetime(timeout)
    lastUse time.Time       // The last use time

    inTransaction bool
    //db    *autorc.Conn    // MySQL connection
    db          *sql.DB     // MySQL connection
    tx          *sql.Tx     // MySQL connection for Transaction
	autoCommit   bool
    initCmds     []string   // MySQL commands/queries executed after connect
    cacheStore   *sync.Map
    tablePrefix  string

    schema      *Schema
    Debug        bool       // Debug logging. You may change it at any time.
    logSlowQuery bool
    logSlowTime  int
}

// NewConnection 实例化数据库连接
// (读+写)连接数据库+选择数据库
func NewConnection(name string, dbDsn string, writable bool) *Connection {
    c := &Connection{
        dbDsn      : dbDsn,
        queryCount : 0,
        tablePrefix: "",
    }
    //c.Connect()
    return c
}

// DB Returns *sql.DB
func (c *Connection) DB() *sql.DB {
    return c.db
}

// TX Returns *sql.TX
func (c *Connection) TX() *sql.Tx {
    return c.tx
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
// 数据库的最大连接数，超过请求就只能等待了，所以不要设置，直接用mysql默认100个就好了
func (c *Connection) SetMaxOpenConns(n int) {
    c.db.SetMaxOpenConns(n)
}

// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
// 连接池中的保持连接的最大连接数
func (c *Connection) SetMaxIdleConns(n int) {
    c.db.SetMaxIdleConns(n)
}

// SetConnMaxIdleTime sets the maximum amount of time a connection may be reused.
func (c *Connection) SetConnMaxIdleTime(d time.Duration) {
    c.db.SetConnMaxIdleTime(d)
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
// 设置为0的话意味着没有最大生命周期，连接总是可重用(默认行为)
func (c *Connection) SetConnMaxLifetime(d time.Duration) {
    c.db.SetConnMaxLifetime(d)
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

// Stats Returns database statistics
// 待实现...
func (c *Connection) Stats() (err error) {
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
func (c *Connection) Query(queryType int, sqlStr string, asObject interface{}) *Result {
    // var stacktrace []map[string]string   // 储存所有调用函数，一层一层的
    // benchmark := Profiler.Start(c.name, sqlstr, stacktrace)
    // Profiler.Stop(benchmark)

    // golang 好像不需要处理执行时因为mysql链接闲置超过8小时的问题：MySQL server has gone away

    // Set the last query
    c.lastQuery = sqlStr

    r := &Result{
        sqlStr: sqlStr,
        //result              sql.Result      // row result resource
        //totalRows           int             // total number of rows
        //currentRow          int             // current row number
        //asObject            interface{}     // return rows as an object or associative array
        //sanitizationEnabled bool            // If this is a records data will be anitized on get
        //rows                *sql.Rows
    }

    if queryType == SELECT {
        // if Config.Get("enable_cache") { return cached() }
        // return Result{result:sql.Result, sqlstr: sqlstr, asObject: asObject}
        //r.result =
        rows, err := c.db.Query(sqlStr) // 查询多条
        fmt.Printf("Connection.Query() = %v %v\n", rows, err)
    } else if queryType == INSERT {
        // return []string{connection.insertID, connection.affectedRows}
    } else if queryType == UPDATE || queryType == DELETE {
        // return connection.affectedRows
    }

    return r
}
//func (c *Connection) Query(queryType int, sqlStr string, asObject interface{}) (*sql.Rows, error) {
    //c.lastQuery = sqlStr
    //rows, err := c.db.Query(query) // 查询多条
    //return rows, err
//}

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
        sqlStr += "SHOW TABLES LIKE " + c.Quote(like)
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
    table = c.QuoteTable(table)

    var sqlStr string
    if  like != "" {
        sqlStr += "SHOW FULL COLUMNS FROM " + table + " LIKE " + c.Quote(like)
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
    table = c.QuoteTable(table)

    var sqlStr string
    if  like != "" {
        sqlStr += "SHOW INDEX FROM " + table + " WHERE " + c.QuoteIdentifier("Key_name") + " LIKE " + c.Quote(like)
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

// LastQuery Returns the last query sql
func (c *Connection) LastQuery() string {
    return c.lastQuery
}
// UserInfo is a test
type UserInfo struct  {
    ID int
    Name string
}

// TablePrefix Return the table prefix defined in the current configuration.
func (c *Connection) TablePrefix(table string) string {
    // 循环输出 stuct 结构
    //t := reflect.TypeOf(c)
    //v := reflect.ValueOf(c)
    //for k := 0; k < t.NumField(); k++ {
        //fmt.Printf("%s -- %v \n", t.Field(k).Name, v.Field(k).Interface())   
    //}

    //fmt.Printf("TablePrefix === %T = %p\n", c, c)
    return c.tablePrefix + table
}

// Quote a value for an SQL query.
func (c *Connection) Quote(values interface{}) string {
    switch vals := values.(type) {
    case string:
        return c.Escape(vals)
    case []string:
        for k, v := range vals {
            vals[k] = c.Escape(v)
        }
        return "(" + strings.Join(vals, ", ") + ")"
    case *Query:
        // Create a sub-query
        return "(" + vals.Compile() + ")"
    default:
        return vals.(string)
    }
}

// QuoteTable Quote a database table name and adds the table prefix if needed.
//table = strings.Replace(table, "#DB#", "lrs", 1 )
// 表名添加引用符号(`)
// 添加表前缀
func (c *Connection) QuoteTable(table string) string {
    table = c.TablePrefix(table)
    table = c.QuoteIdentifier(table)
    return table
}

// QuoteIdentifier Quote a database identifier, such as a column name. Adds the
// table prefix to the identifier if a table name is present.
// 字段名添加引用符号(`)
func (c *Connection) QuoteIdentifier(values interface{}) string {
    switch vals := values.(type) {
    case string:
        if vals == "*" || strings.Index(vals, "`") != -1 {
            // * 不需要变成 `*`，已经有 `` 包含着的直接返回
            return vals
        } else if strings.Index(vals, ".") != -1 {
            // table.column 的写法，变成 `table`.`column`
            parts := regexp.MustCompile(`\.`).Split(vals, 2)
            return c.QuoteIdentifier( c.QuoteTable(parts[0]) ) + "." + c.QuoteIdentifier(parts[1])
        } else {
            return "`" + vals + "`"
        }
    case []string:
        // Separate the column and alias
        value := vals[0]
        alias := vals[1]
        return c.QuoteIdentifier(value) + " AS " + c.QuoteIdentifier(alias)
    default:
        return vals.(string)
    }
}

// Escape is use for Escapes special characters in the txt, so it is safe to place returned string
func (c *Connection) Escape(sql string) string {
    dest := make([]byte, 0, 2*len(sql))
    var escape byte
    for i := 0; i < len(sql); i++ {
        c := sql[i]

        escape = 0

        switch c {
        case 0: /* Must be escaped for 'mysql' */
            escape = '0'
            break
        case '\n': /* Must be escaped for logs */
            escape = 'n'
            break
        case '\r':
            escape = 'r'
            break
        case '\\':
            escape = '\\'
            break
        case '\'':
            escape = '\''
            break
        case '"': /* Better safe than sorry */
            escape = '"'
            break
        case '\032': //十进制26,八进制32,十六进制1a, /* This gives problems on Win32 */
            escape = 'Z'
        }

        if escape != 0 {
            dest = append(dest, '\\', escape)
        } else {
            dest = append(dest, c)
        }
    }

    // SQL standard is to use single-quotes for all values
    return "'" + string(dest) + "'"
}

// StartTransaction is ...
func (c *Connection) StartTransaction() bool {
    c.inTransaction = true
    return true
}

// CommitTransaction is ...
func (c *Connection) CommitTransaction() bool {
    c.inTransaction = false
    return true
}

// RollbackTransaction is ...
func (c *Connection) RollbackTransaction(rollbackAll bool) bool {
    c.inTransaction = false
    return true
}
