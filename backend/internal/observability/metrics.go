package observability

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricsOnce sync.Once

	requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "saas_api_requests_total",
			Help: "Total number of API requests",
		},
		[]string{"method", "path", "status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "saas_api_request_duration_seconds",
			Help:    "API request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
	loginTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "saas_login_total",
			Help: "Total login attempts",
		},
		[]string{"result"},
	)
	approvalDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "saas_task_approval_duration_seconds",
			Help:    "Time from task creation to approval",
			Buckets: []float64{60, 300, 900, 1800, 3600, 7200, 14400, 28800, 86400},
		},
	)
)

func InitMetrics() {
	metricsOnce.Do(func() {
		prometheus.MustRegister(requestTotal, requestDuration, loginTotal, approvalDuration)
	})
}

func MetricsHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func RecordRequest(method, path string, status int, duration time.Duration) {
	statusLabel := http.StatusText(status)
	if statusLabel == "" {
		statusLabel = "unknown"
	}
	requestTotal.WithLabelValues(method, path, statusLabel).Inc()
	requestDuration.WithLabelValues(method, path, statusLabel).Observe(duration.Seconds())
}

func RecordLogin(success bool) {
	if success {
		loginTotal.WithLabelValues("success").Inc()
		return
	}
	loginTotal.WithLabelValues("failure").Inc()
}

func ObserveApprovalDuration(duration time.Duration) {
	if duration < 0 {
		return
	}
	approvalDuration.Observe(duration.Seconds())
}
