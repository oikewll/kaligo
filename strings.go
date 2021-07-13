package kaligo

import(
    //"fmt"
    "strconv"
    //"reflect"
)

// StrToInt is the...
func StrToInt(str string) int{
    int, _ := strconv.Atoi(str)
    return int
}

// StrToInt64 is the...
func StrToInt64(str string) int64{
    int64, _ := strconv.ParseInt(str, 10, 64)
    return int64
}

// IntToStr is the...
func IntToStr(val int) string{
    str := strconv.Itoa(val)
    return str
}

// Int64ToStr is the...
func Int64ToStr(val int64) string{
    str := strconv.FormatInt(val, 10)
    return str
}

// StrToBool is the...
func StrToBool(val string) bool{
    str, _ := strconv.ParseBool(val)
    return str
}
