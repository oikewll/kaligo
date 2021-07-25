package mysql

import (
    "bytes"
    "fmt"
    "io"
	"runtime"
	"regexp"
	"strings"
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

// 字段值处理
func quote(values interface{}) string {
    switch vals := values.(type) {
    case string:
        return Escape(vals)
    case []string:
        for k, v := range vals {
            vals[k] = Escape(v)
        }
        return "(" + strings.Join(vals, ", ") + ")"
    case *Query:
        // Create a sub-query
        return "(" + vals.Compile() + ")"
    default:
        return vals.(string)
    }
}

// 表名添加引用符号(`)
// 添加表前缀
func quoteTable(table string) string {
    //table = strings.Replace(table, "#DB#", "lrs", 1 )

    // 表名前缀
    tablePrefix := ""
    table = tablePrefix + quoteIdentifier(table)

    return table
}

// 字段名添加引用符号(`)
func quoteIdentifier(values interface{}) string {
    switch vals := values.(type) {
    case string:
        if vals == "*" || strings.Index(vals, "`") != -1 {
            // * 不需要变成 `*`，已经有 `` 包含着的直接返回
            return vals
        } else if strings.Index(vals, ".") != -1 {
            // table.column 的写法，变成 `table`.`column`
            parts := regexp.MustCompile(`\.`).Split(vals, 2)
            return quoteIdentifier(parts[0]) + "." + quoteIdentifier(parts[1])
        } else {
            return "`" + vals + "`"
        }
    case []string:
        // Separate the column and alias
        value := vals[0]
        alias := vals[1]
        return quoteIdentifier(value) + " AS " + quoteIdentifier(alias)
    default:
        return vals.(string)
    }
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

// MysqlRealEscapeString is ...
//func MysqlRealEscapeString(value string) string {
    //var sb strings.Builder
    //for i := 0; i < len(value); i++ {
        //c := value[i]
        //switch c {
        //case '\\', 0, '\n', '\r', '\'', '"':
            //sb.WriteByte('\\')
            //sb.WriteByte(c)
        //case '\032':
            //sb.WriteByte('\\')
            //sb.WriteByte('Z')
        //default:
            //sb.WriteByte(c)
        //}
    //}
    //return sb.String()
//}

// Escape is use for Escapes special characters in the txt, so it is safe to place returned string
func Escape(sql string) string {
    dest := make([]byte, 0, 2*len(sql))
    var escape byte
    for i := 0; i < len(sql); i++ {
        c := sql[i]

        escape = 0

        switch c {
        case 0: /* Must be escaped for 'mysql' */
            escape = '0'
            break
        case '\n': /* Must be escaped for logs */
            escape = 'n'
            break
        case '\r':
            escape = 'r'
            break
        case '\\':
            escape = '\\'
            break
        case '\'':
            escape = '\''
            break
        case '"': /* Better safe than sorry */
            escape = '"'
            break
        case '\032': //十进制26,八进制32,十六进制1a, /* This gives problems on Win32 */
            escape = 'Z'
        }

        if escape != 0 {
            dest = append(dest, '\\', escape)
        } else {
            dest = append(dest, c)
        }
    }

    // SQL standard is to use single-quotes for all values
    return "'" + string(dest) + "'"
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

 //Escape Escapes special characters in the txt, so it is safe to place returned string
 //to Query method.
//func Escape(c Conn, txt string) string {
    //if c.Status()&SERVER_STATUS_NO_BACKSLASH_ESCAPES != 0 {
        //return escapeQuotes(txt)
    //}
    //return escapeString(txt)
//}
