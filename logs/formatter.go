package logs

import "fmt"

type Formatter interface {
    Printf(level Level, format string, a ...any) string
}

type PlainFormatter struct {
}

func (f *PlainFormatter) Printf(level Level, format string, a ...any) string {
    return fmt.Sprintf(format, a...) + "\n"
}

type ConsoleFormatter struct {
}

func (f *ConsoleFormatter) Printf(level Level, format string, a ...any) string {
    levelStr := ""
    switch level {
    case LevelDebug:
        levelStr = BlackLight + "Ⓓ DEBUG " + Reset
    case LevelInfo:
        levelStr = "ⓘ INFO  "
    case LevelWarn:
        levelStr = Yellow + "ⓦ WARN  " + Reset
    case LevelError:
        levelStr = Red + "ⓧ ERROR " + Reset
    }
    return levelStr + fmt.Sprintf(format, a...)
}
