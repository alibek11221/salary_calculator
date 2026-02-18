package file

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"salary_calculator/internal/pkg/logging"
)

type Cache[K comparable, V any] struct {
	dir    string
	ttl    time.Duration
	buf    sync.Pool
	logger logging.Logger
}

func New[K comparable, V any](dir string, ttl time.Duration, logger logging.Logger) *Cache[K, V] {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		logger.Error().Err(err).Str("dir", dir).Msg("failed to create cache directory")
	}

	return &Cache[K, V]{
		dir: dir,
		ttl: ttl,
		buf: sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
		logger: logger,
	}
}

type cacheEntry[V any] struct {
	Value V `json:"value"`
}

func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	var zero V
	if c == nil {
		return zero, false
	}
	k := fmt.Sprintf("%v", key)
	filename := filepath.Join(c.dir, k+".json.gz")

	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return zero, false
	}
	if err != nil {
		return zero, false
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return zero, false
	}
	defer func(gzReader *gzip.Reader) {
		_ = gzReader.Close()
	}(gzReader)

	var entry cacheEntry[V]
	if err := json.NewDecoder(gzReader).Decode(&entry); err != nil {
		return zero, false
	}

	return entry.Value, true
}

func (c *Cache[K, V]) Put(key K, value V) error {
	if c == nil {
		return fmt.Errorf("cache is nil")
	}
	k := fmt.Sprintf("%v", key)
	filename := filepath.Join(c.dir, k+".json.gz")

	buf := c.buf.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		c.buf.Put(buf)
	}()

	entry := cacheEntry[V]{
		Value: value,
	}

	if err := json.NewEncoder(buf).Encode(entry); err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	if _, err := buf.WriteTo(gzWriter); err != nil {
		return fmt.Errorf("failed to write compressed data: %w", err)
	}

	return nil
}
