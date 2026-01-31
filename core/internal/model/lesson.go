package model

import (
	"gorm.io/gorm"
)

type Lesson struct {
	gorm.Model
	Title       string `gorm:"type:text;not null"`
	Description string `gorm:"type:text;not null"`
}
