package driver

import (
    "database/sql"
    //"fmt"
    "reflect"
    "regexp"
    //"strconv"
    "strings"
    _ "github.com/mattn/go-sqlite3"     // need comment
    "github.com/owner888/kaligo/database"
)

// DriverName is the default driver name for SQLite.
const DriverName = "sqlite3"

// Dialector is the driver struct
type Dialector struct {
    DriverName string
    DSN        string
    stdDB      *sql.DB      // Connection
    stdTX      *sql.Tx      // Connection for Transaction
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

// ListDatabases If a database name is given it will return the database name with the configured
// prefix. If not, then just the prefix is returnd
func (dialector Dialector) ListDatabases(like string, db *database.DB) []string {
    sqlStr := "PRAGMA database_list"

    type DatabaseList struct {
        Seq   int64
        Name  string
        File  string
    }

    //results := []map[string]interface{}{}
    //db.Query(sqlStr).Scan(&results).Execute()
    //fmt.Printf("%v\n", results)

    rows, err := db.Rows(sqlStr)
    if err != nil {
        db.AddError(err)
    }

    var databases []string
    for rows.Next() {
        database := DatabaseList{}
        o := reflect.ValueOf(&database).Elem()
        numCols := o.NumField()
        columns := make([]interface{}, numCols)
        for i := 0; i < numCols; i++ {
            field := o.Field(i)
            columns[i] = field.Addr().Interface()
        }
        
        if err := rows.Scan(columns...); err != nil {
            db.AddError(err)
        } else {
            databases = append(databases, database.Name)
        }
    }
    return databases
}

// ListTables If a table name is given it will return the table name with the configured
// prefix. If not, then just the prefix is returnd
func (dialector Dialector) ListTables(like string, db *database.DB) []string {
    sqlStr := `SELECT name FROM sqlite_master WHERE type = "table" AND name != "sqlite_sequence" AND name != "geometry_columns" AND name != "spatial_ref_sys"
              UNION ALL SELECT name FROM sqlite_temp_master 
              WHERE type = "table"`
    if  like != "" {
        sqlStr += " AND name LIKE " + db.Quote("%" + like + "%")
    }
    sqlStr += " ORDER BY name"

    var tables []string
    db.Query(sqlStr).Scan(&tables).Execute()
    return tables
}

// ListColumns Lists all of the columns in a table. Optionally, a LIKE string can be
// used to search for specific fields.
func (dialector Dialector) ListColumns(table string, like string, db *database.DB) []database.Column {
    sqlStr := "PRAGMA table_info(" + db.QuoteTable(table) + ")"
    rows, err := db.Rows(sqlStr)
    if err != nil {
        db.AddError(err)
    }

    type Column struct {
        Cid         string
        Name        string
        Type        string
        NotNull     int64
        DfltValue   sql.NullString
        PK          string
    }

    listColumns := []database.Column{}
    for rows.Next() {
        columnTmp := Column{}
        o := reflect.ValueOf(&columnTmp).Elem()
        numCols := o.NumField()
        columns := make([]interface{}, numCols)
        for i := 0; i < numCols; i++ {
            field := o.Field(i)
            columns[i] = field.Addr().Interface()
        }
        
        if err := rows.Scan(columns...); err != nil {
            db.AddError(err)
        } else {
            column := database.Column{
                Field     : columnTmp.Name,
                Type      : columnTmp.Type,
                Default   : columnTmp.DfltValue.String,
                Key       : columnTmp.PK,
                Extra     : columnTmp.Cid,
            }

            column.Null = "NO"
            if columnTmp.NotNull == 1 {
                column.Null = "YES"
            }

            listColumns = append(listColumns, column)
        }
    }
    return listColumns
}

// ListIndexes Lists all of the idexes in a table. Optionally, a LIKE string can be
// used to search for specific indexes by name.
func (dialector Dialector) ListIndexes(table string, like string, db *database.DB) []database.Indexes {
    sqlStr := "SELECT `name`, `tbl_name`, `sql` FROM `sqlite_master` WHERE `type` = 'index' AND `tbl_name` = " + db.Quote(table)
    if  like != "" {
        sqlStr += " LIKE " + db.Quote("%" + like + "%")
    }
    
    rows, err := db.Rows(sqlStr)
    if err != nil {
        db.AddError(err)
    }

    type IndexTmp struct {
        KeyName string  `field:"name"`
        Table   string  `field:"tbl_name"` 
        SQLStr  string  `field:"sql"`
    }

    indexes := []database.Indexes{}
    for rows.Next() {
        indexTmp := IndexTmp{}
        o := reflect.ValueOf(&indexTmp).Elem()
        numCols := o.NumField()
        columns := make([]interface{}, numCols)
        for i := 0; i < numCols; i++ {
            field := o.Field(i)
            columns[i] = field.Addr().Interface()
        }
        
        if err := rows.Scan(columns...); err != nil {
            db.AddError(err)
        } else {
            index := database.Indexes{
                Table     : indexTmp.Table,
                Name      : indexTmp.KeyName,
            }

            index.Unique = false
            if strings.Index(indexTmp.SQLStr, "UNIQUE INDEX") != -1 {
                index.Unique = true
            }

            index.Ascend = false
            if strings.Index(indexTmp.SQLStr, "ASC") != -1 {
                index.Ascend = true
            }

            reg := regexp.MustCompile("ON \""+table+"\" \\(\"(.*?)\"\\)")
            arr := reg.FindStringSubmatch(indexTmp.SQLStr)
            if len(arr) > 1 {
                index.Column = arr[1]
            }
            indexes = append(indexes, index)
        }
    }
    return indexes
}

// DropIndex Drop an index from a table.
func (dialector Dialector) DropIndex(table string, indexName string, db *database.DB) (err error) {
    // sqlite 所有表的索引都不能重复，因为他不是依赖表的
    sqlStr := "DROP INDEX " + db.QuoteIdentifier(indexName)

    _, err = db.Exec(sqlStr)
    if err != nil {
        db.AddError(err)
    }
    return err
}
