package ai

import (
	"errors"
	"net/http"

	"github.com/aipanel/aipanel/server/config"
	"github.com/aipanel/aipanel/server/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db  *gorm.DB
	cfg *config.Config
}

type configRequest struct {
	Name     string `json:"name"`
	Provider string `json:"provider" binding:"required"`
	BaseURL  string `json:"base_url" binding:"required"`
	APIKey   string `json:"api_key"`
	Model    string `json:"model" binding:"required"`
}

func NewHandler(db *gorm.DB, cfg *config.Config) *Handler {
	return &Handler{db: db, cfg: cfg}
}

func (h *Handler) GetConfig(c *gin.Context) {
	cfg, err := h.getStoredConfig()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, h.defaultConfig())
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get ai config"})
		return
	}

	c.JSON(http.StatusOK, cfg)
}

func (h *Handler) SaveConfig(c *gin.Context) {
	var req configRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider, base_url and model are required"})
		return
	}

	name := req.Name
	if name == "" {
		name = "default"
	}

	var existing models.AIConfig
	err := h.db.Order("id asc").First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save ai config"})
		return
	}

	existing.Name = name
	existing.Provider = req.Provider
	existing.BaseURL = req.BaseURL
	existing.APIKey = req.APIKey
	existing.Model = req.Model

	if existing.ID == 0 {
		if err := h.db.Create(&existing).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save ai config"})
			return
		}
	} else if err := h.db.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save ai config"})
		return
	}

	c.JSON(http.StatusOK, existing)
}

func (h *Handler) Test(c *gin.Context) {
	cfg, err := h.getStoredConfig()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			defaultCfg := h.defaultConfig()
			cfg = &models.AIConfig{
				Name:     defaultCfg.Name,
				Provider: defaultCfg.Provider,
				BaseURL:  defaultCfg.BaseURL,
				APIKey:   defaultCfg.APIKey,
				Model:    defaultCfg.Model,
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load ai config"})
			return
		}
	}

	provider := NewProvider(cfg.Provider, cfg.BaseURL, cfg.APIKey, cfg.Model)
	resp, err := provider.Chat([]Message{{Role: "user", Content: "ping"}})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "model": cfg.Model, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "model": resp.Model})
}

func (h *Handler) getStoredConfig() (*models.AIConfig, error) {
	var cfg models.AIConfig
	if err := h.db.Order("id asc").First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (h *Handler) defaultConfig() models.AIConfig {
	return models.AIConfig{
		Name:     "default",
		Provider: h.cfg.AI.Provider,
		BaseURL:  h.cfg.AI.BaseURL,
		APIKey:   h.cfg.AI.APIKey,
		Model:    h.cfg.AI.Model,
	}
}
