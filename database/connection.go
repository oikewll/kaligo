package database

//"context"
//"database/sql"
//"fmt"
//"log"
//"regexp"
//"reflect"
//"strings"
//"time"

// Connection is the struct for MySQL connection handler
//type Connection struct {
//*DB
//active string               // instance name

//timeout time.Duration       // Timeout for connect SetConnMaxLifetime(timeout)
//lastUse time.Time           // The last use time

//inTransaction bool
//stdDB         *sql.DB       // MySQL connection
//stdTX         *sql.Tx       // MySQL connection for Transaction
//autoCommit    bool
//initCmds      []string      // MySQL commands/queries executed after connect
//tablePrefix   string

//Dest          any   // Object：&User{}
//schema        *Schema
//Debug         bool          // Debug logging. You may change it at any time.
//logSlowQuery  bool
//logSlowTime   int

//Context context.Context
//}

// Stats Returns database statistics
//func (c *Connection) Stats() (err error) { return err }
// Caching Per connection cache controller setter/getter
//func (c *Connection) Caching() bool { return false }

//// Query is the function for query multi rows
//func (c *Connection) Query(queryType int, sqlStr string, asObject any) *Result {
//// var stacktrace []map[string]string   // 储存所有调用函数，一层一层的
//// benchmark := Profiler.Start(c.name, sqlstr, stacktrace)
//// Profiler.Stop(benchmark)

//// golang 好像不需要处理执行时因为mysql链接闲置超过8小时的问题：MySQL server has gone away

//// Set the last query
//c.lastQuery = sqlStr

//r := &Result{
//sqlStr: sqlStr,
////result              sql.Result      // row result resource
////totalRows           int             // total number of rows
////currentRow          int             // current row number
////asObject            any     // return rows as an object or associative array
////sanitizationEnabled bool            // If this is a records data will be anitized on get
////rows                *sql.Rows
//}

//if queryType == SELECT {
//// if Config.Get("enable_cache") { return cached() }
//// return Result{result:sql.Result, sqlstr: sqlstr, asObject: asObject}
////r.result =
//rows, err := c.stdDB.Query(sqlStr) // 查询多条
//fmt.Printf("Connection.Query() = %v %v\n", rows, err)
//} else if queryType == INSERT {
//// return []string{connection.insertID, connection.affectedRows}
//} else if queryType == UPDATE || queryType == DELETE {
//// return connection.affectedRows
//}

//return r
//}
