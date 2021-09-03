package database

import (
    "fmt"
    "strings"
)

// Insert is the struct for MySQL DATE type
type Insert struct {
    table    string
    columns  []string
    values   [][]string
    subQuery string     // 子查询
}

// Columns Set the columns that will be inserted.
func (q *Query) Columns(columns []string) *Query {
    q.I.columns = append(q.I.columns, columns...)   // append和后面的...用法，相当于 php里面的array_meage函数
    return q
}

// Values Adds values. Multiple value sets can be added.
func (q *Query) Values(values interface{}) *Query {
    switch vals := values.(type) {
    case []string:
        q.I.values = append(q.I.values, vals)
    case [][]string:
        for _, v := range vals {
            q.I.values = append(q.I.values, v)
        }
    default:
        //fmt.Println("Unknow Type")
    }
    return q
}

// SetValues is a warpper function for calling Columns() and Values().
func (q *Query) SetValues(pairs map[string]string) *Query {
    var keys []string    
    var vals []string    
    for k, v := range pairs {
        keys = append(keys, k)
        vals = append(vals, v)
    }
    q.Columns(keys)
    q.Values( vals)
    return q
}

// SubSelect the query parameters
func (q *Query) SubSelect(query *Query) *Query {
    if query.queryType != SELECT {
        panic("Only SELECT queries can be combined with INSERT queries")
    }
    q.I.subQuery = query.sqlStr
    return q
}

// InsertCompile Compile the SQL query and return it.
func (q *Query) InsertCompile() string {
    var sqlStr string    
    table   := q.I.table
    columns := q.I.columns

    // Start and update query
    sqlStr = "INSERT INTO " + q.QuoteTable(table)

    if len(columns) != 0 {
        columns = arrayUnique(columns)
        //fmt.Printf("table = %v; columns = %v; cryptFields = %v\n", table, FormatJSON(columns), q.cryptFields)
        // Add the column names
        for k, v := range columns {
            columns[k] = q.QuoteIdentifier(v)
        }
        sqlStr += " (" + strings.Join(columns, ", ") + ") "
    } else {
        sqlStr += " "
    }

    if q.I.subQuery == "" {
        var groups []string    
        for _, group := range q.I.values {
            for k, v := range group {
                if q.parameters[v] != "" {
                    // Use the parameter value
                    group[k] = q.parameters[v]
                }

                column := q.I.columns[k]
                // Is the column need encrypt ???
                if cryptFields, ok := q.cryptFields[table]; ok && q.cryptKey != "" && InSlice(column, &cryptFields) {
                    group[k] = fmt.Sprintf("AES_ENCRYPT(%s, \"%s\")", q.Quote(v), q.cryptKey)
                } else {
                    group[k] = q.Quote(v)
                }
            }
            
            groups = append(groups, "(" + strings.Join(group, ", ") + ")")
        }

        // Add the values
        sqlStr += "VALUES " + strings.Join(groups, ", ")
    } else {
        // Add the sub-query
        sqlStr += q.I.subQuery
    }

    //fmt.Printf("InsertCompile === %v\n", sqlStr)
    q.sqlStr = sqlStr

    return sqlStr
}

// InsertReset the query parameters
func (q *Query) InsertReset() *Query {
    //fmt.Println("InsertReset")
    q.I.table    = ""
    q.I.columns  = nil // gc 回收原有数据，len(), cap() 都为0，序列化成 json 的时候，为 null，如果是 columns[:0] 则 gc 不回收，len() 为0， cap() 不变，json 为 []
    q.I.values   = nil

    q.parameters = nil

    return q
}

