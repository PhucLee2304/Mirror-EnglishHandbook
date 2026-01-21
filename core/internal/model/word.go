package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Word struct {
	gorm.Model
	Word       string         `gorm:"type:text;not null;uniqueIndex"`
	SourceUrls datatypes.JSON `gorm:"type:jsonb"`
	Phonetics  []Phonetic     `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
	Meanings   []Meaning      `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
}
