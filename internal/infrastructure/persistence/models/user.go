package models

import "time"

type User struct {
	ID                 string `gorm:"primaryKey;type:varchar(36)"`
	Email              string `gorm:"uniqueIndex;type:varchar(255);not null"`
	PasswordHash       string `gorm:"type:text;not null"`
	Role               string `gorm:"type:varchar(20);not null"`
	MustChangePassword bool   `gorm:"not null;default:true"`
	CreatedAt          time.Time
}
