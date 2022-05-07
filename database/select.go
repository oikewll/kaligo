package database

import (
    "fmt"
    "strconv"
    "strings"
)

// Select is the struct for MySQL DATE type
type Select struct {
    selects   []any
    distinct  bool
    froms     []string
    groupBys  []string
    havings   map[string][]WhereParam   // map["AND"][] WhereParam{column, op, value}
    offset    int
    forUpdate bool
}

// Select Choose the columns to select from.
func (q *Query) Select(columns []any) *Query {
    for _, value := range columns {
        q.S.selects = append(q.S.selects, value)
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

func (q *Query) ForUpdate(forUpdate bool) *Query {
    q.S.forUpdate = forUpdate

    return q
}

// Join Adds addition tables to "JOIN ...".
// @param string joinType type of JOIN: INNER, RIGHT, LEFT, etc
func (q *Query) Join(table string, joinType string) *Query {
    //s.lastJoin = &Join{table: table, joinType: joinType}
    q.lastJoin = NewJoin(table, joinType)
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
    // q.S.havings["AND"] = append(q.S.havings["AND"], []string{column, op, value})

    if q.S.havings == nil {
        q.S.havings = make(map[string][]WhereParam)
    }
    q.S.havings["AND"] = append(q.S.havings["AND"], WhereParam{column, op, value})

    return q
}

// OrHaving Creates a new "AND HAVING" condition for the query.
func (q *Query) OrHaving(column string, op string, value string) *Query {
    // q.S.havings["OR"] = append(q.S.havings["OR"], []string{column, op, value})
    q.S.havings["OR"] = append(q.S.havings["OR"], WhereParam{column, op, value})

    return q
}

// HavingOpen Alias of AndHavingOpen
func (q *Query) HavingOpen() *Query {
    return q.AndHavingOpen()
}

// AndHavingOpen Opens a new "AND HAVING (...)" grouping.
func (q *Query) AndHavingOpen() *Query {
    // q.S.havings["AND"] = append(q.S.havings["AND"], []string{"("})
    q.S.havings["AND"] = append(q.S.havings["AND"], WhereParam{column: "("})

    return q
}

// OrHavingOpen Opens a new "OR HAVING (...)" grouping.
func (q *Query) OrHavingOpen() *Query {
    // q.S.havings["OR"] = append(q.S.havings["OR"], []string{"("})
    q.S.havings["OR"] = append(q.S.havings["OR"], WhereParam{column: ")"})

    return q
}

// HavingClose Alias of AndHavingClose
func (q *Query) HavingClose() *Query {
    return q.AndHavingClose()
}

// AndHavingClose Opens a new "AND HAVING (...)" grouping.
func (q *Query) AndHavingClose() *Query {
    // q.S.havings["AND"] = append(q.S.havings["AND"], []string{")"})
    q.S.havings["AND"] = append(q.S.havings["AND"], WhereParam{column: ")"})

    return q
}

// OrHavingClose Opens a new "OR HAVING (...)" grouping.
func (q *Query) OrHavingClose() *Query {
    // q.S.havings["OR"] = append(q.S.havings["OR"], []string{")"})
    q.S.havings["OR"] = append(q.S.havings["OR"], WhereParam{column: ")"})

    return q
}

// Offset Start returning results after "OFFSET ...".
func (q *Query) Offset(number int) *Query {
    q.S.offset = number

    return q
}

// SelectCompile Set the value of a single column.
func (q *Query) SelectCompile() string {
    if len(q.S.froms) == 0 {
        q.AddError(ErrInvalidValue)
    }
    froms := arrayUnique(q.S.froms)

    // Start a selection query
    sqlStr := "SELECT "

    if q.S.distinct {
        // Select only unique results
        sqlStr += "DISTINCT"
    }

    if len(q.S.selects) == 0 {
        sqlStr += "*"
    } else {
        var selects []string
        for _, v := range q.S.selects {
            switch c := v.(type) {
            case string:
                selects = append(selects, c)
            case Expression: // type Expression string
                selects = append(selects, string(c))
            default:
                logs.Warn("Select Columns Type Error: %T = %v", c, c)
            }
        }
        // 去重
        selects = arrayUnique(selects)

        for k, v := range selects {
            // Is the column need decrypt ???
            for _, table := range q.S.froms {
                if cryptFields, ok := q.cryptFields[table]; ok && q.Dialector.Name() == "mysql" && q.cryptKey != "" && InSlice(v, &cryptFields) {
                    selects[k] = fmt.Sprintf("AES_DECRYPT(%s, \"%s\") AS %v", q.QuoteIdentifier(v), q.cryptKey, q.QuoteIdentifier(v))
                } else {
                    selects[k] = q.QuoteIdentifier(v)
                }
            }
        }
        sqlStr += strings.Join(selects, ", ")
    }

    // Set tables to select froms
    for k, v := range froms {
        froms[k] = q.QuoteTable(v)
    }
    sqlStr += " FROM " + strings.Join(froms, ", ")

    if len(q.joinObjs) != 0 {
        // Add tables to join
        sqlStr += " " + q.CompileJoin(q.joinObjs)
    }

    // select 没有 set 用法
    // Add the columns to update
    // Builder.CompileSet()
    //sqlStr += s.CompileSet(db, s.sets)

    // if len(q.W.wheres) != 0 {
    if len(q.W.params) != 0 {
        // Add selection conditions
        // Builder.CompileConditions()
        // Where.wheres 参数
        // sqlStr += " WHERE " + q.CompileConditions(q.W.wheres)
        conditionsStr, values := q.CompileConditions(q.W.params)
        sqlStr += " WHERE " + conditionsStr
        q.W.values = append(q.W.values, values...)
    }

    if len(q.S.groupBys) != 0 {
        q.S.groupBys = arrayUnique(q.S.groupBys)
        // Add sorting
        for k, v := range q.S.groupBys {
            q.S.groupBys[k] = q.QuoteIdentifier(v)
        }
        sqlStr += " GROUP BY " + strings.Join(q.S.groupBys, ", ")
    }

    if len(q.S.havings) != 0 {
        // Add filtering conditions
        // Builder.CompileConditions()
        // Where.havings 参数
        // sqlStr += " HAVING " + q.CompileConditions(q.S.havings)
        conditionsStr, values := q.CompileConditions(q.S.havings)
        sqlStr += " HAVING " + conditionsStr
        // 值全部捅到 Where.values
        q.W.values = append(q.W.values, values...)
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

    if q.S.forUpdate {
        // Add for update
        if q.InTransaction == true {
            sqlStr += " FOR UPDATE"
        } else {
            logs.Warn("SELECT ... FOR UPDATE can't use for non-transactional environment")
        }
    }

    //fmt.Printf("SelectCompile === %v\n", sqlStr)
    //fmt.Printf("SelectCompile === %v\n", sqlStr)
    q.sqlStr = sqlStr

    return sqlStr
}

// SelectReset the query parameters
func (q *Query) SelectReset() *Query {
    //fmt.Println("SelectReset")
    q.S.selects   = nil
    q.S.distinct  = false
    q.S.froms     = nil
    q.S.groupBys  = nil
    q.S.havings   = nil
    q.S.offset    = 0
    q.S.forUpdate = false

    // q.W.wheres    = nil
    q.W.orderBys  = nil
    q.W.limit     = 0

    q.joinObjs    = nil
    q.lastJoin    = nil
    q.parameters  = nil

    return q
}

/* vim: set expandtab: */
