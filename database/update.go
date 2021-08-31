package database

import (
    "strconv"
)

// Update is the struct for MySQL DATE type
type Update struct {
    table    string
    sets     [][]string  // 多维slice
}

// Set the values to update with an associative array
func (q *Query) Set(pairs map[string]string) *Query {
    for column, value := range pairs {
        q.Value(column, value)
    }
    return q
}

// Value Set the value of a single column.
func (q *Query) Value(column string, value string) *Query {
    var sets []string    
    sets = append(sets, column, value)
    q.U.sets = append(q.U.sets, sets)
    return q
}

// UpdateCompile Set the value of a single column.
func (q *Query) UpdateCompile() string {
    // Start and update query
    sqlStr := "UPDATE " + q.QuoteTable(q.U.table)

    if len(q.joinObjs) != 0 {
        sqlStr += " " + q.CompileJoin(q.joinObjs)
    }

    // Add the columns to update
    sqlStr += " SET " + q.CompileSet(q.U.sets)

    if len(q.W.wheres) != 0 {
        // Add selection conditions
        sqlStr += " WHERE " + q.CompileConditions(q.W.wheres)
    }

    if len(q.W.orderBys) != 0 {
        // Add sorting
        sqlStr += " " + q.CompileOrderBy(q.W.orderBys)
    }

    if q.W.limit != 0 {
        // Add limiting
        sqlStr += " LIMIT " + strconv.Itoa(q.W.limit)
    }

    q.sqlStr = sqlStr

    return sqlStr
}

// UpdateReset the query parameters
func (q *Query) UpdateReset() *Query {
    q.U.table    = ""
    q.U.sets     = nil

    q.W.wheres   = nil
    q.W.orderBys = nil
    q.W.limit    = 0

    q.joinObjs   = nil
    q.lastJoin   = nil
    q.parameters = nil

    return q
}
