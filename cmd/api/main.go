package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmlc643/twitbet-backend/internal/platform/config"
	"github.com/jmlc643/twitbet-backend/internal/platform/database"
	"github.com/jmlc643/twitbet-backend/internal/platform/redis"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Falló la inicialización de PostgreSQL: %v", err)
	}

	rdb, err := redis.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Falló la inicialización de Redis: %v", err)
	}

	router := gin.Default()

	router.SetTrustedProxies(nil)

	router.GET("/healthcheck", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		sqlDB, err := db.DB()
		postgresStatus := "connected"
		if err != nil || sqlDB.PingContext(ctx) != nil {
			postgresStatus = "disconnected"
		}

		redisStatus := "connected"
		if err := rdb.Ping(ctx).Err(); err != nil {
			redisStatus = "disconnected"
		}

		httpStatus := http.StatusOK
		if postgresStatus != "connected" || redisStatus != "connected" {
			httpStatus = http.StatusServiceUnavailable
		}

		c.JSON(httpStatus, gin.H{
			"status": "UP",
			"timestamp": time.Now().Format(time.RFC3339),
			"services": gin.H{
				"postgres": postgresStatus,
				"redis": redisStatus,
			},
		})
	})

	log.Printf("Servidor corriendo en el puerto %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Falló el inicio del servidor HTTP: %v", err)
	}
}