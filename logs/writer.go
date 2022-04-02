package logs

import (
    "fmt"
    "io"
    "os"
)

// Writer log writer interface
type Writer interface {
    Write(string string)
}

type Writers []Writer

func (w Writers) Write(string string) {
    for _, v := range w {
        v.Write(string)
    }
}

type WriterFunc func(string string)

func (w WriterFunc) Write(string string) {
    w(string)
}

type DefaultWriter struct {
    writer io.Writer
}

func (w *DefaultWriter) Write(string string) {
    w.writer.Write([]byte(string))
}

type ConsoleWriter struct {
}

func (w *ConsoleWriter) Write(string string) {
    fmt.Print(string)
}

type FileWriter struct {
    DefaultWriter
}

func NewFileWriter(path string) *FileWriter {
    f, _ := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
    return &FileWriter{DefaultWriter: DefaultWriter{f}}
}
