package database

import (
    "context"
    "database/sql"
    "fmt"
    "reflect"
    "regexp"
    "strings"
    "time"

    // "github.com/owner888/kaligo/util"
)

// Query is the struct for MySQL DATE type
// SQL Builder
type Query struct {
    *DB

    Error     error // Global error
    Context   context.Context
    Schema   *Schema

    RowsAffected int64 // For select、update、insert
    LastInsertId int64 // Only for insert

    S *Select
    W *Where
    J *Join
    B *Builder
    I *Insert
    U *Update
    D *Delete
    //R *Result

    sqlStr    string            // SQL statement
    queryType QueryType         // Query type

    Dest         any            // var user User、var users []User、var result map[string]any、var results []map[string]any、var ages []int64
    Model        any            // Object：&User{}
    ReflectValue reflect.Value  // reflect.ValueOf(Dest)
    lifeTime     int            // Cache lifetime
    cacheKey     string         // Cache key
    cacheAll     bool           // boolean Cache all results

    joinObjs   []*Join          // join objects
    lastJoin   *Join            // last join statement
    parameters map[string]any   // query bind parameters, use for eg: bind(":id", 1) 
}

// QueryType get the type of the query
func (q *Query) QueryType() QueryType {
    return q.queryType
}

// AddError add error to db
// XXX 这个全局有点多余
func (q *Query) AddError(err error) error {
    if q.Error == nil {
        q.Error = err
    } else if err != nil {
        q.Error = fmt.Errorf("%v; %w", q.Error, err)
    }
    return q.Error
}

// Cached Enables the query to be cached for a specified amount of time.
func (q *Query) Cached(lifeTime int, cacheKey string, cacheAll bool) *Query {
    q.lifeTime = lifeTime
    q.cacheKey = cacheKey
    q.cacheAll = cacheAll
    return q
}

// Bind a variable to a parameter in the query.
func (q *Query) Bind(param string, value any) *Query {
    if q.parameters == nil {
        q.parameters = make(map[string]any)
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

// SetCryptKey Set crypt key.
func (q *Query) SetCryptKey(cryptKey string) *Query {
    q.cryptKey = cryptKey
    return q
}

// SetCryptFields Set crypt fields.
func (q *Query) SetCryptFields(cryptFields map[string][]string) *Query {
    q.cryptFields = cryptFields
    return q
}

// Alias Compile method
func (q *Query) String() string {
    return q.sqlStr
}

// Compile the SQL query and return it. Raplaces and parameters with their
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
    q.sqlStr = strings.TrimSpace(q.sqlStr)
    sqlStr = q.sqlStr

    // 只用于 Query 请求的 Bind 方法
    // _, err := db.Query("SELECT `age` FROM `user` WHERE `id` = :id").Bind(":id", 1).Scan(&ages).Execute()
    if q.parameters != nil {
        // Quote all of the values
        values := make(map[string]any, len(q.parameters))
        for k, v := range q.parameters {
            // 如果前面没有:，前面加 :，用于替换
            if k[0:1] != ":" {
                k = ":" + k
            }
            // values[k] = q.Quote(v)
            values[k] = "?"
            q.W.values = append(q.W.values, v)
        }
        // Replace the values in the SQL
        // 把 :id 换成 ?
        sqlStr = Strtr(sqlStr, values)
    }

    return sqlStr
}

// Execute the current query on the given database.
func (q *Query) Execute() (*Query, error) {
    var err error
    var sqlStr string // Compile SQL

    // 当前函数结束时如果有错误则打印日志
    defer func() {
        if err != nil {
            logs.Error(err)
        }
    }()

    curTime := time.Now()   // current timestamp
    sqlStr   = q.Compile()  // Compile the SQL query

    var needReset bool = true   
    // make sure we have a SQL type to work with
    if q.queryType == 0 && len(sqlStr) >= 11 {
        // get the SQL statement type without having to duplicate the entire statement
        stmt := regexp.MustCompile(`[\s]+`).Split(strings.TrimLeft(sqlStr[0:11], "("), 2)
        switch strings.ToUpper(stmt[0]) {
        case "DESCRIBE", "EXECUTE", "EXPLAIN", "SHOW":
            q.queryType = SELECT
            needReset = false
        case "SELECT":
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

    //fmt.Printf("Execute sqlStr = %v; queryType = %v\n", sqlStr, q.queryType)

    // parse model values
    if q.Model != nil {
        if q.Schema, err = Parse(q.Model, q.cacheStore); err != nil {
            return q, err
        }
    }

    // 处理查询缓存
    // cacheObj = cache.Forge(cacheKey)
    // if conn.Caching() && q.lifeTime != 0 && q.queryType == SELECT {
    //     var cacheKey string
    //     if q.cacheKey == "" {
    //         h := md5.New()
    //         io.WriteString(h, "Connection.Query(\"" + sqlStr + "\")")
    //         cacheKey += fmt.Sprintf("db.%x", h.Sum(nil))
    //         return db.Cache(cache.Get(cacheKey), sqlStr, q.AsObject)
    //     } else {
    //         cacheKey += q.cacheKey
    //     }
    // }

    // Execute the query
    q.queryCount++

    logs.Debug(sqlStr)

    // 生成预处理 SQL
    var stmt *sql.Stmt
    stmt, err = q.StdDB.Prepare(sqlStr)
    if err != nil {
        return q, err
    }
    defer stmt.Close()

    // 执行预处理语句
    if q.queryType == SELECT {  // 执行 Rows()
        var rowVars []any    // Prepare(sqlstr).Exec(Vars...)    // Bind(":id", "1")
        if q.W != nil {
            rowVars = q.W.values
        }

        var rows *sql.Rows
        rows, err = q.Rows(sqlStr, rowVars...)
        if err != nil {
            return q, err
        }
        defer rows.Close()
        Scan(rows, q)
    } else {    // 执行 Exec()，// INSERT & DELETE & UPDATE
        var rowsAffected int64 = 0
        var lastInsertID int64 = 0
        var result sql.Result

        if q.queryType == INSERT {
            var insVars [][]any    
            for _, group := range q.I.values {
                var Vars []any    
                for k, v := range group {
                    column := q.I.columns[k]
                    // Is the column need encrypt ???
                    if cryptFields, ok := q.cryptFields[q.I.table]; ok && q.Dialector.Name() == "mysql" && q.cryptKey != "" && InSlice(column, &cryptFields) {
                        Vars = append(Vars, v, q.cryptKey)  // 加密的方式会多一个 key
                    } else {
                        Vars = append(Vars, v)
                    }
                }
                insVars = append(insVars, Vars)
            }
            for _, value := range insVars {
                result, err = stmt.Exec(value...)
            }

            lastInsertID, err = result.LastInsertId()
            rowsAffected, err = result.RowsAffected()

            
            // logs.Trace(q.DB, curTime, func() (string, int64) {
            //     return Explain(sqlStr, util.CastSlice[[]any, any](insVars)...), q.RowsAffected
            // }, q.Error)

            // 日志需要支持 [][]any，目前只支持 []any 类型
            for _, Vars := range insVars {
                logs.Trace(q.DB, curTime, func() (string, int64) {
                    return Explain(sqlStr, Vars...), q.RowsAffected
                }, q.Error)
            }
        } else {    // DELETE & UPDATE
            var Vars []any    
            if q.W != nil {
                Vars = q.W.values
            }

            result, err = stmt.Exec(Vars...)
            lastInsertID, err = result.LastInsertId()
            rowsAffected, err = result.RowsAffected()

            logs.Trace(q.DB, curTime, func() (string, int64) {
                return Explain(sqlStr, Vars...), q.RowsAffected
            }, q.Error)
        }

        q.RowsAffected = rowsAffected
        q.LastInsertId = lastInsertID

        if err != nil {
            return q, err
        }
    }

    // Cache the result if needed
    // if  cacheObj != nil && (q.cacheAll || result.count() != 0) {
    //     cacheObj.setExpiration(q.lifeTime).SetContents(result.asArray()).Set()
    // }

    // 记录日志
    // 3.388208ms [0] INSERT INTO `user` (`name`, `age`) VALUES ( AES_ENCRYPT( "test", "aaa" ), AES_ENCRYPT( "20", "aaa" ) )
    // db.Logger.Trace(stmt.Context, curTime, func() (string, int64) {
    //     return db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...), db.RowsAffected
    // }, db.Error)

    // "DESCRIBE", "EXECUTE", "EXPLAIN", "SHOW" 这几个不需要 Reset，都是完整的 SQL
    q.Reset(needReset)

    return q, nil
}

func (q *Query) Count(value any) *Query {
    return q
}

// Scan is ...
func (q *Query) Scan(value any) *Query {
    if value == nil {
        return q
    }

    q.Dest = value
    // assign query.ReflectValue
    q.ReflectValue = reflect.ValueOf(q.Dest)
    //fmt.Printf("11111 ---> %T = %v\n", q.ReflectValue, q.ReflectValue)
    //fmt.Printf("22222 ---> %T = %v\n", q.ReflectValue.Kind(), q.ReflectValue.Kind())
    //fmt.Printf("33333 ---> %T = %v\n", q.ReflectValue.Elem(), q.ReflectValue.Elem())
    //fmt.Printf("44444 ---> %T = %v\n", q.ReflectValue.Elem().Kind(), q.ReflectValue.Elem().Kind())
    //fmt.Printf("55555 ---> %T = %v\n", q.ReflectValue.Type().Elem().Kind(), q.ReflectValue.Type().Elem().Kind())
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
            // var ages    []int64
            // var users   []User
            // var results []map[string]any
            // Slice 子元素的类型
            if q.ReflectValue.Type().Elem().Kind() == reflect.Struct {
                // var users []User
                q.Model = q.Dest
            } else if q.ReflectValue.Type().Elem().Kind() != reflect.Map {
                // var ages []int64、var names []string 这些slice类型会到这里来
                // 这里先初始化，因为不会去到 Parse() 了，scan.go 里面会报错
                q.Schema = &Schema{}
            }
        }
    }
    if !q.ReflectValue.IsValid() {
        q.AddError(ErrInvalidValue)
    }
    return q
}

// First is the First record
// 按照主键顺序的第一条记录
func (q *Query) First(value any) *Query {
    return q
}

// Last is the Last record
// 按照主键顺序的最后一条记录
func (q *Query) Last(value any) *Query {
    return q
}

// Find is the all records
// 所有记录
func (q *Query) Find(value any) *Query {
    return q
}

// Reset the query parameters
func (q *Query) Reset(needReset bool) *Query {
    if needReset {
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
    }

    //q.Dest      = nil
    q.sqlStr = ""
    q.queryType = 0
    q.lifeTime = 0
    q.cacheKey = ""
    q.cacheAll = false
    q.cryptKey = ""
    q.cryptFields = nil

    // 这里不需要清除，由调用的去清除，SELECT、INSERT、UPDATE、DELETE 这些 Reset() 去清除
    //q.joinObjs   = nil
    //q.lastJoin   = nil
    //q.parameters = nil

    return q
}

/* vim: set expandtab: */
