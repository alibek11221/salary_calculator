package work_calendar

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

const (
	baseURL          = "https://production-calendar.ru/get-period"
	cacheKeyTemplate = "workdays_%d_%d"
)

type Service struct {
	client httpClient
	cache  cache
	token  string
}

func New(client httpClient, cache cache, token string) *Service {
	return &Service{
		cache:  cache,
		token:  token,
		client: client,
	}
}

func (s *Service) GetWorkdaysForYear(ctx context.Context) (map[int]WorkdayResponse, error) {
	start := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	year := 2026
	results := make(map[int]WorkdayResponse, 12)
	monthsToFetch := make([]int, 0, 12)
	urls := make([]string, 0, 12)

	for month := 1; month <= 12; month++ {
		cacheKey := fmt.Sprintf(cacheKeyTemplate, year, month)

		if cached, ok := s.cache.Get(cacheKey); ok {
			results[month] = cached
			continue
		}

		monthsToFetch = append(monthsToFetch, month)
		urls = append(urls, s.generateUrl(month, year))
	}

	if len(monthsToFetch) == 0 {
		log.Info().Str("duration", time.Since(start).String()).Msg("all data from cache")
		return results, nil
	}

	var mu sync.Mutex
	eg, ctx := errgroup.WithContext(ctx)

	for idx, month := range monthsToFetch {
		url := urls[idx]
		month := month
		eg.Go(
			func() error {
				response, err := s.Request(ctx, url)
				if err != nil {
					return err
				}

				cacheKey := fmt.Sprintf(cacheKeyTemplate, year, month)
				if cacheErr := s.cache.Put(cacheKey, *response); cacheErr != nil {
					log.Warn().Err(cacheErr).Int("month", month).Msg("failed to cache month workdays")
				}

				mu.Lock()
				results[month] = *response
				mu.Unlock()

				return nil
			},
		)
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	log.Info().Str("duration", time.Since(start).String()).Int("fetched_months", len(monthsToFetch)).Msg("fetch completed")
	return results, nil
}

func (s *Service) GetWorkdaysForMonth(ctx context.Context, month int, year int) (*WorkdayResponse, error) {
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month: %d", month)
	}

	cacheKey := fmt.Sprintf(cacheKeyTemplate, year, month)

	if cached, ok := s.cache.Get(cacheKey); ok {
		return &cached, nil
	}

	url := s.generateUrl(month, year)
	response, err := s.Request(ctx, url)
	if err != nil {
		return nil, err
	}

	err = s.cache.Put(cacheKey, *response)
	if err != nil {
		log.Warn().Err(err).Msg("failed to cache month workdays")
	}

	return response, nil
}

func (s *Service) Request(ctx context.Context, url string) (*WorkdayResponse, error) {
	start := time.Now()
	defer func() {
		log.Debug().Str("url", url).Dur("duration", time.Since(start)).Msg("request completed")
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return parseResponse(resp)
}

func parseResponse(resp *http.Response) (*WorkdayResponse, error) {
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response read error: %w", err)
	}

	var response WorkdayResponse
	if err = json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}

	return &response, nil
}

func (s *Service) generateUrl(month int, year int) string {
	days := daysInMonth(month, year)
	period := createPeriod(days, month, year)

	return fmt.Sprintf("%s/%s/ru/%s/json", baseURL, s.token, period)
}

func createPeriod(endDay, month, year int) string {
	return fmt.Sprintf(
		"%02d.%02d.%d-%02d.%02d.%d",
		1, month, year,
		endDay, month, year,
	)
}

func daysInMonth(month, year int) int {
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
}
