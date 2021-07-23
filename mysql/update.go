package mysql

import (
    //"errors"
    //"fmt"
    "strconv"
    //"strings"
    //"time"
    //"github.com/owner888/kaligo/util"
)

// Update is the struct for MySQL DATE type
type Update struct {
    table    string
    sets     [][]string  // slice
    joinObjs []*Join  // 多维slice
    lastJoin *Join       // last join statement

    Where
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
func (u *Update) Compile(db *DB) string {
    // db 链接是否存在 ？

    // Start and update query
    //sqlStr := "UPDATE " + db.QuoteTable(u.table)
    sqlStr := "UPDATE " + u.table

    if len(u.joinObjs) != 0 {
        // Builder.CompileJoin()
        sqlStr += u.CompileJoin(db, u.joinObjs)
    }

    // Add the columns to update
    // Builder.CompileSet()
    sqlStr += u.CompileSet(db, u.sets)

    if len(u.wheres) != 0 {
        // Add selection conditions
        // Builder.CompileConditions()
        // Where.wheres 参数
        sqlStr += "WHERE " + u.CompileConditions(db, u.wheres)
    }

    if len(u.orderBys) != 0 {
        // Add sorting
        // Builder.CompileOrderBy()
        // Where.orderBys 参数
        sqlStr += " " + u.CompileOrderBy(db, u.orderBys)
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
func (u *Update) Join(table string, joinType string) *Update {
    //u.lastJoin = new(QueryJoin)
    //u.joins = append(u.joins, u.lastJoin)
    return u
}

