package tests

import (
    "testing"
    // "github.com/owner888/kaligo/logs"
    // "github.com/owner888/kaligo/database"
    "github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
    _, err := db.Delete("demo_user").Where("id", "=", 2).Execute()
    assert.NoError(t, err)
}

// 通过解密字段删除
func TestDeleteDecrypt(t *testing.T) {
    _, err := db.Delete("demo_user").Where("realname", "=", "demoname").
    SetCryptKey("aaa").
    SetCryptFields(map[string][]string{
        "demo_user": {"realname"},
    }).
    Execute()
    assert.NoError(t, err)
}

// 通过 JOIN 条件删除
func TestDeleteJoin(t *testing.T) {
    // 暂时不支持DELETE JOIN写法
    //sqlStr := db.Delete("user").Join("player", "LEFT").On("user.uid", "=", "player.uid").Where("player.id", "=", "test").Compile()
}
