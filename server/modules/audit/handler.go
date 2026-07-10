package audit

import (
	"net/http"
	"strconv"

	"github.com/aipanel/aipanel/server/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Logs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	var logs []models.AuditLog
	if err := h.db.Order("created_at desc").Limit(limit).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get audit logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}
