package tests

import (
    "testing"
    "github.com/owner888/kaligo/logs"
    // "github.com/owner888/kaligo/database"
    "github.com/stretchr/testify/assert"
)

func TestSelectDecrypt(t *testing.T) {
    logs.Info("TestSelectDecrypt")

    var gender int64
    _, err := db.Select("gender").From("demo_user").Where("id", "=", 1).
    SetCryptKey("aaa").
    SetCryptFields(map[string][]string{
        "user": {"realname"},
    }).
    Scan(&gender).Execute()
    assert.NoError(t, err)
    assert.NotNil(t, gender)
}
