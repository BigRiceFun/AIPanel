package models

import "time"

type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Username  string    `gorm:"size:64;not null" json:"username"`
	Action    string    `gorm:"size:128;not null" json:"action"`
	Target    string    `gorm:"size:255" json:"target"`
	IP        string    `gorm:"size:64" json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}
