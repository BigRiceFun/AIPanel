package routes

import (
	"log"
	"net/http"

	"github.com/aipanel/aipanel/server/config"
	"github.com/aipanel/aipanel/server/middleware"
	authmodule "github.com/aipanel/aipanel/server/modules/auth"
	dockermodule "github.com/aipanel/aipanel/server/modules/docker"
	systemmodule "github.com/aipanel/aipanel/server/modules/system"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	api := router.Group("/api")

	authHandler := authmodule.NewHandler(db, cfg)
	api.POST("/auth/login", authHandler.Login)

	protected := api.Group("")
	protected.Use(middleware.AuthRequired(cfg.JWT.Secret))

	systemHandler := systemmodule.NewHandler(cfg)
	protected.GET("/system/status", systemHandler.Status)

	dockerHandler, err := dockermodule.NewHandler()
	if err != nil {
		log.Printf("docker client init failed: %v", err)
	} else {
		protected.GET("/docker/containers", dockerHandler.Containers)
		protected.POST("/docker/start/:id", dockerHandler.Start)
		protected.POST("/docker/stop/:id", dockerHandler.Stop)
		protected.POST("/docker/restart/:id", dockerHandler.Restart)
		protected.DELETE("/docker/remove/:id", dockerHandler.Remove)
	}

	protected.GET("/settings", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "settings module is reserved for phase 2"})
	})
}
