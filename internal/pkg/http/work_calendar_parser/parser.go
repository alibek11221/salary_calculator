package work_calendar_parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"salary_calculator/internal/pkg/logging"
	"slices"

	lru "github.com/hashicorp/golang-lru/v2"
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

func (p *Parser) Parse(year, month int) (*WorkdayResponse, error) {
	if month < 1 || month > 12 {
		return nil, ErrInvalidMonth
	}
	if year < 1900 || year > 2100 {
		return nil, ErrInvalidYear
	}

	if yearEntries, ok := p.cache.Get(year); ok {
		if val, ok := yearEntries[month]; ok {
			return p.cloneResponse(&val), nil
		}
		return nil, ErrInvalidMonth
	}

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

	p.cache.Add(year, entries)

	if val, ok := entries[month]; ok {
		return p.cloneResponse(&val), nil
	}

	return nil, ErrInvalidMonth
}

func (p *Parser) cloneResponse(res *WorkdayResponse) *WorkdayResponse {
	if res == nil {
		return nil
	}

	return &WorkdayResponse{
		Statistics: res.Statistics,
		Days:       slices.Clone(res.Days),
	}
}
