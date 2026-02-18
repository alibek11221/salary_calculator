package file

import (
	"os"
	"testing"
	"time"

	"salary_calculator/internal/pkg/logging"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cache_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	l := logging.New(false)
	ttl := 1 * time.Hour
	cache := New[string, string](tempDir, ttl, l)

	key := "test_key"
	val := "test_value"

	v, ok := cache.Get(key)
	assert.False(t, ok)
	assert.Empty(t, v)

	err = cache.Put(key, val)
	assert.NoError(t, err)

	v, ok = cache.Get(key)
	assert.True(t, ok)
	assert.Equal(t, val, v)

	shortTTL := 10 * time.Millisecond
	cacheShort := New[string, string](tempDir, shortTTL, l)
	err = cacheShort.Put("short_key", "short_val")
	assert.NoError(t, err)

	time.Sleep(20 * time.Millisecond)
	v, ok = cacheShort.Get("short_key")
	assert.False(t, ok)
	assert.Empty(t, v)
}

func TestCache_Errors(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cache_errors")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	l := logging.New(false)
	cache := New[string, string](tempDir, 1*time.Hour, l)

	t.Run("invalid_gzip", func(t *testing.T) {
		err := os.WriteFile(tempDir+"/invalid.json.gz", []byte("not gzip"), 0o644)
		assert.NoError(t, err)

		v, ok := cache.Get("invalid")
		assert.False(t, ok)
		assert.Empty(t, v)
	})
}

func TestCache_IntKey(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cache_test_int")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	l := logging.New(false)
	cache := New[int, int](tempDir, 1*time.Hour, l)
	err = cache.Put(123, 456)
	assert.NoError(t, err)

	v, ok := cache.Get(123)
	assert.True(t, ok)
	assert.Equal(t, 456, v)
}
