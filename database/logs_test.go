package database

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
    assert.Equal(t, "SELECT * FROM test WHERE id = 1", Explain("SELECT * FROM test WHERE id = ?", 1))
}
