package mysql

import (
    //"errors"
    //"fmt"
    "strconv"
    //"strings"
    //"time"
    //"github.com/owner888/kaligo/util"
)

// Delete is the struct for MySQL DATE type
type Delete struct {
    table   string

    Where   // 把Where 嵌套进来，它的参数和函数就可以直接使用了，Where又嵌套了Builder，Builder的参数和函数也都可以用
}

// Table Sets the table to delete from.
func (d *Delete) Table(table string) *Delete {
    d.table = table
    return d
}

// Compile the SQL query and return it.
func (d *Delete) Compile(db *DB) string {
    db.Conn.Connect()   // 这里要返回database connection instance

    // Start a deletion query
    sqlStr := "DELETE FROM " + quoteTable(d.table)

    if len(d.wheres) == 0 {
        //sqlStr += " WHERE " + d.compileConditions(db, d.wheres)
    }

    if len(d.orderBys) == 0 {
        //sqlStr += " WHERE " + d.compileOrderBy(db, d.orderBys)
    }

    // TODO limit 条件拼接
    if d.limit != 0 {
        // Add limiting
        sqlStr += "LIMIT " + strconv.Itoa(d.limit) 
    }
    return sqlStr
}

// Reset the query parameters
func (d *Delete) Reset() *Delete {
    d.table  = ""
    d.wheres   = nil
    d.orderBys = nil
    d.parameters = nil
    d.limit   = 0
    return d
}

