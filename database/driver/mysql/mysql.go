package driver

import (
    "database/sql"
    "reflect"
    //"strconv"
    "strings"
    _ "github.com/go-sql-driver/mysql"   // need comment
    "github.com/owner888/kaligo/database"
)

// DriverName is the default driver name for SQLite.
const DriverName = "mysql"

// Dialector is the driver struct
type Dialector struct {
    DriverName string
    DSN        string
    StdDB      *sql.DB      // Connection
    StdTX      *sql.Tx      // Connection for Transaction
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

// ListDatabases If a database name is given it will return the database name with the configured
// prefix. If not, then just the prefix is returnd
func (dialector Dialector) ListDatabases(like string, db *database.DB) []string {
    sqlStr := "SHOW DATABASES"
    if  like != "" {
        sqlStr += " LIKE " + db.Quote("%" + like + "%")
    }

    var tables []string
    db.Query(sqlStr).Scan(&tables).Execute()
    return tables
}

// ListTables If a table name is given it will return the table name with the configured
// prefix. If not, then just the prefix is returnd
func (dialector Dialector) ListTables(like string, db *database.DB) []string {
    sqlStr := "SHOW TABLES"
    if  like != "" {
        sqlStr += " LIKE " + db.Quote("%" + like + "%")
    }

    var tables []string
    db.Query(sqlStr).Scan(&tables).Execute()
    return tables
}

// ListColumns Lists all of the columns in a table. Optionally, a LIKE string can be
// used to search for specific fields.
func (dialector Dialector) ListColumns(table string, like string, db *database.DB) []database.Column {
    sqlStr := "SHOW FULL COLUMNS FROM " + db.QuoteTable(table)
    if  like != "" {
        sqlStr += " LIKE " + db.Quote("%" + like + "%")
    }

    listColumns := []database.Column{}
    db.Query(sqlStr).Scan(&listColumns).Execute()
    return listColumns
}

// ListIndexes Lists all of the idexes in a table. Optionally, a LIKE string can be
// used to search for specific indexes by name.
func (dialector Dialector) ListIndexes(table string, like string, db *database.DB) []database.Indexes {
    sqlStr := "SHOW INDEX FROM " + db.QuoteTable(table)
    if  like != "" {
        sqlStr += " WHERE " + db.QuoteIdentifier("Key_name") + " LIKE " + db.Quote("%" + like + "%")
    }
    
    rows, err := db.Rows(sqlStr)
    if err != nil {
        db.AddError(err)
    }

    type IndexTmp struct {
        Table           string          `field:"Table"` 
        NonUnique       int64           `field:"Non_unique"`
        KeyName         string          `field:"Key_name"`
        SeqInIndex      int64           `field:"Seq_in_index"`
        ColumnName      string          `field:"Column_name"`
        Collation       sql.NullString  `field:"Collation"`
        Cardinality     int64           `field:"Cardinality"`
        SubPart         sql.NullString  `field:"Sub_part"`
        Packed          sql.NullString  `field:"Packed"`
        Null            string          `field:"Null"`
        IndexType       string          `field:"Index_type"`
        Comment         string          `field:"Comment"`
        IndexComment    string          `field:"Index_comment"`
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
                Table       : indexTmp.Table,
                Name        : indexTmp.KeyName,
                Column      : indexTmp.ColumnName,
                Order       : indexTmp.SeqInIndex,
                Type        : indexTmp.IndexType,
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

// DropIndex Drop an index from a table.
func (dialector Dialector) DropIndex(table string, indexName string, db *database.DB) (err error) {
    var sqlStr string    
    if strings.ToUpper(indexName) == "PRIMARY" {
        sqlStr = "ALTER TABLE " + db.QuoteTable(table)
        sqlStr += " DROP PRIMARY KEY"
    } else {
        sqlStr = "DROP INDEX " + db.QuoteIdentifier(indexName)
        sqlStr += " ON "+ db.QuoteTable(table)
    }

    _, err = db.Exec(sqlStr)
    if err != nil {
        db.AddError(err)
    }
    return err
}
