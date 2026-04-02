package models

import "time"

type Comment struct {
	ID        string `gorm:"primaryKey;type:varchar(36)"`
	TaskID    string `gorm:"type:varchar(36);not null;index"`
	AuthorID  string `gorm:"type:varchar(36);not null"`
	Body      string `gorm:"type:text;not null"`
	CreatedAt time.Time

	Task Task `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
