package model

import (
	"gorm.io/gorm"
)

type Phonetic struct {
	gorm.Model
	Text      string `gorm:"not null"`
	Audio     *string
	SourceUrl *string
	Order     int  `gorm:"not null"`
	WordID    uint `gorm:"not null;index"`
	Word      Word `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
}
