package httpapi

import (
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"system-monitoring-demo/internal/metrics"
)

type Server struct {
	metrics *metrics.AppMetrics
}

func New(m *metrics.AppMetrics) *Server {
	return &Server{metrics: m}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", s.metrics.Instrument("/healthz", s.healthz))
	mux.HandleFunc("/api/v1/ping", s.metrics.Instrument("/api/v1/ping", s.ping))

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/swagger/doc.json", s.swaggerDoc)
	mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("list"),
		),
	)

	return mux
}

func (s *Server) healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	q := r.URL.Query()
	name := q.Get("name")
	if name == "" {
		name = "monitoring"
	}

	if q.Get("fail") == "true" {
		s.metrics.PingRequestsTotal.WithLabelValues("error").Inc()
		s.metrics.PingErrorsTotal.Inc()
		s.metrics.PingProcessing.WithLabelValues("error").Observe(time.Since(started).Seconds())
		http.Error(w, `{"error":"forced failure"}`+"\n", http.StatusInternalServerError)
		return
	}

	// random latency helps to visualize histogram/summary in Prometheus queries
	time.Sleep(time.Duration(20+rand.IntN(180)) * time.Millisecond)

	response := map[string]any{
		"message":     "pong, " + name,
		"server_time": time.Now().UTC().Format(time.RFC3339),
		"duration_ms": time.Since(started).Milliseconds(),
	}
	payload, err := json.Marshal(response)
	if err != nil {
		s.metrics.PingRequestsTotal.WithLabelValues("error").Inc()
		s.metrics.PingErrorsTotal.Inc()
		http.Error(w, `{"error":"encode failed"}`+"\n", http.StatusInternalServerError)
		return
	}

	s.metrics.PingRequestsTotal.WithLabelValues("success").Inc()
	s.metrics.LastPingUnix.Set(float64(time.Now().Unix()))
	s.metrics.PingProcessing.WithLabelValues("success").Observe(time.Since(started).Seconds())
	s.metrics.PayloadSummaryBytes.WithLabelValues("success").Observe(float64(len(payload)))

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Demo-Request", strconv.FormatInt(time.Now().UnixNano(), 10))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}

func (s *Server) swaggerDoc(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(swaggerDocJSON))
}
