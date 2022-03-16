package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	cache := New("redis")
	assert.IsType(t, &Redis{}, cache)
	cache = New("memcache")
	assert.IsType(t, &Memcache{}, cache)
	cache = New("memory")
	assert.IsType(t, &Memory{}, cache)
	cache = New("")
	assert.IsType(t, &Memory{}, cache)
}
