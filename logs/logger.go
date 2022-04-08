package logs

import (
    "errors"
    "fmt"
    "os"
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
    // LevelError  is the critical log level
    LevelCritical
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
    Critical(string, ...any)
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

func (l *logger) Critical(msg string, data ...any) {
    if l.Level >= LevelCritical {
        l.getWriter().Write(l.getFormatter().Printf(l.Prefix, LevelCritical, msg, data...))
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

func Debug(data ...any) {
    root.Debug("", data...)
}

func Info(data ...any) {
    root.Info("", data...)
}

func Warn(data ...any) {
    root.Warn("", data...)
}

func Error(data ...any) {
    root.Error("", data...)
}

func Critical(data ...any) {
    root.Critical("", data...)
}

func Debugf(msg string, data ...any) {
    root.Debug(msg, data...)
}

func Infof(msg string, data ...any) {
    root.Info(msg, data...)
}

func Warnf(msg string, data ...any) {
    root.Warn(msg, data...)
}

func Errorf(msg string, data ...any) {
    root.Error(msg, data...)
}

func Criticalf(msg string, data ...any) {
    root.Critical(msg, data...)
}

// Panic 输出 Critical 日志并 panic
func Panic(msg string, data ...any) {
    Criticalf(msg, data...)
    panic(fmt.Sprintf(msg, data...))
}

// Fatal 输出 Critical 日志并退出程序
func Fatal(msg string, data ...any) {
    Criticalf(msg, data...)
    os.Exit(1)
}
