package database

import (
	"errors"
	"log"

	"core/config"
	"core/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	if db == nil {
		return errors.New("database connection is not established")
	}

	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&model.Definition{},
		&model.Meaning{},
		&model.Phonetic{},
		&model.User{},
		&model.Word{},
	)
	if err != nil {
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}
