package work_calendar_parser

import (
	"fmt"
	"os"
	"path/filepath"
	"salary_calculator/internal/pkg/logging"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	logger := logging.New(false)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testDir := filepath.Join(wd, "testdata")

	p := New(testDir, 10, logger)

	t.Run("success", func(t *testing.T) {
		res, err := p.Parse(2025, 1)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, res)
		assert.Equal(t, 20, res.Statistics.WorkDays)
		assert.Equal(t, 31, res.Statistics.CalendarDays)
	})

	t.Run("success_month_2", func(t *testing.T) {
		res, err := p.Parse(2025, 2)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, res)
		assert.Equal(t, 19, res.Statistics.WorkDays)
	})

	t.Run("error_invalid_month", func(t *testing.T) {
		res, err := p.Parse(2025, 13)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.ErrorIs(t, err, ErrInvalidMonth)
	})

	t.Run("error_file_not_exist", func(t *testing.T) {
		res, err := p.Parse(2024, 1)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("caching", func(t *testing.T) {
		// First call - loads from file
		res1, err := p.Parse(2025, 1)
		if err != nil {
			t.Fatal(err)
		}

		// Delete file to ensure second call uses cache
		filePath := filepath.Join(testDir, "workdays_2025.json")
		tempPath := filePath + ".bak"
		err = os.Rename(filePath, tempPath)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			_ = os.Rename(tempPath, filePath)
		}()

		// Second call - should use cache
		res2, err := p.Parse(2025, 1)
		assert.NoError(t, err)
		assert.Equal(t, res1, res2)
	})

	t.Run("concurrency_race", func(t *testing.T) {
		// This test is designed to trigger race conditions when run with -race flag
		const goroutines = 20
		const iterations = 50
		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					// Use different months to test inner map access
					month := (j % 2) + 1
					_, _ = p.Parse(2025, month)
				}
			}(i)
		}
		wg.Wait()
	})

	t.Run("cache_isolation", func(t *testing.T) {
		res1, err := p.Parse(2025, 1)
		assert.NoError(t, err)
		assert.NotEmpty(t, res1.Days)

		// Modify an EXISTING element of the slice
		res1.Days[0].TypeText = "Modified"

		res2, err := p.Parse(2025, 1)
		assert.NoError(t, err)
		assert.NotEmpty(t, res2.Days)
		assert.NotEqual(t, "Modified", res2.Days[0].TypeText, "Cache element was modified!")
	})

	t.Run("eviction", func(t *testing.T) {
		logger := logging.New(false)
		wd, _ := os.Getwd()
		testDir := filepath.Join(wd, "testdata")

		// Create files for different years to test eviction
		years := []int{2025, 2026, 2027}
		for _, y := range years {
			data, _ := os.ReadFile(filepath.Join(testDir, "workdays_2025.json"))
			_ = os.WriteFile(filepath.Join(testDir, fmt.Sprintf("workdays_%d.json", y)), data, 0o644)
		}
		defer func() {
			for _, y := range years {
				if y != 2025 {
					_ = os.Remove(filepath.Join(testDir, fmt.Sprintf("workdays_%d.json", y)))
				}
			}
		}()

		// Cap = 2
		p := New(testDir, 2, logger)

		// Load 2025 and 2026
		_, _ = p.Parse(2025, 1)
		_, _ = p.Parse(2026, 1)
		assert.Equal(t, 2, p.cache.Len())

		// Access 2025 again (move to front)
		_, _ = p.Parse(2025, 1)

		// Load 2027 (should evict 2026)
		_, _ = p.Parse(2027, 1)
		assert.Equal(t, 2, p.cache.Len())
		assert.False(t, p.cache.Contains(2026))
		assert.True(t, p.cache.Contains(2025))
		assert.True(t, p.cache.Contains(2027))
	})

	t.Run("invalid_month", func(t *testing.T) {
		res, err := p.Parse(2025, 0)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.ErrorIs(t, err, ErrInvalidMonth)

		res, err = p.Parse(2025, 13)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.ErrorIs(t, err, ErrInvalidMonth)
	})

	t.Run("invalid_year", func(t *testing.T) {
		res, err := p.Parse(1899, 1)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.ErrorIs(t, err, ErrInvalidYear)

		res, err = p.Parse(2101, 1)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.ErrorIs(t, err, ErrInvalidYear)
	})
}

func BenchmarkParser_Parse_CacheHit(b *testing.B) {
	logger := logging.New(false)
	wd, _ := os.Getwd()
	testDir := filepath.Join(wd, "testdata")
	p := New(testDir, 10, logger)

	// Pre-populate cache
	_, _ = p.Parse(2025, 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.Parse(2025, 1)
	}
}

func BenchmarkParser_Parse_CacheMiss(b *testing.B) {
	logger := logging.New(false)
	wd, _ := os.Getwd()
	testDir := filepath.Join(wd, "testdata")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		p := New(testDir, 10, logger)
		b.StartTimer()
		_, _ = p.Parse(2025, 1)
	}
}

func BenchmarkParser_Parse_Concurrent(b *testing.B) {
	logger := logging.New(false)
	wd, _ := os.Getwd()
	testDir := filepath.Join(wd, "testdata")
	p := New(testDir, 10, logger)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = p.Parse(2025, 1)
		}
	})
}
