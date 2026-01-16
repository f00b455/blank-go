package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/f00b455/blank-go/internal/config"
	"github.com/f00b455/blank-go/internal/database"
	"github.com/f00b455/blank-go/internal/handlers"
	"github.com/f00b455/blank-go/pkg/dax"
	"github.com/f00b455/blank-go/pkg/stocks"
	"github.com/f00b455/blank-go/pkg/task"
	"github.com/f00b455/blank-go/pkg/weather"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	startTime := time.Now()
	cfg := config.Load()

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to PostgreSQL
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate DAX schema
	if err := dax.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	router := setupRouter(cfg, db, startTime)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting API server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(cfg *config.Config, db interface{}, startTime time.Time) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", handlers.HealthCheck)

	// Initialize task service and handler
	taskRepo := task.NewInMemoryRepository()
	taskService := task.NewService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	api := router.Group("/api/v1")
	{
		api.GET("/ping", handlers.Ping)
		api.GET("/health/detailed", handlers.DetailedHealthCheck(startTime))

		// Task routes
		api.POST("/tasks", taskHandler.CreateTask)
		api.GET("/tasks", taskHandler.ListTasks)
		api.GET("/tasks/:id", taskHandler.GetTask)
		api.PUT("/tasks/:id", taskHandler.UpdateTask)
		api.DELETE("/tasks/:id", taskHandler.DeleteTask)

		// DAX routes (only if database is available)
		if db != nil {
			daxRepo := dax.NewPostgresRepository(db.(*gorm.DB))
			daxService := dax.NewService(daxRepo)
			daxHandler := handlers.NewDAXHandler(daxService)

			daxGroup := api.Group("/dax")
			{
				daxGroup.POST("/import", daxHandler.ImportCSV)
				daxGroup.GET("", daxHandler.GetByFilters)
				daxGroup.GET("/metrics", daxHandler.GetMetrics)
			}
		}

		// Weather routes
		weatherClient := weather.NewClient()
		weatherService := weather.NewService(weatherClient)
		weatherHandler := handlers.NewWeatherHandler(weatherService)

		api.GET("/weather", weatherHandler.GetCurrentWeather)
		api.GET("/weather/forecast", weatherHandler.GetForecast)
		api.GET("/weather/cities/:city", weatherHandler.GetWeatherByCity)

		// Stocks routes
		stocksClient := stocks.NewClient()
		stocksService := stocks.NewService(stocksClient)
		stocksHandler := handlers.NewStocksHandler(stocksService)

		stocksGroup := api.Group("/stocks")
		{
			stocksGroup.GET("/:ticker/summary", stocksHandler.GetStockSummary)
			stocksGroup.GET("/summary", stocksHandler.GetBatchSummary)
		}
	}

	return router
}
