package tests

import (
    "testing"
    // "github.com/owner888/kaligo/logs"
    "github.com/stretchr/testify/assert"
)


func TestUpdate(t *testing.T) {
    sets := map[string]string{"username":"TestUpdate"}
    _, err := db.Update("demo_user").Set(sets).Where("id", "=", 1).Execute()
    assert.NoError(t, err)

}

func TestUpdateLeftJoin(t *testing.T) {
    // sqlStr := db.Update("user").Join("player", "LEFT").On("user.id", "=", "player.uid").Set(sets).Where("player.id", "=", 10).Compile()
    // logs.Debug("sqlStr = ", sqlStr)
}


