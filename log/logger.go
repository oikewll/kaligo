package log

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
    standardLogger Logger = logger{formatter: &ConsoleFormatter{}, writer: &ConsoleWriter{}, Level: LevelDebug}
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
}

type logger struct {
    formatter Formatter
    writer    Writer
    Level     Level
}

func (l logger) LogMode(Level) Logger {
    return l
}

func (l logger) Debug(msg string, data ...any) {
    if l.Level >= LevelDebug {
        l.writer.Write(l.formatter.Printf(LevelDebug, msg, data...))
    }
}

func (l logger) Info(msg string, data ...any) {
    if l.Level >= LevelInfo {
        l.writer.Write(l.formatter.Printf(LevelInfo, msg, data...))
    }
}

func (l logger) Warn(msg string, data ...any) {
    if l.Level >= LevelWarn {
        l.writer.Write(l.formatter.Printf(LevelWarn, msg, data...))
    }
}

func (l logger) Error(msg string, data ...any) {
    if l.Level >= LevelError {
        l.writer.Write(l.formatter.Printf(LevelError, msg, data...))
    }
}

func (l logger) Trace(begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

}

func Debug(msg string, data ...any) {
    standardLogger.Debug(msg, data...)
}

func Info(msg string, data ...any) {
    standardLogger.Info(msg, data...)
}

func Warn(msg string, data ...any) {
    standardLogger.Warn(msg, data...)
}

func Error(msg string, data ...any) {
    standardLogger.Error(msg, data...)
}
