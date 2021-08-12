/**
 * Realize a database operation class
 *
 * @copyright   (C) 2014  seatle
 * @lastmodify  2021-07-06
 *
 */

package mysql

import (
    //"fmt"
    "strconv"
    "strings"
)

// Select is the struct for MySQL DATE type
type Select struct {
    selects  []string
    distinct bool
    froms    []string
    groupBys []string
    havings  map[string][][]string
    offset   int
    //joinObjs []*Join     // join objects
    //lastJoin *Join       // last join statement

    //*Where
}

// Select Choose the columns to select from.
//func (s *Select) Select(columns []string) *Select {
    //s.selects = append(s.selects, columns...)

    //return s
//}

// Select Choose the columns to select from.
func (q *Query) Select(columns []string) *Query {
    q.S.selects = append(q.S.selects, columns...)

    return q
}

// SelectArray Choose the columns to select from.
func (q *Query) SelectArray(columns []string, reset bool) *Query {
    // 重设查询栏目
    if reset {
        q.S.selects = columns
    } else {
        q.S.selects = append(q.S.selects, columns...)
    }

    return q
}

// Distinct the query parameters
func (q *Query) Distinct(value bool) *Query {
    q.S.distinct = value

    return q
}

// From Choose the tables to select "FROM ...".
func (q *Query) From(tables ...string) *Query {
    q.S.froms = append(q.S.froms, tables...)

    return q
}

// Join Adds addition tables to "JOIN ...".
// @param string joinType type of JOIN: INNER, RIGHT, LEFT, etc
func (q *Query) Join(table string, joinType string) *Query {
    //s.lastJoin = &Join{table: table, joinType: joinType}
    q.lastJoin = NewJoin(table, joinType);
    q.joinObjs = append(q.joinObjs, q.lastJoin)

    return q
}

// On Adds "ON ..." conditions for the last created JOIN statement.
func (q *Query) On(c1 string, op string, c2 string) *Query {
    q.lastJoin.On(c1, op, c2)

    return q
}

// AndOn Adds "AND ON ..." conditions for the last created JOIN statement.
func (q *Query) AndOn(c1 string, op string, c2 string) *Query {
    q.lastJoin.AndOn(c1, op, c2)

    return q
}

// OrOn Adds "OR ON ..." conditions for the last created JOIN statement.
func (q *Query) OrOn(c1 string, op string, c2 string) *Query {
    q.lastJoin.OrOn(c1, op, c2)

    return q
}

// OnOpen Adds an opening bracket the last created JOIN statement.
func (q *Query) OnOpen() *Query {
    q.lastJoin.OnOpen()

    return q
}

// OnClose Adds an closing bracket the last created JOIN statement.
func (q *Query) OnClose() *Query {
    q.lastJoin.OnClose()

    return q
}

// GroupBy Creates a "GROUP BY ..." filter.
// @param  columns  []string  column name or []string{column, column} or object
func (q *Query) GroupBy(columns []string) *Query {
    //for idx, column := range columns {
        // 如果column是 []string，再循环一边
        //if i, ok := column.([][]string); ok {
            //for k, v := range column {
                //columns = append(columns, v)
            //}
            //columns[idx] = nil
        //}
    //}

    q.S.groupBys = append(q.S.groupBys, columns...)

    return q
}

// Having Alias of AndHaving
func (q *Query) Having(column string, op string, value string) *Query {
    return q.AndHaving(column, op, value)
}

// AndHaving Creates a new "AND HAVING" condition for the query.
func (q *Query) AndHaving(column string, op string, value string) *Query {
    q.S.havings["AND"] = append(q.S.havings["AND"], []string{column, op, value})

    return q
}

// OrHaving Creates a new "AND HAVING" condition for the query.
func (q *Query) OrHaving(column string, op string, value string) *Query {
    q.S.havings["OR"] = append(q.S.havings["OR"], []string{column, op, value})

    return q
}

// HavingOpen Alias of AndHavingOpen
func (q *Query) HavingOpen() *Query {
    return q.AndHavingOpen()
}

// AndHavingOpen Opens a new "AND HAVING (...)" grouping.
func (q *Query) AndHavingOpen() *Query {
    q.S.havings["AND"] = append(q.S.havings["AND"], []string{"("})

    return q
}

// OrHavingOpen Opens a new "OR HAVING (...)" grouping.
func (q *Query) OrHavingOpen() *Query {
    q.S.havings["OR"] = append(q.S.havings["OR"], []string{"("})

    return q
}

// HavingClose Alias of AndHavingClose
func (q *Query) HavingClose() *Query {
    return q.AndHavingClose()
}

// AndHavingClose Opens a new "AND HAVING (...)" grouping.
func (q *Query) AndHavingClose() *Query {
    q.S.havings["AND"] = append(q.S.havings["AND"], []string{")"})

    return q
}

// OrHavingClose Opens a new "OR HAVING (...)" grouping.
func (q *Query) OrHavingClose() *Query {
    q.S.havings["OR"] = append(q.S.havings["OR"], []string{")"})

    return q
}

// Offset Start returning results after "OFFSET ...".
func (q *Query) Offset(number int) *Query {
    q.S.offset = number

    return q
}

func arrayUnique(arr []string) []string{
	size := len(arr)
	result := make([]string, 0, size)
	temp := map[string]struct{}{}
	for i:=0; i < size; i++ {
		if _,ok := temp[arr[i]]; ok != true {
			temp[arr[i]] = struct{}{}
			result = append(result, arr[i])
		}
	}
	return result
}

// SelectCompile Set the value of a single column.
func (q *Query) SelectCompile() string {
    // Start a selection query
    sqlStr := "SELECT "

    if q.S.distinct {
        // Select only unique results
        sqlStr += "DISTINCT"
    }

    if len(q.S.selects) == 0 {
        // Select all columns
        sqlStr += "*"
    } else {
        q.S.selects = arrayUnique(q.S.selects)
        for k, v := range q.S.selects {
            q.S.selects[k] = q.C.QuoteIdentifier(v)
        }
        sqlStr += strings.Join(q.S.selects, ", ")
    }

    if len(q.S.froms) != 0 {
        // Set tables to select from
        q.S.froms = arrayUnique(q.S.froms)
        for k, v := range q.S.froms {
            q.S.froms[k] = q.C.QuoteTable(v)
        }
        sqlStr += " FROM " + strings.Join(q.S.froms, ", ")
    }

    if len(q.joinObjs) != 0 {
        // Add tables to join
        sqlStr += " " + q.CompileJoin(q.joinObjs)
    }

    // select 没有 set 用法
    // Add the columns to update
    // Builder.CompileSet()
    //sqlStr += s.CompileSet(db, s.sets)

    if len(q.W.wheres) != 0 {
        // Add selection conditions
        // Builder.CompileConditions()
        // Where.wheres 参数
        sqlStr += " WHERE " + q.CompileConditions(q.W.wheres)
    }

    if len(q.S.groupBys) != 0 {
        // Add sorting
        q.S.groupBys = arrayUnique(q.S.groupBys)
        for k, v := range q.S.groupBys {
            q.S.groupBys[k] = q.C.QuoteIdentifier(v)
        }
        sqlStr += " GROUP BY " + strings.Join(q.S.groupBys, ", ")
    }

    if len(q.S.havings) != 0 {
        // Add filtering conditions
        // Builder.CompileConditions()
        // Where.havings 参数
        sqlStr += " HAVING " + q.CompileConditions(q.S.havings)
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

    if q.S.offset != 0 {
        // Add offsets
        sqlStr += " OFFSET " + strconv.Itoa(q.S.offset)
    }

    //fmt.Printf("SelectCompile === %v\n", sqlStr)
    q.sqlStr = sqlStr

    return sqlStr
}

// SelectReset the query parameters
func (q *Query) SelectReset() *Query {
    //fmt.Println("SelectReset")
    q.S.selects  = nil
    q.S.distinct = false
    q.S.froms    = nil
    q.S.groupBys = nil
    q.S.havings  = nil
    q.S.offset   = 0

    q.W.wheres   = nil
    q.W.orderBys = nil
    q.W.limit    = 0

    q.joinObjs   = nil
    q.lastJoin   = nil
    q.parameters = nil

    return q
}
