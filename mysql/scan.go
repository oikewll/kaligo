/**
 * Realize a database operation class
 *
 * @copyright   (C) 2014  seatle
 * @lastmodify  2021-07-06
 *
 */

package mysql

import (
    "database/sql"
    "database/sql/driver"
    //"fmt"
    "reflect"
    //"strconv"
    "strings"
    "time"
)

func prepareValues(values []interface{}, db *DB, columnTypes []*sql.ColumnType, columns []string) {
	if db.query.Schema != nil {
		for idx, name := range columns {
			if field := db.query.Schema.LookUpField(name); field != nil {
				values[idx] = reflect.New(reflect.PtrTo(field.FieldType)).Interface()
				continue
			}
			values[idx] = new(interface{})
		}
	} else if len(columnTypes) > 0 {
		for idx, columnType := range columnTypes {
			if columnType.ScanType() != nil {
				values[idx] = reflect.New(reflect.PtrTo(columnType.ScanType())).Interface()
			} else {
				values[idx] = new(interface{})
			}
		}
	} else {
		for idx := range columns {
			values[idx] = new(interface{})
		}
	}
}

func scanIntoMap(mapValue map[string]interface{}, values []interface{}, columns []string) {
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

// Scan is...
func Scan(rows *sql.Rows, db *DB) {

    columns, _     := rows.Columns()

    //for rows.Next() {
        //var name string
        //var age int
        //err := rows.Scan(&name, &age)
        //if err == nil {
            //fmt.Printf("name = %v\nage  = %v\n", name, age)
        //}
    //}

    //fmt.Printf("Dest 44444 = %T = %v\n", columns, columns)
    //fmt.Printf("Dest 44444 = %T = %v\n", columnTypes, columnTypes)

    values := make([]interface{}, len(columns))
    db.RowsAffected = 0

    switch dest := db.query.Dest.(type) {
    case map[string]interface{}, *map[string]interface{}:
        if rows.Next() {
            // 如果字段类型不为空
            columnTypes, _ := rows.ColumnTypes()
            prepareValues(values, db, columnTypes, columns)

            db.RowsAffected++
            db.AddError(rows.Scan(values...))

            mapValue, ok := dest.(map[string]interface{})
            if !ok {
                if v, ok := dest.(*map[string]interface{}); ok {
                    mapValue = *v
                }
            }
            scanIntoMap(mapValue, values, columns)
        }
    case *[]map[string]interface{}:
        columnTypes, _ := rows.ColumnTypes()
        for rows.Next() {
            prepareValues(values, db, columnTypes, columns)

            db.RowsAffected++
            db.AddError(rows.Scan(values...))

            mapValue := map[string]interface{}{}
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
            db.RowsAffected++
            db.AddError(rows.Scan(dest))
        }
    default:
        Schema := db.query.Schema
        //fmt.Printf("Schema = %v\n", Schema)
        switch db.query.ReflectValue.Kind() {
        case reflect.Slice, reflect.Array:

        case reflect.Struct, reflect.Ptr:
            //fmt.Printf("11111 %v = %v\n", db.query.ReflectValue.Type(), Schema.ModelType)
            if db.query.ReflectValue.Type() != Schema.ModelType {
                Schema, _ = Parse(db.query.Dest, db.cacheStore)
            }
            if rows.Next() {
                for idx, column := range columns {
                    // 从 db.query.Schema 里面去查找，因为在 query.Execute() 方法里面已经执行了 Parse()，或者上面也会执行，所以这里肯定是有的
                    if field := Schema.LookUpField(column); field != nil {
                        values[idx] = reflect.New(reflect.PtrTo(field.IndirectFieldType)).Interface()
                        //fmt.Printf("22222 %v = %v = %T\n", idx, column, values[idx])
                    } else if names := strings.Split(column, "__"); len(names) > 1 {
                        //if rel, ok := Schema.Relationships.Relations[names[0]]; ok {
                            //if field := rel.FieldSchema.LookUpField(strings.Join(names[1:], "__")); field != nil && field.Readable {
                                //values[idx] = reflect.New(reflect.PtrTo(field.IndirectFieldType)).Interface()
                                //continue
                            //}
                        //}
                        values[idx] = &sql.RawBytes{}
                    } else if len(columns) == 1 {
                        values[idx] = dest
                    } else {
                        values[idx] = &sql.RawBytes{}
                    }
                }

                db.RowsAffected++
                db.AddError(rows.Scan(values...))
                for idx, column := range columns {
                    if field := Schema.LookUpField(column); field != nil {
                        field.Set(db.query.ReflectValue, values[idx])
                    } else if names := strings.Split(column, "__"); len(names) > 1 {
                        //if rel, ok := Schema.query.Relations[names[0]]; ok {
                            //if field := rel.FieldSchema.LookUpField(strings.Join(names[1:], "__")); field != nil && field.Readable {
                                //relValue := rel.Field.ReflectValueOf(db.query.ReflectValue)
                                //value := reflect.ValueOf(values[idx]).Elem()

                                //if relValue.Kind() == reflect.Ptr && relValue.IsNil() {
                                    //if value.IsNil() {
                                        //continue
                                    //}
                                    //relValue.Set(reflect.New(relValue.Type().Elem()))
                                //}

                                //field.Set(relValue, values[idx])
                            //}
                        //}
                    }
                }
            }
        default:
            db.AddError(rows.Scan(dest))
        }
    }

    if err := rows.Err(); err != nil && err != db.Error {
        db.AddError(err)
    }

    if db.RowsAffected == 0 && db.Error == nil {
        db.AddError(ErrRecordNotFound)
    }
}
