package database

import (
    "strconv"
)

type set struct {
    column string
    value  any
}

// Update is the struct for MySQL DATE type
type Update struct {
    table string
    sets  []set // 多维slice
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
    q.U.sets = append(q.U.sets, set{column: column, value: value})
    return q
}

// Value Set the value of a single column.
func (q *Query) ValueExpression(column string, value Expression) *Query {
    q.U.sets = append(q.U.sets, set{column: column, value: value})
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

    // if len(q.W.wheres) != 0 {
    if len(q.W.params) != 0 {
        // Add selection conditions
        // sqlStr += " WHERE " + q.CompileConditions(q.W.wheres)
        conditionsStr, values := q.CompileConditions(q.W.params)
        sqlStr += " WHERE " + conditionsStr
        q.W.values = append(q.W.values, values)
    }

    if len(q.W.orderBys) != 0 {
        // Add sorting
        sqlStr += " " + q.CompileOrderBy(q.W.orderBys)
    }

    // SQLite does not support LIMIT for DELETE、UPDATE
    if q.W.limit != 0 && q.Dialector.Name() != "sqlite" {
        // Add limiting
        sqlStr += " LIMIT " + strconv.Itoa(q.W.limit)
    }

    q.sqlStr = sqlStr

    return sqlStr
}

// UpdateReset the query parameters
func (q *Query) UpdateReset() *Query {
    q.U.table = ""
    q.U.sets = nil

    // q.W.wheres = nil
    q.W.params = nil
    q.W.orderBys = nil
    q.W.limit = 0

    q.joinObjs = nil
    q.lastJoin = nil
    q.parameters = nil

    return q
}
