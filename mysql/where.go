package mysql

import (
    //"errors"
    //"fmt"
    //"strconv"
    //"strings"
    //"time"
)

// Where is the struct for MySQL DATE type
type Where struct {
    //wheres [][3]string
    wheres   map[string][][]string
    orderBys [][2]string
    limit    int

    //*Builder
}

// Where Alias of andWhere
func (q *Query) Where(column string, op string, value string) *Query {
    return q.AndWhere(column, op, value)
}

// AndWhere Creates a new "AND WHERE" condition for the query.
// @title andWhere
// @description 查询条件函数
// @auth   seatle 2021/07/22 11:40
// @param  column string 字段名
// @param  op     string 操作符 >、=、<、>=、<=
// @param  value  string 查询值
// @return w     *Where Where对象
func (q *Query) AndWhere(column string, op string, value string) *Query {
    //w.wheres["AND"] = []string{column, op, value}
    // 需要加锁吗？不用吧，因为一个查询是一个对象啊
    if q.W.wheres == nil {
        q.W.wheres = make(map[string][][]string)
    }
    q.W.wheres["AND"] = append(q.W.wheres["AND"], []string{column, op, value})
    return q
}

// OrWhere Creates a new "OR WHERE" condition for the query.
func (q *Query) OrWhere(column string, op string, value string) *Query {
    q.W.wheres["OR"] = append(q.W.wheres["OR"], []string{column, op, value})
    return q
}

// WhereOpen Alias of andWhereOpen
func (q *Query) WhereOpen() *Query {
    return q.AndWhereOpen()
}

// AndWhereOpen Opens a new "AND WHERE (...)" grouping.
func (q *Query) AndWhereOpen() *Query {
    //w.wheres["AND"] = []string{"("}
    q.W.wheres["AND"] = append(q.W.wheres["AND"], []string{"("})
    return q
}

// OrWhereOpen Opens a new "OR WHERE (...)" grouping.
func (q *Query) OrWhereOpen() *Query {
    q.W.wheres["OR"] = append(q.W.wheres["OR"], []string{"("})
    return q
}

// WhereClose Alias of andWhereClose
func (q *Query) WhereClose() *Query {
    return q.AndWhereClose()
}

// AndWhereClose Closes an open "AND WHERE (...)" grouping.
func (q *Query) AndWhereClose() *Query {
    q.W.wheres["AND"] = append(q.W.wheres["AND"], []string{")"})
    return q
}

// OrWhereClose Closes an open "OR WHERE (...)" grouping.
func (q *Query) OrWhereClose() *Query {
    q.W.wheres["OR"] = append(q.W.wheres["OR"], []string{")"})
    return q
}

// OrderBy Applies sorting with "ORDER By ..."
func (q *Query) OrderBy(column string, direction string) *Query {
    q.W.orderBys = append(q.W.orderBys, [2]string{column, direction})
    return q
}

// Limit Return up to "LIMIT ..." results
func (q *Query) Limit(value int) *Query {
    q.W.limit = value
    return q
}

// Reset the query parameters
//func (w *Where) Reset() *Where {
    //return w
//}
