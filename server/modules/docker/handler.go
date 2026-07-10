package docker

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aipanel/aipanel/server/service"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	client *client.Client
	audit  *service.AuditService
}

type ContainerResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

type LogsResponse struct {
	Container string   `json:"container"`
	Logs      []string `json:"logs"`
}

func NewHandler(audit *service.AuditService) (*Handler, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Handler{client: cli, audit: audit}, nil
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
	h.audit.Record(c, "docker.start", c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h *Handler) Stop(c *gin.Context) {
	timeout := 10
	if err := h.client.ContainerStop(c.Request.Context(), c.Param("id"), container.StopOptions{Timeout: &timeout}); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	h.audit.Record(c, "docker.stop", c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h *Handler) Restart(c *gin.Context) {
	timeout := 10
	if err := h.client.ContainerRestart(c.Request.Context(), c.Param("id"), container.StopOptions{Timeout: &timeout}); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	time.Sleep(300 * time.Millisecond)
	h.audit.Record(c, "docker.restart", c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h *Handler) Remove(c *gin.Context) {
	if err := h.client.ContainerRemove(c.Request.Context(), c.Param("id"), container.RemoveOptions{Force: false}); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	h.audit.Record(c, "docker.remove", c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h *Handler) Logs(c *gin.Context) {
	id := c.Param("id")
	follow := c.Query("follow") == "true"

	ctx := c.Request.Context()
	if follow {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 8*time.Second)
		defer cancel()
	}

	reader, err := h.client.ContainerLogs(ctx, id, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     follow,
		Tail:       "100",
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	defer reader.Close()

	logs, err := decodeDockerLogs(reader)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LogsResponse{Container: id, Logs: logs})
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

func decodeDockerLogs(reader io.Reader) ([]string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdout, &stderr, reader); err != nil {
		return nil, err
	}

	text := strings.TrimSpace(stdout.String() + stderr.String())
	if text == "" {
		return []string{}, nil
	}

	lines := strings.Split(text, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines, nil
}
