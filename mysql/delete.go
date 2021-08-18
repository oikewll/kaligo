package mysql

import (
    //"fmt"
    "strconv"
)

// Delete is the struct for MySQL DATE type
type Delete struct {
    table   string

    //*Where   // 把Where 嵌套进来，它的参数和函数就可以直接使用了，Where又嵌套了Builder，Builder的参数和函数也都可以用
}

// Table Sets the table to delete from.
// 不需要了，db.go Delete(table string) 已经有了
//func (d *Delete) Table(table string) *Delete {
    //d.table = table
    //return d
//}

// DeleteCompile the SQL query and return it.
func (q *Query) DeleteCompile() string {
    // Start a deletion query
    sqlStr := "DELETE FROM " + q.QuoteTable(q.D.table)

    if len(q.W.wheres) != 0 {
        // Add deletion conditions
        sqlStr += " WHERE " + q.CompileConditions(q.W.wheres)
    }

    if len(q.W.orderBys) != 0 {
        // Add sorting
        sqlStr += " WHERE " + q.CompileOrderBy(q.W.orderBys)
    }

    if q.W.limit != 0 {
        // Add limiting
        sqlStr += "LIMIT " + strconv.Itoa(q.W.limit) 
    }

    //fmt.Printf("DeleteCompile === %v\n", sqlStr)
    q.sqlStr = sqlStr

    return sqlStr
}

// DeleteReset the query parameters
func (q *Query) DeleteReset() *Query {
    //fmt.Println("DeleteReset")
    q.D.table    = ""

    q.W.wheres   = nil
    q.W.orderBys = nil
    q.W.limit    = 0

    q.parameters = nil

    return q
}

