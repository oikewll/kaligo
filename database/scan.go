package database

import (
    "database/sql"
    "database/sql/driver"
    "reflect"
    "time"
)

// prepareValues 准备可以被rows.Scan进去数据 的 数据类型指针
func prepareValues(values []any, q *Query, columnTypes []*sql.ColumnType, columns []string) {
    if q.Schema != nil {
        for idx, name := range columns {
            if field := q.Schema.LookUpField(name); field != nil {
                values[idx] = reflect.New(reflect.PtrTo(field.FieldType)).Interface()
                continue
            }
            values[idx] = new(any)
        }
    } else if len(columnTypes) > 0 {
        for idx, columnType := range columnTypes {
            if columnType.ScanType() != nil {
                values[idx] = reflect.New(reflect.PtrTo(columnType.ScanType())).Interface()
            } else {
                values[idx] = new(any)
            }
        }
    } else {
        for idx := range columns {
            values[idx] = new(any)
        }
    }
}

// scanIntoMap 扫描进 Map
func scanIntoMap(mapValue map[string]any, values []any, columns []string) {
    for idx, column := range columns {
        if reflectValue := reflect.Indirect(reflect.Indirect(reflect.ValueOf(values[idx]))); reflectValue.IsValid() {
            mapValue[column] = reflectValue.Interface()
            if valuer, ok := mapValue[column].(driver.Valuer); ok {
                mapValue[column], _ = valuer.Value()
            } else if b, ok := mapValue[column].(sql.RawBytes); ok {
                mapValue[column] = string(b)
            }
        } else {
            mapValue[column] = nil
        }
    }
}

// Scan 扫描数据行
// 将数据行中的数据解析到 q.Dest 中, q.Dest 可以是 基础类型、time.Time类型、结构体、数组类型的指针
// 扫描的行数: q.RowsAffected
func Scan(rows *sql.Rows, q *Query) {
    columns, _ := rows.Columns()
    values := make([]any, len(columns))
    q.RowsAffected = 0

    switch dest := q.Dest.(type) {
    case map[string]any, *map[string]any:
        if rows.Next() {
            // 如果字段类型不为空
            columnTypes, _ := rows.ColumnTypes()
            prepareValues(values, q, columnTypes, columns)

            q.RowsAffected++
            q.AddError(rows.Scan(values...))

            mapValue, ok := dest.(map[string]any)
            if !ok {
                if v, ok := dest.(*map[string]any); ok {
                    mapValue = *v
                }
            }
            scanIntoMap(mapValue, values, columns)
        }
    case *[]map[string]any:
        columnTypes, _ := rows.ColumnTypes()
        for rows.Next() {
            prepareValues(values, q, columnTypes, columns)

            q.RowsAffected++
            q.AddError(rows.Scan(values...))

            mapValue := map[string]any{}
            scanIntoMap(mapValue, values, columns)
            *dest = append(*dest, mapValue)
        }
    case *int, *int8, *int16, *int32, *int64,
        *uint, *uint8, *uint16, *uint32, *uint64, *uintptr,
        *float32, *float64,
        *bool, *string, *time.Time,
        *sql.NullInt32, *sql.NullInt64, *sql.NullFloat64,
        *sql.NullBool, *sql.NullString, *sql.NullTime:
        for rows.Next() {
            q.RowsAffected++
            q.AddError(rows.Scan(dest))
        }
    default:
        Schema := q.Schema
        switch q.ReflectValue.Kind() {
        case reflect.Slice, reflect.Array:
            var (
                reflectValueType = q.ReflectValue.Type().Elem()
                isPtr            = reflectValueType.Kind() == reflect.Ptr
                fields           = make([]*Field, len(columns))
            )

            if isPtr {
                reflectValueType = reflectValueType.Elem()
            }

            q.ReflectValue.Set(reflect.MakeSlice(q.ReflectValue.Type(), 0, 20))

            //if q.Model != nil && reflectValueType != Schema.ModelType && reflectValueType.Kind() == reflect.Struct {
            //Schema, _ = Parse(q.Model, q.cacheStore)
            //}

            for idx, column := range columns {
                //fmt.Printf("idx = [%v] ---> column = [%v]\n", idx, column)
                // query.Execute() 方法里面已经执行了 Parse()
                if field := Schema.LookUpField(column); field != nil {
                    //fmt.Printf("column = [%v] ---> field = [%v]\n", column, field)
                    fields[idx] = field
                } else {
                    // var ages []int64 会跑这里来
                    values[idx] = &sql.RawBytes{}
                }
            }

            // pluck values into slice of data
            isPluck := false
            if len(fields) == 1 {
                if _, ok := reflect.New(reflectValueType).Interface().(sql.Scanner); ok || // is scanner
                    reflectValueType.Kind() != reflect.Struct || // is not struct
                    Schema.ModelType.ConvertibleTo(TimeReflectType) { // is time
                    isPluck = true
                }
            }

            for rows.Next() {
                q.RowsAffected++
                elem := reflect.New(reflectValueType)
                if isPluck {
                    q.AddError(rows.Scan(elem.Interface()))
                } else {
                    // 准备 values
                    for idx, field := range fields {
                        if field != nil {
                            values[idx] = reflect.New(reflect.PtrTo(field.IndirectFieldType)).Interface()
                        }
                    }

                    q.AddError(rows.Scan(values...))

                    // 赋值
                    for idx, field := range fields {
                        field.Set(elem, values[idx])
                    }
                }

                if isPtr {
                    q.ReflectValue.Set(reflect.Append(q.ReflectValue, elem))
                } else {
                    q.ReflectValue.Set(reflect.Append(q.ReflectValue, elem.Elem()))
                }
            }

        case reflect.Struct, reflect.Ptr:
            // 这里应该不会进入，因为 Execute() 里面执行了 Parse() 才到这里来的
            //if q.ReflectValue.Type() != Schema.ModelType && q.ReflectValue.Type().Kind() == reflect.Struct {
            //var err error
            //if q.Schema, err = Parse(q.Model, q.cacheStore); err != nil {
            //q.AddError(err)
            //}
            //}
            if rows.Next() {
                for idx, column := range columns {
                    if field := Schema.LookUpField(column); field != nil {
                        values[idx] = reflect.New(reflect.PtrTo(field.IndirectFieldType)).Interface()
                    } else if len(columns) == 1 {
                        values[idx] = dest
                    } else {
                        values[idx] = &sql.RawBytes{}
                    }
                }

                q.RowsAffected++
                // 给 values 填充数据
                q.AddError(rows.Scan(values...))
                for idx, column := range columns {
                    if field := Schema.LookUpField(column); field != nil {
                        field.Set(q.ReflectValue, values[idx])
                    }
                }
            }
        default:
            q.AddError(rows.Scan(dest))
        }
    }

    if err := rows.Err(); err != nil && err != q.Error {
        q.AddError(err)
    }

    if q.RowsAffected == 0 && q.Error == nil {
        q.AddError(ErrRecordNotFound)
    }
}
