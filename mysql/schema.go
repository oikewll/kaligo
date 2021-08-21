package mysql

import (
    "fmt"
    "sync"
    "reflect"
    //"strings"
    "go/ast"
)

// Schema is the struct for MySQL DATE type
type Schema struct {
    *DB

    Name                      string
    ModelType                 reflect.Type
    Table                     string
    Fields                    []*Field
    FieldsByName              map[string]*Field     // 通过 struct 字段名查询
    FieldsByColumn            map[string]*Field     // 通过 数据库 字段名查询
    err                       error
    initialized               chan struct{}
    cacheStore                *sync.Map
}

func (s Schema) String() string {
    if s.ModelType.Name() == "" {
        return fmt.Sprintf("%s(%s)", s.Name, s.Table)
    }
    return fmt.Sprintf("%s.%s", s.ModelType.PkgPath(), s.ModelType.Name())
}

// LookUpField is 通过表名 或者 数据库名 查询字段
func (s Schema) LookUpField(name string) *Field {
    if field, ok := s.FieldsByColumn[name]; ok {
        return field
    }

    if field, ok := s.FieldsByName[name]; ok {
        return field
    }
    return nil
}

// Parse get data type from dialector
// 1、传入 struct，比如 &User{}，生成 User对应的 []*Field
// 2、当查询了数据 rows.Next 循环的时候，去上面的 []*Field 找对应的字段，通过 field.Set 去设置值
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
        FieldsByColumn: map[string]*Field{},
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

    // 循环字段，给 FieldsByColumn、FieldsByName 赋值
    for _, field := range schema.Fields {
        if field.Column == "" && field.DataType != "" {
            // 驼峰转下划线，数据库字段都是下划线写法
            field.Column = TransFieldName(field.Name)
        }

        if field.Column != "" {
            if _, ok := schema.FieldsByColumn[field.Column]; !ok {
                schema.FieldsByColumn[field.Column] = field
                schema.FieldsByName[field.Name]     = field
            }
        }

        field.setupValuerAndSetter()
    }

    return schema, schema.err
}

// ListTables If a table name is given it will return the table name with the configured
// prefix. If not, then just the prefix is returnd
func (s *Schema) ListTables(like string) []string {
    var sqlStr string
    if  like != "" {
        sqlStr += "SHOW TABLES LIKE " + s.Quote(like)
    } else {
        sqlStr += "SHOW TABLES"
    }

    var tables []string
    tables = append(tables, "111")
    return tables
}

// ListColumns Lists all of the columns in a table. Optionally, a LIKE string can be
// used to search for specific fields.
func (s *Schema) ListColumns(table string, like string) map[string] map[string]string {
    table = s.QuoteTable(table)

    var sqlStr string
    if  like != "" {
        sqlStr += "SHOW FULL COLUMNS FROM " + table + " LIKE " + s.Quote(like)
    } else {
        sqlStr += "SHOW FULL COLUMNS FROM " + table
    }

    var column map[string]string
    column["name"] = "Field"
    var columns map[string] map[string]string
    columns["Field"] = column

    return columns
}

// ListIndexes Lists all of the idexes in a table. Optionally, a LIKE string can be
// used to search for specific indexes by name.
func (s *Schema) ListIndexes(table string, like string) []map[string]string {
    table = s.QuoteTable(table)

    var sqlStr string
    if  like != "" {
        sqlStr += "SHOW INDEX FROM " + table + " WHERE " + s.QuoteIdentifier("Key_name") + " LIKE " + s.Quote(like)
    } else {
        sqlStr += "SHOW INDEX FROM " + table
    }

    var indexes []map[string]string
    mapName := map[string]string {
        "name"      : "Key_name",
        "column"    : "Column_name",
        "order"     : "Seq_in_index",
        "type"      : "Index_type",
        "primary"   : "true",
        "unique"    : "Non_unique",
        "null"      : "YES",
        "ascending" : "Collation",
    }
    indexes = append(indexes, mapName)

    return indexes
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
