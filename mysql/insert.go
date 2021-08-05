package mysql

import (
    //"errors"
    //"fmt"
    //"strconv"
    //"strings"
    //"time"
    //"github.com/owner888/kaligo/util"
)

// Insert is the struct for MySQL DATE type
type Insert struct {
    table   string
    columns []string    // slice
    values  [][]string  // 多维slice
    //values  []map[string]string // map slice
    params  []string

    *Builder
}

// Table Set the table to insert into.
func (i *Insert) Table(table string) *Insert {
    i.table = table
    return i
}

// Columns Set the columns that will be inserted.
func (i *Insert) Columns(columns []string) *Insert {
    i.columns = append(i.columns, columns...)   // append和后面的...用法，相当于 php里面的array_meage函数
    return i
}

// Values Adds values. Multiple value sets can be added.
func (i *Insert) Values(values interface{}) *Insert {
    switch vals := values.(type) {
    case []string:
        i.values = append(i.values, vals)
    case [][]string:
        for _, v := range vals {
            i.values = append(i.values, v)
        }
    default:
        //fmt.Println("Unknow Type")
    }
    return i
}

// Set is a warpper function for calling Columns() and Values().
func (i *Insert) Set(pairs map[string]string) *Insert {
    var keys []string    
    var vals []string    
    for k, v := range pairs {
        keys = append(keys, k)
        vals = append(vals, v)
    }
    i.Columns(keys)
    i.Values( vals)
    return i
}

// Reset the query parameters
func (i *Insert) reset() *Insert {
    i.table = ""
    i.columns   = nil // 垃圾回收器会自动回收原有的数据，len(), cap() 都为0，序列化成 json 的时候，为 null
    //i.columns = i.columns[:0] // len() 为0， cap() 是原来的长度，json 为 []
    i.values    = nil
    i.params    = nil
    return i
}

