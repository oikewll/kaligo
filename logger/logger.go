package logger

import (
    "errors"
    "time"
    //"context"
)

// ErrRecordNotFound record not found error
var ErrRecordNotFound = errors.New("record not found")

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// LogLevel ...
type LogLevel int

const (
    // Silent is the default log level
    Silent LogLevel = iota + 1
    // Error is the error log level
    Error
    // Warn is the warn log level
    Warn
    // Info is the lower log level
    Info
)

// Writer log writer interface
type Writer interface {
	Printf(string, ...interface{})
}

// Interface logger interface
type Interface interface {
    LogMode(LogLevel) Interface
    //Info(context.Context, string, ...interface{})
    Info(string, ...interface{})
    Warn(string, ...interface{})
    Error(string, ...interface{})
    Trace(begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

type logger struct {
    Writer
    //Config
    LogLevel LogLevel
    infoStr, warnStr, errStr            string
    traceStr, traceErrStr, traceWarnStr string
}

func (l logger) Info(msg string, data ...interface{}) {
    if l.LogLevel >= Info {
        //l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
    }
}

func (l logger) Warn(msg string, data ...interface{}) {
    if l.LogLevel >= Warn {
        //l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
    }
}

func (l logger) Error(msg string, data ...interface{}) {
    if l.LogLevel >= Error {
        //l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
    }
}

