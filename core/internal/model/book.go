package model

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title string `gorm:"type:text;not null;uniqueIndex"`
}
