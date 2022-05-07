package tests

import (
    "testing"
    // "github.com/owner888/kaligo/logs"
    "github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
    sets := map[string]string{"username": "TestUpdate"}
    _, err := db.Update("demo_user").Set(sets).Where("id", "=", 1).Execute()
    assert.NoError(t, err)
}

// 通过解密字段删除
func TestUpdateDecrypt(t *testing.T) {
    sets := map[string]string{"username": "TestUpdateDecrypt"}
    _, err := db.Update("demo_user").Set(sets).Where("realname", "=", "demoname").
    SetCryptKey("aaa").
    SetCryptFields(map[string][]string{
        "demo_user": {"realname"},
    }).
    Execute()
    assert.NoError(t, err)
}

// 通过 JOIN 条件修改字段内容
func TestUpdateJoin(t *testing.T) {
    // sqlStr := db.Update("user").Join("player", "LEFT").On("user.id", "=", "player.uid").Set(sets).Where("player.id", "=", 10).Compile()
    // logs.Debug("sqlStr = ", sqlStr)
}


