package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/f00b455/blank-go/internal/version"
	"github.com/gin-gonic/gin"
)

const (
	bytesToKB = 1024
	bytesToMB = bytesToKB * 1024
)

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
			"version":        version.Version,
			"uptime_seconds": uptime,
			"system": gin.H{
				"go_version":      runtime.Version(),
				"goroutines":      runtime.NumGoroutine(),
				"memory_alloc_mb": float64(memStats.Alloc) / bytesToMB,
				"memory_sys_mb":   float64(memStats.Sys) / bytesToMB,
				"gc_runs":         memStats.NumGC,
			},
			"checks": gin.H{
				"api": "ok",
			},
		})
	}
}
