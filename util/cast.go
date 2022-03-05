// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package util provides easy and safe casting in Go.
package util

import "time"

// To is use for all type
func To[T any](i any) (t T) {
	switch any(t).(type) {
	case bool:
		i = ToBool(i)
	case time.Time:
		i = ToTime(i)
	case time.Duration:
		i = ToDuration(i)
	case float64:
		i = ToFloat64(i)
	case float32:
		i = ToFloat32(i)
	case int64:
		i = ToInt64(i)
	case int32:
		i = ToInt32(i)
	case int16:
		i = ToInt16(i)
	case int8:
		i = ToInt8(i)
	case int:
		i = ToInt(i)
	case uint:
		i = ToUint(i)
	case uint64:
		i = ToUint64(i)
	case uint32:
		i = ToUint32(i)
	case uint16:
		i = ToUint16(i)
	case uint8:
		i = ToUint8(i)
	case string:
		i = ToString(i)
	case map[string]string:
		i = ToStringMapString(i)
	case map[string][]string:
		i = ToStringMapStringSlice(i)
	case map[string]bool:
		i = ToStringMapBool(i)
	case map[string]int:
		i = ToStringMapInt(i)
	case map[string]int64:
		i = ToStringMapInt64(i)
	case map[string]any:
		i = ToStringMap(i)
	case []bool:
		i = ToBoolSlice(i)
	case []string:
		i = ToStringSlice(i)
	case []int:
		i = ToIntSlice(i)
	case []time.Duration:
		i = ToDurationSlice(i)
	case []any:
		i = ToSlice(i)
	}
	return i.(T)
}

// ToBool casts an interface to a bool type.
func ToBool(i interface{}) bool {
	v, _ := ToBoolE(i)
	return v
}

// ToTime casts an interface to a time.Time type.
func ToTime(i interface{}) time.Time {
	v, _ := ToTimeE(i)
	return v
}

// ToTimeInDefaultLocation casts an interface to a *time.Location type.
func ToTimeInDefaultLocation(i interface{}, location *time.Location) time.Time {
	v, _ := ToTimeInDefaultLocationE(i, location)
	return v
}

// ToDuration casts an interface to a time.Duration type.
func ToDuration(i interface{}) time.Duration {
	v, _ := ToDurationE(i)
	return v
}

// ToFloat64 casts an interface to a float64 type.
func ToFloat64(i interface{}) float64 {
	v, _ := ToFloat64E(i)
	return v
}

// ToFloat32 casts an interface to a float32 type.
func ToFloat32(i interface{}) float32 {
	v, _ := ToFloat32E(i)
	return v
}

// ToInt64 casts an interface to an int64 type.
func ToInt64(i interface{}) int64 {
	v, _ := ToInt64E(i)
	return v
}

// ToInt32 casts an interface to an int32 type.
func ToInt32(i interface{}) int32 {
	v, _ := ToInt32E(i)
	return v
}

// ToInt16 casts an interface to an int16 type.
func ToInt16(i interface{}) int16 {
	v, _ := ToInt16E(i)
	return v
}

// ToInt8 casts an interface to an int8 type.
func ToInt8(i interface{}) int8 {
	v, _ := ToInt8E(i)
	return v
}

// ToInt casts an interface to an int type.
func ToInt(i interface{}) int {
	v, _ := ToIntE(i)
	return v
}

// ToUint casts an interface to a uint type.
func ToUint(i interface{}) uint {
	v, _ := ToUintE(i)
	return v
}

// ToUint64 casts an interface to a uint64 type.
func ToUint64(i interface{}) uint64 {
	v, _ := ToUint64E(i)
	return v
}

// ToUint32 casts an interface to a uint32 type.
func ToUint32(i interface{}) uint32 {
	v, _ := ToUint32E(i)
	return v
}

// ToUint16 casts an interface to a uint16 type.
func ToUint16(i interface{}) uint16 {
	v, _ := ToUint16E(i)
	return v
}

// ToUint8 casts an interface to a uint8 type.
func ToUint8(i interface{}) uint8 {
	v, _ := ToUint8E(i)
	return v
}

// ToString casts an interface to a string type.
func ToString(i interface{}) string {
	v, _ := ToStringE(i)
	return v
}

// ToStringMapString casts an interface to a map[string]string type.
func ToStringMapString(i interface{}) map[string]string {
	v, _ := ToStringMapStringE(i)
	return v
}

// ToStringMapStringSlice casts an interface to a map[string][]string type.
func ToStringMapStringSlice(i interface{}) map[string][]string {
	v, _ := ToStringMapStringSliceE(i)
	return v
}

// ToStringMapBool casts an interface to a map[string]bool type.
func ToStringMapBool(i interface{}) map[string]bool {
	v, _ := ToStringMapBoolE(i)
	return v
}

// ToStringMapInt casts an interface to a map[string]int type.
func ToStringMapInt(i interface{}) map[string]int {
	v, _ := ToStringMapIntE(i)
	return v
}

// ToStringMapInt64 casts an interface to a map[string]int64 type.
func ToStringMapInt64(i interface{}) map[string]int64 {
	v, _ := ToStringMapInt64E(i)
	return v
}

// ToStringMap casts an interface to a map[string]interface{} type.
func ToStringMap(i interface{}) map[string]interface{} {
	v, _ := ToStringMapE(i)
	return v
}

// ToSlice casts an interface to a []interface{} type.
func ToSlice(i interface{}) []interface{} {
	v, _ := ToSliceE(i)
	return v
}

// ToBoolSlice casts an interface to a []bool type.
func ToBoolSlice(i interface{}) []bool {
	v, _ := ToBoolSliceE(i)
	return v
}

// ToStringSlice casts an interface to a []string type.
func ToStringSlice(i interface{}) []string {
	v, _ := ToStringSliceE(i)
	return v
}

// ToIntSlice casts an interface to a []int type.
func ToIntSlice(i interface{}) []int {
	v, _ := ToIntSliceE(i)
	return v
}

// ToDurationSlice casts an interface to a []time.Duration type.
func ToDurationSlice(i interface{}) []time.Duration {
	v, _ := ToDurationSliceE(i)
	return v
}
