package service

import (
	"github.com/aipanel/aipanel/server/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditService struct {
	db *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

func (s *AuditService) Record(c *gin.Context, action string, target string) {
	if s == nil || s.db == nil {
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	log := models.AuditLog{
		Action: action,
		Target: target,
		IP:     c.ClientIP(),
	}
	if id, ok := userID.(uint); ok {
		log.UserID = id
	}
	if name, ok := username.(string); ok {
		log.Username = name
	}
	if log.Username == "" {
		log.Username = "unknown"
	}

	_ = s.db.Create(&log).Error
}
