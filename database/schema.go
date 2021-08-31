package database

import (
    "errors"
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
            field.Column = toDBName(field.Name)
        }

        if field.Column != "" {
            if _, ok := schema.FieldsByColumn[field.Column]; !ok {
                schema.FieldsByColumn[field.Column] = field
                schema.FieldsByName[field.Name]     = field
                //fmt.Printf("66666 ---> %v\n", field)
            }
        }

        field.setupValuerAndSetter()
    }

    return schema, schema.err
}

// ListTables If a table name is given it will return the table name with the configured
// prefix. If not, then just the prefix is returnd
func (s *Schema) ListTables(args ...string) []string {
    var like, sqlStr string
    if len(args) > 0 {
        like = args[0]
    }

    if  like != "" {
        sqlStr += "SHOW TABLES LIKE " + s.Quote("%" + like + "%")
    } else {
        sqlStr += "SHOW TABLES"
    }

    var tables []string
    s.Query.Query(sqlStr).Scan(&tables).Execute()
    return tables
}

// ListColumns Lists all of the columns in a table. Optionally, a LIKE string can be
// used to search for specific fields.
func (s *Schema) ListColumns(table string, args ...string) []Column {
    table = s.QuoteTable(table)
    var like, sqlStr string
    if len(args) > 0 {
        like = args[0]
    }

    if  like != "" {
        sqlStr += "SHOW FULL COLUMNS FROM " + table + " LIKE " + s.Quote("%" + like + "%")
    } else {
        sqlStr += "SHOW FULL COLUMNS FROM " + table
    }

    columns := []Column{}
    s.Query.Query(sqlStr).Scan(&columns).Execute()
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
func (s *Schema) CreateDatabase(database string, args ...interface{}) (err error) {
    database = s.QuoteTable(database)

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

    sqlStr := "CREATE DATABASE "
    if ifNotExists {
        sqlStr += "IF NOT EXISTS "
    }
    sqlStr += s.QuoteIdentifier(database) + s.processCharset(charset, true)

    _, err = s.Query.Exec(sqlStr)
    if err != nil {
        s.Query.DB.AddError(err)
    }
    return err
}

// DropDatabase Drops a database. Will throw a Database Exception if it cannot.
func (s *Schema) DropDatabase(database string) (err error) {
    sqlStr := "DROP DATABASE " + s.QuoteIdentifier(database)

    _, err = s.Query.Exec(sqlStr)
    if err != nil {
        s.Query.DB.AddError(err)
    }
    return err
}

// DropTable Drops a table. Will throw a Database Exception if it cannot.
func (s *Schema) DropTable(table string) (err error) {
    sqlStr := "DROP TABLE IF EXISTS "
    sqlStr += s.QuoteIdentifier(s.Query.DB.TablePrefix(table))
    
    _, err = s.Query.Exec(sqlStr)
    if err != nil {
        s.Query.DB.AddError(err)
    }
    return err
}

// RenameTable Renames a table. Will throw a Database Exception if it cannot.
func (s *Schema) RenameTable(table string, newTable string) (err error) {
    sqlStr := "RENAME TABLE "
    sqlStr += s.QuoteIdentifier(s.Query.DB.TablePrefix(table))
    sqlStr += " TO "
    sqlStr += s.QuoteIdentifier(s.Query.DB.TablePrefix(newTable))

    _, err = s.Query.Exec(sqlStr)
    if err != nil {
        s.Query.DB.AddError(err)
    }
    return err
}

// CreateTable Creates a table.
func (s *Schema) CreateTable(table string, fields map[string] map[string]interface{}, args ...interface{}) (err error) {
    var (
        primaryKeys     []string
        ifNotExists     bool = true
        engine          string = "InnoDB"
        charset         string = "utf8"
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

    sqlStr := "CREATE TABLE ";
    if ifNotExists {
        sqlStr += "IF NOT EXISTS "
    }

    sqlStr += s.QuoteIdentifier(s.Query.DB.TablePrefix(table)) + " ("
    sqlStr += s.processFields(fields, "")

    if len(primaryKeys) > 0 {
        for k, v := range primaryKeys {
            primaryKeys[k] = s.QuoteIdentifier(v)
        }
        sqlStr += ",\n\tPRIMARY KEY (" + strings.Join(primaryKeys, ", ") + ")"
    }

    if len(foreignKeys) > 0 {
        sqlStr += s.processForeignKeys(foreignKeys)
    }

    sqlStr += "\n)"
    if engine != "" {
        sqlStr += " ENGINE = " + engine + " "
    }

    sqlStr += s.processCharset(charset, true, "") + ";"

    _, err = s.Query.Exec(sqlStr)
    if err != nil {
        s.Query.DB.AddError(err)
    }
    return err
}

// TruncateTable Truncates a table.
func (s *Schema) TruncateTable(table string) (err error) {
    sqlStr := "TRUNCATE TABLE "
    sqlStr += s.QuoteIdentifier(s.Query.DB.TablePrefix(table))

    _, err = s.Query.Exec(sqlStr)
    if err != nil {
        s.Query.DB.AddError(err)
    }
    return err
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
    sqlStr := operation + " " + s.QuoteIdentifier(s.Query.DB.TablePrefix(table))
    rows, err := s.Query.Rows(sqlStr)
    if err != nil {
        s.Query.DB.AddError(err)
    }

    // 字段名记得要大写开头，才是public，否则访问不到
    type Operation struct {
        Table   string `field:"Table"` 
        Op      string `field:"Op"`
        MsgType string `field:"Msg_type"`
        MsgText string `field:"Msg_text"`
    }

    ops := []Operation{}
    for rows.Next() {
        op := Operation{}
        o := reflect.ValueOf(&op).Elem()
        numCols := o.NumField()
        columns := make([]interface{}, numCols)
        for i := 0; i < numCols; i++ {
            field := o.Field(i)
            columns[i] = field.Addr().Interface()
        }
        err := rows.Scan(columns...)
        if err != nil {
            s.AddError(err)
        } else {
            ops = append(ops, op)
        }
    }

    msgType := ops[0].MsgType
    msgText := ops[0].MsgText
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
    sqlStr := "SELECT * FROM "
    sqlStr += s.QuoteIdentifier(s.Query.DB.TablePrefix(table))
    sqlStr += " LIMIT 1"

    _, err := s.Query.Rows(sqlStr)
    if err != nil {
        s.Query.DB.AddError(err)
        return false
    }
    return true
}

// FieldExists Checks if given field(s) in a given table exists.
func (s *Schema) FieldExists(table string, value interface{}) bool {
    var columns []string    
    switch value.(type) {
    case string:
        columns = append(columns, value.(string))
    default:
        columns = value.([]string)
    }

    for k, v := range columns {
        columns[k] = s.QuoteIdentifier(v)
    }

    sqlStr := "SELECT "
    sqlStr += strings.Join(columns, ", ")
    sqlStr += " FROM "
    sqlStr += s.QuoteIdentifier(s.Query.DB.TablePrefix(table))
    sqlStr += " LIMIT 1"

    _, err := s.Query.Rows(sqlStr)
    if err != nil {
        s.Query.DB.AddError(err)
        return false
    }
    return true
}

// CreateIndex Creates an index on that table.
func (s *Schema) CreateIndex(table string, indexColumns interface{}, indexName string, index string) (err error) {
    var sqlStr string    
    acceptedIndex := []string{"UNIQUE", "FULLTEXT", "SPATIAL", "NONCLUSTERED", "PRIMARY"}

    // make sure the index type is uppercase
    if index != "" {
        index = strings.ToUpper(index)
    }

    if indexName == "" {
        switch indexColumns.(type) {
        case []string:
            for k, v := range indexColumns.([]string) {
                if indexName == "" {
                    indexName += ""
                } else {
                    indexName += "_"
                }

                key := ToString(k)
                if IsNumeric(key) {
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

    if index == "PRIMARY" {
        sqlStr = "ALTER TABLE "
        sqlStr += s.QuoteIdentifier(s.Query.DB.TablePrefix(table))
        sqlStr += " ADD PRIMARY KEY "

        switch indexColumns.(type) {
        case []string:
            columns := ""
            for k, v := range indexColumns.([]string) {
                if columns == "" {
                    columns += ""
                } else {
                    columns += ", "
                }

                key := ToString(k)
                if IsNumeric(key) {
                    columns += s.QuoteIdentifier(v)
                } else {
                    columns += s.QuoteIdentifier(key) + " " + strings.ToUpper(v)
                }
            }
            sqlStr += " (" + columns + ")"
        }
    } else {
        sqlStr = "CREATE "

        if index != "" {
            if InSlice(index, &acceptedIndex) {
                sqlStr += index + " "
            }
        }

        sqlStr += "INDEX "
        sqlStr += s.QuoteIdentifier(indexName)
        sqlStr += " ON "
        sqlStr += s.QuoteIdentifier(s.Query.DB.TablePrefix(table))

        switch indexColumns.(type) {
        case []string:
            columns := ""
            for k, v := range indexColumns.([]string) {
                if columns == "" {
                    columns += ""
                } else {
                    columns += ", "
                }

                key := ToString(k)
                if IsNumeric(key) {
                    columns += s.QuoteIdentifier(v)
                } else {
                    columns += s.QuoteIdentifier(key) + " " + strings.ToUpper(v)
                }
            }
            sqlStr += " (" + columns + ")"
        default:
            sqlStr += " (" + s.QuoteIdentifier(indexColumns.(string)) + ")"
        }
    }

    _, err = s.Query.Exec(sqlStr)
    if err != nil {
        s.Query.DB.AddError(err)
    }
    return err
}

// DropIndex Drop an index from a table.
func (s *Schema) DropIndex(table string, indexName string) (err error) {
    var sqlStr string    
    if strings.ToUpper(indexName) == "PRIMARY" {
        sqlStr = "ALTER TABLE " + s.QuoteIdentifier(s.Query.DB.TablePrefix(table))
        sqlStr += " DROP PRIMARY KEY"
    } else {
        sqlStr = "DROP INDEX " + s.QuoteIdentifier(indexName)
        sqlStr += " ON "+ s.QuoteIdentifier(s.Query.DB.TablePrefix(table))
    }

    _, err = s.Query.Exec(sqlStr)
    if err != nil {
        s.Query.DB.AddError(err)
    }
    return err
}

// AddForeignKey Adds a single foreign key to a table
func (s *Schema) AddForeignKey(table string, foreignKey []map[string]interface{}) (err error) {
    sqlStr := "ALTER TABLE " + s.QuoteIdentifier(s.Query.DB.TablePrefix(table)) + " ADD " + strings.TrimLeft(s.processForeignKeys(foreignKey), ",")

    res, err := s.Query.Exec(sqlStr)
    if err != nil {
        s.Query.DB.AddError(err)
    }

    _, err = res.RowsAffected()
    if err != nil {
        s.Query.DB.AddError(err)
    }
    return err
}

// DropForeignKey Drops a foreign key from a table
func (s *Schema) DropForeignKey(table string, fkName string) (err error) {
    sqlStr := "ALTER TABLE " + s.QuoteIdentifier(s.Query.DB.TablePrefix(table)) + " DROP FOREIGN KEY " + s.QuoteIdentifier(fkName)

    _, err = s.Query.Exec(sqlStr)
    if err != nil {
        s.Query.DB.AddError(err)
    }
    return err
}

func (s *Schema) processForeignKeys(foreignKeys []map[string]interface{}) string {
    var fkList []string    
    for _, definition := range foreignKeys {
        // some sanity checks
        if _, ok := definition["key"]; !ok {
            s.AddError(errors.New("Foreign keys on processForeignKeys() must specify a foreign key name"))
            return ""
        }

        if _, ok := definition["reference"]; !ok {
            s.AddError(errors.New("Foreign keys on processForeignKeys() must specify a foreign key reference"))
            return ""
        }

        reference := definition["reference"].(map[string]string)

        var referenceTable, referenceColumn string    
        if table, ok := reference["table"]; ok {
            referenceTable = table
        } else {
            s.AddError(errors.New("Foreign keys on processForeignKeys() must specify a reference table name"))
            return ""
        }

        if column, ok := reference["column"]; ok {
            referenceColumn = column
        } else {
            s.AddError(errors.New("Foreign keys on processForeignKeys() must specify a reference column name"))
            return ""
        }

        var sqlStr string    
        if table, ok := definition["constraint"]; ok {
            sqlStr += " CONSTRAINT " + s.QuoteIdentifier(s.Query.DB.TablePrefix(table.(string))) 
        }

        sqlStr += " FOREIGN KEY (" + s.QuoteIdentifier(definition["key"].(string)) + ")"
        sqlStr += " REFERENCES " + s.QuoteIdentifier(s.Query.DB.TablePrefix(referenceTable)) + " ("
        referenceColumnArr := strings.Split(referenceColumn, ",")
        for k, v := range referenceColumnArr {
            referenceColumnArr[k] = s.QuoteIdentifier(v)
        }
        sqlStr += strings.Join(referenceColumnArr, ", ")
        sqlStr += ")"

        if val, ok := definition["on_update"]; ok {
            sqlStr += " ON UPDATE " + val.(string) 
        }
        if val, ok := definition["on_delete"]; ok {
            sqlStr += " ON DELETE " + val.(string)
        }

        fkList = append(fkList, "\n\t" + strings.TrimLeft(sqlStr, " "))
    }

    return ", " + strings.Join(fkList, ",")
}

// AddFields adds fields to a table.
func (s *Schema) AddFields(table string, fields map[string] map[string]interface{}) error {
    return s.alterFields("ADD", table, fields)
}

// DropFields drops fields from a table.
func (s *Schema) DropFields(table string, value interface{}) error {
    var fields []string
    switch value.(type) {
    case string:
        fields = append(fields, value.(string))
    default:
        fields = value.([]string)
    }
    return s.alterFields("DROP", table, fields)
}

// ModifyFields alters fields in a table.
func (s *Schema) ModifyFields(table string, fields map[string] map[string]interface{}) error {
    return s.alterFields("MODIFY", table, fields)
}

// 返回的 rowsAffected 为 0，但是实际上是修改成功了，不知道为什么
func (s *Schema) alterFields(alterType string, table string, fields interface{}) (err error) {
    sqlStr := "ALTER TABLE " + s.QuoteIdentifier(s.Query.DB.TablePrefix(table)) + " "

    if alterType == "DROP" {
        var dropFields, sqlDropFields []string
        dropFields = fields.([]string)
        for _, field := range dropFields {
            sqlDropFields = append(sqlDropFields, "DROP " + s.QuoteIdentifier(field))
        }
        sqlStr += strings.Join(sqlDropFields, ", ")
    } else {
        fieldMaps := fields.(map[string] map[string]interface{})

        useBrackets := true    
        if InSlice(alterType, &[]string{"ADD", "CHANGE", "MODIFY"}) {
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
        sqlStr += s.processFields(fieldMaps, prefix)
        if useBrackets {
            sqlStr += ")"
        }
    }

    _, err = s.Query.Exec(sqlStr)
    if err != nil {
        s.Query.DB.AddError(err)
    }
    return err
}

func (s *Schema) processCharset(charset string, isDefault bool, args ...string) string {

    var collation string    
    if len(args) > 0 {
        collation = args[0]
    }

    // utf8_unicode_ci
    charsets := strings.Split(charset, "_")
    if collation == "" && len(charsets) > 1 {
        collation = charset     // utf8_unicode_ci
        charset   = charsets[0] // utf8
    }

    charset = " CHARACTER SET " + charset
    if isDefault {
        charset = " DEFAULT " + charset
    }

    if collation != "" {
        if isDefault {
            charset += " DEFAULT"
        }
        charset += " COLLATE " + collation
    }

    return charset
}

func (s *Schema) processFields(fields map[string] map[string]interface{}, prefix string) string {
    var sqlFields []string    

    for field, dict := range fields {
        dict = MapChangeKeyCase(dict, true)
        tmpPrefix := prefix
        if value, ok := dict["NAME"]; ok && field != value && tmpPrefix == "MODIFY " {
            tmpPrefix = "CHANGE "
        }
        sqlStr := "\n\t" + tmpPrefix
        sqlStr += s.QuoteIdentifier(field)

        if value, ok := dict["NAME"]; ok && field != value {
            sqlStr += " " + s.QuoteIdentifier(value.(string)) + " "
        }

        if value, ok := dict["TYPE"]; ok {
            sqlStr += " " + value.(string)
        }

        if value, ok := dict["CONSTRAINT"]; ok {
            sqlStr += "(" + ToString(value) + ")"
        }

        if value, ok := dict["CHARSET"]; ok {
            sqlStr += s.processCharset(value.(string), false)
        }

        if value, ok := dict["UNSIGNED"]; ok && value.(bool) == true {
            sqlStr += " UNSIGNED"
        }

        if value, ok := dict["DEFAULT"]; ok {
            sqlStr += " DEFAULT " + s.Quote(value.(string))
        }

        if value, ok := dict["NULL"]; ok && value.(bool) == true {
            sqlStr += " NULL"
        } else {
            sqlStr += " NOT NULL"
        }

        if value, ok := dict["AUTO_INCREMENT"]; ok && value.(bool) == true {
            sqlStr += " AUTO_INCREMENT"
        }

        if value, ok := dict["PRIMARY_KEY"]; ok && value.(bool) == true {
            sqlStr += " PRIMARY_KEY"
        }

        if value, ok := dict["COMMENT"]; ok {
            sqlStr += " COMMENT " + s.Escape(value.(string))
        }

        if value, ok := dict["FIRST"]; ok && value.(bool) == true {
            sqlStr += " FIRST"
        }

        if value, ok := dict["AFTER"]; ok {
            sqlStr += " AFTER " + s.QuoteIdentifier(value.(string))
        }

        sqlFields = append(sqlFields, sqlStr)
    }

    return strings.Join(sqlFields, ", ")
}

