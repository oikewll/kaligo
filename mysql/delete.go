package mysql

import (
    "strconv"
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
func (d *Delete) Compile() string {
    // Start a deletion query
    sqlStr := "DELETE FROM " + d.connection.QuoteTable(d.table)

    if len(d.wheres) == 0 {
        //sqlStr += " WHERE " + d.compileConditions(d.wheres)
    }

    if len(d.orderBys) == 0 {
        //sqlStr += " WHERE " + d.compileOrderBy(d.orderBys)
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
    d.table      = ""
    d.wheres     = nil
    d.orderBys   = nil
    d.parameters = nil
    d.limit      = 0
    return d
}

