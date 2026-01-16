package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Word struct {
	gorm.Model
	Word        string         `gorm:"not null;index:idx_word_content,unique"`
	SourceUrls  datatypes.JSON `gorm:"type:jsonb"`
	ContentHash string         `gorm:"type:char(64);not null;index:idx_word_content,unique" json:"-"`
	Phonetics   []Phonetic     `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
	Meanings    []Meaning      `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
}
