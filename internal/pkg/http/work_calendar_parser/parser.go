package work_calendar_parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"salary_calculator/internal/pkg/logging"

	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/sync/singleflight"
)

var (
	ErrInvalidMonth = errors.New("invalid month")
	ErrInvalidYear  = errors.New("invalid year")
)

type Parser struct {
	logger   logging.Logger
	dir      string
	cacheCap int

	cache *lru.Cache[int, map[int]WorkdayResponse]
	group singleflight.Group
}

func New(dir string, cacheCap int, logger logging.Logger) *Parser {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		logger.Error().Err(err).Str("dir", dir).Msg("failed to create cache directory")
	}

	cache, err := lru.New[int, map[int]WorkdayResponse](cacheCap)
	if err != nil {
		logger.Fatal().Err(err).Int("cacheCap", cacheCap).Msg("failed to create LRU cache")
	}

	return &Parser{
		logger:   logger,
		dir:      dir,
		cacheCap: cacheCap,
		cache:    cache,
	}
}

// Parse возвращает календарь месяца из production-calendar файла года.
//
// Возвращённый WorkdayResponse read-only: слайс Days шарится с кэшем —
// не мутировать. При кэш-промахе загрузка файла года выполняется через
// singleflight: конкурентные запросы одного года читают файл один раз.
func (p *Parser) Parse(year, month int) (*WorkdayResponse, error) {
	if month < 1 || month > 12 {
		return nil, ErrInvalidMonth
	}
	if year < 1900 || year > 2100 {
		return nil, ErrInvalidYear
	}

	if yearEntries, ok := p.cache.Get(year); ok {
		if val, ok := yearEntries[month]; ok {
			return &val, nil
		}
		return nil, ErrInvalidMonth
	}

	v, err, _ := p.group.Do(strconv.Itoa(year), func() (any, error) {
		if yearEntries, ok := p.cache.Get(year); ok {
			return yearEntries, nil
		}

		entries, err := p.loadYear(year)
		if err != nil {
			return nil, err
		}

		p.cache.Add(year, entries)

		return entries, nil
	})
	if err != nil {
		return nil, err
	}

	entries := v.(map[int]WorkdayResponse)
	if val, ok := entries[month]; ok {
		return &val, nil
	}

	return nil, ErrInvalidMonth
}

func (p *Parser) loadYear(year int) (map[int]WorkdayResponse, error) {
	fileName := fmt.Sprintf("workdays_%d.json", year)
	path := filepath.Join(p.dir, fileName)

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			p.logger.Info().Str("file", fileName).Msg("file does not exist")
		}
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	var entries map[int]WorkdayResponse
	if err := json.NewDecoder(file).Decode(&entries); err != nil {
		return nil, fmt.Errorf("decode workdays file: %w", err)
	}

	return entries, nil
}
