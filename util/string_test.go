package util

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestFilterInjectionsWords(t *testing.T) {
    testCases := []struct{ original, flitered string }{
        {"select * from x", "  from x"},
    }
    for _, tC := range testCases {
        t.Run(tC.original, func(t *testing.T) {
            assert.Equal(t, tC.flitered, FilterInjectionsWords(tC.original))
        })
    }
}
