package mysql

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"time"
)

// Row is a type for result row. It contains values for any column of received row.
//
// If row is a result of ordinary text query, its element can be
// []byte slice, contained result text or nil if NULL is returned.
//
// If it is a result of prepared statement execution, its element field can be:
// intX, uintX, floatX, []byte, Date, Time, time.Time (in Local location) or nil.
type Row []interface{}

// Bin gets the nn-th value and returns it as []byte ([]byte{} if NULL).
func (tr Row) Bin(nn int) (bin []byte) {
	switch data := tr[nn].(type) {
	case nil:
		// bin = []byte{}
	case []byte:
		bin = data
	default:
		buf := new(bytes.Buffer)
		fmt.Fprint(buf, data)
		bin = buf.Bytes()
	}
	return
}

// Str gets the nn-th value and returns it as string ("" if NULL).
func (tr Row) Str(nn int) (str string) {
	switch data := tr[nn].(type) {
	case nil:
		// str = ""
	case []byte:
		str = string(data)
	case time.Time:
		str = TimeString(data)
	case time.Duration:
		str = DurationString(data)
	default:
		str = fmt.Sprint(data)
	}
	return
}

// MaxInt is...
const MaxInt = int64(int(^uint(0) >> 1))
// MinInt is...
const MinInt = -MaxInt - 1

// IntErr gets the nn-th value and returns it as int (0 if NULL). Returns error if
// conversion is impossible.
func (tr Row) IntErr(nn int) (val int, err error) {
	switch data := tr[nn].(type) {
	case nil:
		// nop
	case int32:
		val = int(data)
	case int16:
		val = int(data)
	case uint16:
		val = int(data)
	case int8:
		val = int(data)
	case uint8:
		val = int(data)
	case []byte:
		val, err = strconv.Atoi(string(data))
	case int64:
		if data >= MinInt && data <= MaxInt {
			val = int(data)
		} else {
			err = strconv.ErrRange
		}
	case uint32:
		if int64(data) <= MaxInt {
			val = int(data)
		} else {
			err = strconv.ErrRange
		}
	case uint64:
		if data <= uint64(MaxInt) {
			val = int(data)
		} else {
			err = strconv.ErrRange
		}
	default:
		err = os.ErrInvalid
	}
	return
}

// Int gets the nn-th value and returns it as int (0 if NULL). Panics if conversion is
// impossible.
func (tr Row) Int(nn int) (val int) {
	val, err := tr.IntErr(nn)
	if err != nil {
		panic(err)
	}
	return
}

// ForceInt gets the nn-th value and returns it as int. Returns 0 if value is NULL or
// conversion is impossible.
func (tr Row) ForceInt(nn int) (val int) {
	val, _ = tr.IntErr(nn)
	return
}

// MaxUint is...
const MaxUint = uint64(^uint(0))

// UintErr gets the nn-th value and return it as uint (0 if NULL). Returns error if
// conversion is impossible.
func (tr Row) UintErr(nn int) (val uint, err error) {
	switch data := tr[nn].(type) {
	case nil:
		// nop
	case uint32:
		val = uint(data)
	case uint16:
		val = uint(data)
	case uint8:
		val = uint(data)
	case []byte:
		var v uint64
		v, err = strconv.ParseUint(string(data), 0, 0)
		val = uint(v)
	case uint64:
		if data <= MaxUint {
			val = uint(data)
		} else {
			err = strconv.ErrRange
		}
	case int8, int16, int32, int64:
		v := reflect.ValueOf(data).Int()
		if v >= 0 && uint64(v) <= MaxUint {
			val = uint(v)
		} else {
			err = strconv.ErrRange
		}
	default:
		err = os.ErrInvalid
	}
	return
}

// Uint gets the nn-th value and returns it as uint (0 if NULL). Panics if conversion is
// impossible.
func (tr Row) Uint(nn int) (val uint) {
	val, err := tr.UintErr(nn)
	if err != nil {
		panic(err)
	}
	return
}

// ForceUint gets the nn-th value and returns it as uint. Returns 0 if value is NULL or
// conversion is impossible.
func (tr Row) ForceUint(nn int) (val uint) {
	val, _ = tr.UintErr(nn)
	return
}

// DateErr gets the nn-th value and returns it as Date (0000-00-00 if NULL). Returns error
// if conversion is impossible.
func (tr Row) DateErr(nn int) (val Date, err error) {
	switch data := tr[nn].(type) {
	case nil:
		// nop
	case Date:
		val = data
	case []byte:
		val, err = ParseDate(string(data))
	}
	return
}

// Date is like DateErr but panics if conversion is impossible.
func (tr Row) Date(nn int) (val Date) {
	val, err := tr.DateErr(nn)
	if err != nil {
		panic(err)
	}
	return
}

// ForceDate is like DateErr but returns 0000-00-00 if conversion is impossible.
func (tr Row) ForceDate(nn int) (val Date) {
	val, _ = tr.DateErr(nn)
	return
}

// TimeErr gets the nn-th value and returns it as time.Time in loc location (zero if NULL)
// Returns error if conversion is impossible. It can convert Date to time.Time.
func (tr Row) TimeErr(nn int, loc *time.Location) (t time.Time, err error) {
	switch data := tr[nn].(type) {
	case nil:
		// nop
	case time.Time:
		if loc == time.Local {
			t = data
		} else {
			y, mon, d := data.Date()
			h, m, s := data.Clock()
			t = time.Date(y, mon, d, h, m, s, t.Nanosecond(), loc)
		}
	case Date:
		t = data.Time(loc)
	case []byte:
		t, err = ParseTime(string(data), loc)
	}
	return
}

// Time is like TimeErr but panics if conversion is impossible.
func (tr Row) Time(nn int, loc *time.Location) (val time.Time) {
	val, err := tr.TimeErr(nn, loc)
	if err != nil {
		panic(err)
	}
	return
}

// ForceTime is like TimeErr but returns 0000-00-00 00:00:00 if conversion is
// impossible.
func (tr Row) ForceTime(nn int, loc *time.Location) (val time.Time) {
	val, _ = tr.TimeErr(nn, loc)
	return
}

// LocaltimeErr gets the nn-th value and returns it as time.Time in Local location
// (zero if NULL). Returns error if conversion is impossible.
// It can convert Date to time.Time.
func (tr Row) LocaltimeErr(nn int) (t time.Time, err error) {
	switch data := tr[nn].(type) {
	case nil:
		// nop
	case time.Time:
		t = data
	case Date:
		t = data.Time(time.Local)
	case []byte:
		t, err = ParseTime(string(data), time.Local)
	}
	return
}

// Localtime is like LocaltimeErr but panics if conversion is impossible.
func (tr Row) Localtime(nn int) (val time.Time) {
	val, err := tr.LocaltimeErr(nn)
	if err != nil {
		panic(err)
	}
	return
}

// ForceLocaltime is like LocaltimeErr but returns 0000-00-00 00:00:00 if conversion is
// impossible.
func (tr Row) ForceLocaltime(nn int) (val time.Time) {
	val, _ = tr.LocaltimeErr(nn)
	return
}

// DurationErr gets the nn-th value and returns it as time.Duration (0 if NULL). Returns error
// if conversion is impossible.
func (tr Row) DurationErr(nn int) (val time.Duration, err error) {
	switch data := tr[nn].(type) {
	case nil:
	case time.Duration:
		val = data
	case []byte:
		val, err = ParseDuration(string(data))
	default:
		err = fmt.Errorf("Can't convert `%v` to time.Duration", data)
	}
	return
}

// Duration is like DurationErr but panics if conversion is impossible.
func (tr Row) Duration(nn int) (val time.Duration) {
	val, err := tr.DurationErr(nn)
	if err != nil {
		panic(err)
	}
	return
}

// ForceDuration is like DurationErr but returns 0 if conversion is impossible.
func (tr Row) ForceDuration(nn int) (val time.Duration) {
	val, _ = tr.DurationErr(nn)
	return
}

// BoolErr gets the nn-th value and returns it as bool. Returns error
// if conversion is impossible.
func (tr Row) BoolErr(nn int) (val bool, err error) {
	switch data := tr[nn].(type) {
	case nil:
		// nop
	case int8:
		val = (data != 0)
	case int32:
		val = (data != 0)
	case int16:
		val = (data != 0)
	case int64:
		val = (data != 0)
	case uint8:
		val = (data != 0)
	case uint32:
		val = (data != 0)
	case uint16:
		val = (data != 0)
	case uint64:
		val = (data != 0)
	case []byte:
		var v int64
		v, err = strconv.ParseInt(string(data), 0, 64)
		val = (v != 0)
	default:
		err = os.ErrInvalid
	}
	return
}

// Bool is like BoolErr but panics if conversion is impossible.
func (tr Row) Bool(nn int) (val bool) {
	val, err := tr.BoolErr(nn)
	if err != nil {
		panic(err)
	}
	return
}

// ForceBool is like BoolErr but returns false if conversion is impossible.
func (tr Row) ForceBool(nn int) (val bool) {
	val, _ = tr.BoolErr(nn)
	return
}

// Int64Err gets the nn-th value and returns it as int64 (0 if NULL). Returns error if
// conversion is impossible.
func (tr Row) Int64Err(nn int) (val int64, err error) {
	switch data := tr[nn].(type) {
	case nil:
		// nop
	case int64, int32, int16, int8:
		val = reflect.ValueOf(data).Int()
	case uint64, uint32, uint16, uint8:
		u := reflect.ValueOf(data).Uint()
		if u > math.MaxInt64 {
			err = strconv.ErrRange
		} else {
			val = int64(u)
		}
	case []byte:
		val, err = strconv.ParseInt(string(data), 10, 64)
	default:
		err = os.ErrInvalid
	}
	return
}

// Int64 gets the nn-th value and returns it as int64 (0 if NULL).
// Panics if conversion is impossible.
func (tr Row) Int64(nn int) (val int64) {
	val, err := tr.Int64Err(nn)
	if err != nil {
		panic(err)
	}
	return
}

// ForceInt64 gets the nn-th value and returns it as int64. Returns 0 if value is NULL or
// conversion is impossible.
func (tr Row) ForceInt64(nn int) (val int64) {
	val, _ = tr.Int64Err(nn)
	return
}

// Uint64Err gets the nn-th value and returns it as uint64 (0 if NULL). Returns error if
// conversion is impossible.
func (tr Row) Uint64Err(nn int) (val uint64, err error) {
	switch data := tr[nn].(type) {
	case nil:
		// nop
	case uint64, uint32, uint16, uint8:
		val = reflect.ValueOf(data).Uint()
	case int64, int32, int16, int8:
		i := reflect.ValueOf(data).Int()
		if i < 0 {
			err = strconv.ErrRange
		} else {
			val = uint64(i)
		}
	case []byte:
		val, err = strconv.ParseUint(string(data), 10, 64)
	default:
		err = os.ErrInvalid
	}
	return
}

// Uint64 gets the nn-th value and returns it as uint64 (0 if NULL).
// Panic if conversion is impossible.
func (tr Row) Uint64(nn int) (val uint64) {
	val, err := tr.Uint64Err(nn)
	if err != nil {
		panic(err)
	}
	return
}

// ForceUint64 gets the nn-th value and returns it as uint64. Returns 0 if value is NULL or
// conversion is impossible.
func (tr Row) ForceUint64(nn int) (val uint64) {
	val, _ = tr.Uint64Err(nn)
	return
}

// FloatErr gets the nn-th value and returns it as float64 (0 if NULL). Returns error if
// conversion is impossible.
func (tr Row) FloatErr(nn int) (val float64, err error) {
	switch data := tr[nn].(type) {
	case nil:
		// nop
	case float64, float32:
		val = reflect.ValueOf(data).Float()
	case int64, int32, int16, int8:
		i := reflect.ValueOf(data).Int()
		if i >= 2<<53 || i <= -(2<<53) {
			err = strconv.ErrRange
		} else {
			val = float64(i)
		}
	case uint64, uint32, uint16, uint8:
		u := reflect.ValueOf(data).Uint()
		if u >= 2<<53 {
			err = strconv.ErrRange
		} else {
			val = float64(u)
		}
	case []byte:
		val, err = strconv.ParseFloat(string(data), 64)
	default:
		err = os.ErrInvalid
	}
	return
}

// Float gets the nn-th value and returns it as float64 (0 if NULL).
// Panics if conversion is impossible.
func (tr Row) Float(nn int) (val float64) {
	val, err := tr.FloatErr(nn)
	if err != nil {
		panic(err)
	}
	return
}

// ForceFloat gets the nn-th value and returns it as float64. Returns 0 if value is NULL or
// if conversion is impossible.
func (tr Row) ForceFloat(nn int) (val float64) {
	val, _ = tr.FloatErr(nn)
	return
}