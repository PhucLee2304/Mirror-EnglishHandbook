package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"core/config"
	"core/database"

	"core/internal/handler"
	"core/internal/repo"
	"core/internal/service"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
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

	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			log.Fatalf("failed to create logs directory: %v", err)
		}
	}

	writer, err := rotatelogs.New(
		"logs/access_%Y-%m-%d.log",
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(7*24*time.Hour),
	)
	if err != nil {
		log.Fatalf("failed to create rotatelogs: %v", err)
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout, writer)

	if cfg.AppMode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	wordRepo := repo.NewWordRepository(db)

	wordService := service.NewWordService(cfg, wordRepo)

	wordHandler := handler.NewWordHandler(wordService)
	wordHandler.SetupRouter(api, cfg)

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.AppPort
	}
	addr := fmt.Sprintf(":%s", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	log.Printf("Dictionary server started on %s", addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)

	_ = writer.Close()
	log.Println("Server stopped gracefully")
}
