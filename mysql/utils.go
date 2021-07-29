package mysql

import (
    "bytes"
    "fmt"
    "io"
	"runtime"
	"regexp"
	"strings"
    "strconv"
)

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

var reg = regexp.MustCompile(`\B[A-Z]`)

// transFieldName 转换字段名称，驼峰写法转下划线写法
func transFieldName(name string) string {
	return strings.ToLower(reg.ReplaceAllString(name, "_$0"))
}

// Strtr strtr()
//
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

// 转义可能导致 SQL 注入攻击的字符
func escapeString(txt string) string {
    var (
        esc string
        buf bytes.Buffer
    )
    last := 0
    for ii, bb := range txt {
        switch bb {
        case 0:
            esc = `\0`
        case '\n':
            esc = `\n`
        case '\r':
            esc = `\r`
        case '\\':
            esc = `\\`
        case '\'':
            esc = `\'`
        case '"':
            esc = `\"`
        case '\032':
            esc = `\Z`
        default:
            continue
        }
        io.WriteString(&buf, txt[last:ii])
        io.WriteString(&buf, esc)
        last = ii + 1
    }
    io.WriteString(&buf, txt[last:])
    return buf.String()
}

// 转义可能导致 SQL 注入攻击的 引用字符
func escapeQuotes(txt string) string {
    var buf bytes.Buffer
    last := 0
    for ii, bb := range txt {
        if bb == '\'' {
            io.WriteString(&buf, txt[last:ii])
            io.WriteString(&buf, `''`)
            last = ii + 1
        }
    }
    io.WriteString(&buf, txt[last:])
    return buf.String()
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
