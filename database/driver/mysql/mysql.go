package driver

import (
	"database/sql"
	"fmt"
	"reflect"

	//"regexp"
	//"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql" // need comment
	"github.com/owner888/kaligo/database"
)

// DriverName is the default driver name for SQLite.
const DriverName = "mysql"

// Dialector is the driver struct
type Dialector struct {
	DriverName string
	DSN        string
	StdDB      *sql.DB // Connection
	StdTX      *sql.Tx // Connection for Transaction
}

// Open is
func Open(dsn string) database.Dialector {
	return &Dialector{DSN: dsn}
}

// Name is dialector name
func (dialector Dialector) Name() string {
	return "mysql"
}

// Initialize is
func (dialector Dialector) Initialize(db *database.DB) (err error) {
	if dialector.DriverName == "" {
		dialector.DriverName = DriverName
	}

	if dialector.StdDB != nil {
		db.StdDB = dialector.StdDB
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
	db.Row("SELECT DATABASE()").Scan(&name)
	return
}

// ListDatabases If a database name is given it will return the database name with the configured
// prefix. If not, then just the prefix is returnd
func (dialector Dialector) ListDatabases(like string, db *database.DB) []string {
	sqlStr := "SHOW DATABASES"
	if like != "" {
		sqlStr += " LIKE " + db.Quote("%"+like+"%")
	}

	var names []string
	db.Query(sqlStr).Scan(&names).Execute()
	return names
}

// ListTables If a table name is given it will return the table name with the configured
// prefix. If not, then just the prefix is returnd
func (dialector Dialector) ListTables(like string, db *database.DB) []string {
	sqlStr := "SHOW TABLES"
	if like != "" {
		sqlStr += " LIKE " + db.Quote("%"+like+"%")
	}

	var names []string
	db.Query(sqlStr).Scan(&names).Execute()
	return names
}

// ListColumns Lists all of the columns in a table. Optionally, a LIKE string can be
// used to search for specific fields.
func (dialector Dialector) ListColumns(table, like string, db *database.DB) []database.Column {
	sqlStr := "SHOW FULL COLUMNS FROM " + db.QuoteTable(table)
	if like != "" {
		sqlStr += " LIKE " + db.Quote("%"+like+"%")
	}
	rows, err := db.Rows(sqlStr)
	if err != nil {
		db.AddError(err)
	}

	var (
		columnField      string
		columnType       string
		columnCollation  sql.NullString
		columnNull       string
		columnKey        string
		columnDefault    sql.NullString
		columnExtra      string
		columnPrivileges string
		columnComment    string
	)

	listColumns := []database.Column{}
	for rows.Next() {
		if err := rows.Scan(&columnField, &columnType, &columnCollation, &columnNull, &columnKey, &columnDefault, &columnExtra, &columnPrivileges, &columnComment); err != nil {
			db.AddError(err)
		} else {
			//if strings.Index(columnType, "unsigned") != -1 { }
			//if strings.Index(columnType, "zerofill") != -1 { }
			var readable, creatable, updatable, notNull, unique, primaryKey, autoIncrement bool = true, true, true, false, false, false, false
			if strings.Index(columnPrivileges, "select") != -1 {
				readable = true
			}
			if strings.Index(columnPrivileges, "insert") != -1 {
				creatable = true
			}
			if strings.Index(columnPrivileges, "update") != -1 {
				updatable = true
			}
			if columnNull == "NO" { // Allow Null
				notNull = true
			}
			if columnKey == "PRI" {
				primaryKey = true
			} else if columnKey == "UNI" {
				unique = true
			}
			if columnExtra == "auto_increment" {
				autoIncrement = true
			}

			dataType, dataSize := database.ParseType(columnType)
			column := database.Column{
				Name:          database.ToSchemaName(columnField),
				DBName:        columnField,
				DataType:      dataType,
				Size:          dataSize,
				Precision:     0,
				NotNull:       notNull,
				DefaultValue:  columnDefault.String,
				Unique:        unique,
				PrimaryKey:    primaryKey,
				AutoIncrement: autoIncrement,
				Comment:       columnComment,
				Readable:      readable,
				Creatable:     creatable,
				Updatable:     updatable,
				Extra:         columnExtra,
			}

			listColumns = append(listColumns, column)
		}
	}

	return listColumns
}

// ListIndexes Lists all of the idexes in a table. Optionally, a LIKE string can be
// used to search for specific indexes by name.
func (dialector Dialector) ListIndexes(table, like string, db *database.DB) []database.Indexes {
	sqlStr := "SHOW INDEX FROM " + db.QuoteTable(table)
	if like != "" {
		sqlStr += " WHERE " + db.QuoteIdentifier("Key_name") + " LIKE " + db.Quote("%"+like+"%")
	}

	rows, err := db.Rows(sqlStr)
	if err != nil {
		db.AddError(err)
	}

	type IndexTmp struct {
		Table        string         `field:"Table"`
		NonUnique    int64          `field:"Non_unique"`
		KeyName      string         `field:"Key_name"`
		SeqInIndex   int64          `field:"Seq_in_index"`
		ColumnName   string         `field:"Column_name"`
		Collation    sql.NullString `field:"Collation"`
		Cardinality  int64          `field:"Cardinality"`
		SubPart      sql.NullString `field:"Sub_part"`
		Packed       sql.NullString `field:"Packed"`
		Null         string         `field:"Null"`
		IndexType    string         `field:"Index_type"`
		Comment      string         `field:"Comment"`
		IndexComment string         `field:"Index_comment"`
	}

	indexes := []database.Indexes{}
	for rows.Next() {
		indexTmp := IndexTmp{}
		o := reflect.ValueOf(&indexTmp).Elem()
		numCols := o.NumField()
		columns := make([]any, numCols)
		for i := 0; i < numCols; i++ {
			field := o.Field(i)
			columns[i] = field.Addr().Interface()
		}

		if err := rows.Scan(columns...); err != nil {
			db.AddError(err)
		} else {
			index := database.Indexes{
				Table:  indexTmp.Table,
				Name:   indexTmp.KeyName,
				Column: indexTmp.ColumnName,
				Order:  indexTmp.SeqInIndex,
				Type:   indexTmp.IndexType,
			}

			index.Primary = false
			if indexTmp.KeyName == "PRIMARY" {
				index.Primary = true
			}

			index.Unique = false
			if indexTmp.NonUnique == 0 {
				index.Unique = true
			}

			index.Null = false
			if indexTmp.Null == "YES" {
				index.Null = true
			}

			index.Ascend = false
			if indexTmp.Collation.String == "A" {
				index.Ascend = true
			}

			indexes = append(indexes, index)
		}
	}
	return indexes
}

// CreateDatabase Creates a database. Will throw a Database Exception if it cannot.
func (dialector Dialector) CreateDatabase(database, charset string, ifNotExists bool, db *database.DB) (err error) {
	database = db.QuoteTable(database)

	sqlStr := "CREATE DATABASE "
	if ifNotExists {
		sqlStr += "IF NOT EXISTS "
	}
	sqlStr += db.QuoteIdentifier(database) + db.ProcessCharset(charset, true)

	_, err = db.Exec(sqlStr)
	if err != nil {
		db.AddError(err)
	}
	return err
}

// DropDatabase Drops a database. Will throw a Database Exception if it cannot.
func (dialector Dialector) DropDatabase(database string, db *database.DB) (err error) {
	sqlStr := "DROP DATABASE " + db.QuoteIdentifier(database)

	_, err = db.Exec(sqlStr)
	if err != nil {
		db.AddError(err)
	}
	return err
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

	if len(foreignKeys) > 0 {
		sqlStr += db.ProcessForeignKeys(foreignKeys)
	}

	sqlStr += "\n)"
	if engine != "" {
		sqlStr += " ENGINE = " + engine + " "
	}

	sqlStr += db.ProcessCharset(charset, true, "") + ";"

	_, err = db.Exec(sqlStr)
	if err != nil {
		db.AddError(err)
	}
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
	if err != nil {
		db.AddError(err)
	}
	return err
}

// TruncateTable Truncates a table.
func (dialector Dialector) TruncateTable(table string, db *database.DB) (err error) {
	sqlStr := "TRUNCATE TABLE "
	sqlStr += db.QuoteTable(table)

	_, err = db.Exec(sqlStr)
	if err != nil {
		db.AddError(err)
	}
	return err
}

// TableExists Generic check if a given table exists.
func (dialector Dialector) TableExists(table string, db *database.DB) bool {
	sqlStr := "SELECT * FROM "
	sqlStr += db.QuoteTable(table)
	sqlStr += " LIMIT 1"

	_, err := db.Rows(sqlStr)
	if err != nil {
		db.AddError(err)
		return false
	}
	return true
}

// FieldExists Checks if given field(s) in a given table exists.
func (dialector Dialector) FieldExists(table string, value any, db *database.DB) bool {
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
		db.AddError(err)
		return false
	}
	return true
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

	if index == "PRIMARY" {
		sqlStr = "ALTER TABLE "
		sqlStr += db.QuoteTable(table)
		sqlStr += " ADD PRIMARY KEY "
	} else {
		sqlStr = "CREATE "
		if index != "" {
			sqlStr += index + " "
		}
		sqlStr += "INDEX "
		sqlStr += db.QuoteIdentifier(indexName)
		sqlStr += " ON "
		sqlStr += db.QuoteTable(table)
	}
	sqlStr += " (" + columns + ")"

	_, err = db.Exec(sqlStr)
	if err != nil {
		db.AddError(err)
	}
	return err
}

// RenameIndex Rename an index from a table.
func (dialector Dialector) RenameIndex(table, oldName, newName string, db *database.DB) (err error) {
	sqlStr := "ALTER TABLE " + db.QuoteTable(table) + " RENAME INDEX " + db.QuoteIdentifier(oldName) + " TO " + db.QuoteIdentifier(newName)
	_, err = db.Exec(sqlStr)
	if err != nil {
		db.AddError(err)
	}
	return
}

// DropIndex Drop an index from a table.
func (dialector Dialector) DropIndex(table, indexName string, db *database.DB) (err error) {
	var sqlStr string
	if strings.ToUpper(indexName) == "PRIMARY" {
		sqlStr = "ALTER TABLE " + db.QuoteTable(table)
		sqlStr += " DROP PRIMARY KEY"
	} else {
		sqlStr = "DROP INDEX " + db.QuoteIdentifier(indexName)
		sqlStr += " ON " + db.QuoteTable(table)
	}

	_, err = db.Exec(sqlStr)
	if err != nil {
		db.AddError(err)
	}
	return err
}

// AddForeignKey Adds a single foreign key to a table
func (dialector Dialector) AddForeignKey(table string, foreignKey []map[string]any, db *database.DB) (err error) {
	sqlStr := "ALTER TABLE " + db.QuoteTable(table) + " ADD " + strings.TrimLeft(db.ProcessForeignKeys(foreignKey), ",")

	_, err = db.Exec(sqlStr)
	if err != nil {
		db.AddError(err)
	}
	return err
}

// DropForeignKey Drops a foreign key from a table
func (dialector Dialector) DropForeignKey(table string, fkName string, db *database.DB) (err error) {
	sqlStr := "ALTER TABLE " + db.QuoteTable(table) + " DROP FOREIGN KEY " + db.QuoteIdentifier(fkName)

	_, err = db.Exec(sqlStr)
	if err != nil {
		db.AddError(err)
	}
	return err
}

// AddFields adds fields to a table.
func (dialector Dialector) AddFields(table string, fields []map[string]any, db *database.DB) error {
	return dialector.AlterFields("ADD", table, fields, db)
}

// DropFields drops fields from a table.
func (dialector Dialector) DropFields(table string, value any, db *database.DB) error {
	var fields []string
	switch value.(type) {
	case string:
		fields = append(fields, value.(string))
	default:
		fields = value.([]string)
	}
	return dialector.AlterFields("DROP", table, fields, db)
}

// ModifyFields alters fields in a table.
func (dialector Dialector) ModifyFields(table string, fields []map[string]any, db *database.DB) error {
	return dialector.AlterFields("MODIFY", table, fields, db)
}

// AlterFields is ...
func (dialector Dialector) AlterFields(alterType string, table string, fields any, db *database.DB) (err error) {
	sqlStr := "ALTER TABLE " + db.QuoteTable(table) + " "

	if alterType == "DROP" {
		var dropFields, sqlDropFields []string
		dropFields = fields.([]string)
		for _, field := range dropFields {
			sqlDropFields = append(sqlDropFields, "DROP "+db.QuoteIdentifier(field))
		}
		sqlStr += strings.Join(sqlDropFields, ", ")
	} else {
		fieldMaps := fields.([]map[string]any)

		useBrackets := true
		if database.InSlice(alterType, &[]string{"ADD", "CHANGE", "MODIFY"}) {
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
		if prefix == "MODIFY " {
			sqlStr += "\n\tCHANGE "
			if value, ok := dict["OLDNAME"]; ok {
				sqlStr += db.QuoteIdentifier(value.(string))
			} else if value, ok := dict["NAME"]; ok {
				sqlStr += db.QuoteIdentifier(value.(string))
			}
		} else if prefix == "ADD " {
			sqlStr += "\n\tADD "
		} else if prefix == "DROP " {
			sqlStr += "\n\tDROP "
		}

		if value, ok := dict["NAME"]; ok {
			sqlStr += " " + db.QuoteIdentifier(value.(string))
		}

		if value, ok := dict["TYPE"]; ok {
			sqlStr += " " + value.(string)
		}

		if value, ok := dict["CONSTRAINT"]; ok {
			sqlStr += "(" + database.ToString(value) + ")"
		}

		if value, ok := dict["CHARSET"]; ok {
			sqlStr += db.ProcessCharset(value.(string), false)
		}

		if value, ok := dict["UNSIGNED"]; ok && value.(bool) == true {
			sqlStr += " UNSIGNED"
		}

		if value, ok := dict["DEFAULT"]; ok {
			sqlStr += " DEFAULT " + db.Quote(value.(string))
		}

		if value, ok := dict["NOTNULL"]; ok && value.(bool) == true {
			sqlStr += " NOT NULL"
		} else {
			sqlStr += " NULL"
		}

		if value, ok := dict["AUTO_INCREMENT"]; ok && value.(bool) == true {
			sqlStr += " AUTO_INCREMENT"
		}

		if value, ok := dict["PRIMARY_KEY"]; ok && value.(bool) == true {
			sqlStr += " PRIMARY_KEY"
		}

		if value, ok := dict["COMMENT"]; ok {
			sqlStr += " COMMENT " + db.Escape(value.(string))
		}

		if value, ok := dict["FIRST"]; ok && value.(bool) == true {
			sqlStr += " FIRST"
		}

		if value, ok := dict["AFTER"]; ok {
			sqlStr += " AFTER " + db.QuoteIdentifier(value.(string))
		}

		sqlFields = append(sqlFields, sqlStr)
	}

	return strings.Join(sqlFields, ", ")
}
