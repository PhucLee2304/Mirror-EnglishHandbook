package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null;uniqueIndex:idx_users_email_when_not_deleted,where:deleted_at IS NULL"`
	Avatar       *string
	Platform     string        `gorm:"not null"`
	DeviceTokens []DeviceToken `gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
}
