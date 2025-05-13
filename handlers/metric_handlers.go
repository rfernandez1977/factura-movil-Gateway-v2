package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricHandlers maneja las rutas de métricas
type MetricHandlers struct {}

func NewMetricHandlers() *MetricHandlers {
    return &MetricHandlers{}
}

// PrometheusHandler expone las métricas para Prometheus
func (h *MetricHandlers) PrometheusHandler() gin.HandlerFunc {
    handler := promhttp.Handler()
    return func(c *gin.Context) {
        handler.ServeHTTP(c.Writer, c.Request)
    }
}