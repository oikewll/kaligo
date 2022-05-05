package tests

import (
    "testing"
    "github.com/owner888/kaligo/logs"
    // "github.com/owner888/kaligo/database"
    "github.com/stretchr/testify/assert"
)

func TestSelectDecrypt(t *testing.T) {
    logs.Info("TestSelectDecrypt")

    var ages int64
    _, err := db.Select("age").From("user").Where("id", "=", 1).
    SetCryptKey("aaa").
    SetCryptFields(map[string][]string{
        "user": {"name", "age"},
    }).
    Scan(&ages).Execute()
    assert.NoError(t, err)
    assert.NotNil(t, ages)
}
