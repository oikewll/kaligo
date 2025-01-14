package database

type WhereParam struct {
    column  string
    op      string
    value   any
}

// Where is the struct for MySQL DATE type
type Where struct {
    params   map[string][]WhereParam   // map["AND"][] WhereParam{column, op, value}
    values   []any
    orderBys [][2]string
    limit    int
}

func (q *Query) WhereWrapper(wrapper func(*Query)) *Query {
    wrapper(q)
    return q
}

// Where Alias of andWhere
func (q *Query) Where(column string, op string, value any) *Query {
    return q.AndWhere(column, op, value)
}

// AndWhere Creates a new "AND WHERE" condition for the query.
// @title andWhere
// @description 查询条件函数
// @auth   seatle 2021/07/22 11:40
// @param  column string 字段名
// @param  op     string 操作符 >、=、<、>=、<=、IN、LIKE
// @param  value  string 查询值
// @return w     *Where Where对象
func (q *Query) AndWhere(column string, op string, value any) *Query {
    if q.W.params == nil {
        q.W.params = make(map[string][]WhereParam)
    }
    q.W.params["AND"] = append(q.W.params["AND"], WhereParam{column, op, value})
    return q
}

// OrWhere Creates a new "OR WHERE" condition for the query.
func (q *Query) OrWhere(column string, op string, value string) *Query {
    q.W.params["OR"] = append(q.W.params["OR"], WhereParam{column, op, value})
    return q
}

// WhereOpen Alias of andWhereOpen
func (q *Query) WhereOpen() *Query {
    return q.AndWhereOpen()
}

// AndWhereOpen Opens a new "AND WHERE (...)" grouping.
func (q *Query) AndWhereOpen() *Query {
    q.W.params["AND"] = append(q.W.params["AND"], WhereParam{column: "("})
    return q
}

// OrWhereOpen Opens a new "OR WHERE (...)" grouping.
func (q *Query) OrWhereOpen() *Query {
    q.W.params["OR"] = append(q.W.params["OR"], WhereParam{column: "("})
    return q
}

// WhereClose Alias of andWhereClose
func (q *Query) WhereClose() *Query {
    return q.AndWhereClose()
}

// AndWhereClose Closes an open "AND WHERE (...)" grouping.
func (q *Query) AndWhereClose() *Query {
    q.W.params["AND"] = append(q.W.params["AND"], WhereParam{column: ")"})
    return q
}

// OrWhereClose Closes an open "OR WHERE (...)" grouping.
func (q *Query) OrWhereClose() *Query {
    q.W.params["OR"] = append(q.W.params["OR"], WhereParam{column: ")"})
    return q
}

// OrderBy Applies sorting with "ORDER By ..."
func (q *Query) OrderBy(column string, args ...string) *Query {
    var direction string    
    if len(args) > 0 {
        direction = args[0]
    } else {
        direction = "ASC"
    }
    q.W.orderBys = append(q.W.orderBys, [2]string{column, direction})
    return q
}

// Limit Return up to "LIMIT ..." results
func (q *Query) Limit(value int) *Query {
    q.W.limit = value
    return q
}

/* vim: set expandtab: */
