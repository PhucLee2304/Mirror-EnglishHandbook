package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"core/config"
	"core/database"
	"core/internal/seed"
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

	if cfg.DataJsonPath == "" {
		log.Fatalf("DATA_JSON_PATH is not set")
	}
	entries, err := os.ReadDir(cfg.DataJsonPath)
	if err != nil {
		log.Fatalf("Failed to read data directory: %v", err)
	}
	log.Printf("Found %d files in %s", len(entries), cfg.DataJsonPath)

	for _, entry := range entries {
		log.Printf("Found file: %s", entry.Name())
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		path := filepath.Join(cfg.DataJsonPath, entry.Name())
		log.Printf("Seeding data from file: %s", path)

		if err := seed.LoadFile(db, path); err != nil {
			log.Fatalf("Failed to load seed file %s: %v", path, err)
		}
	}
	log.Println("Seeding completed successfully")
}
