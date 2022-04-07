package driver

import (
    "database/sql"
    "fmt"

    //"reflect"
    "regexp"
    //"strconv"
    "strings"

    _ "github.com/mattn/go-sqlite3" // use for call go-sqlite3 init() method
    "github.com/owner888/kaligo/database"
)

// DriverName is the default driver name for SQLite.
const DriverName = "sqlite3"

// Dialector is the driver struct
type Dialector struct {
    DriverName string
    DSN        string
    stdDB      *sql.DB // Connection
    stdTX      *sql.Tx // Connection for Transaction
}

// Open is
func Open(dsn string) database.Dialector {
    return &Dialector{DSN: dsn}
}

// Name is
func (dialector Dialector) Name() string {
    return "sqlite"
}

// Initialize is
func (dialector Dialector) Initialize(db *database.DB) (err error) {
    if dialector.DriverName == "" {
        dialector.DriverName = DriverName
    }

    if dialector.stdDB != nil {
        db.StdDB = dialector.stdDB
    } else {
        db.StdDB, err = sql.Open(dialector.DriverName, dialector.DSN)
        if err != nil {
            return err
        }
    }
    return
}

// CurrentDatabase is Current Database
func (dialector Dialector) CurrentDatabase(db *database.DB) (name string) {
    var null any
    db.Row("PRAGMA database_list").Scan(&null, &name, &null)
    return
}

// ListDatabases If a database name is given it will return the database name with the configured
// prefix. If not, then just the prefix is returnd
func (dialector Dialector) ListDatabases(like string, db *database.DB) ([]string, error) {
    rows, err := db.Rows("PRAGMA database_list")
    if err != nil {
        return nil, err
    }

    var (
        null  any
        name  string
        names []string
    )

    for rows.Next() {
        if err := rows.Scan(&null, &name, &null); err != nil {
            return nil, err
        } else {
            names = append(names, name)
        }
    }
    return names, nil
}

// ListTables If a table name is given it will return the table name with the configured
// prefix. If not, then just the prefix is returnd
func (dialector Dialector) ListTables(like string, db *database.DB) ([]string, error) {
    sqlStr := `SELECT name FROM sqlite_master WHERE type = "table" AND name != "sqlite_sequence" AND name != "geometry_columns" AND name != "spatial_ref_sys"
              UNION ALL SELECT name FROM sqlite_temp_master 
              WHERE type = "table"`
    if like != "" {
        sqlStr += " AND name LIKE " + db.Quote("%"+like+"%")
    }
    sqlStr += " ORDER BY name"

    var names []string
    db.Query(sqlStr).Scan(&names).Execute()
    return names, nil
}

// ListColumns Lists all of the columns in a table. Optionally, a LIKE string can be
// used to search for specific fields.
func (dialector Dialector) ListColumns(table, like string, db *database.DB) ([]database.Column, error) {
    sqlStr := "PRAGMA table_info(" + db.QuoteTable(table) + ")"
    rows, err := db.Rows(sqlStr)
    if err != nil {
        db.AddError(err)
    }
    var createSQL string
    // CREATE UNIQUE INDEX "id_idx" ON "user" ("id")
    db.Row("SELECT sql FROM sqlite_master WHERE type = ? AND tbl_name = ?", "index", table).Scan(&createSQL)
    var unique bool = false
    if strings.Index(createSQL, "UNIQUE INDEX") != -1 {
        unique = true
    }

    var (
        columnCid, columnNotNull, columnPK int64
        columnName, columnType             string
        columnDfltValue                    sql.NullString
    )

    listColumns := []database.Column{}
    for rows.Next() {
        if err := rows.Scan(&columnCid, &columnName, &columnType, &columnNotNull, &columnDfltValue, &columnPK); err != nil {
            return nil, err
        } else {
            var notNull, primaryKey, autoIncrement bool = true, false, false
            if columnNotNull == 0 {
                notNull = false
            }
            if columnPK == 1 {
                primaryKey = true
                autoIncrement = true
            }

            dataType, dataSize := database.ParseType(columnType)
            column := database.Column{
                Name:          database.ToSchemaName(columnName),
                DBName:        columnName,
                DataType:      dataType,
                Size:          dataSize,
                Precision:     0,
                NotNull:       notNull,
                DefaultValue:  columnDfltValue.String,
                Unique:        unique,
                PrimaryKey:    primaryKey,
                AutoIncrement: autoIncrement,
                Comment:       "",
                Readable:      true,
                Creatable:     true,
                Updatable:     true,
                Extra:         database.ToString(columnCid),
            }

            listColumns = append(listColumns, column)
        }
    }
    return listColumns, nil
}

func (dialector Dialector) getPrimaryKeyColumn(createSQL string) string {
    arr := strings.Split(createSQL, "\n")
    var primaryKeyStr, primaryKeyColumn string
    for _, v := range arr {
        if strings.Index(v, "PRIMARY KEY") != -1 {
            primaryKeyStr = v
        }
    }
    if primaryKeyStr != "" {
        arr = strings.Split(primaryKeyStr, " ")
        primaryKeyColumn = strings.TrimSpace(arr[0])
        primaryKeyColumn = strings.Replace(primaryKeyColumn, "\"", "", -1)
    }
    return primaryKeyColumn
}

// ListIndexes Lists all of the idexes in a table. Optionally, a LIKE string can be
// used to search for specific indexes by name.
func (dialector Dialector) ListIndexes(table, like string, db *database.DB) ([]database.Indexes, error) {
    sqlStr := "SELECT `name`, `tbl_name`, `sql` FROM `sqlite_master` WHERE `type` = 'index' AND `tbl_name` = " + db.Quote(table)
    if like != "" {
        sqlStr += " LIKE " + db.Quote("%"+like+"%")
    }
    rows, err := db.Rows(sqlStr)
    if err != nil {
        return nil, err
    }

    var (
        name    string
        tblName string
        sql     string
        indexes []database.Indexes
    )
    for rows.Next() {
        if err := rows.Scan(&name, &tblName, &sql); err != nil {
            db.AddError(err)
        } else {
            index := database.Indexes{
                Table: tblName,
                Name:  name,
            }

            index.Unique = false
            if strings.Index(sql, "UNIQUE INDEX") != -1 {
                index.Unique = true
            }

            index.Ascend = false
            if strings.Index(sql, "ASC") != -1 {
                index.Ascend = true
            }

            // 对于有ASC的提取不到，因为 ( 后面换行了，要修复一下
            arr := regexp.MustCompile("ON \"" + table + "\" \\(\"(.*?)\"\\)").FindStringSubmatch(sql)
            if len(arr) > 1 {
                index.Column = arr[1]
            }
            indexes = append(indexes, index)
        }
    }
    return indexes, nil
}

// CreateDatabase Creates a database. Will throw a Database Exception if it cannot.
func (dialector Dialector) CreateDatabase(database, charset string, ifNotExists bool, db *database.DB) (err error) {
    return ErrCreateDatabaseNotImplemented
}

// DropDatabase Drops a database. Will throw a Database Exception if it cannot.
func (dialector Dialector) DropDatabase(database string, db *database.DB) (err error) {
    return ErrDropDatabaseNotImplemented
}

// CreateTable Creates a table.
func (dialector Dialector) CreateTable(table string, fields []map[string]any, primaryKeys []string, ifNotExists bool, engine, charset string, foreignKeys []map[string]any, db *database.DB) (err error) {

    sqlStr := "CREATE TABLE "
    if ifNotExists {
        sqlStr += "IF NOT EXISTS "
    }

    sqlStr += db.QuoteTable(table) + " ("
    sqlStr += dialector.ProcessFields(fields, "", db)

    if len(primaryKeys) > 0 {
        for k, v := range primaryKeys {
            primaryKeys[k] = db.QuoteIdentifier(v)
        }
        sqlStr += ",\n\tPRIMARY KEY (" + strings.Join(primaryKeys, ", ") + ")"
    }

    // 要测试一下
    if len(foreignKeys) > 0 {
        sqlStr += db.ProcessForeignKeys(foreignKeys)
    }

    sqlStr += "\n)"
    //if engine != "" {
    //sqlStr += " ENGINE = " + engine + " "
    //}

    //sqlStr += db.ProcessCharset(charset, true, "") + ";"

    _, err = db.Exec(sqlStr)
    return err
}

// RenameTable Renames a table. Will throw a Database Exception if it cannot.
func (dialector Dialector) RenameTable(oldTable, newTable string, db *database.DB) (err error) {
    sqlStr := "ALTER TABLE "
    sqlStr += db.QuoteTable(oldTable)
    sqlStr += " RENAME TO "
    sqlStr += db.QuoteTable(newTable)

    _, err = db.Exec(sqlStr)
    if err != nil {
        db.AddError(err)
    }
    return err
}

// DropTable Drops a table. Will throw a Database Exception if it cannot.
func (dialector Dialector) DropTable(table string, db *database.DB) (err error) {
    sqlStr := "DROP TABLE IF EXISTS "
    sqlStr += db.QuoteTable(table)

    _, err = db.Exec(sqlStr)
    return err
}

// TruncateTable Truncates a table.
func (dialector Dialector) TruncateTable(table string, db *database.DB) (err error) {
    sqlStr := "DELETE FROM " // sqlite DELETE FROM 可以让自增字段从0开始，mysql 要用 TRUNCATE TABLE 才行
    sqlStr += db.QuoteTable(table)

    _, err = db.Exec(sqlStr)
    return err
}

// TableExists Generic check if a given table exists.
func (dialector Dialector) TableExists(table string, db *database.DB) (bool, error) {
    var count int
    err := db.Row("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count)
    return count > 0, err
}

// FieldExists Checks if given field(s) in a given table exists.
func (dialector Dialector) FieldExists(table string, value any, db *database.DB) (bool, error) {
    var columns []string
    switch value.(type) {
    case string:
        columns = append(columns, value.(string))
    default:
        columns = value.([]string)
    }

    for k, v := range columns {
        columns[k] = db.QuoteIdentifier(v)
    }

    sqlStr := "SELECT "
    sqlStr += strings.Join(columns, ", ")
    sqlStr += " FROM "
    sqlStr += db.QuoteTable(table)
    sqlStr += " LIMIT 1"

    _, err := db.Rows(sqlStr)
    if err != nil {
        return false, err
    }
    return true, nil
}

// CreateIndex Creates an index on that table.
func (dialector Dialector) CreateIndex(table string, indexColumns any, indexName, index string, db *database.DB) (err error) {
    var sqlStr string
    acceptedIndex := []string{"UNIQUE", "PRIMARY"} // Can't support FULLTEXT, SPATIAL, NONCLUSTERED

    // make sure the index type is uppercase
    if index != "" {
        index = strings.ToUpper(index)
        if !database.InSlice(index, &acceptedIndex) {
            return fmt.Errorf("failed to create index with name %s", index)
        }
    }

    // 索引名为空，如果是多个字段则以下划线拼接，单个字段名以字段名为索引名
    if indexName == "" {
        switch indexColumns.(type) {
        case []string:
            for k, v := range indexColumns.([]string) {
                if indexName == "" {
                    indexName += ""
                } else {
                    indexName += "_"
                }

                key := database.ToString(k)
                if database.IsNumeric(key) {
                    indexName += v
                } else {
                    key = strings.Replace(key, "(", "", -1)
                    key = strings.Replace(key, ")", "", -1)
                    key = strings.Replace(key, " ", "", -1)
                    indexName += key
                }
            }
        default:
            indexName = indexColumns.(string)
        }
    }

    columns := ""
    switch indexColumns.(type) {
    case []string:
        for k, v := range indexColumns.([]string) {
            if columns == "" {
                columns += ""
            } else {
                columns += ", "
            }

            key := database.ToString(k)
            if database.IsNumeric(key) {
                columns += db.QuoteIdentifier(v)
            } else {
                columns += db.QuoteIdentifier(key) + " " + strings.ToUpper(v)
            }
        }
    default:
        columns = db.QuoteIdentifier(indexColumns.(string))
    }

    if index != "PRIMARY" {
        sqlStr = "CREATE "
        if index != "" {
            sqlStr += index + " "
        }
        sqlStr += "INDEX "
        sqlStr += db.QuoteIdentifier(indexName)
        sqlStr += " ON "
        sqlStr += db.QuoteTable(table)
        sqlStr += " (" + columns + ")"

        _, err = db.Exec(sqlStr)
        if err != nil {
            db.AddError(err)
        }
    } else {
        err = ErrPrimaryKeyNotImplemented
    }
    return err
}

// RenameIndex Rename an index from a table.
func (dialector Dialector) RenameIndex(table, oldName, newName string, db *database.DB) (err error) {
    var sql string
    db.Row("SELECT sql FROM sqlite_master WHERE type = ? AND tbl_name = ? AND name = ?", "index", table, oldName).Scan(&sql)
    if sql != "" {
        // Drop old index
        dialector.DropIndex(table, oldName, db)
        // Create a new index
        _, err = db.Exec(strings.Replace(sql, oldName, newName, 1))
    } else {
        err = fmt.Errorf("failed to find index with name %v", oldName)
    }
    return err
}

// DropIndex Drop an index from a table.
func (dialector Dialector) DropIndex(table, indexName string, db *database.DB) (err error) {
    // sqlite 所有表的索引都不能重复，因为他不是依赖表的
    sqlStr := "DROP INDEX " + db.QuoteIdentifier(indexName)

    _, err = db.Exec(sqlStr)
    return err
}

// AddForeignKey Adds a single foreign key to a table
func (dialector Dialector) AddForeignKey(table string, foreignKey []map[string]any, db *database.DB) (err error) {
    return ErrForeignKeyNotImplemented
}

// DropForeignKey Drops a foreign key from a table
func (dialector Dialector) DropForeignKey(table string, fkName string, db *database.DB) (err error) {
    return ErrForeignKeyNotImplemented
}

// AddFields adds fields to a table.
func (dialector Dialector) AddFields(table string, fields []map[string]any, db *database.DB) error {
    // SQLite 不支持一条语句添加多个字段
    if len(fields) > 1 {
        return ErrAlertTableMultipleAddNotImplemented
    }
    return dialector.AlterFields("ADD", table, fields, db)
}

// DropFields drops fields from a table.
func (dialector Dialector) DropFields(table string, value any, db *database.DB) error {
    return ErrAlertTableNotImplemented
}

// ModifyFields alters fields in a table.
func (dialector Dialector) ModifyFields(table string, fields []map[string]any, db *database.DB) error {
    return ErrAlertTableNotImplemented
}

// AlterFields is ...
func (dialector Dialector) AlterFields(alterType string, table string, fields any, db *database.DB) (err error) {
    sqlStr := "ALTER TABLE " + db.QuoteTable(table) + " "

    fieldMaps := fields.([]map[string]any)

    useBrackets := true
    if database.InSlice(alterType, &[]string{"ADD"}) {
        useBrackets = false
    }
    var prefix string
    if useBrackets {
        sqlStr += alterType + " "
        sqlStr += "("
        prefix = ""
    } else {
        prefix = alterType + " "
    }
    sqlStr += dialector.ProcessFields(fieldMaps, prefix, db)
    if useBrackets {
        sqlStr += ")"
    }

    _, err = db.Exec(sqlStr)
    if err != nil {
        db.AddError(err)
    }
    return err
}

// ProcessFields is
func (dialector Dialector) ProcessFields(fields []map[string]any, prefix string, db *database.DB) string {
    var sqlFields []string

    for _, dict := range fields {
        dict = database.MapChangeKeyCase(dict, true)
        sqlStr := ""
        // ALTER TABLE statement to modify、drop、rename a column in SQLite does not support.
        if prefix == "ADD " {
            sqlStr += "\n\tADD "
        }

        if value, ok := dict["NAME"]; ok {
            sqlStr += " " + db.QuoteIdentifier(value.(string))
        }

        if value, ok := dict["UNSIGNED"]; ok && value.(bool) == true {
            sqlStr += " UNSIGNED"
        }

        if value, ok := dict["TYPE"]; ok {
            sqlStr += " " + value.(string)
        }

        if value, ok := dict["CONSTRAINT"]; ok {
            sqlStr += "(" + database.ToString(value) + ")"
        }

        if value, ok := dict["NOTNULL"]; ok && value.(bool) == true {
            sqlStr += " NOT NULL"
        } else {
            sqlStr += " NULL"
        }

        if value, ok := dict["DEFAULT"]; ok {
            sqlStr += " DEFAULT " + db.Quote(value.(string))
        }

        // Sqlite does not support COMMENT、FIRST、AFTER

        sqlFields = append(sqlFields, sqlStr)
    }

    return strings.Join(sqlFields, ", ")
}
