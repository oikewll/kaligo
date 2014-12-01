package util

import (
    "time"
    "math/rand"
)

// 用法如下：
// sql = util.Substr(sql, 0, len(sql)-2)
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

func RandNum() int{
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    return r.Intn(100) 
}
