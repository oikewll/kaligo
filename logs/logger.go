package logs

import (
    "errors"
    "fmt"
    "os"
    "strings"
)

// ErrRecordNotFound record not found error
var ErrRecordNotFound = errors.New("record not found")

// Level ...
type Level int

const (
    LevelDefault Level = iota
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
    root Logger = &logger{formatter: &ConsoleFormatter{}, writer: &ConsoleWriter{}, level: LevelDebug}
)

// Logger logger interface
type Logger interface {
    Level() Level
    LogMode(Level) Logger
    //Info(context.Context, string, ...any)
    Debug(...any)
    Info(...any)
    Warn(...any)
    Error(...any)
    Critical(...any)
    Debugf(string, ...any)
    Infof(string, ...any)
    Warnf(string, ...any)
    Errorf(string, ...any)
    Criticalf(string, ...any)
    Panic(...any)
    Panicf(string, ...any)
    Fatal(...any)
    Fatalf(string, ...any)
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
    level      Level
    Prefix     string
    TimeFormat string
    parent     Logger // formatter 和 writer 可以继承自 parent
}

func New(prefix string, level Level, parant Logger) Logger {
    if parant == nil {
        parant = root
    }
    return &logger{Prefix: prefix, level: level, parent: parant}
}

func (l *logger) LogMode(level Level) Logger {
    l.level = level
    return l
}

func (l *logger) Debug(data ...any) {
    l.Debugf("", data...)
}

func (l *logger) Info(data ...any) {
    l.Infof("", data...)
}

func (l *logger) Warn(data ...any) {
    l.Warnf("", data...)
}

func (l *logger) Error(data ...any) {
    l.Errorf("", data...)
}

func (l *logger) Critical(data ...any) {
    l.Criticalf("", data...)
}

func (l *logger) Debugf(msg string, data ...any) {
    if l.Level() >= LevelDebug {
        l.getWriter().Write(l.getFormatter().Printf(l.Prefix, LevelDebug, msg, data...))
    }
}

func (l *logger) Infof(msg string, data ...any) {
    if l.Level() >= LevelInfo {
        l.getWriter().Write(l.getFormatter().Printf(l.Prefix, LevelInfo, msg, data...))
    }
}

func (l *logger) Warnf(msg string, data ...any) {
    if l.Level() >= LevelWarn {
        l.getWriter().Write(l.getFormatter().Printf(l.Prefix, LevelWarn, msg, data...))
    }
}

func (l *logger) Errorf(msg string, data ...any) {
    if l.Level() >= LevelError {
        l.getWriter().Write(l.getFormatter().Printf(l.Prefix, LevelError, msg, data...))
    }
}

func (l *logger) Criticalf(msg string, data ...any) {
    if l.Level() >= LevelCritical {
        l.getWriter().Write(l.getFormatter().Printf(l.Prefix, LevelCritical, msg, data...))
    }
}

// Panic 输出 Critical 日志并 panic
func (l *logger) Panic(data ...any) {
    l.Critical(data...)
    panic(fmt.Sprint(data...))
}

// Fatal 输出 Critical 日志并退出程序
func (l *logger) Fatal(data ...any) {
    l.Critical(data...)
    os.Exit(1)
}

// Panic 输出 Critical 日志并 panic
func (l *logger) Panicf(msg string, data ...any) {
    l.Criticalf(msg, data...)
    panic(fmt.Sprintf(msg, data...))
}

// Fatal 输出 Critical 日志并退出程序
func (l *logger) Fatalf(msg string, data ...any) {
    l.Criticalf(msg, data...)
    os.Exit(1)
}

func (l *logger) Level() Level {
    if l.level == LevelDefault {
        return l.parent.(*logger).Level()
    }
    return l.level
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

// ====== Public 工具类方法 ======

// ParseLevel string 转 Level debug > info > warn > error > critical > silent
func ParseLevel(level string) Level {
    level = strings.ToLower(level)
    switch level {
    case "silent":
        return LevelSilent
    case "critical":
        return LevelCritical
    case "error":
        return LevelError
    case "warn":
        return LevelWarn
    case "info":
        return LevelInfo
    case "debug":
        return LevelDebug
    default:
        return LevelDefault
    }
}

// ====== 以下是 root logger 的快捷方式 ======

func LogMode(level Level) Logger {
    root.LogMode(level)
    return root
}

func Debug(data ...any) {
    root.Debug(data...)
}

func Info(data ...any) {
    root.Info(data...)
}

func Warn(data ...any) {
    root.Warn(data...)
}

func Error(data ...any) {
    root.Error(data...)
}

func Critical(data ...any) {
    root.Critical(data...)
}

func Debugf(msg string, data ...any) {
    root.Debugf(msg, data...)
}

func Infof(msg string, data ...any) {
    root.Infof(msg, data...)
}

func Warnf(msg string, data ...any) {
    root.Warnf(msg, data...)
}

func Errorf(msg string, data ...any) {
    root.Errorf(msg, data...)
}

func Criticalf(msg string, data ...any) {
    root.Criticalf(msg, data...)
}

// Panic 输出 Critical 日志并 panic
func Panic(data ...any) {
    root.Panic(data...)
}

// Panic 输出 Critical 日志并 panic
func Panicf(msg string, data ...any) {
    root.Panicf(msg, data...)
}

// Fatal 输出 Critical 日志并退出程序
func Fatal(data ...any) {
    root.Fatal(data...)
}

// Fatal 输出 Critical 日志并退出程序
func Fatalf(msg string, data ...any) {
    root.Fatalf(msg, data...)
}
