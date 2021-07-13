package mysql

import (
    "bytes"
    "fmt"
    "io"
)

// Version returns version string
func Version() string {
    return "1.5.3"
}

func syntaxError(ln int) error {
    return fmt.Errorf("syntax error at line: %d", ln)
}

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

// Escape: Escapes special characters in the txt, so it is safe to place returned string
// to Query method.
//func Escape(c Conn, txt string) string {
    //if c.Status()&SERVER_STATUS_NO_BACKSLASH_ESCAPES != 0 {
        //return escapeQuotes(txt)
    //}
    //return escapeString(txt)
//}
