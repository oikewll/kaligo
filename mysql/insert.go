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
    //values  [][]string // 不要用二维slice，非常蛋疼，改map把
    values  [][]string // 不要用二维slice，非常蛋疼，改map把
    params  []string
}

// Set the table to insert into.
func (i *Insert) setTable(table string) *Insert {
    i.table = table
    return i
}

// Set the columns that will be inserted.
func (i *Insert) setColumns(columns []string) *Insert {
    i.columns = append(i.columns, columns...)   // append和后面的...用法，相当于 php里面的array_meage函数
    return i
}

// Adds values. Multiple value sets can be added.
func (i *Insert) setValues(values [][]string) *Insert {
    i.values = append(i.values, values...)
    return i
}

func (i *Insert) set(pairs map[string]string) *Insert {
    //i.setColumns(util.Arr.MapKeys(pairs))
    return i
}

// Reset the query parameters
func (i *Insert) reset() *Insert {
    i.table = ""
    i.columns = nil // 垃圾回收器会自动回收原有的数据，len(), cap() 都为0，序列化成 json 的时候，为 null
    //i.columns = i.columns[:0] // len() 为0， cap() 是原来的长度，json 为 []
    i.values = nil
    i.params = nil
    return i
}

