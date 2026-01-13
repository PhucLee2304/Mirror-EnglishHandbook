package model

import (
	"gorm.io/gorm"
)

type Meaning struct {
	gorm.Model
	PartOfSpeech string       `gorm:"not null;index"`
	Synonyms     []string     `gorm:"serialize:json"`
	Antonyms     []string     `gorm:"serialize:json"`
	Definitions  []Definition `gorm:"foreignKey:MeaningID;constraint:OnDelete:CASCADE;"`
	WordID       uint         `gorm:"not null;index"`
	Word         Word         `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
}
