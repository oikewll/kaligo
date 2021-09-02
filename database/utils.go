package database

import (
    "encoding/json"
    "fmt"
    "io"
    "runtime"
    "reflect"
    "regexp"
    "strings"
    "strconv"
    "time"
)

var kaliSourceDir string

func init() {
    _, file, _, _ := runtime.Caller(0)
    // compatible solution to get gorm source directory with various operating systems
    kaliSourceDir = regexp.MustCompile(`utils.utils\.go`).ReplaceAllString(file, "")
}

// FileWithLineNum return the file name and line number of the current file
func FileWithLineNum() string {
    // the second caller usually from gorm internal, so set i start from 2
    for i := 2; i < 15; i++ {
        _, file, line, ok := runtime.Caller(i)
        if ok && (!strings.HasPrefix(file, kaliSourceDir) || strings.HasSuffix(file, "_test.go")) {
            return file + ":" + strconv.FormatInt(int64(line), 10)
        }
    }

    return ""
}

// Version returns version string
func Version() string {
    return "1.5.3"
}

func syntaxError(ln int) error {
    return fmt.Errorf("syntax error at line: %d", ln)
}

// CatchError is...
func catchError(err *error) {
    if pv := recover(); pv != nil {
        switch e := pv.(type) {
        case runtime.Error:
            panic(pv)
        case error:
            if e == io.EOF {
                *err = io.ErrUnexpectedEOF
            } else {
                *err = e
            }
        default:
            panic(pv)
        }
    }
}

// FormatJSON is ...
func FormatJSON(jsonObj interface{}) string {
    jsonStr, _ := json.MarshalIndent(jsonObj, "", "    ")
    return string(jsonStr)
}

// ParseType is Extracts the text between parentheses, if any.
//     Returns: []interface{}{"CHAR", 6}
//     columnType, columnSize := parseType("CHAR(6)");
func ParseType(columnType string) (retType DataType, retSize int64) {
    var strOpen int    
    if strOpen = strings.Index(columnType, "("); strOpen != -1 {
        // Closing parenthesis
        strClose := strings.Index(columnType, ")")

        // Length without parentheses
        strLength := columnType[strOpen + 1 : strClose]

        // Type without the length
        columnType = columnType[0:strOpen]
        retSize = ToInt64(strLength)
    } else {
        // No length specified
        retSize = 0
    }

    columnType = GetDataType(columnType)
    retType    = DataType(columnType)
    return
}

// GetDataType is ...
func GetDataType(value string) (dataType string) {
    switch strings.ToLower(value) {
    case "boolean":
        dataType = "bool"
    case "int", "integer", "smallint", "bigint":
        dataType = "int"
    case "dec", "decimal", "double precision", "float", "numeric", "real":
        dataType = "float"
    case "date", "time", "time with time zone", "timestamp", "timestamp with time zone":
        dataType = "time"
    case "bit", "bit varying", "char", "char varying", "character", "character varying", "interval", "national char", "national char varying",
        "national character", "national character varying", "nchar", "nchar varying", "varchar", "char large object", "character large object",
        "clob", "national character large object", "nchar large object", "nclob":
        dataType = "string"
    case "binary large object", "blob", "binary", "binary varying", "varbinary":
        dataType = "bytes"
    default:
        dataType = "bytes"
    }
    return
}

var commonInitialisms = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}

// ToSchemaName 转换字段名称，下划线写法转驼峰写法
func ToSchemaName(name string) string {
    result := strings.Replace(strings.Title(strings.Replace(name, "_", " ", -1)), " ", "", -1)
    for _, initialism := range commonInitialisms {
        result = regexp.MustCompile(strings.Title(strings.ToLower(initialism))+"([A-Z]|$|_)").ReplaceAllString(result, initialism+"$1")
    }
    return result
}

// ToDBName 转换字段名称，驼峰写法转下划线写法
func ToDBName(name string) string {
    if name == "" {
        return ""
    }

    var (
        value                          = name
        buf                            strings.Builder
        lastCase, nextCase, nextNumber bool // upper case == true
        curCase                        = value[0] <= 'Z' && value[0] >= 'A'
    )

    for i, v := range value[:len(value)-1] {
        nextCase   = value[i+1] <= 'Z' && value[i+1] >= 'A'
        nextNumber = value[i+1] >= '0' && value[i+1] <= '9'

        if curCase {
            if lastCase && (nextCase || nextNumber) {
                buf.WriteRune(v + 32)
            } else {
                if i > 0 && value[i-1] != '_' && value[i+1] != '_' {
                    buf.WriteByte('_')
                }
                buf.WriteRune(v + 32)
            }
        } else {
            buf.WriteRune(v)
        }

        lastCase = curCase
        curCase  = nextCase
    }

    if curCase {
        if !lastCase && len(value) > 1 {
            buf.WriteByte('_')
        }
        buf.WriteByte(value[len(value)-1] + 32)
    } else {
        buf.WriteByte(value[len(value)-1])
    }
    ret := buf.String()
    return ret
}

// ParseTagSetting is 解析标签
func ParseTagSetting(str string, sep string) map[string]string {
    settings := map[string]string{}
    names := strings.Split(str, sep)

    for i := 0; i < len(names); i++ {
        j := i
        if len(names[j]) > 0 {
            for {
                if names[j][len(names[j])-1] == '\\' {
                    i++
                    names[j] = names[j][0:len(names[j])-1] + sep + names[i]
                    names[i] = ""
                } else {
                    break
                }
            }
        }

        values := strings.Split(names[j], ":")
        k := strings.TrimSpace(strings.ToUpper(values[0]))

        if len(values) >= 2 {
            settings[k] = strings.Join(values[1:], ":")
        } else if k != "" {
            settings[k] = k
        }
    }

    return settings
}

// Strtr strtr()
// If the parameter length is 1, type is: map[string]string
// Strtr("baab", map[string]string{"ab": "01"}) will return "ba01"
// If the parameter length is 2, type is: string, string
// Strtr("baab", "ab", "01") will return "1001", a => 0; b => 1.
func Strtr(haystack string, params ...interface{}) string {
    ac := len(params)
    if ac == 1 {
        pairs := params[0].(map[string]string)
        length := len(pairs)
        if length == 0 {
            return haystack
        }
        oldnew := make([]string, length*2)
        for o, n := range pairs {
            if o == "" {
                return haystack
            }
            oldnew = append(oldnew, o, n)
        }
        return strings.NewReplacer(oldnew...).Replace(haystack)
    } else if ac == 2 {
        from := params[0].(string)
        to := params[1].(string)
        trlen, lt := len(from), len(to)
        if trlen > lt {
            trlen = lt
        }
        if trlen == 0 {
            return haystack
        }

        str := make([]uint8, len(haystack))
        var xlat [256]uint8
        var i int
        var j uint8
        if trlen == 1 {
            for i = 0; i < len(haystack); i++ {
                if haystack[i] == from[0] {
                    str[i] = to[0]
                } else {
                    str[i] = haystack[i]
                }
            }
            return string(str)
        }
        // trlen != 1
        for {
            xlat[j] = j
            if j++; j == 0 {
                break
            }
        }
        for i = 0; i < trlen; i++ {
            xlat[from[i]] = to[i]
        }
        for i = 0; i < len(haystack); i++ {
            str[i] = xlat[haystack[i]]
        }
        return string(str)
    }

    return haystack
}

// CheckTruth is 检查是否为空
func CheckTruth(val interface{}) bool {
    if v, ok := val.(bool); ok {
        return v
    }

    if v, ok := val.(string); ok {
        v = strings.ToLower(v)
        return v != "false"
    }

    return !reflect.ValueOf(val).IsZero()
}

// ToSlice is interface{} to slice interface{}
func ToSlice(value interface{}) []interface{} {
    reflectValue := reflect.ValueOf(value)
    if reflectValue.Kind() == reflect.Ptr {
        reflectValue = reflectValue.Elem()
    }

    if reflectValue.Kind() != reflect.Slice {
        panic("ToSlice value not slice")
    }

    len := reflectValue.Len()
    ret := make([]interface{}, len)
    for i := 0; i < len; i++ {
        ret[i] = reflectValue.Index(i).Interface()
    }

    return ret
}

// ToInt64 is int to string
func ToInt64(value interface{}) int64 {
    switch v := value.(type) {
    case int:
        return int64(v)
    case int8:
        return int64(v)
    case int16:
        return int64(v)
    case int32:
        return int64(v)
    case int64:
        return v
    case uint:
        return int64(v)
    case uint8:
        return int64(v)
    case uint16:
        return int64(v)
    case uint32:
        return int64(v)
    case uint64:
        return int64(v)
    case float32:
        return int64(v)
    case float64:
        return int64(v)
    case time.Time:
        return v.Unix()
    case *time.Time:
        if v != nil {
            return v.Unix()
        } 
    case string:
        if i, err := strconv.Atoi(v); err == nil {
            return int64(i)
        }
    }
    return 0
}

// ToUint64 is int to string
func ToUint64(value interface{}) uint64 {
    switch v := value.(type) {
    case int:
        return uint64(v)
    case int8:
        return uint64(v)
    case int16:
        return uint64(v)
    case int32:
        return uint64(v)
    case int64:
        return uint64(v)
    case uint:
        return uint64(v)
    case uint8:
        return uint64(v)
    case uint16:
        return uint64(v)
    case uint32:
        return uint64(v)
    case uint64:
        return v
    case float32:
        return uint64(v)
    case float64:
        return uint64(v)
    case time.Time:
        return uint64(v.Unix())
    case string:
        if i, err := strconv.Atoi(v); err == nil {
            return uint64(i)
        }
    }
    return 0
}

// ToFloat is int to string
func ToFloat(value interface{}) float64 {
    switch v := value.(type) {
    case int:
        return float64(v)
    case int8:
        return float64(v)
    case int16:
        return float64(v)
    case int32:
        return float64(v)
    case int64:
        return float64(v)
    case uint:
        return float64(v)
    case uint8:
        return float64(v)
    case uint16:
        return float64(v)
    case uint32:
        return float64(v)
    case uint64:
        return float64(v)
    case float32:
        return float64(v)
    case float64:
        return v
    case time.Time:
        return float64(v.Unix())
    case string:
        if float, err := strconv.ParseFloat(v, 64); err == nil {
            return float
        }
    }
    return 0
}

// ToString is int to string
func ToString(value interface{}) string {
    switch v := value.(type) {
    case string:
        return v
    case int:
        return strconv.FormatInt(int64(v), 10)
    case int8:
        return strconv.FormatInt(int64(v), 10)
    case int16:
        return strconv.FormatInt(int64(v), 10)
    case int32:
        return strconv.FormatInt(int64(v), 10)
    case int64:
        return strconv.FormatInt(v, 10)
    case uint:
        return strconv.FormatUint(uint64(v), 10)
    case uint8:
        return strconv.FormatUint(uint64(v), 10)
    case uint16:
        return strconv.FormatUint(uint64(v), 10)
    case uint32:
        return strconv.FormatUint(uint64(v), 10)
    case uint64:
        return strconv.FormatUint(v, 10)
    case float32:
        return strconv.FormatFloat(float64(v), 'E', -1, 32)
    case float64:
        return strconv.FormatFloat(v, 'E', -1, 64)
    }
    return ""
}

// InSlice like php in_array()
func InSlice(a string, list *[]string) bool {
    if list == nil {
        return false
    }
    for _, b := range *list {
        if b == a {
            return true
        }
    }
    return false
}

// IsNumeric like php isNumeric()
func IsNumeric(s string) bool {
    _, err := strconv.ParseFloat(s, 64)
    return err == nil
}

// StructToMap 将一个结构体所有字段(包括通过组合得来的字段)到一个map中
// value:结构体的反射值
// data:存储字段数据的map
func StructToMap(value reflect.Value, data map[string]interface{}) {
    if value.Kind() != reflect.Struct {
        return
    }

    for i := 0; i < value.NumField(); i++ {
        var fieldValue = value.Field(i)
        if fieldValue.CanInterface() {
            var fieldType = value.Type().Field(i)
            if fieldType.Anonymous {
                // 匿名组合字段,进行递归解析
                StructToMap(fieldValue, data)
            } else {
                // 非匿名字段
                var fieldName = fieldType.Tag.Get("db")
                if fieldName == "-" {
                    continue
                }
                if fieldName == "" {
                    fieldName = ToDBName(fieldType.Name)
                }
                data[fieldName] = fieldValue.Interface()
                //t.Log(fieldName + ":" + fieldValue.Interface().(string))
            }
        }
    }
}

// MapChangeKeyCase is change key case
func MapChangeKeyCase(values map[string]interface{}, caseUpper bool) map[string]interface{}{
    valueMaps := make(map[string]interface{}, len(values))
    for k, v := range values {
        if caseUpper {
            valueMaps[strings.ToUpper(k)] = v
        } else {
            valueMaps[strings.ToLower(k)] = v
        }
    }
    return valueMaps
}

//// AddSlashes is ...
//// 转义：引号、双引号添加反斜杠
//func (db *DB) AddSlashes(val string) string {
//val = strings.Replace(val, "\"", "\\\"", -1)
//val = strings.Replace(val, "'", "\\'", -1)
//return val
//}

//// StripSlashes is ...
//// 反转义：引号、双引号去除反斜杠
//func (db *DB) StripSlashes(val string) string {
//val = strings.Replace(val, "\\\"", "\"", -1)
//val = strings.Replace(val, "\\'", "'", -1)
//return val
//}

//// GetSafeParam is ...
//// 防止XSS跨站攻击
//func (db *DB) GetSafeParam(val string) string {
//val = strings.Replace(val, "&", "&amp;", -1)
//val = strings.Replace(val, "<", "&lt;", -1)
//val = strings.Replace(val, ">", "&gt;", -1)
//val = strings.Replace(val, "\"", "&quot;", -1)
//val = strings.Replace(val, "'", "&#039;", -1)
//return val
//}
