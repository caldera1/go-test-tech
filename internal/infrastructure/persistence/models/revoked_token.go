package models

import "time"

type RevokedToken struct {
	ID        string `gorm:"primaryKey;type:varchar(36)"`
	RevokedAt time.Time
}
