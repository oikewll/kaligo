package mysql

import (
    //"errors"
    //"fmt"
    //"strconv"
    //"strings"
    //"time"
    //"github.com/owner888/kaligo/util"
)

// Where is the struct for MySQL DATE type
type Where struct {
    //wheres [][3]string
    wheres  map[string][][]string
    orderBys [][2]string
    limit   int

    *Builder
}

// Where Alias of andWhere
func (w *Where) Where(column string, op string, value string) *Where {
    return w.AndWhere(column, op, value)
}

// AndWhere Creates a new "AND WHERE" condition for the query.
// @title andWhere
// @description 查询条件函数
// @auth   seatle 2021/07/22 11:40
// @param  column string 字段名
// @param  op     string 操作符 >、=、<、>=、<=
// @param  value  string 查询值
// @return w     *Where Where对象
func (w *Where) AndWhere(column string, op string, value string) *Where {
    //w.wheres["AND"] = []string{column, op, value}
    w.wheres["AND"] = append(w.wheres["AND"], []string{column, op, value})
    return w
}

// OrWhere Creates a new "OR WHERE" condition for the query.
func (w *Where) OrWhere(column string, op string, value string) *Where {
    w.wheres["OR"] = append(w.wheres["OR"], []string{column, op, value})
    return w
}

// WhereOpen Alias of andWhereOpen
func (w *Where) WhereOpen() *Where {
    return w.AndWhereOpen()
}

// AndWhereOpen Opens a new "AND WHERE (...)" grouping.
func (w *Where) AndWhereOpen() *Where {
    //w.wheres["AND"] = []string{"("}
    w.wheres["AND"] = append(w.wheres["AND"], []string{"("})
    return w
}

// OrWhereOpen Opens a new "OR WHERE (...)" grouping.
func (w *Where) OrWhereOpen() *Where {
    w.wheres["OR"] = append(w.wheres["OR"], []string{"("})
    return w
}

// WhereClose Alias of andWhereClose
func (w *Where) WhereClose() *Where {
    return w.AndWhereClose()
}

// AndWhereClose Closes an open "AND WHERE (...)" grouping.
func (w *Where) AndWhereClose() *Where {
    w.wheres["AND"] = append(w.wheres["AND"], []string{")"})
    return w
}

// OrWhereClose Closes an open "OR WHERE (...)" grouping.
func (w *Where) OrWhereClose() *Where {
    w.wheres["OR"] = append(w.wheres["OR"], []string{")"})
    return w
}

// OrderBy Applies sorting with "ORDER By ..."
func (w *Where) OrderBy(column string, direction string) *Where {
    w.orderBys = append(w.orderBys, [2]string{column, direction})
    return w
}

// Limit Return up to "LIMIT ..." results
func (w *Where) Limit(value int) *Where {
    w.limit = value
    return w
}

// Reset the query parameters
//func (w *Where) Reset() *Where {
    //return w
//}
