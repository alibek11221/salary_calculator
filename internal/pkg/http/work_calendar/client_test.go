package work_calendar_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"salary_calculator/internal/pkg/http/work_calendar"
	"salary_calculator/internal/pkg/logging"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const token = "token"

func TestClient_GetWorkdaysForYear(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(cache *Mockcache, client *MockhttpClient)
		wantErr    bool
		wantErrMsg string
		wantCount  int
	}{
		{
			name: "all_cached",
			setupMocks: func(cache *Mockcache, client *MockhttpClient) {
				for m := 1; m <= 12; m++ {
					key := fmt.Sprintf("workdays_2026_%d", m)
					cache.EXPECT().Get(key).Return(work_calendar.WorkdayResponse{}, true)
				}
			},
			wantCount: 12,
		},
		{
			name: "fetch_all_success",
			setupMocks: func(cache *Mockcache, client *MockhttpClient) {
				for m := 1; m <= 12; m++ {
					key := fmt.Sprintf("workdays_2026_%d", m)
					cache.EXPECT().Get(key).Return(work_calendar.WorkdayResponse{}, false)

					resp := &work_calendar.WorkdayResponse{}
					respData, _ := json.Marshal(resp)
					client.EXPECT().Do(gomock.Any()).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(respData)),
					}, nil)

					cache.EXPECT().Put(key, *resp).Return(nil)
				}
			},
			wantCount: 12,
		},
		{
			name: "fetch_error",
			setupMocks: func(cache *Mockcache, client *MockhttpClient) {
				cache.EXPECT().Get(gomock.Any()).Return(work_calendar.WorkdayResponse{}, false).AnyTimes()
				client.EXPECT().Do(gomock.Any()).Return(nil, fmt.Errorf("network error")).AnyTimes()
			},
			wantErr:    true,
			wantErrMsg: "network error",
		},
		{
			name: "fetch_some_cache_put_error",
			setupMocks: func(cache *Mockcache, client *MockhttpClient) {
				for m := 1; m <= 12; m++ {
					key := fmt.Sprintf("workdays_2026_%d", m)
					cache.EXPECT().Get(key).Return(work_calendar.WorkdayResponse{}, false)

					resp := &work_calendar.WorkdayResponse{}
					respData, _ := json.Marshal(resp)
					client.EXPECT().Do(gomock.Any()).Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(respData)),
					}, nil)

					cache.EXPECT().Put(key, *resp).Return(fmt.Errorf("cache error"))
				}
			},
			wantCount: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			cache := NewMockcache(ctrl)
			client := NewMockhttpClient(ctrl)
			l := logging.New(false)

			service := work_calendar.New(client, cache, token, l)

			tt.setupMocks(cache, client)

			res, err := service.GetWorkdaysForYear(ctx, 2026)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, res)
				assert.Contains(t, err.Error(), tt.wantErrMsg)
			} else {
				assert.Nil(t, err)
				assert.Len(t, res, tt.wantCount)
			}
		})
	}
}

func TestClient_GetWorkdaysForMonth(t *testing.T) {
	tests := []struct {
		name       string
		month      int
		year       int
		setupMocks func(cache *Mockcache, client *MockhttpClient)
		wantErr    bool
	}{
		{
			name:  "cached",
			month: 1,
			year:  2024,
			setupMocks: func(cache *Mockcache, client *MockhttpClient) {
				cache.EXPECT().Get("workdays_2024_1").Return(work_calendar.WorkdayResponse{}, true)
			},
		},
		{
			name:  "fetch_success",
			month: 1,
			year:  2024,
			setupMocks: func(cache *Mockcache, client *MockhttpClient) {
				cache.EXPECT().Get("workdays_2024_1").Return(work_calendar.WorkdayResponse{}, false)
				resp := &work_calendar.WorkdayResponse{}
				respData, _ := json.Marshal(resp)
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(respData)),
				}, nil)
				cache.EXPECT().Put("workdays_2024_1", *resp).Return(nil)
			},
		},
		{
			name:    "invalid_month",
			month:   13,
			year:    2024,
			wantErr: true,
			setupMocks: func(cache *Mockcache, client *MockhttpClient) {
			},
		},
		{
			name:  "fetch_success_cache_put_error",
			month: 1,
			year:  2024,
			setupMocks: func(cache *Mockcache, client *MockhttpClient) {
				cache.EXPECT().Get("workdays_2024_1").Return(work_calendar.WorkdayResponse{}, false)
				resp := &work_calendar.WorkdayResponse{}
				respData, _ := json.Marshal(resp)
				client.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(respData)),
				}, nil)
				cache.EXPECT().Put("workdays_2024_1", *resp).Return(fmt.Errorf("cache error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			cache := NewMockcache(ctrl)
			client := NewMockhttpClient(ctrl)
			l := logging.New(false)
			service := work_calendar.New(client, cache, token, l)

			tt.setupMocks(cache, client)

			res, err := service.GetWorkdaysForMonth(ctx, tt.month, tt.year)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}

func TestClient_Request_Errors(t *testing.T) {
	t.Run("unexpected_status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockhttpClient(ctrl)
		l := logging.New(false)
		service := work_calendar.New(client, nil, token, l)

		client.EXPECT().Do(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewReader([]byte("error"))),
		}, nil)

		res, err := service.Request(context.Background(), "http://example.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code: 500")
		assert.Nil(t, res)
	})

	t.Run("invalid_json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockhttpClient(ctrl)
		l := logging.New(false)
		service := work_calendar.New(client, nil, token, l)

		client.EXPECT().Do(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		res, err := service.Request(context.Background(), "http://example.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "JSON parse error")
		assert.Nil(t, res)
	})
}
