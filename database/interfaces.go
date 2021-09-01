package database

// Dialector database dialector
type Dialector interface {
	Name() string
	Initialize(*DB) error
    CurrentDatabase(db *DB) string
    ListDatabases(like string, db *DB) []string
    ListTables(like string, db *DB) []string
    ListColumns(table, like string, db *DB) []Column
    ListIndexes(table, like string, db *DB) []Indexes
    CreateDatabase(database, charset string, ifNotExists bool, db *DB) (err error)
    DropDatabase(database string, db *DB) (err error)
    CreateTable(table string, fields []map[string]interface{}, primaryKeys []string, ifNotExists bool, engine, charset string, foreignKeys []map[string]interface{}, db *DB) (err error)
    RenameTable(oldTable, newTable string, db *DB) (err error)
    DropTable(table string, db *DB) (err error)
    TruncateTable(table string, db *DB) (err error) 
    TableExists(table string, db *DB) bool
    FieldExists(table string, value interface{}, db *DB) bool
    CreateIndex(table string, indexColumns interface{}, indexName, index string, db *DB) (err error)
    RenameIndex(table, oldName, newName string, db *DB) (err error)
    DropIndex(table, indexName string, db *DB) (err error)
    AddForeignKey(table string, foreignKey []map[string]interface{}, db *DB) (err error)
    DropForeignKey(table string, fkName string, db *DB) (err error) 
    AddFields(table string, fields []map[string]interface{}, db *DB) error
    DropFields(table string, value interface{}, db *DB) error
    ModifyFields(table string, fields []map[string]interface{}, db *DB) error
    AlterFields(alterType string, table string, fields interface{}, db *DB) (err error) 
    ProcessFields(fields []map[string]interface{}, prefix string, db *DB) string

	//Migrator(db *DB) Migrator
	//DataTypeOf(*schema.Field) string
	//DefaultValueOf(*schema.Field) clause.Expression
	//BindVarTo(writer clause.Writer, stmt *Statement, v interface{})
	//QuoteTo(clause.Writer, string)
	//Explain(sql string, vars ...interface{}) string
}

