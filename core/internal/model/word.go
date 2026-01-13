package model

import (
	"gorm.io/gorm"
)

type Word struct {
	gorm.Model
	Word       string     `gorm:"not null;index"`
	SourceUrls []string   `gorm:"serialize:json"`
	Phonetics  []Phonetic `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
	Meanings   []Meaning  `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
}
