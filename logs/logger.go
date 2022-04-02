package logs

import (
    "errors"
    "time"
)

// ErrRecordNotFound record not found error
var ErrRecordNotFound = errors.New("record not found")

// Level ...
type Level int

const (
    _ Level = iota
    // LevelSilent is the default log level
    LevelSilent
    // LevelError is the error log level
    LevelError
    // LevelWarn is the warn log level
    LevelWarn
    // LevelInfo is the lower log level
    LevelInfo
    // LevelDebug is the lower log level
    LevelDebug
)

var (
    root Logger = &logger{formatter: &ConsoleFormatter{}, writer: &ConsoleWriter{}, Level: LevelDebug}
)

// Logger logger interface
type Logger interface {
    LogMode(Level) Logger
    //Info(context.Context, string, ...any)
    Debug(string, ...any)
    Info(string, ...any)
    Warn(string, ...any)
    Error(string, ...any)
    Trace(begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

type Log struct {
    Message string
    Level   Level
    file    string
    line    int
}

type logger struct {
    formatter  Formatter
    writer     Writer
    Level      Level
    Prefix     string
    TimeFormat string
    parent     Logger // formatter 和 writer 可以继承自 parent
}

func New(prefix string, level Level, parant Logger) Logger {
    if parant == nil {
        parant = root
    }
    return &logger{Prefix: prefix, Level: level, parent: parant}
}

func (l *logger) LogMode(Level) Logger {
    return l
}

func (l *logger) Debug(msg string, data ...any) {
    if l.Level >= LevelDebug {
        l.getWriter().Write(l.getFormatter().Printf(l.Prefix, LevelDebug, msg, data...))
    }
}

func (l *logger) Info(msg string, data ...any) {
    if l.Level >= LevelInfo {
        l.getWriter().Write(l.getFormatter().Printf(l.Prefix, LevelInfo, msg, data...))
    }
}

func (l *logger) Warn(msg string, data ...any) {
    if l.Level >= LevelWarn {
        l.getWriter().Write(l.getFormatter().Printf(l.Prefix, LevelWarn, msg, data...))
    }
}

func (l *logger) Error(msg string, data ...any) {
    if l.Level >= LevelError {
        l.getWriter().Write(l.getFormatter().Printf(l.Prefix, LevelError, msg, data...))
    }
}

func (l *logger) getWriter() Writer {
    if l.writer == nil {
        return l.parent.(*logger).getWriter()
    }
    return l.writer
}

func (l *logger) getFormatter() Formatter {
    if l.formatter == nil {
        return l.parent.(*logger).getFormatter()
    }
    return l.formatter
}

func (l *logger) Trace(begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

}

func Debug(msg string, data ...any) {
    root.Debug(msg, data...)
}

func Info(msg string, data ...any) {
    root.Info(msg, data...)
}

func Warn(msg string, data ...any) {
    root.Warn(msg, data...)
}

func Error(msg string, data ...any) {
    root.Error(msg, data...)
}
