package kaligo

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
    timer := NewTimer()
    assert.NotNil(t, timer)
}
