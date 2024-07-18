package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tiant-developer/go-tiant/base"
	"time"
)

func PromMiddleware(appName string) gin.HandlerFunc {
	m := base.PromMonitor
	return func(c *gin.Context) {
		relativePath := c.Request.URL.Path
		start := time.Now()
		c.Next()
		code := fmt.Sprintf("%d", c.Writer.Status())
		m.APIRequestsCounter.With(prometheus.Labels{"handler": relativePath, "appName": appName, "method": c.Request.Method, "code": code}).Inc()
		m.RequestDuration.With(prometheus.Labels{"handler": relativePath, "appName": appName, "method": c.Request.Method, "code": code}).Observe(getRequestCost(start, time.Now()))
	}
}

func getRequestCost(start, end time.Time) float64 {
	return float64(end.Sub(start).Nanoseconds()/1e4) / 100.0
}
