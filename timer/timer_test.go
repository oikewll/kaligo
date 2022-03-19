package timer

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
    timer := New()
    assert.NotNil(t, timer)
}
