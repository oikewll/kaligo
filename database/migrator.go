package database

import (
    "strings"
)

type Migrator struct {
    *DB
}

// CurrentDatabase is Current Database
func (m *Migrator) CurrentDatabase() (name string) {
    return m.Dialector.CurrentDatabase(m.DB)
}

// ListDatabases If a database name is given it will return the database name with the configured
// prefix. If not, then just the prefix is returnd
func (m *Migrator) ListDatabases(args ...string) ([]string, error) {
    like := ""
    if len(args) > 0 {
        like = args[0]
    }
    return m.Dialector.ListDatabases(like, m.DB)
}

// ListTables If a table name is given it will return the table name with the configured
// prefix. If not, then just the prefix is returnd
func (m *Migrator) ListTables(args ...string) ([]string, error) {
    like := ""
    if len(args) > 0 {
        like = args[0]
    }
    return m.Dialector.ListTables(like, m.DB)
}

// ListColumns Lists all of the columns in a table. Optionally, a LIKE string can be
// used to search for specific fields.
func (m *Migrator) ListColumns(table string, args ...string) ([]Column, error) {
    like := ""
    if len(args) > 0 {
        like = args[0]
    }
    return m.Dialector.ListColumns(table, like, m.DB)
}

// ListIndexes Lists all of the idexes in a table. Optionally, a LIKE string can be
// used to search for specific indexes by name.
func (m *Migrator) ListIndexes(table string, args ...string) ([]Indexes, error) {
    like := ""
    if len(args) > 0 {
        like = args[0]
    }
    return m.Dialector.ListIndexes(table, like, m.DB)
}

// CreateDatabase Creates a database. Will throw a Database Exception if it cannot.
func (m *Migrator) CreateDatabase(database string, args ...any) (err error) {
    var (
        ifNotExists bool   = true
        charset     string = "utf8"
    )

    if len(args) > 1 {
        charset = args[0].(string)
        ifNotExists = args[1].(bool)
    } else if len(args) > 0 {
        ifNotExists = args[0].(bool)
    }

    return m.Dialector.CreateDatabase(database, charset, ifNotExists, m.DB)
}

// DropDatabase Drops a database. Will throw a Database Exception if it cannot.
func (m *Migrator) DropDatabase(database string) (err error) {
    return m.Dialector.DropDatabase(database, m.DB)
}

// CreateTable Creates a table.
func (m *Migrator) CreateTable(table string, fields []map[string]any, args ...any) (err error) {
    var (
        primaryKeys []string
        ifNotExists bool   = true
        engine      string = "InnoDB"
        charset     string = "utf8_general_ci"
        foreignKeys []map[string]any
    )

    switch len(args) {
    case 5:
        primaryKeys = args[0].([]string)
        ifNotExists = args[1].(bool)
        engine = args[2].(string)
        charset = args[3].(string)
        foreignKeys = args[4].([]map[string]any)
    case 4:
        primaryKeys = args[0].([]string)
        ifNotExists = args[1].(bool)
        engine = args[2].(string)
        charset = args[3].(string)
    case 3:
        primaryKeys = args[0].([]string)
        ifNotExists = args[1].(bool)
        engine = args[2].(string)
    case 2:
        primaryKeys = args[0].([]string)
        ifNotExists = args[1].(bool)
    case 1:
        primaryKeys = args[0].([]string)
    default:
    }

    err = m.Dialector.CreateTable(table, fields, primaryKeys, ifNotExists, engine, charset, foreignKeys, m.DB)
    m.AddError(err)
    return
}

// DropTable Drops a table. Will throw a Database Exception if it cannot.
func (m *Migrator) DropTable(table string) (err error) {
    return m.Dialector.DropTable(table, m.DB)
}

// RenameTable Renames a table. Will throw a Database Exception if it cannot.
func (m *Migrator) RenameTable(oldTable, newTable string) (err error) {
    return m.Dialector.RenameTable(oldTable, newTable, m.DB)
}

// TruncateTable Truncates a table.
func (m *Migrator) TruncateTable(table string) (err error) {
    return m.Dialector.TruncateTable(table, m.DB)
}

// AnalyzeTable Analyzes a table.
func (m *Migrator) AnalyzeTable(table string) bool {
    return m.tableMaintenance("ANALYZE TABLE", table)
}

// CheckTable Check a table.
func (m *Migrator) CheckTable(table string) bool {
    return m.tableMaintenance("CHECK TABLE", table)
}

// OptimizeTable Optimize a table.
func (m *Migrator) OptimizeTable(table string) bool {
    return m.tableMaintenance("OPTIMIZE TABLE", table)
}

// RepairTable Repair a table.
func (m *Migrator) RepairTable(table string) bool {
    return m.tableMaintenance("REPAIR TABLE", table)
}

func (m *Migrator) tableMaintenance(operation string, table string) bool {
    // SQLite does't support table maintenance
    if m.Dialector.Name() == "sqlite" {
        return false
    }

    sqlStr := operation + " " + m.QuoteTable(table)
    rows, err := m.Rows(sqlStr)
    if err != nil {
        m.AddError(err)
    }

    var null any
    var msgType, msgText string

    for rows.Next() {
        err := rows.Scan(null, null, msgType, msgText)
        if err != nil {
            m.AddError(err)
        }
    }

    if msgType == "status" && InSlice(strings.ToLower(msgText), &[]string{"ok", "table is already up to date"}) {
        return true
    }

    if InSlice(msgType, &[]string{"info", "warning", "error"}) {
        msgType = strings.ToUpper(msgType)
    } else {
        msgType = "INFO"
    }

    //logger(msgType, "Table: " + table + ", Operation: " + ops[0].Op + ", Message: " + msgText, 'Schema.tableMaintenance');

    return false
}

// TableExists Generic check if a given table exists.
func (m *Migrator) TableExists(table string) (bool, error) {
    return m.Dialector.TableExists(table, m.DB)
}

// FieldExists Checks if given field(s) in a given table exists.
func (m *Migrator) FieldExists(table string, value any) (bool, error) {
    return m.Dialector.FieldExists(table, value, m.DB)
}

// CreateIndex Creates an index on that table.
func (m *Migrator) CreateIndex(table string, indexColumns any, indexName, index string) (err error) {
    return m.Dialector.CreateIndex(table, indexColumns, indexName, index, m.DB)
}

// RenameIndex Rename an index from a table.
func (m *Migrator) RenameIndex(table, oldName, newName string) (err error) {
    return m.Dialector.RenameIndex(table, oldName, newName, m.DB)
}

// DropIndex Drop an index from a table.
func (m *Migrator) DropIndex(table string, indexName string) (err error) {
    return m.Dialector.DropIndex(table, indexName, m.DB)
}

// AddForeignKey Adds a single foreign key to a table
// player.userid(fk_userid foreign key) -> user.id，user.id needed index(primary key or unique key)
func (m *Migrator) AddForeignKey(table string, foreignKey []map[string]any) (err error) {
    return m.Dialector.AddForeignKey(table, foreignKey, m.DB)
}

// DropForeignKey Drops a foreign key from a table
func (m *Migrator) DropForeignKey(table string, fkName string) (err error) {
    return m.Dialector.DropForeignKey(table, fkName, m.DB)
}

// AddFields adds fields to a table.
func (m *Migrator) AddFields(table string, fields []map[string]any) error {
    return m.Dialector.AddFields(table, fields, m.DB)
}

// DropFields drops fields from a table.
func (m *Migrator) DropFields(table string, value any) error {
    return m.Dialector.DropFields(table, value, m.DB)
}

// ModifyFields alters fields in a table.
func (m *Migrator) ModifyFields(table string, fields []map[string]any) error {
    return m.Dialector.ModifyFields(table, fields, m.DB)
}

// AlterFields is ...
//func (m *Migrator) AlterFields(alterType, table string, fields any) (err error) {
//return m.Dialector.AlterFields(alterType, table, fields, m.DB)
//}

// ProcessFields is ...
func (m *Migrator) ProcessFields(fields []map[string]any, prefix string) string {
    return m.Dialector.ProcessFields(fields, prefix, m.DB)
}

/* vim: set expandtab: */
