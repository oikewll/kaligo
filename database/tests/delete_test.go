package tests

import (
    "testing"
    "github.com/owner888/kaligo/logs"
    // "github.com/owner888/kaligo/database"
    // "github.com/stretchr/testify/assert"
)


func TestDelete(t *testing.T) {
    // 暂时不支持DELETE JOIN写法
    //sqlStr := db.Delete("user").Join("player", "LEFT").On("user.uid", "=", "player.uid").Where("player.id", "=", "test").Compile()
    sqlStr := db.Delete("user").Where("nickname", "=", "test").Compile()
    logs.Debug("sqlStr = ", sqlStr)
}

