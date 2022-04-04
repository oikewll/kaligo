package logs

import (
    "testing"
)

func TestLogLevel(t *testing.T) {
    Critical("Critical message")
    Error("Error message")
    Warn("Warn message")
    Info("Info message")
    Debug("Debug message")
}

func TestLogFile(t *testing.T) {
    root = &logger{formatter: &PlainFormatter{}, writer: Writers{&ConsoleWriter{}, NewFileWriter("logs.log")}, Level: LevelDebug}
    TestLogLevel(t)
    log := New("KALI", LevelDebug, nil)
    log.Error("Error with prefix")
}
