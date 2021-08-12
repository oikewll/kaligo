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
    "strings"
    "strconv"
)

// Builder is the struct for MySQL DATE type
type Builder struct {
    //*Query
}

// CompileJoin Compiles an array of JOIN statements into an SQL partial.
//func (q *Query) CompileJoin(c *Connection, joins []*Join) string {
func (q *Query) CompileJoin(joins []*Join) string {
    var statements []string    

    for _, join := range joins {
        statements = append(statements, join.Compile(q.C))
    }

    return strings.Join(statements, " ")
}

// CompileConditions Compiles an array of conditions into an SQL partial. Used for WHERE and HAVING
func (q *Query) CompileConditions(conditions map[string][][]string) string {
    var lastCondition string    
    var sqlStr string    

    for _, group := range conditions {
        // Process groups of conditions
        for logic, condition := range group {
            conditionStr := strings.Join(condition, "")
            if conditionStr == "(" {
                if sqlStr != "" && lastCondition != "(" {
                    // Include logic operator
                    sqlStr += " " + strconv.Itoa(logic) + " "
                }

                sqlStr += "("
            } else if conditionStr == ")" {
                sqlStr += ")"
            } else {
                if sqlStr != "" && lastCondition != "(" {
                    // Add the logic operator
                    sqlStr += " " + strconv.Itoa(logic) + " "
                }

                // Split the condition 
                // 'name', '=', 'John'
                // 'age', 'BETWEEN', '10,20'
                column := condition[0]
                op     := condition[1]
                value  := condition[2]
                // value 传 NULL 字符串过来，要把等号(=)换成 IS
                if value == "NULL" {
                    if op == "=" {
                        // Convert "val = NULL" to "val IS NULL"
                        op = "IS"
                    } else if op == "!=" {
                        // Convert "val != NULL" to "val IS NOT NULL"
                        op = "IS NOT"
                    }
                }

                // Database operators are always uppercase
                op = strings.ToUpper(op)

                //if (op == "BETWEEN" || op == "NOT BETWEEN") && is_array(value) {
                if op == "BETWEEN" || op == "NOT BETWEEN" {
                    ////BETWEEN always has exactly two arguments
                    valueArr := strings.Split(value, ",")
                    // trim一下兼容有空格写法：10,20 和 10, 20 都兼容
                    min := strings.TrimSpace(valueArr[0])
                    max := strings.TrimSpace(valueArr[1])
                    if q.parameters[min] != "" {
                        // Set the parameter as the minimum
                        min = q.parameters[min]
                    }

                    if q.parameters[max] != "" {
                        // Set the parameter as the maximum
                        max = q.parameters[max]
                    }

                    value = q.C.Quote(min) + " AND " + q.C.Quote(max)
                } else {
                    if q.parameters[value] != "" {
                        // Set the parameter as the value
                        value = q.parameters[value]
                    }
                    
                    // Quote the entire value normally
                    value = q.C.Quote(value)
                }

                // Append the statement to the query
                sqlStr += q.C.QuoteIdentifier(column) + " " + op + " " + value
            }

            lastCondition = conditionStr
        }
    }

    return sqlStr
}

// CompileSet Compiles an array of set values into an SQL partial. Used for UPDATE
func (q *Query) CompileSet(values [][]string) string {
    var sqlStr string    

    dict := make(map[string]string)
    for _, group := range values {
        // Split the set
        column := group[0]
        value  := group[1]

        // Quote the column name
        column = q.C.QuoteIdentifier(column)

        if val, ok := dict[value]; ok {
            // Use the parameter value
            value = val
        }
        dict[column] = column + "=" + q.C.Quote(value)
    }

    var sets []string    
    for _, v := range dict {
        sets = append(sets, v)
    }
    sqlStr = strings.Join(sets, ", ")
    return sqlStr
}

// CompileOrderBy Compiles an array of ORDER BY statements into an SQL partial..
func (q *Query) CompileOrderBy(columns [][2]string) string {
    var sorts []string    

    for _, group := range columns {
        // Split the orderby
        column    := group[0]
        direction := group[1]

        direction = strings.ToUpper(direction)
        if direction != "" {
            // Make the direction uppercase
            direction = " " + direction
        }

        sorts = append(sorts, q.C.QuoteIdentifier(column) + direction)
    }

    return "ORDER BY " + strings.Join(sorts, ", ")
}

