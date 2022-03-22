package util

import (
    "math/rand"
    "strconv"
    "time"
    "html"
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

// FilterInjections include SQL and XSS injections
func FilterInjections(str string) string {
    if str == "" {
        return str
    }

    str = FilterInjectionsWords(str)
    str = html.EscapeString(str)
    str = Addslashes(str)
    return str
}

// FilterInjectionsWords include SQL and XSS words
func FilterInjectionsWords(str string) string {
    // arr := [3]string{
    //     "/<(\\/?)(script|i?frame|style|html|body|title|link|meta|object|\\?|\\%)([^>]*?)>/isU",
    //     "/(<[^>]*)on[a-zA-Z]+\s*=([^>]*>)/isU",
    //     "/select|insert|update|delete|\'|\/\*|\*|\.\.\/|\.\/|union|into|load_file|outfile|dump/is"
    // }
    //
    // for _, value := range arr {
    //     re := regexp.MustCompile(value)
    //     re.ReplaceAllString(str, "")
    // }

    return str
}

// Addslashes 函数返回在预定义字符之前添加反斜杠的字符串
// 在防止被注入攻击时，常会用到两个函数：htmlspecialchars()和addslashes() 、trim() 函数
// 配合 html.EscapeString(hstr) 可以防止 XSS，XSS 实际上是往数据库添加 <script></script> 内容，利用当前用户cookie权限去调用接口做坏事
// 预定义字符是：
// 单引号（'）
// 双引号（"）
// 反斜杠（\）
func Addslashes(str string) string {
    tmpRune := []rune{}
    strRune := []rune(str)
    for _, ch := range strRune {
        switch ch {
        case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
            tmpRune = append(tmpRune, []rune{'\\'}[0])
            tmpRune = append(tmpRune, ch)
        default:
            tmpRune = append(tmpRune, ch)
        }
    }
    return string(tmpRune)
}

// Stripslashes 函数删除由 Addslashes 函数添加的反斜杠
func Stripslashes(str string) string {
    dstRune := []rune{}
    strRune := []rune(str)
    strLenth := len(strRune)
    for i := 0; i < strLenth; i++ {
        if strRune[i] == []rune{'\\'}[0] {
            i++
        }
        dstRune = append(dstRune, strRune[i])
    }
    return string(dstRune)
}
