package mysql

import (
    //"fmt"
)

// Schema is the struct for MySQL DATE type
type Schema struct {
    name string              // database connection config name
    C    *Connection         // database connection instance, Include *sql.DB
}

// CreateDatabase Creates a database. Will throw a Database Exception if it cannot.
func (s *Schema) CreateDatabase(database string, charset string, ifNotExists bool) int {
    sqlStr := "CREATE DATABASE "
    if ifNotExists {
        sqlStr += "IF NOT EXISTS "
    }
    //sqlStr += quoteIdentifier(database) + processCharset(charset, true)
    //return s.connection.query(0, sqlStr, false)
    return 10
}

// DropDatabase Drops a database. Will throw a Database Exception if it cannot.
func (s *Schema) DropDatabase(database string) int {
    //sqlStr := "DROP DATABASE "
    //sqlStr += quoteIdentifier(database) + processCharset(charset, true)
    //return s.connection.query(0, sqlStr, false)
    return 10
}

// DropTable Drops a table. Will throw a Database Exception if it cannot.
func (s *Schema) DropTable(database string) int {
    //sqlStr := "DROP TABLE IF EXISTS "
    //sqlStr += quoteIdentifier(database) + processCharset(charset, true)
    //return s.connection.query(0, sqlStr, false)
    return 10
}
