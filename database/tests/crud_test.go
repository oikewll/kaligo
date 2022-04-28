package tests

import (
    // "testing"
)

type User struct {
    ID   uint   `db:"id"`
    Name string `db:"name"`
    Age  uint   `db:"age"`
    Sex  uint   `db:"sex"`
}


