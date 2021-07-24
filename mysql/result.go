/**
 * Realize a database operation class
 *
 * @copyright   (C) 2014  seatle
 * @lastmodify  2021-07-06
 *
 */

package mysql

import (
    "database/sql"
    //"fmt"
    //"log"
    //"time"
    //"regexp"
    //"strings"
    //"sort"
    //"container/list"
    //"sync"
    //"reflect"
)

// Result is the struct for MySQL stores the result locally.
type Result struct {
    sqlStr              string          // string Executed SQL for this result
    result              sql.Result      // row result resource
    totalRows           int             // total number of rows
    currentRow          int             // current row number
    asObject            interface{}     // return rows as an object or associative array
    sanitizationEnabled bool            // If this is a records data will be anitized on get

    //rows []mysql.Row
	//res  mysql.Result
    //rows sql.Rows
    //row *Row
}

// NewResult 实例化数据库连接
// (读+写)连接数据库+选择数据库
func NewResult() *Result {
    r := new(Result)
    return r
}
