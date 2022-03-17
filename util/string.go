package util

import (
    "math/rand"
    "strconv"
    "time"
)

// Substr 返回一个字符串中从指定位置开始到指定字符数的字符
// ep: 
//   sql = util.Substr(sql, 0, len(sql)-2)
func Substr(str string, start, length int) string {
    rs := []rune(str)
    rl := len(rs)
    end := 0

    if start < 0 {
        start = rl - 1 + start
    }
    end = start + length

    if start > end {
        start, end = end, start
    }

    if start < 0 {
        start = 0
    }
    if start > rl {
        start = rl
    }
    if end < 0 {
        end = 0
    }
    if end > rl {
        end = rl
    }

    return string(rs[start:end])
}

// RandNum 生成随机数
func RandNum() int{
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    return r.Intn(100) 
}

// RandomStr 生成随机字符串
func RandomStr(length int) string {
    str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    bytes := []byte(str)
    var result []byte
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < length; i++ {
        result = append(result, bytes[r.Intn(len(bytes))])
    }
    return string(result)
}

// StrToInt is the...
func StrToInt(str string) int{
    var int, _ = strconv.Atoi(str)
    return int
}

// StrToInt64 is the...
func StrToInt64(str string) int64{
    var int64, _ = strconv.ParseInt(str, 10, 64)
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
