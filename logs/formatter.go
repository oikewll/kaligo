package logs

import (
    "fmt"
    "runtime"
    "strconv"
    "strings"
    "time"

    "github.com/owner888/kaligo/logs/ansi"
)

var (
    Seperator       = " "
    Terminator      = "\n"
    showFunction    = false
    StacktraceDepth = 4
)

// Formatter 格式化日志输出
type Formatter interface {
    Printf(prefix string, level Level, format string, a ...any) string
}

type PlainFormatter struct {
}

func (f *PlainFormatter) Printf(prefix string, level Level, format string, a ...any) string {
    var builder strings.Builder
    builder.WriteString(formatTime(time.Now()))
    builder.WriteString(Seperator)
    if len(prefix) > 0 {
        builder.WriteString("[")
        builder.WriteString(prefix)
        builder.WriteString("]")
        builder.WriteString(Seperator)
    }
    if len(format) > 0 {
        builder.WriteString(fmt.Sprintf(format, a...))
    } else {
        builder.WriteString(fmt.Sprint(a...))
    }
    builder.WriteString(Terminator)
    return builder.String()
}

type ConsoleFormatter struct {
}

func (f *ConsoleFormatter) Printf(prefix string, level Level, format string, a ...any) string {
    var builder strings.Builder
    builder.WriteString(formatTime(time.Now()))
    builder.WriteString(Seperator)
    switch level {
    case LevelDebug:
        builder.WriteString(ansi.White)
        builder.WriteString("Ⓓ DEBUG")
        builder.WriteString(ansi.Reset)
    case LevelInfo:
        builder.WriteString(ansi.Cyan)
        builder.WriteString("ⓘ INFO ")
        builder.WriteString(ansi.Reset)
    case LevelWarn:
        builder.WriteString(ansi.Yellow)
        builder.WriteString("ⓦ WARN ")
        builder.WriteString(ansi.Reset)
    case LevelError:
        builder.WriteString(ansi.Red)
        builder.WriteString("ⓧ ERROR")
        builder.WriteString(ansi.Reset)
    case LevelCritical:
        builder.WriteString(ansi.RedBold)
        builder.WriteString("🅧 CRITICAL")
    }
    builder.WriteString(Seperator)
    if len(prefix) > 0 {
        builder.WriteString("[")
        builder.WriteString(prefix)
        builder.WriteString("]")
        builder.WriteString(Seperator)
    }
    if fun, file, line, ok := runtime.Caller(StacktraceDepth); ok {
        if showFunction {
            builder.WriteString(fmt.Sprintf("[%v]", extractName(runtime.FuncForPC(fun).Name())))
            builder.WriteString(Seperator)
        }
        builder.WriteString(fmt.Sprintf("(%v:%v)", extractName(file), line))
        builder.WriteString(Seperator)
    }
    if len(format) > 0 {
        builder.WriteString(fmt.Sprintf(format, a...))
    } else {
        builder.WriteString(fmt.Sprint(a...))
    }
    if level == LevelCritical {
        builder.WriteString(ansi.Reset)
    }
    builder.WriteString(Terminator)
    return builder.String()
}

func formatTime(t time.Time) string {
    var builder strings.Builder
    year, month, day := t.Date()
    hour, minite, second := t.Clock()
    components := []string{
        strconv.Itoa(year), "-", strconv.Itoa(int(month)), "-", strconv.Itoa(day), " ",
        strconv.Itoa(hour), ":", strconv.Itoa(minite), ":", strconv.Itoa(second),
    }
    for _, v := range components {
        builder.WriteString(v)
    }
    return builder.String()
}

func extractName(s string) string {
    parts := strings.Split(s, "/")
    return parts[len(parts)-1]
}
