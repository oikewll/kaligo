package log

import (
    "testing"
)

func TestLogLevel(t *testing.T) {
    Error("Error message")
    Warn("Warning message")
    Info("Info message")
    Debug("Debug message")
}
