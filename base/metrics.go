package base

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tiant-developer/go-tiant/zlog"
)

var PromMonitor *PrometheusMonitor

type PrometheusMonitor struct {
	APIRequestsCounter *prometheus.CounterVec
	RequestDuration    *prometheus.HistogramVec
}

func RegistryMetrics(engine *gin.Engine, appName string) {
	APIRequestsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "",
			Subsystem: appName,
			Name:      "http_requests_total",
			Help:      "A counter for requests to the wrapped handler.",
		},
		[]string{"handler", "appName", "method", "code"},
	)

	RequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "",
			Subsystem: appName,
			Name:      "http_request_duration_milliseconds",
			Help:      "A histogram of latencies for requests.",
			Buckets:   []float64{1, 5, 30, 50, 70, 100, 150, 200, 300, 500, 1000, 3000},
		},
		[]string{"handler", "method", "code", "appName"},
	)
	PromMonitor = &PrometheusMonitor{
		APIRequestsCounter: APIRequestsCounter,
		RequestDuration:    RequestDuration,
	}

	runtimeMetricsRegister := prometheus.NewRegistry()
	runtimeMetricsRegister.MustRegister(collectors.NewGoCollector())
	runtimeMetricsRegister.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	runtimeMetricsRegister.MustRegister(APIRequestsCounter, RequestDuration)

	engine.GET("/metrics", func(ctx *gin.Context) {
		// 避免metrics打点输出过多无用日志
		zlog.SetNoLogFlag(ctx)

		httpHandler := promhttp.InstrumentMetricHandler(
			runtimeMetricsRegister, promhttp.HandlerFor(runtimeMetricsRegister, promhttp.HandlerOpts{}),
		)
		httpHandler.ServeHTTP(ctx.Writer, ctx.Request)
	})
}
