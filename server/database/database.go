package database

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/aipanel/aipanel/server/models"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Open(path string) (*gorm.DB, error) {
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	return gorm.Open(sqlite.Open(path), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.AIConfig{}, &models.AuditLog{})
}

func SeedAdmin(db *gorm.DB, username string, password string) error {
	var existing models.User
	err := db.Where("username = ?", username).First(&existing).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return db.Create(&models.User{
		Username:     username,
		PasswordHash: string(hash),
	}).Error
}
