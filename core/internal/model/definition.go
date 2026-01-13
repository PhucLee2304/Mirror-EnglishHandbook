package model

import (
	"gorm.io/gorm"
)

type Definition struct {
	gorm.Model
	DefinitionText string `gorm:"not null;type:text"`
	Example        *string
	Antonyms       []string `gorm:"serialize:json"`
	Synonyms       []string `gorm:"serialize:json"`
	MeaningID      uint     `gorm:"not null;index"`
	Meaning        Meaning  `gorm:"foreignKey:MeaningID;constraint:OnDelete:CASCADE;"`
}
