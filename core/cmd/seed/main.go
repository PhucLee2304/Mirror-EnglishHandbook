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

	if cfg.HandbookDataJsonPath == "" {
		log.Fatalf("HANDBOOK_DATA_JSON_PATH is not set")
		return
	}
	entries, err := os.ReadDir(cfg.HandbookDataJsonPath)
	if err != nil {
		log.Fatalf("Failed to read data directory: %v", err)
		return
	}
	log.Printf("Found %d files in %s", len(entries), cfg.HandbookDataJsonPath)

	for _, entry := range entries {
		log.Printf("Found file: %s", entry.Name())
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		path := filepath.Join(cfg.HandbookDataJsonPath, entry.Name())
		log.Printf("Seeding data from file: %s", path)

		if err := seed.LoadHandbookFile(db, path); err != nil {
			log.Fatalf("Failed to load handbook seed file %s: %v", path, err)
			return
		}
	}

	if cfg.BookDataJsonPath == "" {
		log.Fatalf("BOOK_DATA_JSON_PATH is not set")
		return
	}
	bookEntries, err := os.ReadDir(cfg.BookDataJsonPath)
	if err != nil {
		log.Fatalf("Failed to read data directory: %v", err)
		return
	}
	log.Printf("Found %d files in %s", len(bookEntries), cfg.BookDataJsonPath)
	for _, entry := range bookEntries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		path := filepath.Join(cfg.BookDataJsonPath, entry.Name())
		log.Printf("Seeding data from file: %s", path)
		if err := seed.LoadBookFile(db, path); err != nil {
			log.Fatalf("Failed to load book seed file %s: %v", path, err)
			return
		}
	}
	log.Println("Seeding completed successfully")
}
