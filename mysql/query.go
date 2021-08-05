package mysql

import (
    "fmt"
    "strings"
    "regexp"
)

// Query is the struct for MySQL DATE type
type Query struct {
    sqlStr      string              // SQL statement
    queryType   int                 // Query type
    lifeTime    int                 // Cache lifetime
    cacheKey    string              // Cache key
    cacheAll    bool                // boolean Cache all results
    parameters  map[string]string   // Quoted query parameters
    asObject    interface{}         // Return results as associative arrays(map[string]string || []map[string]string) or objects(&User{ID:1, Name: "sam"})
    connection  *Connection         // db connection, Include *sql.DB

    // select、insert、update、delete、builder、join 都有嵌入其他类，只有这个 Query 是独立的
}

// QueryType get the type of the query
func (q *Query) QueryType() int {
    return q.queryType
}

// Cached Enables the query to be cached for a specified amount of time.
func (q *Query) Cached(lifeTime int, cacheKey string, cacheAll bool) *Query {
    q.lifeTime = lifeTime
    q.cacheKey = cacheKey
    q.cacheAll = cacheAll

    return q
}

// AsAssoc Returns results as associative arrays
func (q *Query) AsAssoc() *Query {
    q.asObject = false

    return q
}

// AsObject Returns results as objects.
func (q *Query) AsObject(value interface{}) *Query {
    q.asObject = value

    return q
}

// Param Set the value of a parameter in the query.
func (q *Query) Param(param string, value string) *Query {
    // Add or overload a new parameter
    q.parameters[param] = value

    return q
}

// Bind a variable to a parameter in the query.
func (q *Query) Bind(param string, value string) *Query {
    // Bind a value to a variable
    q.parameters[param] = value

    return q
}

// Parameters Add multiple parameters to the query.
func (q *Query) Parameters(params map[string]string) *Query {
    // Merge the new parameters in
    for param, value := range params {
        q.parameters[param] = value
    }

    return q
}

// SetConnection Set a DB Connection to use when compiling the SQL.
func (q *Query) SetConnection(c *Connection) *Query {
    q.connection = c
    return q
}

// Compile the SQL query and return it. Raplaces and parameters with their
// @return result Result DatabaseResult for SELECT queries
// @return result interface{} the insert id for INSERT queries
// @return result integer number of affected rows for all other queries
func (q *Query) Compile(args ...*Connection) string {
    var conn *Connection    
    if len(args) != 0 {
        conn = args[0]
    } else {
        // Get the database instance
        //db := New()
        //conn = db.C
        conn = q.connection
    }
    //fmt.Printf("Query Compile === %T = %p\n", conn, conn)

    // Import the SQL locally
    sqlStr := q.sqlStr
    fmt.Printf("Query Compile sqlStr === %v\n", sqlStr)

    if q.parameters != nil {
        // Quote all of the values
        values := q.parameters
        for k, v := range values {
            // 前面加 :
            if k[0:1] != ":" {
                k = ":" + k
            }
            values[k] = conn.Quote(v)
        }

        // Replace the values in the SQL
        sqlStr = Strtr(sqlStr, values)
    }

    return strings.TrimSpace(sqlStr)
}

// Execute the current query on the given database.
func (q *Query) Execute(args ...*Connection) string {
//func (q *Query) Execute(args ...*Connection) interface{} {
    var conn *Connection    
    if len(args) != 0 {
        conn = args[0]
    } else {
        // Get the database instance
        //db := New()
        //conn = db.C
        conn = q.connection
    }
    //fmt.Printf("Query Execute === %T = %p\n", conn, conn)

    //fmt.Printf("Execute sqlStr111 = %v\n", q.sqlStr)
    // Compile the SQL query
    sqlStr := q.Compile(conn)
    //fmt.Printf("Execute sqlStr222 = %v\n", sqlStr)

    // make sure we have a SQL type to work with
    if q.queryType == 0 {
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
    q.connection.queryCount++
    // Connection.Query()
    //result := conn.Query(q.queryType, sqlStr, q.AsObject)

    //Cache the result if needed
    //if  cacheObj != nil && (q.cacheAll || result.count() != 0) {
        //cacheObj.setExpiration(q.lifeTime).SetContents(result.asArray()).Set()
    //}

    return "hello"
    //return result
}

