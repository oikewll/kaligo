package tests

import (
    "testing"
    // "github.com/owner888/kaligo/logs"
    // "github.com/owner888/kaligo/database"
    "github.com/stretchr/testify/assert"
)

// 查询数据 DECODE
func TestSelectDecrypt(t *testing.T) {
    var realname string
    _, err := db.Select("realname").From("demo_user").Where("id", "=", 4).
    SetCryptKey("aaa").
    SetCryptFields(map[string][]string{
        "demo_user": {"realname"},
    }).
    Scan(&realname).Execute()
    assert.NoError(t, err)
    assert.NotNil(t, realname)
}


// 查询条件 DECODE
func TestSelectWhereDecrypt(t *testing.T) {
    var gender int64
    _, err := db.Select("gender").From("demo_user").Where("realname", "=", "demoname").
    SetCryptKey("aaa").
    SetCryptFields(map[string][]string{
        "demo_user": {"realname"},
    }).
    Scan(&gender).Execute()
    assert.NoError(t, err)
    assert.NotNil(t, gender)
}
