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
    //"sync"
    //"reflect"
)

// Result is the struct for MySQL stores the result locally.
type Result struct {
    sqlStr              string          // string Executed SQL for this result
    result              sql.Result      // row result resource
    totalRows           int             // total number of rows
    currentRow          int             // current row number
    rows                *sql.Rows       // For select
    row                 *Row

    fields              []*Field        // Fields table
    affectedRows        int             // For update and delete
    insertID            int             // Primary key value ( useful for AUTO_INCREMENT primary keys)
}

// NewResult 实例化数据库连接
// (读+写)连接数据库+选择数据库
func NewResult() *Result {
    r := new(Result)
    return r
}

