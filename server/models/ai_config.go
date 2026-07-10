package models

import "time"

type AIConfig struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Provider  string    `gorm:"size:64;not null" json:"provider"`
	BaseURL   string    `gorm:"size:255;not null" json:"base_url"`
	APIKey    string    `gorm:"size:255" json:"api_key"`
	Model     string    `gorm:"size:128;not null" json:"model"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
