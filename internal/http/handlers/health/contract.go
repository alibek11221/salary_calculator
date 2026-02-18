package health

import (
	"context"
	"fmt"
	"time"
)

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

type Checker interface {
	HealthCheck(ctx context.Context) error
	Name() string
}

type Status struct {
	Name    string    `json:"name"`
	Status  string    `json:"status"`
	Message string    `json:"message,omitempty"`
	Time    time.Time `json:"timestamp"`
}

type SystemInfo struct {
	Version      string `json:"go_version"`
	NumCPU       int    `json:"cpu_cores"`
	NumGoroutine int    `json:"goroutines"`
	Uptime       string `json:"uptime,omitempty"`
}

type MemoryInfo struct {
	Alloc        string `json:"alloc"`
	TotalAlloc   string `json:"total_alloc"`
	Sys          string `json:"sys"`
	HeapAlloc    string `json:"heap_alloc"`
	HeapSys      string `json:"heap_sys"`
	HeapIdle     string `json:"heap_idle"`
	HeapInuse    string `json:"heap_inuse"`
	HeapReleased string `json:"heap_released"`
	StackInuse   string `json:"stack_inuse"`
	StackSys     string `json:"stack_sys"`
	GCSys        string `json:"gc_sys"`
	NextGC       string `json:"next_gc"`

	Lookups     uint64 `json:"lookups"`
	Mallocs     uint64 `json:"mallocs"`
	Frees       uint64 `json:"frees"`
	HeapObjects uint64 `json:"heap_objects"`
	LastGC      uint64 `json:"last_gc_timestamp"`
	NumGC       uint32 `json:"gc_cycles"`
}

type Response struct {
	Status   string     `json:"status"`
	Services []Status   `json:"services,omitempty"`
	System   SystemInfo `json:"system"`
	Memory   MemoryInfo `json:"memory"`
	Time     time.Time  `json:"timestamp"`
}
