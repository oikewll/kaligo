/**
 * Realize a database operation class
 *
 * @copyright   (C) 2014  seatle
 * @lastmodify  2021-07-06
 *
 */

package mysql

import (
    "strings"
)

// Join is the struct for MySQL DATE type
type Join struct {
    joinType  string
    table     string
    alias     string
    onValues  [][4]string

    Builder
}

// NewJoin creates a new Select Object
// @title NewJoin
// @description JOIN条件函数
// @auth   seatle 2021/07/22 11:40
// @param  table     interface{}  column name or []string{column, alias} or object 
// @param  joinType  string       type of JOIN: INNER, RIGHT, LEFT, etc
// @return *Join Join对象
func NewJoin(table string, joinType string) *Join {
    return &Join{
        table   : table,
        joinType: joinType,
    }
}

// On Alias of andOn
func (j *Join) On(c1 string, op string, c2 string) *Join {
    return j.AndOn(c1, op, c2)
}

// AndOn Adds a new AND condition for joining.
func (j *Join) AndOn(c1 string, op string, c2 string) *Join {
    j.onValues = append(j.onValues, [4]string{c1, op, c2, "AND"}) 
    return j
}

// OrOn Adds a new OR condition for joining.
func (j *Join) OrOn(c1 string, op string, c2 string) *Join {
    j.onValues = append(j.onValues, [4]string{c1, op, c2, "OR"}) 
    return j
}

// OnOpen Adds a opening bracket.
func (j *Join) OnOpen() *Join {
    j.onValues = append(j.onValues, [4]string{"", "", "", "("}) 
    return j
}

// OnClose Adds a closing bracket.
func (j *Join) OnClose() *Join {
    j.onValues = append(j.onValues, [4]string{"", "", "", ")"}) 
    return j
}

// Compile the SQL partial for a JOIN statement and return it.
func (j *Join) Compile(db *DB) string {

    var sqlStr string    

    // 操作句柄不存在的情况下
    //if db == nil {
        //db = dbConnection.instance(db)
    //}

    // j.joinType = "LEFT"
    if j.joinType == "" {
        sqlStr += strings.ToUpper(j.joinType) + " JOIN"
    } else {
        sqlStr += "JOIN"
    }

    // 子查询先不实现，太难了啊
    // JOIN (SELECT ks, COUNT(*) AS '# Tasks' FROM Table GROUP BY ks) t1) ON (t1.ks = t2.ks)
    //s := new(Select)
    //if s != nil {
        //// Compile the subquery and add it
        //sqlStr += " (" + s.compile() + ")"
    //} elseif expression != nil {
        // Compile the expression and add its value
        //sqlStr += " " + trim(expression.value(), " ()") + ")"
    //} else {
        // Quote the table name that is being joined
        //sqlStr += " " + quoteTable(j.table)
    //}
    sqlStr += " " + quoteTable(j.table)

    // Add the alias if needed
    if j.alias != "" {
        sqlStr += " AS " + quoteTable(j.alias)
    }

    var conditions []string    

    for _, condition := range j.onValues {
        c1 := condition[0]
        op := condition[1]
        c2 := condition[2]
        ch := condition[3]  // chaining

        cString := c1 + op + c2

        // Just a chaining character?
        if cString == "" {
            conditions = append(conditions, ch)
        } else {
            // Check if we have a pending bracket open
            if conditions[len(conditions)-1] == "(" {
                // Update the chain type
                conditions[len(conditions)-1] = " " + ch + " ("
            } else {
                // Just add chain type
                conditions = append(conditions, " " + ch + " ")
            }

            if op != "" {
                // Make the operator uppercase and spaced
                op = " " + strings.ToUpper(op)
            }

            // Quote each of the identifiers used for the condition
            var c2Str string    
            if c2 == "" {
                c2Str += "NULL"
            } else {
                //c2Str += db.quoteIdentifier(c2)
                c2Str += c2
            }
            conditions = append(conditions, c1 + op + " " + c2Str)
        }
    }
    
    // remove the first chain type
    conditions = conditions[1:] // 删除开头1个元素，直接移动数据指针方式
    //conditions = append(conditions[:0], conditions[1:]...)    // 后面的数据向开头移动

    // if there are conditions, concat the conditions "... AND ..." and glue them on...
    if conditions != nil {
        sqlStr += " ON (" + strings.Join(conditions, "") + ")"
    }

    return sqlStr
}

// Reset the query parameters
//func (j *Join) Reset() *Join {
    //return j
//}