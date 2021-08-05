package mysql

import (
    //"fmt"
    "strconv"
)

// Update is the struct for MySQL DATE type
type Update struct {
    table    string
    sets     [][]string  // slice
    joinObjs []*Join  // 多维slice
    lastJoin *Join       // last join statement

    *Where
}

// Table Sets the table to update.
func (u *Update) Table(table string) *Update {
    u.table = table
    return u
}

// Set the values to update with an associative array
func (u *Update) Set(pairs map[string]string) *Update {
    var sets []string    
    for column, value := range pairs {
        sets = append(sets, column, value)
        u.sets = append(u.sets, sets)
    }
    return u
}

// Value Set the value of a single column.
func (u *Update) Value(column string, value string) *Update {
    var sets []string    
    sets = append(sets, column, value)
    u.sets = append(u.sets, sets)
    return u
}

// Compile Set the value of a single column.
func (u *Update) Compile(args ...*Connection) string {
    var conn *Connection    
    if len(args) != 0 {
        conn = args[0]
    } else {
        // Get the database instance
        db := New()
        conn = db.C
    }
    //fmt.Printf("Update Compile === %T = %p\n", conn, conn)

    // Start and update query
    sqlStr := "UPDATE " + u.connection.QuoteTable(u.table)

    if len(u.joinObjs) != 0 {
        // Builder.CompileJoin()
        sqlStr += u.CompileJoin(conn, u.joinObjs)
    }

    // Add the columns to update
    // Builder.CompileSet()
    sqlStr += u.CompileSet(conn, u.sets)

    if len(u.wheres) != 0 {
        // Add selection conditions
        // Builder.CompileConditions()
        // Where.wheres 参数
        sqlStr += "WHERE " + u.CompileConditions(conn, u.wheres)
    }

    if len(u.orderBys) != 0 {
        // Add sorting
        // Builder.CompileOrderBy()
        // Where.orderBys 参数
        sqlStr += " " + u.CompileOrderBy(conn, u.orderBys)
    }

    if u.limit != 0 {
        // Add limiting
        sqlStr += " LIMIT " + strconv.Itoa(u.limit)
    }

    return sqlStr
}

// Reset the query parameters
func (u *Update) Reset() *Update {
    u.table      = ""
    u.joinObjs   = nil
    u.sets       = nil
    u.wheres     = nil
    u.orderBys   = nil
    u.limit      = 0
    u.lastJoin   = nil  // 清空对象
    u.parameters = nil

    return u
}

// Join Adds addition tables to "JOIN ...".
// @param joinType string join type (LEFT, RIGHT, INNER, etc)
func (u *Update) Join(table string, joinType string) *Update {
    u.lastJoin = &Join{
        table: table,
        joinType: joinType,
    }
    u.joinObjs = append(u.joinObjs, u.lastJoin)
    return u
}

// On Adds "ON ..." condition for the last created JOIN statement.
// @param joinType string join type (LEFT, RIGHT, INNER, etc)
func (u *Update) On(c1 string, op string, c2 string) *Update {
    u.lastJoin.On(c1, op, c2)
    return u
}

// AndOn Adds "ON ..." condition for the last created JOIN statement.
// @param joinType string join type (LEFT, RIGHT, INNER, etc)
func (u *Update) AndOn(c1 string, op string, c2 string) *Update {
    u.lastJoin.AndOn(c1, op, c2)
    return u
}

// OrOn Adds "OR ON ..." condition for the last created JOIN statement.
// @param joinType string join type (LEFT, RIGHT, INNER, etc)
func (u *Update) OrOn(c1 string, op string, c2 string) *Update {
    u.lastJoin.OrOn(c1, op, c2)
    return u
}

// OnOpen Adds an opening bracket the last created JOIN statement.
// @param joinType string join type (LEFT, RIGHT, INNER, etc)
func (u *Update) OnOpen() *Update {
    u.lastJoin.OnOpen()
    return u
}

// OnClose Adds an closing bracket the last created JOIN statement.
// @param joinType string join type (LEFT, RIGHT, INNER, etc)
func (u *Update) OnClose() *Update {
    u.lastJoin.OnClose()
    return u
}
