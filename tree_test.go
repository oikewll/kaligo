package kaligo

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestParams(t *testing.T) {
    params := Params{
        Param{"key1", "value1"},
        Param{"key2", "value2"},
        Param{"key3", ""},
    }
    assert.Equal(t, "value1", params.ByName("key1"))
    assert.Equal(t, "", params.ByName("key3", "defaultValue"))
    assert.Equal(t, "defaultValue", params.ByName("unknown_key", "defaultValue"))
    assert.Equal(t, "", params.ByName("unknown_key"))
}
