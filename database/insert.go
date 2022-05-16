package database

import (
    "fmt"
    "strings"
)

// Insert is the struct for MySQL DATE type
type Insert struct {
    table    string
    columns  []string
    values   [][]any
    updates  []map[string]any
    subQuery string // 子查询
}

// Columns Set the columns that will be inserted.
func (q *Query) Columns(columns []string) *Query {
    q.I.columns = append(q.I.columns, columns...) // append 和后面的 ... 用法，相当于 php 里面的 array_meage() 函数

    return q
}

// Values Adds values. Multiple value sets can be added.
func (q *Query) Values(values any) *Query {
    switch vals := values.(type) {
    case []any:
        q.I.values = append(q.I.values, vals)
    case [][]any:
        q.I.values = append(q.I.values, vals...)
    default:
        logs.Error("Insert Values: Unknow Type")
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
    q.Values(vals)

    return q
}

// SubSelect the query parameters
func (q *Query) SubSelect(query *Query) *Query {
    if query.queryType != SELECT {
        logs.Panic("Only SELECT queries can be combined with INSERT queries")
    }
    q.I.subQuery = query.sqlStr

    return q
}

// ON DUPLICATE KEY UPDATE
func (q *Query) OnDuplicateKeyUpdate(updates any) *Query {
    switch vals := updates.(type) {
    case map[string]any:
        q.I.updates = append(q.I.updates, vals)
    case []map[string]any:
        q.I.updates = append(q.I.updates, vals...)
    default:
        logs.Error("Insert OnDuplicateKeyUpdate: Unknow Type")
    }

    return q
}

// InsertCompile Compile the SQL query and return it.
func (q *Query) InsertCompile() (sqlStr string) {
    var placeholders []string

    table   := q.I.table
    columns := q.I.columns

    // Start and insert query
    sqlStr = "INSERT INTO " + q.QuoteTable(table)

    if len(columns) != 0 {
        columns = arrayUnique(columns)
        for k, v := range columns {
            columns[k] = q.QuoteIdentifier(v)

            // Is the column need encrypt ???
            if cryptFields, ok := q.cryptFields[table]; ok && q.Dialector.Name() == "mysql" && q.cryptKey != "" && InSlice(v, &cryptFields) {
                placeholders = append(placeholders, "AES_ENCRYPT( ?, ? )")
            } else {
                placeholders = append(placeholders, "?")
            }
        }
        sqlStr += " (" + strings.Join(columns, ", ") + ")"
    }

    // logs.Debugf("table = %v; columns = %v; cryptFields = %v\n", table, FormatJSON(columns), q.cryptFields)

    // 子查询为空
    if q.I.subQuery == "" {
        sqlStr += " VALUES ( " + strings.Join(placeholders, ", ") + " )"
    } else {
        // Add the sub-query
        // INSERT INTO table1 ( `column1` ) SELECT `col1` FROM `table2`
        sqlStr += " " + q.I.subQuery
    }

    // INSERT INTO table(`UniqueKeyField`, `field1`, `field2) VALUES ("UniqueKeyFieldVal", "field1Val", "field2Val") ON DUPLICATE KEY UPDATE `field1` = "field1Val", `field2` = "field2Val";
    if q.I.updates != nil {
        var updates []string
        for _, sliceMap := range q.I.updates {
            for key, value := range sliceMap {
                updates = append(updates, fmt.Sprintf("`%s` = \"%s\"", key, value))
            }
        }
        sqlStr += " ON DUPLICATE KEY UPDATE " + strings.Join(updates, ", ")
    }

    q.sqlStr = sqlStr

    return sqlStr
}

// InsertReset the query parameters
func (q *Query) InsertReset() *Query {
    q.I.table    = ""
    q.I.columns  = nil // gc 回收原有数据，len(), cap() 都为0，序列化成 json 的时候，为 null，如果是 columns[:0] 则 gc 不回收，len() 为0， cap() 不变，json 为 []
    q.I.values   = nil
    q.I.updates  = nil
    q.parameters = nil

    return q
}
