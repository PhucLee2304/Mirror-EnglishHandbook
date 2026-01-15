package main

import (
	"log"
	"net/http"

	"core/config"
	"core/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if cfg.AppMode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// wordRepo := repo.NewWordRepo(db)
	// wordService := service.NewWordService(wordRepo)
	// wordHandler := handler.NewWordHandler(wordService)

	// api := r.Group("/api")
	// wordHandler.Register(api)

	addr := ":" + cfg.AppPort
	log.Printf("Starting server on %s in %s mode", addr, cfg.AppMode)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
