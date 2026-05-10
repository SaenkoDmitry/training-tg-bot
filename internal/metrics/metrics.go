package metrics

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"gorm.io/gorm"
)

var (
	// Системные (собираем сами через gopsutil)

	CpuUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "system_cpu_usage_percent",
		Help: "Current CPU usage percentage",
	}, []string{"mode"})

	MemUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "system_memory_usage_bytes",
		Help: "Current memory usage in bytes",
	}, []string{"type"})

	// Бизнес-метрики

	TotalUsers = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "app_users_total",
		Help: "Total number of registered users",
	})

	WorkoutsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "app_workouts_total",
		Help: "Total number of workouts created",
	}, []string{"status"})

	ApiRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "app_api_requests_total",
		Help: "Total number of API requests",
	}, []string{"method", "endpoint", "status"},
	)

	ApiRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "app_api_request_duration_seconds",
		Help:    "API request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "endpoint"})
)

func Init() *prometheus.Registry {
	registry := prometheus.NewRegistry()

	// Стандартные Go-метрики (GC, goroutines, memory)
	registry.MustRegister(collectors.NewGoCollector())
	// Процессы (RSS, CPU time, open FDs)
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// Наши кастомные
	registry.MustRegister(CpuUsage, MemUsage, TotalUsers, WorkoutsTotal, ApiRequestsTotal, ApiRequestDuration)

	return registry
}

// StartSystemMetricsCollector Запускаем сбор системных метрик в фоне
func StartSystemMetricsCollector() {
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for range ticker.C {
			// CPU
			if percent, err := cpu.Percent(0, false); err == nil && len(percent) > 0 {
				CpuUsage.WithLabelValues("total").Set(percent[0])
			}

			// Memory
			if vmStat, err := mem.VirtualMemory(); err == nil {
				MemUsage.WithLabelValues("used").Set(float64(vmStat.Used))
				MemUsage.WithLabelValues("total").Set(float64(vmStat.Total))
				MemUsage.WithLabelValues("available").Set(float64(vmStat.Available))
			}

			// Go runtime stats
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			// Уже собирается через NewGoCollector, но можно добавить кастомные
		}
	}()
}

func StartUsersCollector(db *gorm.DB) {
	updateTotalUsers(db)

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			updateTotalUsers(db)
		}
	}()
}

func updateTotalUsers(db *gorm.DB) {
	var count int64
	if err := db.Table("users").Count(&count).Error; err == nil {
		TotalUsers.Set(float64(count))
	}
}
