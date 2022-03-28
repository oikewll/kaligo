package logs

import (
	"testing"
)

func TestLogLevel(t *testing.T) {
	Error("Error message")
	Warn("Warning message")
	Info("Info message")
	Debug("Debug message")
}

func TestLogFile(t *testing.T) {
	standardLogger = logger{formatter: &PlainFormatter{}, writer: NewFileWriter("logs.log"), Level: LevelDebug}
	TestLogLevel(t)
}
