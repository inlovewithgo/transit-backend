package metrics

import (
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
)

var (
    httpRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"path", "method", "status"},
    )

    httpDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )

    httpActiveRequests = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "http_active_requests",
            Help: "Number of currently active HTTP requests",
        },
        []string{"method", "path"},
    )
)

func Init() {
    prometheus.MustRegister(httpRequests, httpDuration, httpActiveRequests)
}

func GinMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        path := c.FullPath()
        if path == "" {
            path = c.Request.URL.Path
        }
        
        method := c.Request.Method
        
        httpActiveRequests.WithLabelValues(method, path).Inc()
        c.Next()
        
        duration := time.Since(start)
        
        status := strconv.Itoa(c.Writer.Status())
        
        httpRequests.WithLabelValues(path, method, status).Inc()
        httpDuration.WithLabelValues(method, path).Observe(duration.Seconds())
        
        httpActiveRequests.WithLabelValues(method, path).Dec()
    }
}

func GetMetrics() (map[string]interface{}, error) {
    metricFamilies, err := prometheus.DefaultGatherer.Gather()
    if err != nil {
        return nil, err
    }
    
    metrics := make(map[string]interface{})
    for _, mf := range metricFamilies {
        metrics[*mf.Name] = mf
    }
    
    return metrics, nil
}

func ResetMetrics() {
    httpRequests.Reset()
    httpDuration.Reset()
    httpActiveRequests.Reset()
}