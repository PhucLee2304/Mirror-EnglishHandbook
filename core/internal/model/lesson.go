package model

import (
	"gorm.io/gorm"
)

type Lesson struct {
	gorm.Model
	Title       string     `gorm:"type:text;not null"`
	Description string     `gorm:"type:text;not null"`
	IsVideo     bool       `gorm:"type:boolean;not null"`
	AudioURL    string     `gorm:"type:text;not null"`
	BookID      uint       `gorm:"not null;index"`
	Book        Book       `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;"`
	Questions   []Question `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE;"`
}
