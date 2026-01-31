package model

import "gorm.io/gorm"

type DeviceToken struct {
	gorm.Model
	Token    string `gorm:"not null;uniqueIndex:idx_device_tokens_when_not_deleted,where:deleted_at IS NULL"`
	Platform string `gorm:"not null"`
	UserID   uint   `gorm:"not null;index:idx_device_tokens_user_id"`
	User     User   `gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
}
