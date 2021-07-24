/**
 * Realize a database operation class
 *
 * @copyright   (C) 2014  seatle
 * @lastmodify  2021-07-06
 *
 */

package mysql

import (
    "strconv"
    "strings"
)

// Select is the struct for MySQL DATE type
type Select struct {
    selects  []string
    distinct bool
    froms    []string
    joinObjs []*Join     // join objects
    groupBys []string
    havings  map[string][][]string
    offset    int
    lastJoin  *Join      // last join statement

    Where
}

// NewSelect creates a new Select Object
func NewSelect(columns []string) *Select {
    NewWhere()  // 要验证一下这里调用以后，当前调用 Where 里面的 limitValue 的值是否有被赋值
    return &Select{
        selects  : columns,
        distinct : false,
        offset   : 0,
    }
}

// Distinct the query parameters
func (s *Select) Distinct(value bool) *Select {
    s.distinct = value

    return s
}

// Select Choose the columns to select from.
func (s *Select) Select(columns []string) *Select {
    s.selects = append(s.selects, columns...)

    return s
}

// SelectArray Choose the columns to select from.
func (s *Select) SelectArray(columns []string, reset bool) *Select {
    // 重设查询栏目
    if reset {
        s.selects = columns
    } else {
        s.selects = append(s.selects, columns...)
    }

    return s
}

// From Choose the tables to select "FROM ...".
func (s *Select) From(tables []string) *Select {
    s.froms = append(s.froms, tables...)

    return s
}

// Join Adds addition tables to "JOIN ...".
// @param string joinType type of JOIN: INNER, RIGHT, LEFT, etc
func (s *Select) Join(table string, joinType string) *Select {
    //s.lastJoin = &Join{table: table, joinType: joinType}
    s.lastJoin = NewJoin(table, joinType);
    s.joinObjs = append(s.joinObjs, s.lastJoin)

    return s
}

// On Adds "ON ..." conditions for the last created JOIN statement.
func (s *Select) On(c1 string, op string, c2 string) *Select {
    s.lastJoin.On(c1, op, c2)

    return s
}

// AndOn Adds "AND ON ..." conditions for the last created JOIN statement.
func (s *Select) AndOn(c1 string, op string, c2 string) *Select {
    s.lastJoin.AndOn(c1, op, c2)

    return s
}

// OrOn Adds "OR ON ..." conditions for the last created JOIN statement.
func (s *Select) OrOn(c1 string, op string, c2 string) *Select {
    s.lastJoin.OrOn(c1, op, c2)

    return s
}

// OnOpen Adds an opening bracket the last created JOIN statement.
func (s *Select) OnOpen() *Select {
    s.lastJoin.OnOpen()

    return s
}

// OnClose Adds an closing bracket the last created JOIN statement.
func (s *Select) OnClose() *Select {
    s.lastJoin.OnClose()

    return s
}

// GroupBy Creates a "GROUP BY ..." filter.
// @param  columns  []string  column name or []string{column, column} or object
func (s *Select) GroupBy(columns []string) *Select {
    //for idx, column := range columns {
        // 如果column是 []string，再循环一边
        //if i, ok := column.([][]string); ok {
            //for k, v := range column {
                //columns = append(columns, v)
            //}
            //columns[idx] = nil
        //}
    //}

    s.groupBys = append(s.groupBys, columns...)

    return s
}

// Having Alias of AndHaving
func (s *Select) Having(column string, op string, value string) *Select {
    return s.AndHaving(column, op, value)
}

// AndHaving Creates a new "AND HAVING" condition for the query.
func (s *Select) AndHaving(column string, op string, value string) *Select {
    s.havings["AND"] = append(s.havings["AND"], []string{column, op, value})

    return s
}

// OrHaving Creates a new "AND HAVING" condition for the query.
func (s *Select) OrHaving(column string, op string, value string) *Select {
    s.havings["OR"] = append(s.havings["OR"], []string{column, op, value})

    return s
}

// HavingOpen Alias of AndHavingOpen
func (s *Select) HavingOpen() *Select {
    return s.AndHavingOpen()
}

// AndHavingOpen Opens a new "AND HAVING (...)" grouping.
func (s *Select) AndHavingOpen() *Select {
    s.havings["AND"] = append(s.havings["AND"], []string{"("})

    return s
}

// OrHavingOpen Opens a new "OR HAVING (...)" grouping.
func (s *Select) OrHavingOpen() *Select {
    s.havings["OR"] = append(s.havings["OR"], []string{"("})

    return s
}

// HavingClose Alias of AndHavingClose
func (s *Select) HavingClose() *Select {
    return s.AndHavingClose()
}

// AndHavingClose Opens a new "AND HAVING (...)" grouping.
func (s *Select) AndHavingClose() *Select {
    s.havings["AND"] = append(s.havings["AND"], []string{")"})

    return s
}

// OrHavingClose Opens a new "OR HAVING (...)" grouping.
func (s *Select) OrHavingClose() *Select {
    s.havings["OR"] = append(s.havings["OR"], []string{")"})

    return s
}

// Offset Start returning results after "OFFSET ...".
func (s *Select) Offset(number int) *Select {
    s.offset = number

    return s
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

// Compile Set the value of a single column.
func (s *Select) Compile(db *DB) string {
    // db 链接是否存在 ？

    // Start a selection query
    sqlStr := "SELECT "

    if s.distinct {
        // Select only unique results
        sqlStr += "DISTINCT"
    }

    if len(s.selects) == 0 {
        // Select all columns
        sqlStr += "*"
    } else {
        s.selects = arrayUnique(s.selects)
        //for k, v := range s.selects {
            //s.selects[k] = db.QuoteIdent(v)
        //}
        sqlStr += strings.Join(s.selects, ", ")
    }

    if len(s.froms) != 0 {
        s.froms = arrayUnique(s.froms)
        // Set tables to select from
        for k, v := range s.froms {
            s.froms[k] = quoteTable(v)
        }
        sqlStr += strings.Join(s.froms, ", ")
    }

    if len(s.joinObjs) != 0 {
        // Add tables to join
        // Builder.CompileJoin()
        sqlStr += s.CompileJoin(db, s.joinObjs)
    }

    // select 没有 set 用法
    // Add the columns to update
    // Builder.CompileSet()
    //sqlStr += s.CompileSet(db, s.sets)

    if len(s.wheres) != 0 {
        // Add selection conditions
        // Builder.CompileConditions()
        // Where.wheres 参数
        sqlStr += "WHERE " + s.CompileConditions(db, s.wheres)
    }

    if len(s.groupBys) != 0 {
        // Add sorting
        s.groupBys = arrayUnique(s.groupBys)
        //for k, v := range s.groupBys {
            //s.groupBys[k] = db.QuoteIdent(v)
        //}
        sqlStr += strings.Join(s.groupBys, ", ")
    }

    if len(s.havings) != 0 {
        // Add filtering conditions
        // Builder.CompileConditions()
        // Where.havings 参数
        sqlStr += "HAVING " + s.CompileConditions(db, s.havings)
    }

    if len(s.orderBys) != 0 {
        // Add sorting
        // Builder.CompileOrderBy()
        // Where.orderBys 参数
        sqlStr += " " + s.CompileOrderBy(db, s.orderBys)
    }

    if s.limit != 0 {
        // Add limiting
        sqlStr += " LIMIT " + strconv.Itoa(s.limit)
    }

    if s.offset != 0 {
        // Add offsets
        sqlStr += " OFFSET " + strconv.Itoa(s.offset)
    }

    return sqlStr
}

// Reset the query parameters
func (s *Select) Reset() *Select {
    s.selects    = nil
    s.froms      = nil
    s.joinObjs   = nil
    s.wheres     = nil
    s.groupBys   = nil
    s.havings    = nil
    s.orderBys   = nil
    s.distinct   = false
    s.limit      = 0
    s.offset     = 0
    s.lastJoin   = nil
    s.parameters = nil

    return s
}
