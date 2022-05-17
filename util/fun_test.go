package util

import (
    "testing"
    // "log"

    // "github.com/stretchr/testify/assert"
)

type User struct {
    Name string
    Age int
}

func TestPrintMemStats(t *testing.T) {
    PrintMemStats()
}

func TestStructPrint(t *testing.T) {
    u := User{
        Name: "demo",
        Age: 10,
    }
    StructPrint(u)
}
