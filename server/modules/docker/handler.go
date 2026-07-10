package docker

import (
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	client *client.Client
}

type ContainerResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

func NewHandler() (*Handler, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Handler{client: cli}, nil
}

func (h *Handler) Containers(c *gin.Context) {
	containers, err := h.client.ContainerList(c.Request.Context(), container.ListOptions{All: true})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	result := make([]ContainerResponse, 0, len(containers))
	for _, item := range containers {
		result = append(result, ContainerResponse{
			ID:     shortID(item.ID),
			Name:   firstName(item.Names),
			Image:  item.Image,
			Status: item.State,
		})
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Start(c *gin.Context) {
	if err := h.client.ContainerStart(c.Request.Context(), c.Param("id"), container.StartOptions{}); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Stop(c *gin.Context) {
	timeout := 10
	if err := h.client.ContainerStop(c.Request.Context(), c.Param("id"), container.StopOptions{Timeout: &timeout}); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Restart(c *gin.Context) {
	timeout := 10
	if err := h.client.ContainerRestart(c.Request.Context(), c.Param("id"), container.StopOptions{Timeout: &timeout}); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	time.Sleep(300 * time.Millisecond)
	c.Status(http.StatusNoContent)
}

func (h *Handler) Remove(c *gin.Context) {
	if err := h.client.ContainerRemove(c.Request.Context(), c.Param("id"), container.RemoveOptions{Force: false}); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func shortID(id string) string {
	if len(id) <= 12 {
		return id
	}
	return id[:12]
}

func firstName(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return strings.TrimPrefix(names[0], "/")
}
