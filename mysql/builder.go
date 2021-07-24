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
    "strconv"
)

// Builder is the struct for MySQL DATE type
type Builder struct {
    Query
}

// CompileJoin Compiles an array of JOIN statements into an SQL partial.
func (b *Builder) CompileJoin(db *DB, joins []*Join) string {
    var statements []string    
    for _, join := range joins {
        statements = append(statements, join.Compile(db))
    }

    return strings.Join(statements, ", ")
}


// CompileConditions Compiles an array of conditions into an SQL partial. Used for WHERE and HAVING
func (b *Builder) CompileConditions(db *DB, conditions map[string][][]string) string {
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
                    if b.parameters[min] != "" {
                        // Set the parameter as the minimum
                        min = b.parameters[min]
                    }

                    if b.parameters[max] != "" {
                        // Set the parameter as the maximum
                        max = b.parameters[max]
                    }

                    value = quote(min) + " AND " + quote(max)
                } else {
                    if b.parameters[value] != "" {
                        // Set the parameter as the value
                        value = b.parameters[value]
                    }
                    
                    // Quote the entire value normally
                    value = quote(value)
                }

                // Append the statement to the query
                sqlStr += quoteIdentifier(column) + " " + op + " " + value
            }

            lastCondition = strings.Join(condition, "")
        }
    }

    return sqlStr
}

// CompileSet Compiles an array of set values into an SQL partial. Used for UPDATE
func (b *Builder) CompileSet(db *DB, values [][]string) string {
    var sets []string    
    for _, group := range values {
        column := group[0]
        value  := group[1]

        column = quoteIdentifier(column)
        value  = quote(value)

        sets = append(sets, column + " = " + value)
    }

    return strings.Join(sets, ", ")
}

// CompileOrderBy Compiles an array of ORDER BY statements into an SQL partial..
func (b *Builder) CompileOrderBy(db *DB, columns [][2]string) string {

    var sorts []string    
    for _, group := range columns {
        column    := group[0]
        direction := group[1]

        column    = quoteIdentifier(column)
        direction = quote(direction)

        direction = strings.ToUpper(direction)
        if direction != "" {
            // Make the direction uppercase
            direction = " " + direction
        }

        sorts = append(sorts, column + direction)
    }

    return "ORDER BY " + strings.Join(sorts, ", ")
}

// Reset the query parameters
//func (b *Builder) Reset() *Builder {
    //return b
//}
