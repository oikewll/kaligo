package log

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
        levelStr = WhiteLight + "[D] " + Reset
    case LevelInfo:
        levelStr = "[I] "
    case LevelWarn:
        levelStr = Yellow + "[W] " + Reset
    case LevelError:
        levelStr = Red + "[E] " + Reset
    }
    return levelStr + fmt.Sprintf(format, a...)
}
