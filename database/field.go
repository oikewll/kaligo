package database

import (
    "database/sql"
    "database/sql/driver"
    "fmt"
    "reflect"
    "strconv"
    "strings"
    "time"

    "github.com/jinzhu/now"
)

// DataType is ...
type DataType string

// TimeType is ...
type TimeType int64

// TimeReflectType is Time反射类型
var TimeReflectType = reflect.TypeOf(time.Time{})

// 模拟enum，因为Golang没有，但是模拟的还是没有enum的性质，比如：
// 1、固定性(你给他一个不存在的值他也吃下去)；2、可轮询；3、独立的命名空间；
const (
    // UnixSecond is ...
    UnixSecond TimeType = 1
    // UnixMillisecond is ...
    UnixMillisecond TimeType = 2
    // UnixNanosecond is ...
    UnixNanosecond TimeType = 3
)

const (
    // Bool is bool type
    Bool DataType = "bool"
    // Int is int type
    Int DataType = "int"
    // Uint is uint type
    Uint DataType = "uint"
    // Float is float type
    Float DataType = "float"
    // String is string type
    String DataType = "string"
    // Time is time type
    Time DataType = "time"
    // Bytes is bytes type
    Bytes DataType = "bytes"
)

// Column is ...
type Column struct {
    Name          string   // 字段名
    DBName        string   // 数据库栏位名
    DataType      DataType // 数据类型，这里模拟了一个enum
    Size          int64    // 栏位长度
    Precision     int64    // 精是，小数点保留多少位，decimal(10, 2) 里面的 2
    NotNull       bool     // 是否允许为空
    DefaultValue  string   // 默认值
    Unique        bool     // 是否唯一索引
    PrimaryKey    bool     // 是否主键
    AutoIncrement bool     // 是否自增
    Comment       string   // 注释
    Readable      bool     // Privileges select
    Creatable     bool     // Privileges insert
    Updatable     bool     // Privileges update
    Extra         string   // 扩展
    CryptKey      string   // 密钥 AES_DECRYPT、AES_ENCRYPT
}

// Field is ...
type Field struct {
    *Column

    FieldType         reflect.Type
    IndirectFieldType reflect.Type
    StructField       reflect.StructField
    StructTag         reflect.StructTag
    TagSettings       map[string]string

    ReflectValueOf func(reflect.Value) reflect.Value
    ValueOf        func(reflect.Value) (value any, zero bool)
    Set            func(reflect.Value, any) error
}

// ParseField is 解析字段
func ParseField(fieldStruct reflect.StructField) *Field {
    //fmt.Printf("ParseField %v\n", FormatJSON(fieldStruct))

    field := &Field{
        FieldType:         fieldStruct.Type,
        IndirectFieldType: fieldStruct.Type,
        StructField:       fieldStruct,
        StructTag:         fieldStruct.Tag,
        TagSettings:       ParseTagSetting(fieldStruct.Tag.Get("db"), ";"),
    }
    field.Name = fieldStruct.Name

    // Array、Slice、Struct 需要通过 Elem() 获取指针指向的值(子元素的值)，其他类型直接拿到的就是他的值
    for field.IndirectFieldType.Kind() == reflect.Ptr {
        field.IndirectFieldType = field.IndirectFieldType.Elem()
    }

    fieldValue := reflect.New(field.IndirectFieldType)
    valuer, isValuer := fieldValue.Interface().(driver.Valuer)
    if isValuer {
        if v, err := valuer.Value(); reflect.ValueOf(v).IsValid() && err == nil {
            fieldValue = reflect.ValueOf(v)
        }

        var getRealFieldValue func(reflect.Value)
        getRealFieldValue = func(v reflect.Value) {
            rv := reflect.Indirect(v)
            if rv.Kind() == reflect.Struct && !rv.Type().ConvertibleTo(TimeReflectType) {
                for i := 0; i < rv.Type().NumField(); i++ {
                    newFieldType := rv.Type().Field(i).Type
                    for newFieldType.Kind() == reflect.Ptr {
                        newFieldType = newFieldType.Elem()
                    }

                    fieldValue = reflect.New(newFieldType)

                    if rv.Type() != reflect.Indirect(fieldValue).Type() {
                        getRealFieldValue(fieldValue)
                    }

                    if fieldValue.IsValid() {
                        return
                    }

                    for key, value := range ParseTagSetting(field.IndirectFieldType.Field(i).Tag.Get("db"), ";") {
                        if _, ok := field.TagSettings[key]; !ok {
                            field.TagSettings[key] = value
                        }
                    }
                }
            }
        }

        getRealFieldValue(fieldValue)
    }
    //fmt.Printf("field.TagSettings = %v\n", field.TagSettings)

    if val, ok := field.TagSettings["COLUMN"]; ok {
        field.DBName = val
    }

    // 是否主键索引
    if val, ok := field.TagSettings["PRIMARYKEY"]; ok && CheckTruth(val) {
        field.PrimaryKey = true
    }

    // 是否自增字段
    if val, ok := field.TagSettings["AUTOINCREMENT"]; ok && CheckTruth(val) {
        field.AutoIncrement = true
    }

    // 默认值
    if val, ok := field.TagSettings["DEFAULT"]; ok {
        field.DefaultValue = val
    }
    // default value is function or null or blank (primary keys)
    field.DefaultValue = strings.TrimSpace(field.DefaultValue)

    // 类型大小
    if num, ok := field.TagSettings["SIZE"]; ok {
        field.Size = ToInt64(num)
    }

    // 是否可以为NULL
    if val, ok := field.TagSettings["NOTNULL"]; ok && CheckTruth(val) {
        field.NotNull = true
    }

    // 是否唯一索引
    if val, ok := field.TagSettings["UNIQUE"]; ok && CheckTruth(val) {
        field.Unique = true
    }

    // 备注
    if val, ok := field.TagSettings["COMMENT"]; ok {
        field.Comment = val
    }

    switch reflect.Indirect(fieldValue).Kind() {
    case reflect.Bool:
        field.DataType = Bool
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        field.DataType = Int
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        field.DataType = Uint
    case reflect.Float32, reflect.Float64:
        field.DataType = Float
    case reflect.String:
        field.DataType = String
    case reflect.Struct:
        if _, ok := fieldValue.Interface().(*time.Time); ok {
            field.DataType = Time
        } else if fieldValue.Type().ConvertibleTo(TimeReflectType) {
            field.DataType = Time
        } else if fieldValue.Type().ConvertibleTo(reflect.TypeOf(&time.Time{})) {
            field.DataType = Time
        }
    case reflect.Array, reflect.Slice:
        if reflect.Indirect(fieldValue).Type().Elem() == reflect.TypeOf(uint8(0)) {
            field.DataType = Bytes
        }
    }

    if val, ok := field.TagSettings["TYPE"]; ok {
        switch DataType(strings.ToLower(val)) {
        case Bool, Int, Uint, Float, String, Time, Bytes:
            field.DataType = DataType(strings.ToLower(val))
        default:
            field.DataType = DataType(val)
        }
    }

    if field.Size == 0 {
        switch reflect.Indirect(fieldValue).Kind() {
        case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64, reflect.Float64:
            field.Size = 64
        case reflect.Int8, reflect.Uint8:
            field.Size = 8
        case reflect.Int16, reflect.Uint16:
            field.Size = 16
        case reflect.Int32, reflect.Uint32, reflect.Float32:
            field.Size = 32
        }
    }

    return field
}

// create valuer, setter when parse struct
// 解析fields值到value中
func (field *Field) setupValuerAndSetter() {
    // ValueOf
    switch {
    case len(field.StructField.Index) == 1:
        field.ValueOf = func(value reflect.Value) (any, bool) {
            fieldValue := reflect.Indirect(value).Field(field.StructField.Index[0])
            return fieldValue.Interface(), fieldValue.IsZero()
        }
    case len(field.StructField.Index) == 2 && field.StructField.Index[0] >= 0:
        field.ValueOf = func(value reflect.Value) (any, bool) {
            fieldValue := reflect.Indirect(value).Field(field.StructField.Index[0]).Field(field.StructField.Index[1])
            return fieldValue.Interface(), fieldValue.IsZero()
        }
    default:
        field.ValueOf = func(value reflect.Value) (any, bool) {
            v := reflect.Indirect(value)
            for _, idx := range field.StructField.Index {
                if idx >= 0 {
                    v = v.Field(idx)
                } else {
                    v = v.Field(-idx - 1)

                    if v.Type().Elem().Kind() != reflect.Struct {
                        return nil, true
                    }

                    if !v.IsNil() {
                        v = v.Elem()
                    } else {
                        return nil, true
                    }
                }
            }
            return v.Interface(), v.IsZero()
        }
    }

    // ReflectValueOf
    switch {
    case len(field.StructField.Index) == 1:
        field.ReflectValueOf = func(value reflect.Value) reflect.Value {
            return reflect.Indirect(value).Field(field.StructField.Index[0])
        }
    case len(field.StructField.Index) == 2 && field.StructField.Index[0] >= 0 && field.FieldType.Kind() != reflect.Ptr:
        field.ReflectValueOf = func(value reflect.Value) reflect.Value {
            return reflect.Indirect(value).Field(field.StructField.Index[0]).Field(field.StructField.Index[1])
        }
    default:
        field.ReflectValueOf = func(value reflect.Value) reflect.Value {
            v := reflect.Indirect(value)
            for idx, fieldIdx := range field.StructField.Index {
                if fieldIdx >= 0 {
                    v = v.Field(fieldIdx)
                } else {
                    v = v.Field(-fieldIdx - 1)
                }

                if v.Kind() == reflect.Ptr {
                    if v.Type().Elem().Kind() == reflect.Struct {
                        if v.IsNil() {
                            v.Set(reflect.New(v.Type().Elem()))
                        }
                    }

                    if idx < len(field.StructField.Index)-1 {
                        v = v.Elem()
                    }
                }
            }
            return v
        }
    }

    // Setter 回调函数
    fallbackSetter := func(value reflect.Value, v any, setter func(reflect.Value, any) error) (err error) {
        if v == nil {
            field.ReflectValueOf(value).Set(reflect.New(field.FieldType).Elem())
        } else {
            reflectV := reflect.ValueOf(v)
            // Optimal value type acquisition for v
            reflectValType := reflectV.Type()

            if reflectValType.AssignableTo(field.FieldType) {
                field.ReflectValueOf(value).Set(reflectV)
                return
            } else if reflectValType.ConvertibleTo(field.FieldType) {
                field.ReflectValueOf(value).Set(reflectV.Convert(field.FieldType))
                return
            } else if field.FieldType.Kind() == reflect.Ptr {
                fieldValue := field.ReflectValueOf(value)
                fieldType := field.FieldType.Elem()

                if reflectValType.AssignableTo(fieldType) {
                    if !fieldValue.IsValid() {
                        fieldValue = reflect.New(fieldType)
                    } else if fieldValue.IsNil() {
                        fieldValue.Set(reflect.New(fieldType))
                    }
                    fieldValue.Elem().Set(reflectV)
                    return
                } else if reflectValType.ConvertibleTo(fieldType) {
                    if fieldValue.IsNil() {
                        fieldValue.Set(reflect.New(fieldType))
                    }

                    fieldValue.Elem().Set(reflectV.Convert(fieldType))
                    return
                }
            }

            if reflectV.Kind() == reflect.Ptr {
                if reflectV.IsNil() {
                    field.ReflectValueOf(value).Set(reflect.New(field.FieldType).Elem())
                } else {
                    err = setter(value, reflectV.Elem().Interface())
                }
            } else if valuer, ok := v.(driver.Valuer); ok {
                if v, err = valuer.Value(); err == nil {
                    err = setter(value, v)
                }
            } else {
                return fmt.Errorf("failed to set value %+v to field %s", v, field.Name)
            }
        }

        return
    }

    // Set
    switch field.FieldType.Kind() {
    case reflect.Bool:
        field.Set = func(value reflect.Value, v any) (err error) {
            switch data := v.(type) {
            case bool:
                field.ReflectValueOf(value).SetBool(data)
            case *bool:
                if data != nil {
                    field.ReflectValueOf(value).SetBool(*data)
                } else {
                    field.ReflectValueOf(value).SetBool(false)
                }
            case int64:
                if data > 0 {
                    field.ReflectValueOf(value).SetBool(true)
                } else {
                    field.ReflectValueOf(value).SetBool(false)
                }
            case string:
                b, _ := strconv.ParseBool(data)
                field.ReflectValueOf(value).SetBool(b)
            default:
                return fallbackSetter(value, v, field.Set)
            }
            return nil
        }
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        field.Set = func(value reflect.Value, v any) (err error) {
            switch data := v.(type) {
            case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, time.Time, *time.Time:
                field.ReflectValueOf(value).SetInt(ToInt64(data))
            case []byte:
                return field.Set(value, string(data))
            case string:
                if i, err := strconv.ParseInt(data, 0, 64); err == nil {
                    field.ReflectValueOf(value).SetInt(i)
                } else {
                    return err
                }
            default:
                return fallbackSetter(value, v, field.Set)
            }
            return err
        }
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        field.Set = func(value reflect.Value, v any) (err error) {
            switch data := v.(type) {
            case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, time.Time:
                field.ReflectValueOf(value).SetUint(ToUint64(data))
            case []byte:
                return field.Set(value, string(data))
            case string:
                if i, err := strconv.ParseUint(data, 0, 64); err == nil {
                    field.ReflectValueOf(value).SetUint(i)
                } else {
                    return err
                }
            default:
                return fallbackSetter(value, v, field.Set)
            }
            return err
        }
    case reflect.Float32, reflect.Float64:
        field.Set = func(value reflect.Value, v any) (err error) {
            switch data := v.(type) {
            case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
                field.ReflectValueOf(value).SetFloat(ToFloat(data))
            case []byte:
                return field.Set(value, string(data))
            case string:
                if i, err := strconv.ParseFloat(data, 64); err == nil {
                    field.ReflectValueOf(value).SetFloat(i)
                } else {
                    return err
                }
            default:
                return fallbackSetter(value, v, field.Set)
            }
            return err
        }
    case reflect.String:
        field.Set = func(value reflect.Value, v any) (err error) {
            switch data := v.(type) {
            case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, string:
                field.ReflectValueOf(value).SetString(ToString(data))
            case float64, float32:
                field.ReflectValueOf(value).SetString(fmt.Sprintf("%."+ToString(field.Precision)+"f", data))
            default:
                return fallbackSetter(value, v, field.Set)
            }
            return err
        }
    default:
        fieldValue := reflect.New(field.FieldType)
        switch fieldValue.Elem().Interface().(type) {
        case time.Time:
            field.Set = func(value reflect.Value, v any) error {
                switch data := v.(type) {
                case time.Time:
                    field.ReflectValueOf(value).Set(reflect.ValueOf(v))
                case *time.Time:
                    if data != nil {
                        field.ReflectValueOf(value).Set(reflect.ValueOf(data).Elem())
                    } else {
                        field.ReflectValueOf(value).Set(reflect.ValueOf(time.Time{}))
                    }
                case string:
                    if t, err := now.Parse(data); err == nil {
                        field.ReflectValueOf(value).Set(reflect.ValueOf(t))
                    } else {
                        return fmt.Errorf("failed to set string %v to time.Time field %s, failed to parse it as time, got error %v", v, field.Name, err)
                    }
                default:
                    return fallbackSetter(value, v, field.Set)
                }
                return nil
            }
        case *time.Time:
            field.Set = func(value reflect.Value, v any) error {
                switch data := v.(type) {
                case time.Time:
                    fieldValue := field.ReflectValueOf(value)
                    if fieldValue.IsNil() {
                        fieldValue.Set(reflect.New(field.FieldType.Elem()))
                    }
                    fieldValue.Elem().Set(reflect.ValueOf(v))
                case *time.Time:
                    field.ReflectValueOf(value).Set(reflect.ValueOf(v))
                case string:
                    if t, err := now.Parse(data); err == nil {
                        fieldValue := field.ReflectValueOf(value)
                        if fieldValue.IsNil() {
                            if v == "" {
                                return nil
                            }
                            fieldValue.Set(reflect.New(field.FieldType.Elem()))
                        }
                        fieldValue.Elem().Set(reflect.ValueOf(t))
                    } else {
                        return fmt.Errorf("failed to set string %v to time.Time field %s, failed to parse it as time, got error %v", v, field.Name, err)
                    }
                default:
                    return fallbackSetter(value, v, field.Set)
                }
                return nil
            }
        default:
            if _, ok := fieldValue.Elem().Interface().(sql.Scanner); ok {
                // pointer scanner
                field.Set = func(value reflect.Value, v any) (err error) {
                    reflectV := reflect.ValueOf(v)
                    if !reflectV.IsValid() {
                        field.ReflectValueOf(value).Set(reflect.New(field.FieldType).Elem())
                    } else if reflectV.Type().AssignableTo(field.FieldType) {
                        field.ReflectValueOf(value).Set(reflectV)
                    } else if reflectV.Kind() == reflect.Ptr {
                        if reflectV.IsNil() || !reflectV.IsValid() {
                            field.ReflectValueOf(value).Set(reflect.New(field.FieldType).Elem())
                        } else {
                            return field.Set(value, reflectV.Elem().Interface())
                        }
                    } else {
                        fieldValue := field.ReflectValueOf(value)
                        if fieldValue.IsNil() {
                            fieldValue.Set(reflect.New(field.FieldType.Elem()))
                        }

                        if valuer, ok := v.(driver.Valuer); ok {
                            v, _ = valuer.Value()
                        }

                        err = fieldValue.Interface().(sql.Scanner).Scan(v)
                    }
                    return
                }
            } else if _, ok := fieldValue.Interface().(sql.Scanner); ok {
                // struct scanner
                field.Set = func(value reflect.Value, v any) (err error) {
                    reflectV := reflect.ValueOf(v)
                    if !reflectV.IsValid() {
                        field.ReflectValueOf(value).Set(reflect.New(field.FieldType).Elem())
                    } else if reflectV.Type().AssignableTo(field.FieldType) {
                        field.ReflectValueOf(value).Set(reflectV)
                    } else if reflectV.Kind() == reflect.Ptr {
                        if reflectV.IsNil() || !reflectV.IsValid() {
                            field.ReflectValueOf(value).Set(reflect.New(field.FieldType).Elem())
                        } else {
                            return field.Set(value, reflectV.Elem().Interface())
                        }
                    } else {
                        if valuer, ok := v.(driver.Valuer); ok {
                            v, _ = valuer.Value()
                        }

                        err = field.ReflectValueOf(value).Addr().Interface().(sql.Scanner).Scan(v)
                    }
                    return
                }
            } else {
                field.Set = func(value reflect.Value, v any) (err error) {
                    return fallbackSetter(value, v, field.Set)
                }
            }
        }
    }
}

// Field is ...
//type Field struct {
//Catalog  string
//Db       string
//Table    string
//OrgTable string
//Name     string
//OrgName  string
//DispLen  uint32
//Charset  uint16
//Flags    uint16
//Type     byte
//Scale    byte
//}

//type Field struct {
//Name                   string
//DBName                 string
//BindNames              []string
//DataType               DataType
//GORMDataType           DataType
//PrimaryKey             bool
//AutoIncrement          bool
//AutoIncrementIncrement int64
//Creatable              bool
//Updatable              bool
//Readable               bool
//HasDefaultValue        bool
////AutoCreateTime         TimeType
////AutoUpdateTime         TimeType
//DefaultValue           string
//DefaultValueInterface  any
//NotNull                bool
//Unique                 bool
//Comment                string
//Size                   int
//Precision              int
//Scale                  int
//FieldType              reflect.Type
//IndirectFieldType      reflect.Type
//StructField            reflect.StructField
//Tag                    reflect.StructTag
//TagSettings            map[string]string
//Schema                 *Schema
//EmbeddedSchema         *Schema
//OwnerSchema            *Schema
//ReflectValueOf         func(reflect.Value) reflect.Value
//ValueOf                func(reflect.Value) (value any, zero bool)
//Set                    func(reflect.Value, any) error
//IgnoreMigration        bool
//}
