package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
)

type Handler struct {
	checkers []Checker
	timeout  time.Duration
}

func New(checkers []Checker, timeout time.Duration) *Handler {
	if timeout == 0 {
		timeout = 5 * time.Second // default timeout
	}

	return &Handler{
		checkers: checkers,
		timeout:  timeout,
	}
}

// ServeHTTP godoc
// @Summary      Проверка здоровья
// @Description  Возвращает текущее состояние сервиса и системные метрики
// @Tags         health
// @Produce      json
// @Success      200  {object}  Response
// @Failure      503  {object}  Response
// @Router       /health/health [get]
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	response := Response{
		Status:   "healthy",
		Services: make([]Status, 0, len(h.checkers)),
		System: SystemInfo{
			Version:      runtime.Version(),
			NumCPU:       runtime.NumCPU(),
			NumGoroutine: runtime.NumGoroutine(),
		},
		Memory: MemoryInfo{
			Alloc:        formatBytes(m.Alloc),
			TotalAlloc:   formatBytes(m.TotalAlloc),
			Sys:          formatBytes(m.Sys),
			HeapAlloc:    formatBytes(m.HeapAlloc),
			HeapSys:      formatBytes(m.HeapSys),
			HeapIdle:     formatBytes(m.HeapIdle),
			HeapInuse:    formatBytes(m.HeapInuse),
			HeapReleased: formatBytes(m.HeapReleased),
			StackInuse:   formatBytes(m.StackInuse),
			StackSys:     formatBytes(m.StackSys),
			GCSys:        formatBytes(m.GCSys),
			NextGC:       formatBytes(m.NextGC),

			Lookups:     m.Lookups,
			Mallocs:     m.Mallocs,
			Frees:       m.Frees,
			HeapObjects: m.HeapObjects,
			LastGC:      m.LastGC,
			NumGC:       m.NumGC,
		},
		Time: time.Now(),
	}

	for _, checker := range h.checkers {
		status := Status{
			Name: checker.Name(),
			Time: time.Now(),
		}

		err := checker.HealthCheck(ctx)
		if err != nil {
			status.Status = "unhealthy"
			status.Message = err.Error()
			response.Status = "unhealthy"

			log.Warn().
				Str("service", checker.Name()).
				Err(err).
				Msg("Health check failed")
		} else {
			status.Status = "healthy"
		}

		response.Services = append(response.Services, status)
	}

	statusCode := http.StatusOK
	if response.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error().Err(err).Msg("Failed to encode health response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

type BasicHealthChecker struct{}

func NewBasicHealthChecker() *BasicHealthChecker {
	return &BasicHealthChecker{}
}

func (b *BasicHealthChecker) HealthCheck(ctx context.Context) error {
	_ = make([]byte, 1024)

	if runtime.NumGoroutine() > 10000 {
		return fmt.Errorf("too many goroutines: %d", runtime.NumGoroutine())
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	const maxAllocBytes = 1 << 30
	if m.Alloc > maxAllocBytes {
		return fmt.Errorf("memory allocation too high: %d bytes (limit: %d bytes)", m.Alloc, maxAllocBytes)
	}

	heapUsageRatio := float64(m.HeapAlloc) / float64(m.HeapSys)
	if heapUsageRatio > 0.9 {
		return fmt.Errorf("heap usage too high: %.2f%%", heapUsageRatio*100)
	}

	return nil
}

func (b *BasicHealthChecker) Name() string {
	return "basic"
}
