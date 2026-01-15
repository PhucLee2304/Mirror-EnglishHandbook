package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Definition struct {
	gorm.Model
	DefinitionText string `gorm:"not null;type:text"`
	Example        *string
	Antonyms       datatypes.JSON `gorm:"type:jsonb"`
	Synonyms       datatypes.JSON `gorm:"type:jsonb"`
	MeaningID      uint           `gorm:"not null;index"`
	Meaning        Meaning        `gorm:"foreignKey:MeaningID;constraint:OnDelete:CASCADE;"`
}
