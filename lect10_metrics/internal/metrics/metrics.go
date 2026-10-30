package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// AppMetrics contains all demo metrics for system monitoring showcase.
type AppMetrics struct {
	PingRequestsTotal   *prometheus.CounterVec
	PingErrorsTotal     prometheus.Counter
	InFlightRequests    prometheus.Gauge
	LastPingUnix        prometheus.Gauge
	RequestDuration     *prometheus.HistogramVec
	ResponseSizeBytes   *prometheus.HistogramVec
	PingProcessing      *prometheus.SummaryVec
	PayloadSummaryBytes *prometheus.SummaryVec
}

func New(reg prometheus.Registerer) *AppMetrics {
	m := &AppMetrics{
		PingRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "demo",
				Name:      "ping_requests_total",
				Help:      "Total number of ping requests split by result.",
			},
			[]string{"result"},
		),
		PingErrorsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "demo",
				Name:      "ping_errors_total",
				Help:      "Total number of ping requests that ended with an error.",
			},
		),
		InFlightRequests: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "demo",
				Name:      "http_inflight_requests",
				Help:      "Current number of in-flight HTTP requests.",
			},
		),
		LastPingUnix: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "demo",
				Name:      "last_ping_unix_timestamp",
				Help:      "Unix timestamp of the last successful ping.",
			},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "demo",
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request latency in seconds.",
				Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.2, 0.5, 1, 2},
			},
			[]string{"path", "method", "status"},
		),
		ResponseSizeBytes: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "demo",
				Name:      "http_response_size_bytes",
				Help:      "HTTP response size in bytes.",
				Buckets:   []float64{64, 128, 256, 512, 1024, 2048, 4096},
			},
			[]string{"path", "method"},
		),
		PingProcessing: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  "demo",
				Name:       "ping_processing_seconds",
				Help:       "Ping business processing time in seconds.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"result"},
		),
		PayloadSummaryBytes: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  "demo",
				Name:       "ping_payload_size_bytes",
				Help:       "Ping response payload size in bytes.",
				Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01},
			},
			[]string{"result"},
		),
	}

	reg.MustRegister(
		m.PingRequestsTotal,
		m.PingErrorsTotal,
		m.InFlightRequests,
		m.LastPingUnix,
		m.RequestDuration,
		m.ResponseSizeBytes,
		m.PingProcessing,
		m.PayloadSummaryBytes,
	)

	return m
}

// Instrument wraps a handler with generic HTTP metrics.
func (m *AppMetrics) Instrument(path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.InFlightRequests.Inc()
		defer m.InFlightRequests.Dec()

		rw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		started := time.Now()

		next(rw, r)

		status := strconv.Itoa(rw.status)
		m.RequestDuration.WithLabelValues(path, r.Method, status).Observe(time.Since(started).Seconds())
		m.ResponseSizeBytes.WithLabelValues(path, r.Method).Observe(float64(rw.bytes))
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}
