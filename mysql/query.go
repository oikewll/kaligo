package mysql

import (
    "context"
    "database/sql"
    //"fmt"
    "strings"
    "regexp"
    "reflect"
    //"os"
)

// Query is the struct for MySQL DATE type
type Query struct {
    *DB
    Context     context.Context
    Schema      *Schema
    stdDB       *sql.DB       // MySQL connection
    stdTX       *sql.Tx       // MySQL connection for Transaction

    S *Select
    W *Where
    J *Join
    B *Builder
    I *Insert
    U *Update
    D *Delete
    R *Result

    TablePrefix   string
    Dest          interface{}        //
    Model         interface{}        // Object：&User{}
    ReflectValue  reflect.Value
    sqlStr        string             // SQL statement
    queryType     QueryType          // Query type
    lifeTime      int                // Cache lifetime
    cacheKey      string             // Cache key
    cacheAll      bool               // boolean Cache all results

    joinObjs      []*Join            // join objects
    lastJoin      *Join              // last join statement
    parameters    map[string]string  // Quoted query parameters
}

// QueryType get the type of the query
func (q *Query) QueryType() QueryType {
    return q.queryType
}

// Cached Enables the query to be cached for a specified amount of time.
func (q *Query) Cached(lifeTime int, cacheKey string, cacheAll bool) *Query {
    q.lifeTime = lifeTime
    q.cacheKey = cacheKey
    q.cacheAll = cacheAll
    return q
}

// Scan is ...
func (q *Query) Scan(value interface{}) *Query {
    q.Dest = value

    // assign query.ReflectValue
    if q.Dest != nil {
        q.ReflectValue = reflect.ValueOf(q.Dest)
        for q.ReflectValue.Kind() == reflect.Ptr {
            if q.ReflectValue.IsNil() && q.ReflectValue.CanAddr() {
                q.ReflectValue.Set(reflect.New(q.ReflectValue.Type().Elem()))
            }

            q.ReflectValue = q.ReflectValue.Elem()
            // assign model values，只有 Struct 才给 q.Model 赋值
            if q.ReflectValue.Kind() == reflect.Struct {
                // var user User
                q.Model = q.Dest
            } else if q.ReflectValue.Kind() == reflect.Slice {
                // var users  []User
                // users   := []User{}
                // results := []map[string]interface{}{}
                if q.ReflectValue.Type().Elem().Kind() == reflect.Struct {
                    // var users []User
                    q.Model = q.Dest
                }
            }
        }
        if !q.ReflectValue.IsValid() {
            q.AddError(ErrInvalidValue)
        }
    }

    return q
}

// AsAssoc Returns results as associative arrays
//func (q *Query) AsAssoc() *Query { }
// AsObject Returns results as objects.
//func (q *Query) AsObject(value interface{}) *Query { }

// Bind a variable to a parameter in the query.
func (q *Query) Bind(param string, value string) *Query {
    if q.parameters == nil {
        q.parameters = make(map[string]string)
    }
    // Bind a value to a variable
    q.parameters[param] = value
    return q
}

// Parameters Add multiple parameters to the query.
func (q *Query) Parameters(params map[string]string) *Query {
    // Merge the new parameters in
    for param, value := range params {
        q.Bind(param, value)
    }
    return q
}

// Compile the SQL query and return it. Raplaces and parameters with their
// @return result Result DatabaseResult for SELECT queries
// @return result interface{} the insert id for INSERT queries
// @return result integer number of affected rows for all other queries
func (q *Query) Compile() string {
    var sqlStr string

    switch q.queryType {
    case SELECT:
        q.SelectCompile()
    case INSERT:
        q.InsertCompile()
    case UPDATE:
        q.UpdateCompile()
    case DELETE:
        q.DeleteCompile()
    default:
    }

    // Import the SQL locally
    sqlStr = q.sqlStr

    if q.parameters != nil {
        // Quote all of the values
        values := make(map[string]string, len(q.parameters))
        for k, v := range q.parameters {
            // 如果前面没有:，前面加 :，用于替换
            if k[0:1] != ":" {
                k = ":" + k
            }
            values[k] = q.Quote(v)
        }
        // Replace the values in the SQL
        sqlStr = Strtr(sqlStr, values)
    }

    // 不需要了, 一个Query()一个对象，db.query = &Query{} 以后之前那个就会被回收掉了
    //q.Reset()

    return strings.TrimSpace(sqlStr)
}

// Execute the current query on the given database.
func (q *Query) Execute() *Query {
    var err error

    // Compile the SQL query
    sqlStr := q.Compile()
    //fmt.Printf("Execute sqlStr = %v\n", sqlStr)

    // make sure we have a SQL type to work with
    if q.queryType == 0 && len(sqlStr) >= 11 {
        // get the SQL statement type without having to duplicate the entire statement
        stmt := regexp.MustCompile(`[\s]+`).Split(strings.TrimLeft(sqlStr[0:11], "("), 2)
        switch strings.ToUpper(stmt[0]) {
        case "DESCRIBE", "EXECUTE", "EXPLAIN", "SELECT", "SHOW":
            q.queryType = SELECT
        case "INSERT":
            q.queryType = INSERT
        case "UPDATE":
            q.queryType = UPDATE
        case "DELETE":
            q.queryType = DELETE
        default:
            q.queryType = 0
        }
    }

    // parse model values
    if q.Model != nil {
        if q.Schema, err = Parse(q.Model, q.cacheStore); err != nil {
            q.AddError(err)
        }
    }

    // 处理查询缓存
    //cacheObj = cache.Forge(cacheKey)
    //if conn.Caching() && q.lifeTime != 0 && q.queryType == SELECT {
        //var cacheKey string
        //if q.cacheKey == "" {
            //h := md5.New()
            //io.WriteString(h, "Connection.Query(\"" + sqlStr + "\")")
            //cacheKey += fmt.Sprintf("db.%x", h.Sum(nil))
            //return db.Cache(cache.Get(cacheKey), sqlStr, q.AsObject)
        //} else {
            //cacheKey += q.cacheKey
        //}
    //}

    // Execute the query
    q.queryCount++

    rows, err := q.stdDB.Query(sqlStr)
    if err != nil {
        q.AddError(err)
        return q
    }
    defer rows.Close();

    Scan(rows, q)

    //Cache the result if needed
    //if  cacheObj != nil && (q.cacheAll || result.count() != 0) {
        //cacheObj.setExpiration(q.lifeTime).SetContents(result.asArray()).Set()
    //}

    // 记录日志
    //db.Logger.Trace(stmt.Context, curTime, func() (string, int64) {
        //return db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...), db.RowsAffected
    //}, db.Error)

    return q
}

// First is the First record
// 按照主键顺序的第一条记录
func (q *Query) First(value interface{}) *Query {
    return q
}

// Last is the Last record
// 按照主键顺序的最后一条记录
func (q *Query) Last(value interface{}) *Query {
    return q
}

// Find is the all records
// 所有记录
func (q *Query) Find(value interface{}) *Query {
    return q
}

// Reset the query parameters
func (q *Query) Reset() *Query {
    switch q.queryType {
    case SELECT:
        q.SelectReset()
    case INSERT:
        q.InsertReset()
    case UPDATE:
        q.UpdateReset()
    case DELETE:
        q.DeleteReset()
    default:
    }

    //q.Dest      = nil
    q.sqlStr    = ""
    q.queryType = 0
    q.lifeTime  = 0
    q.cacheKey  = ""
    q.cacheAll  = false

    // 这里不需要清除，由调用的去清除，SELECT、INSERT、UPDATE、DELETE这些reset去清除
    //q.joinObjs   = nil
    //q.lastJoin   = nil
    //q.parameters = nil

    return q
}
