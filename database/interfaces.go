package database

// Dialector database dialector
type Dialector interface {
	Name() string
	Initialize(*DB) error
    ListDatabases(like string, db *DB) []string
    ListTables(like string, db *DB) []string
    ListColumns(table string, like string, db *DB) []Column
    ListIndexes(table string, like string, db *DB) []Indexes
    DropIndex(table string, indexName string, db *DB) (err error)

	//Migrator(db *DB) Migrator
	//DataTypeOf(*schema.Field) string
	//DefaultValueOf(*schema.Field) clause.Expression
	//BindVarTo(writer clause.Writer, stmt *Statement, v interface{})
	//QuoteTo(clause.Writer, string)
	//Explain(sql string, vars ...interface{}) string
}

