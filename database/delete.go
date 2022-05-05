package database

import (
    //"fmt"
    "strconv"
)

// Delete is the struct for MySQL DATE type
type Delete struct {
    table   string

    //*Where   // 把Where 嵌套进来，它的参数和函数就可以直接使用了，Where又嵌套了Builder，Builder的参数和函数也都可以用
}

// DeleteCompile the SQL query and return it.
func (q *Query) DeleteCompile() string {
    // Start a deletion query
    sqlStr := "DELETE FROM " + q.QuoteTable(q.D.table)

    if len(q.W.params) != 0 {
        // Add selection conditions
        conditionsStr, values := q.CompileConditions(q.W.params)
        sqlStr += " WHERE " + conditionsStr
        q.W.values = append(q.W.values, values...)
    }

    if len(q.W.orderBys) != 0 {
        sqlStr += " WHERE " + q.CompileOrderBy(q.W.orderBys)
    }

    // SQLite does not support LIMIT for DELETE、UPDATE 
    if q.W.limit != 0 && q.Dialector.Name() != "sqlite" {
        sqlStr += "LIMIT " + strconv.Itoa(q.W.limit) 
    }

    q.sqlStr = sqlStr

    return sqlStr
}

// DeleteReset the query parameters
func (q *Query) DeleteReset() *Query {
    q.D.table    = ""

    q.W.params   = nil
    q.W.orderBys = nil
    q.W.limit    = 0

    q.parameters = nil

    return q
}

