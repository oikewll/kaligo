package mysql

import (
    //"errors"
    //"fmt"
    //"strconv"
    //"strings"
    //"time"
    //"github.com/owner888/kaligo/util"
)

// Query is the struct for MySQL DATE type
type Query struct {
    queryType   int                 // Query type
    lifeTime    int                 // Cache lifetime
    cacheKey    string              // Cache key
    cacheAll    bool                // boolean Cache all results
    sqlStr      string              // SQL statement
    parameters  map[string]string   // Quoted query parameters
    asObject    interface{}         // Return results as associative arrays(map[string]string) or objects(&User{ID:1, Name: "sam"})
    connection  *Connection         // db connection, Include *sql.DB

    // select、insert、update、delete、builder、join 都有嵌入其他类，只有这个 Query 是独立的
}

// QueryObj is ...
//var QueryObj Query = Query{
    //queryType: SELECT,
    //cacheAll : false,
    //asObject : false,
//}

func init() {
}

// NewQuery Creates a new SQL query of the specified type.
func NewQuery(sqlStr string, queryType int) *Query {
    return &Query{
        sqlStr   : sqlStr,
        queryType: queryType,
    }
}

// QueryType get the type of the query
func (q *Query) QueryType(c1 string, op string, c2 string) int {
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
func (q *Query) AsObject(class interface{}) *Query {
    q.asObject = class

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
func (q *Query) SetConnection(params map[string]string) *Query {
    // Merge the new parameters in
    for param, value := range params {
        q.parameters[param] = value
    }

    return q
}
// Reset the query parameters
//func (q *Query) reset() *Query {
    //return q
//}
