package observability

import (
	"testing"
	"time"
)

func TestMetricsAndCounters(t *testing.T) {
	InitMetrics()
	RecordRequest("GET", "/healthz", 200, 10*time.Millisecond)
	RecordRequest("GET", "/unknown", 999, 5*time.Millisecond)
	RecordLogin(true)
	RecordLogin(false)
	ObserveApprovalDuration(2 * time.Minute)
	ObserveApprovalDuration(-1 * time.Second)
}
