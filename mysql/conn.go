/**
 * Realize a database operation class
 *
 * @copyright   (C) 2014  seatle
 * @lastmodify  2021-07-06
 *
 */

package mysql

import (
    "fmt"
    "log"
    "time"
    //"regexp"
    //"strings"
    //"sort"
    //"container/list"
    //"sync"
    //"reflect"
    "github.com/owner888/kaligo/conf"
    "github.com/owner888/kaligo/util"
    "database/sql"
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

// Conn is the struct for MySQL connection handler
type Conn struct {
    proto string        // Network protocol
    laddr string        // Local address
    raddr string        // Remote (server) address

    dbuser string       // MySQL username
	dbpass string       // MySQL password
	dbname string       // Database name
    plugin string       // Authentication plugin

    dbDsn string        // MySQL Dsn
    info serverInfo     // MySQL server information
	seq  byte           // MySQL sequence number

    // Timeout for connect
    timeout time.Duration

    maxIdleConns int
    maxOpenConns int
    //Conn *autorc.Conn
    //res mysql.Result
    netConn *sql.DB        // MySQL connection
    rows sql.Rows
    res sql.Result
    //row *Row

    // Debug logging. You may change it at any time.
	Debug bool
}

// NewConn 实例化数据库连接 
// (读+写)连接数据库+选择数据库
func NewConn() *Conn {
    c := Conn{
        proto: "tcp",
        laddr: "",
        raddr:  conf.GetValue("db", "host") + ":" + conf.GetValue("db", "port"),
        dbuser: conf.GetValue("db", "user"),
        dbpass: conf.GetValue("db", "pass"),
        dbname: conf.GetValue("db", "name"),
        maxOpenConns: util.StrToInt(conf.GetValue("db", "max_open_conns")),
        maxIdleConns: util.StrToInt(conf.GetValue("db", "max_idle_conns")),
    }

    return &c
}

func (c *Conn) init() {
	c.seq = 0 // Reset sequence number, mainly for reconnect
	if c.Debug {
		log.Printf("[%2d ->] Init packet:", c.seq)
	}
	//Register("SET NAMES utf8")
}

// SetTimeout sets timeout for Connect and Reconnect
func (c *Conn) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// NetConn return internall net.Conn
//func (my *Conn) NetConn() net.Conn {
	//return my.net_conn
//}

func (c *Conn) connect() (err error) {
    defer util.CatchError(&err)
    // 推出这个函数时不能关闭链接，否则其他调用的函数就无法执行 Query()、Exel() 方法 了
    //defer db.Close()

    c.netConn = nil

    c.dbDsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s",
		c.dbuser,
		c.dbpass,
		c.raddr,
		c.dbname,
		"utf8mb4",
	)

    //fmt.Printf("%v", c.dbDsn)

    if c.netConn == nil {
        // 数据库的抽象（*sql.DB），并不是真的数据库连接
        //my.conn := autorc.New("tcp", "", address, user, pass, name) 
        c.netConn, err = sql.Open("mysql", c.dbDsn)
        if err != nil {
            c.netConn = nil
            return
        }
    }

    c.init()

    // See "Important settings" section.
    c.netConn.SetConnMaxLifetime(time.Minute * 3)   // 连接存活 3分钟
    c.netConn.SetMaxOpenConns(c.maxOpenConns)       // 数据库的最大连接数
    c.netConn.SetMaxIdleConns(c.maxIdleConns)       // 连接池中的保持连接的最大连接数

    // 初始化一个数据库连接，sql.Open 的时候实际上是返回一个数据库的抽象而已，并没有真的和mysql链接上
    //err = c.netConn.Ping()
    err = c.Ping()
    if err != nil {
        fmt.Println("连接数据库失败 --- ", err.Error())
        return
    }

    //my.netConn.Query("set names utf8");

    return
}

// Connect is Establishes a connection with MySQL server version 4.1 or later.
func (c *Conn) Connect() (err error) {
    if c.netConn != nil {
        return ErrAlredyConn
    }

    return c.connect()
}

func (c *Conn) close() (err error) {
    defer util.CatchError(&err)

    // Always close and invalidate connection
    defer func() {
        err = c.netConn.Close()
        c.netConn = nil // Mark that we disconnect
    }()

    return
}

// Close connection to the server
func (c *Conn) Close() (err error) {
    if c.netConn == nil {
        return ErrNotConn
    }
    //if c.unreaded_reply {
		//return ErrUnreadedReply
	//}

    err = c.close()
    return err
}

// Reconnect to Close and reopen connection.
// Ignore unreaded rows, reprepare all prepared statements.
func (c *Conn) Reconnect() (err error) {
	if c.netConn != nil {
		// Close connection, ignore all errors
		c.close()
	}
	// Reopen the connection.
	if err = c.connect(); err != nil {
		return
	}

	//// Reprepare all prepared statements
	//var (
		//new_stmt *Stmt
		//new_map  = make(map[uint32]*Stmt)
	//)
	//for _, stmt := range c.stmt_map {
		//new_stmt, err = c.prepare(stmt.sql)
		//if err != nil {
			//return
		//}
		//// Assume that fields set in new_stmt by prepare() are indentical to
		//// corresponding fields in stmt. Why can they be different?
		//stmt.id = new_stmt.id
		//stmt.rebind = true
		//new_map[stmt.id] = stmt
	//}
	//// Replace the stmt_map
	//c.stmt_map = new_map

	return
}

// Use to Change database
func (c *Conn) Use(dbname string) (err error) {
    defer util.CatchError(&err)

    if c.netConn == nil {
        return ErrNotConn
    }
    //if c.unreaded_reply {
		//return ErrUnreadedReply
	//}

    c.dbname = dbname
    return
}

// Ping is Send MySQL PING to the server.
func (c *Conn) Ping() (err error) {
	defer util.CatchError(&err)

	if c.netConn == nil {
		return ErrNotConn
	}
	//if c.unreaded_reply {
		//return ErrUnreadedReply
	//}

	// Send command
	//my.sendCmd(_COM_PING)
	// Get server response
	//my.getResult(nil, nil)

    err = c.netConn.Ping()
    if err != nil {
        return ErrNotConn
        //fmt.Println("连接数据库失败", err.Error())
        //return
    }

	return
}

// QueryRow is the function for query one row
func (c *Conn) QueryRow(query string) (*sql.Row) {
    row := c.netConn.QueryRow(query, 1) // 查询一条
    return row
}

// Query is the function for query multi rows
func (c *Conn) Query(query string) (*sql.Rows, error) {
    rows, err := c.netConn.Query(query) // 查询多条
    return rows, err
}

// Exec is the function for Insert、Update、Delete
func (c *Conn) Exec(query string) (sql.Result, error) {
    res, err := c.netConn.Exec(query) // 查询多条
    return res, err
}

//func init() {
	//mysql.New = New
//}
