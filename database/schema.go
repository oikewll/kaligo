package database

import (
    "fmt"
    "go/ast"
    "reflect"
    "sync"
)

// Schema is the struct for MySQL DATE type
type Schema struct {
    *Query

    Name           string
    ModelType      reflect.Type
    Table          string
    Fields         []*Field
    FieldsByName   map[string]*Field // 通过 struct 字段名查询
    FieldsByDBName map[string]*Field // 通过 数据库 字段名查询
    err            error
    initialized    chan struct{}
    cacheStore     *sync.Map
}

// Indexes 数据表索引，字段名记得要大写开头，才是public，否则访问不到
type Indexes struct {
    Table   string `field:"Table"`
    Name    string `field:"Key_name"`
    Column  string `field:"Column_name"`
    Order   int64  `field:"Seq_in_index"`
    Type    string `field:"Index_type"`
    Primary bool   `field:"Key_name == 'PRIMARY'"`
    Unique  bool   `field:"Non_unique == 0"`
    Null    bool   `field:"Null == 'YES'"`
    Ascend  bool   `field:"Collation == 'A'"`
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
// 3、map[string]any、[]map[string]any{}、[]int64 都不会到这里来
func Parse(dest any, cacheStore *sync.Map) (*Schema, error) {
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

    // 等待其他协程数据
    if v, ok := cacheStore.Load(modelType); ok {
        s := v.(*Schema)
        // Wait for the initialization of other goroutines to complete
        <-s.initialized
        return s, s.err
    }

    //modelValue := reflect.New(modelType)
    // 其实没必要去弄表名，因为schema又不是一张表
    //tableName  := ToDBName(modelType.Name())

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
                schema.FieldsByName[field.Name] = field
                //fmt.Printf("66666 ---> %v\n", field)
            }
        }

        field.setupValuerAndSetter()
    }

    return schema, schema.err
}

/* vim: set expandtab: */
