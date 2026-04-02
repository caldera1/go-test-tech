package models

import "time"

type Task struct {
	ID              string `gorm:"primaryKey;type:varchar(36)"`
	Title           string `gorm:"type:varchar(255);not null"`
	Description     string `gorm:"type:text"`
	DueDate         time.Time
	Status          string `gorm:"type:varchar(30);not null;index"`
	AssignedUserID  string `gorm:"type:varchar(36);not null;index"`
	CreatedByUserID string `gorm:"type:varchar(36);not null"`
	CreatedAt       time.Time

	AssignedUser  User `gorm:"foreignKey:AssignedUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser User `gorm:"foreignKey:CreatedByUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
