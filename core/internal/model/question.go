package model

import (
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	Content   string  `gorm:"type:text;not null"`
	TimeStart float64 `gorm:"type:float;not null"`
	TimeEnd   float64 `gorm:"type:float;not null"`
	Order     int     `gorm:"type:int;not null"`
	LessonID  uint    `gorm:"not null;index"`
	Lesson    Lesson  `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE;"`
}
