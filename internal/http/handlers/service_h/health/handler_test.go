package health_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"salary_calculator/internal/http/handlers/service_h/health"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Healthy", func(t *testing.T) {
		mockChecker := NewMockChecker(ctrl)
		mockChecker.EXPECT().Name().Return("db").AnyTimes()
		mockChecker.EXPECT().HealthCheck(gomock.Any()).Return(nil)

		handler := health.New([]health.Checker{mockChecker}, 1*time.Second)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp health.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "healthy", resp.Status)
		assert.Len(t, resp.Services, 1)
		assert.Equal(t, "db", resp.Services[0].Name)
		assert.Equal(t, "healthy", resp.Services[0].Status)
	})

	t.Run("Unhealthy", func(t *testing.T) {
		mockChecker := NewMockChecker(ctrl)
		mockChecker.EXPECT().Name().Return("db").AnyTimes()
		mockChecker.EXPECT().HealthCheck(gomock.Any()).Return(errors.New("db down"))

		handler := health.New([]health.Checker{mockChecker}, 1*time.Second)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var resp health.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "unhealthy", resp.Status)
		assert.Len(t, resp.Services, 1)
		assert.Equal(t, "db", resp.Services[0].Name)
		assert.Equal(t, "unhealthy", resp.Services[0].Status)
		assert.Equal(t, "db down", resp.Services[0].Message)
	})
}

func TestBasicHealthChecker(t *testing.T) {
	checker := health.NewBasicHealthChecker()
	assert.Equal(t, "basic", checker.Name())

	err := checker.HealthCheck(context.Background())
	assert.NoError(t, err)
}
