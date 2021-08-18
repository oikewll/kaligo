package mysql

import (
    //"fmt"
    "strconv"
)

// Update is the struct for MySQL DATE type
type Update struct {
    table    string
    sets     [][]string  // 多维slice
    // 转移到query
    //joinObjs []*Join     // slice
    //lastJoin *Join       // last join statement

    //*Where
}

// Table Sets the table to update.
// 不需要了，db.go Update(table string) 已经有了
//func (q *Query) Table(table string) *Query {
    //q.U.table = table
    //return q
//}

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
    var sqlStr string    

    // Start and update query
    sqlStr = "UPDATE " + q.QuoteTable(q.U.table)

    if len(q.joinObjs) != 0 {
        // Builder.CompileJoin()
        sqlStr += " " + q.CompileJoin(q.joinObjs)
    }

    // Add the columns to update
    // Builder.CompileSet()
    sqlStr += " SET " + q.CompileSet(q.U.sets)

    if len(q.W.wheres) != 0 {
        // Add selection conditions
        // Builder.CompileConditions()
        // Where.wheres 参数
        sqlStr += " WHERE " + q.CompileConditions(q.W.wheres)
    }

    if len(q.W.orderBys) != 0 {
        // Add sorting
        // Builder.CompileOrderBy()
        // Where.orderBys 参数
        sqlStr += " " + q.CompileOrderBy(q.W.orderBys)
    }

    if q.W.limit != 0 {
        // Add limiting
        sqlStr += " LIMIT " + strconv.Itoa(q.W.limit)
    }

    //fmt.Printf("UpdateCompile === %v\n", sqlStr)
    q.sqlStr = sqlStr

    return sqlStr
}

// UpdateReset the query parameters
func (q *Query) UpdateReset() *Query {
    //fmt.Println("UpdateReset")
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

// Join Adds addition tables to "JOIN ...".
// @param joinType string join type (LEFT, RIGHT, INNER, etc)
//func (q *Query) Join(table string, joinType string) *Query {
    //q.lastJoin = &Join{
        //table: table,
        //joinType: joinType,
    //}
    //q.joinObjs = append(q.joinObjs, q.lastJoin)
    //return q
//}

// On Adds "ON ..." condition for the last created JOIN statement.
// @param joinType string join type (LEFT, RIGHT, INNER, etc)
//func (q *Query) On(c1 string, op string, c2 string) *Query {
    //q.lastJoin.On(c1, op, c2)
    //return q
//}

// AndOn Adds "ON ..." condition for the last created JOIN statement.
// @param joinType string join type (LEFT, RIGHT, INNER, etc)
//func (q *Query) AndOn(c1 string, op string, c2 string) *Query {
    //q.lastJoin.AndOn(c1, op, c2)
    //return q
//}

// OrOn Adds "OR ON ..." condition for the last created JOIN statement.
// @param joinType string join type (LEFT, RIGHT, INNER, etc)
//func (q *Query) OrOn(c1 string, op string, c2 string) *Query {
    //q.lastJoin.OrOn(c1, op, c2)
    //return q
//}

// OnOpen Adds an opening bracket the last created JOIN statement.
// @param joinType string join type (LEFT, RIGHT, INNER, etc)
//func (q *Query) OnOpen() *Query {
    //q.lastJoin.OnOpen()
    //return q
//}

// OnClose Adds an closing bracket the last created JOIN statement.
// @param joinType string join type (LEFT, RIGHT, INNER, etc)
//func (q *Query) OnClose() *Query {
    //q.lastJoin.OnClose()
    //return q
//}
