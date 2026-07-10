package routes

import (
	"log"
	"net/http"

	"github.com/aipanel/aipanel/server/config"
	"github.com/aipanel/aipanel/server/middleware"
	aimodule "github.com/aipanel/aipanel/server/modules/ai"
	auditmodule "github.com/aipanel/aipanel/server/modules/audit"
	authmodule "github.com/aipanel/aipanel/server/modules/auth"
	dockermodule "github.com/aipanel/aipanel/server/modules/docker"
	logmodule "github.com/aipanel/aipanel/server/modules/log"
	systemmodule "github.com/aipanel/aipanel/server/modules/system"
	terminalmodule "github.com/aipanel/aipanel/server/modules/terminal"
	"github.com/aipanel/aipanel/server/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	api := router.Group("/api")

	authHandler := authmodule.NewHandler(db, cfg)
	api.POST("/auth/login", authHandler.Login)

	protected := api.Group("")
	protected.Use(middleware.AuthRequired(cfg.JWT.Secret))

	auditService := service.NewAuditService(db)

	systemHandler := systemmodule.NewHandler(cfg)
	protected.GET("/system/status", systemHandler.Status)
	protected.GET("/system/logs", logmodule.NewHandler().SystemLogs)

	dockerHandler, err := dockermodule.NewHandler(auditService)
	if err != nil {
		log.Printf("docker client init failed: %v", err)
	} else {
		protected.GET("/docker/containers", dockerHandler.Containers)
		protected.GET("/docker/logs/:id", dockerHandler.Logs)
		protected.POST("/docker/start/:id", dockerHandler.Start)
		protected.POST("/docker/stop/:id", dockerHandler.Stop)
		protected.POST("/docker/restart/:id", dockerHandler.Restart)
		protected.DELETE("/docker/remove/:id", dockerHandler.Remove)
	}

	aiHandler := aimodule.NewHandler(db, cfg)
	protected.GET("/ai/config", aiHandler.GetConfig)
	protected.POST("/ai/config", aiHandler.SaveConfig)
	protected.POST("/ai/test", aiHandler.Test)

	protected.GET("/audit/logs", auditmodule.NewHandler(db).Logs)

	api.GET("/terminal/ws", terminalmodule.NewHandler(cfg, auditService).WS)

	protected.GET("/settings", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "settings module is reserved for phase 2"})
	})
}
