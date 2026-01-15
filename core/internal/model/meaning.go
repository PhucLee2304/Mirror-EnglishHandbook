package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Meaning struct {
	gorm.Model
	PartOfSpeech string         `gorm:"not null;index"`
	Synonyms     datatypes.JSON `gorm:"type:jsonb"`
	Antonyms     datatypes.JSON `gorm:"type:jsonb"`
	Definitions  []Definition   `gorm:"foreignKey:MeaningID;constraint:OnDelete:CASCADE;"`
	WordID       uint           `gorm:"not null;index"`
	Word         Word           `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
}
