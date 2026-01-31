package model

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title   string   `gorm:"type:text;not null;uniqueIndex"`
	Lessons []Lesson `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;"`
}
