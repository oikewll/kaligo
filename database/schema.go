package database

import (
    //"database/sql"
    //"errors"
    "fmt"
    "reflect"
    "strings"
    "sync"
    "go/ast"
)

// Schema is the struct for MySQL DATE type
type Schema struct {
    *Query

    Name                      string
    ModelType                 reflect.Type
    Table                     string
    Fields                    []*Field
    FieldsByName              map[string]*Field     // 通过 struct 字段名查询
    FieldsByDBName            map[string]*Field     // 通过 数据库 字段名查询
    err                       error
    initialized               chan struct{}
    cacheStore                *sync.Map
}

// Indexes 数据表索引，字段名记得要大写开头，才是public，否则访问不到
type Indexes struct {
    Table       string  `field:"Table"` 
    Name        string  `field:"Key_name"`
    Column      string  `field:"Column_name"`
    Order       int64   `field:"Seq_in_index"`
    Type        string  `field:"Index_type"`
    Primary     bool    `field:"Key_name == 'PRIMARY'"`
    Unique      bool    `field:"Non_unique == 0"`
    Null        bool    `field:"Null == 'YES'"`
    Ascend      bool    `field:"Collation == 'A'"`

    //Collation       string   `field:"Collation"`
    //Cardinality     int64    `field:"Cardinality"`
    //SubPart         string   `field:"Sub_part"`
    //Packed          string   `field:"Packed"`
    //Comment         string   `field:"Comment"`
    //IndexComment    string   `field:"Index_comment"`
}

func (s Schema) String() string {
    if s.ModelType.Name() == "" {
        return fmt.Sprintf("%s(%s)", s.Name, s.Table)
    }
    return fmt.Sprintf("%s.%s", s.ModelType.PkgPath(), s.ModelType.Name())
}

// LookUpField is 通过表名 或者 数据库名 查询字段
func (s Schema) LookUpField(name string) *Field {
    if field, ok := s.FieldsByDBName[name]; ok {
        return field
    }

    if field, ok := s.FieldsByName[name]; ok {
        return field
    }
    return nil
}

// Parse get data type from dialector
// 1、传入 struct，比如 &User{}，生成 User对应的 []*Field
// 2、当查询了数据 rows.Next 循环的时候，去上面的[]*Field 找对应的字段，通过 field.Set 去设置值
// 3、map[string]interface{}、[]map[string]interface{}{}、[]int64 都不会到这里来
func Parse(dest interface{}, cacheStore *sync.Map) (*Schema, error) {
    if dest == nil {
        return nil, fmt.Errorf("%w: %+v", ErrUnsupportedDataType, dest)
    }

    // 类型：*reflect.rtype = mysql.User
    modelType := reflect.ValueOf(dest).Type()
    for modelType.Kind() == reflect.Slice || modelType.Kind() == reflect.Array || modelType.Kind() == reflect.Ptr {
        modelType = modelType.Elem()
    }

    // 非 Struct，直接返回错误
    if modelType.Kind() != reflect.Struct {
        if modelType.PkgPath() == "" {
            return nil, fmt.Errorf("%w: %+v", ErrUnsupportedDataType, dest)
        }
        return nil, fmt.Errorf("%w: %s.%s", ErrUnsupportedDataType, modelType.PkgPath(), modelType.Name())
    }

    //fmt.Printf("11111 = cacheStore.Load = %v\n", modelType)

    // 等待其他协程数据
    if v, ok := cacheStore.Load(modelType); ok {
        //fmt.Printf("22222 = cacheStore.Load = %v\n", v)
        s := v.(*Schema)
        // Wait for the initialization of other goroutines to complete
        <-s.initialized
        return s, s.err
    }

    //modelValue := reflect.New(modelType)
    // 其实没必要去弄表名，因为schema又不是一张表
    //tableName  := TransFieldName(modelType.Name())

    // Model
    schema := &Schema{
        Name:           modelType.Name(),
        ModelType:      modelType,
        FieldsByName:   map[string]*Field{},
        FieldsByDBName: map[string]*Field{},
        cacheStore:     cacheStore,
        initialized:    make(chan struct{}),
    }
    // When the schema initialization is completed, the channel will be closed
    defer close(schema.initialized)

    // 等待其他协程数据
    if v, loaded := cacheStore.LoadOrStore(modelType, schema); loaded {
        s := v.(*Schema)
        // Wait for the initialization of other goroutines to complete
        <-s.initialized
        return s, s.err
    }

    // 函数结束时，检查错误，有写日志，并且删除缓存
    defer func() {
        if schema.err != nil {
            //logger.Default.Error(context.Background(), schema.err.Error())
            cacheStore.Delete(modelType)
        }
    }()

    // 从 struct 映射出字段，struct 名当作 Schema，字段作为 Field
    for i := 0; i < modelType.NumField(); i++ {
        // IsExported 是否以大写字母开头
        if fieldStruct := modelType.Field(i); ast.IsExported(fieldStruct.Name) {
            if field := ParseField(fieldStruct); field != nil {
                schema.Fields = append(schema.Fields, field)
            }
        }
    }

    // 循环字段，给 FieldsByDBName、FieldsByName 赋值
    for _, field := range schema.Fields {
        if field.DBName == "" && field.DataType != "" {
            // 驼峰转下划线，数据库字段都是下划线写法
            field.DBName = ToDBName(field.Name)
        }

        if field.DBName != "" {
            if _, ok := schema.FieldsByDBName[field.DBName]; !ok {
                schema.FieldsByDBName[field.DBName] = field
                schema.FieldsByName[field.Name]     = field
                //fmt.Printf("66666 ---> %v\n", field)
            }
        }

        field.setupValuerAndSetter()
    }

    return schema, schema.err
}

// CurrentDatabase is Current Database
func (s *Schema) CurrentDatabase() (name string) {
    return s.Dialector.CurrentDatabase(s.DB)
}

// ListDatabases If a database name is given it will return the database name with the configured
// prefix. If not, then just the prefix is returnd
func (s *Schema) ListDatabases(args ...string) []string {
    like := ""
    if len(args) > 0 {
        like = args[0]
    }
    return s.Dialector.ListDatabases(like, s.DB)
}

// ListTables If a table name is given it will return the table name with the configured
// prefix. If not, then just the prefix is returnd
func (s *Schema) ListTables(args ...string) []string {
    like := ""
    if len(args) > 0 {
        like = args[0]
    }
    return s.Dialector.ListTables(like, s.DB)
}

// ListColumns Lists all of the columns in a table. Optionally, a LIKE string can be
// used to search for specific fields.
func (s *Schema) ListColumns(table string, args ...string) []Column {
    like := ""
    if len(args) > 0 {
        like = args[0]
    }
    return s.Dialector.ListColumns(table, like, s.DB)
}

// ListIndexes Lists all of the idexes in a table. Optionally, a LIKE string can be
// used to search for specific indexes by name.
func (s *Schema) ListIndexes(table string, args ...string) []Indexes {
    like := ""
    if len(args) > 0 {
        like = args[0]
    }
    return s.Dialector.ListIndexes(table, like, s.DB)
}

// CreateDatabase Creates a database. Will throw a Database Exception if it cannot.
func (s *Schema) CreateDatabase(database string, args ...interface{}) (err error) {
    var (
        ifNotExists bool   = true
        charset     string = "utf8"
    )

    if len(args) > 1 {
        charset     = args[0].(string)
        ifNotExists = args[1].(bool)
    } else if len(args) > 0 {
        ifNotExists = args[0].(bool)
    }

    return s.Dialector.CreateDatabase(database, charset, ifNotExists, s.DB)
}

// DropDatabase Drops a database. Will throw a Database Exception if it cannot.
func (s *Schema) DropDatabase(database string) (err error) {
    return s.Dialector.DropDatabase(database, s.DB)
}

// CreateTable Creates a table.
func (s *Schema) CreateTable(table string, fields []map[string]interface{}, args ...interface{}) (err error) {
    var (
        primaryKeys     []string
        ifNotExists     bool = true
        engine          string = "InnoDB"
        charset         string = "utf8_general_ci"
        foreignKeys     []map[string]interface{}
    )

    switch len(args) {
    case 5:
        primaryKeys = args[0].([]string)
        ifNotExists = args[1].(bool)
        engine      = args[2].(string)
        charset     = args[3].(string)
        foreignKeys = args[4].([]map[string]interface{})
    case 4:
        primaryKeys = args[0].([]string)
        ifNotExists = args[1].(bool)
        engine      = args[2].(string)
        charset     = args[3].(string)
    case 3:
        primaryKeys = args[0].([]string)
        ifNotExists = args[1].(bool)
        engine      = args[2].(string)
    case 2:
        primaryKeys = args[0].([]string)
        ifNotExists = args[1].(bool)
    case 1:
        primaryKeys = args[0].([]string)
    default:
    }

    return s.Dialector.CreateTable(table, fields, primaryKeys, ifNotExists, engine, charset, foreignKeys, s.DB)
}

// DropTable Drops a table. Will throw a Database Exception if it cannot.
func (s *Schema) DropTable(table string) (err error) {
    return s.Dialector.DropTable(table, s.DB)
}

// RenameTable Renames a table. Will throw a Database Exception if it cannot.
func (s *Schema) RenameTable(oldTable, newTable string) (err error) {
    return s.Dialector.RenameTable(oldTable, newTable, s.DB)
}

// TruncateTable Truncates a table.
func (s *Schema) TruncateTable(table string) (err error) {
    return s.Dialector.TruncateTable(table, s.DB)
}

// AnalyzeTable Analyzes a table.
func (s *Schema) AnalyzeTable(table string) bool {
    return s.tableMaintenance("ANALYZE TABLE", table)
}

// CheckTable Check a table.
func (s *Schema) CheckTable(table string) bool {
    return s.tableMaintenance("CHECK TABLE", table)
}

// OptimizeTable Optimize a table.
func (s *Schema) OptimizeTable(table string) bool {
    return s.tableMaintenance("OPTIMIZE TABLE", table)
}

// RepairTable Repair a table.
func (s *Schema) RepairTable(table string) bool {
    return s.tableMaintenance("REPAIR TABLE", table)
}

func (s *Schema) tableMaintenance(operation string, table string) bool {
    // SQLite does't support table maintenance
    if s.Dialector.Name() == "sqlite" {
        return false
    }

    sqlStr := operation + " " + s.QuoteTable(table)
    rows, err := s.Rows(sqlStr)
    if err != nil {
        s.AddError(err)
    }

    var null interface{}
    var msgType, msgText string

    for rows.Next() {
        err := rows.Scan(null, null, msgType, msgText)
        if err != nil {
            s.AddError(err)
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
func (s *Schema) TableExists(table string) bool {
    return s.Dialector.TableExists(table, s.DB)
}

// FieldExists Checks if given field(s) in a given table exists.
func (s *Schema) FieldExists(table string, value interface{}) bool {
    return s.Dialector.FieldExists(table, value, s.DB)
}

// CreateIndex Creates an index on that table.
func (s *Schema) CreateIndex(table string, indexColumns interface{}, indexName, index string) (err error) {
    return s.Dialector.CreateIndex(table, indexColumns, indexName, index, s.DB)
}

// RenameIndex Rename an index from a table.
func (s *Schema) RenameIndex(table, oldName, newName string) (err error) {
    return s.Dialector.RenameIndex(table, oldName, newName, s.DB)
}

// DropIndex Drop an index from a table.
func (s *Schema) DropIndex(table string, indexName string) (err error) {
    return s.Dialector.DropIndex(table, indexName, s.DB)
}

// AddForeignKey Adds a single foreign key to a table
func (s *Schema) AddForeignKey(table string, foreignKey []map[string]interface{}) (err error) {
    return s.Dialector.AddForeignKey(table, foreignKey, s.DB)
}

// DropForeignKey Drops a foreign key from a table
func (s *Schema) DropForeignKey(table string, fkName string) (err error) {
    return s.Dialector.DropForeignKey(table, fkName, s.DB)
}

// AddFields adds fields to a table.
func (s *Schema) AddFields(table string, fields []map[string]interface{}) error {
    return s.Dialector.AddFields(table, fields, s.DB)
}

// DropFields drops fields from a table.
func (s *Schema) DropFields(table string, value interface{}) error {
    return s.Dialector.DropFields(table, value, s.DB)
}

// ModifyFields alters fields in a table.
func (s *Schema) ModifyFields(table string, fields []map[string]interface{}) error {
    return s.Dialector.ModifyFields(table, fields, s.DB)
}

// AlterFields is ...
//func (s *Schema) AlterFields(alterType, table string, fields interface{}) (err error) {
    //return s.Dialector.AlterFields(alterType, table, fields, s.DB)
//}

// ProcessFields is ...
func (s *Schema) ProcessFields(fields []map[string]interface{}, prefix string) string {
    return s.Dialector.ProcessFields(fields, prefix, s.DB)
}

