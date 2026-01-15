package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

const Version = "1.0.0"

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func DetailedHealthCheck(startTime time.Time) gin.HandlerFunc {
	return func(c *gin.Context) {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		uptime := time.Since(startTime).Seconds()

		c.JSON(http.StatusOK, gin.H{
			"status":         "healthy",
			"timestamp":      time.Now().UTC().Format(time.RFC3339),
			"version":        Version,
			"uptime_seconds": uptime,
			"system": gin.H{
				"go_version":      runtime.Version(),
				"goroutines":      runtime.NumGoroutine(),
				"memory_alloc_mb": float64(memStats.Alloc) / 1024 / 1024,
				"memory_sys_mb":   float64(memStats.Sys) / 1024 / 1024,
				"gc_runs":         memStats.NumGC,
			},
			"checks": gin.H{
				"api": "ok",
			},
		})
	}
}
