package util

import (
    "strconv"
    "testing"

    "github.com/stretchr/testify/assert"
)

type IntSlice []int

func TestCast(t *testing.T) {
    s := []int{1, 2}
    var anySlice []any = CastSlice[int, any](s)
    assert.Equal(t, []any{1, 2}, anySlice)
}

func TestMap(t *testing.T) {
    s := []int{1, 2}
    var anySlice []string = MapSlice(s, strconv.Itoa)
    assert.Equal(t, []string{"1", "2"}, anySlice)
}

func TestCompactMap(t *testing.T) {
    s := []string{"1", "2"}
    anySlice, err := CompactMapSliceE(s, strconv.Atoi)
    assert.NoError(t, err)
    assert.Equal(t, []int{1, 2}, anySlice)
    s = []string{"a", "2"}
    anySlice, err = CompactMapSliceE(s, strconv.Atoi)
    assert.Error(t, err)
    assert.Equal(t, []int{0, 2}, anySlice)
}

func TestFlat(t *testing.T) {
    s := [][]int{{1, 2}, {3, 4}}
    assert.Equal(t, []int{1, 2, 3, 4}, FlatSlice(s))
    ss := []IntSlice{{1, 2}, {3, 4}}
    assert.Equal(t, IntSlice{1, 2, 3, 4}, FlatSlice(ss))
}
