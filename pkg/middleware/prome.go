package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// DefaultBuckets prometheus buckets in seconds
var DefaultBuckets = []float64{0.5, 1.0, 5.0}

const (
	reqsMetricName    = "http_request_total"
	latencyMetricName = "http_request_duration_seconds"
)

// Metrics in prometheus defination
type Metrics struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

// NewProm returns a new prometheus middleware
func NewProm(app string) *Metrics {
	m := Metrics{}

	m.reqs = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        reqsMetricName,
		Help:        "How many http requests processed, with labels status code, method and path.",
		ConstLabels: prometheus.Labels{"app": app},
	},
		[]string{"code", "method", "path"},
	)

	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyMetricName,
		Help:        "How long it took to process the request, with labels status code, method and path.",
		ConstLabels: prometheus.Labels{"app": app},
		Buckets:     DefaultBuckets,
	},
		[]string{"code", "method", "path"},
	)

	prometheus.MustRegister(m.reqs, m.latency)

	return &m
}

// Register router of this middleware
func (m *Metrics) Register(app *gin.Engine) {
	app.Use(m.Run())
	app.GET("/metrics", func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})
}

// Run start this middleware with context
func (m *Metrics) Run() gin.HandlerFunc {
	return m.serve
}

func (m *Metrics) serve(ctx *gin.Context) {
	if ctx.Request.URL.Path == "/metrics" {
		ctx.Next()
		return
	}

	start := time.Now()
	req := ctx.Request
	code := strconv.Itoa(ctx.Writer.Status())

	ctx.Next()

	m.reqs.WithLabelValues(code, req.Method, req.URL.Path).Inc()
	m.latency.WithLabelValues(code, req.Method, req.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1e9)
}
