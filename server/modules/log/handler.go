package log

import (
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

type LogsResponse struct {
	Type string   `json:"type"`
	Logs []string `json:"logs"`
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) SystemLogs(c *gin.Context) {
	logType := c.DefaultQuery("type", "system")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	logs, err := readSystemLogs(logType, limit)
	if err != nil {
		c.JSON(http.StatusOK, LogsResponse{Type: logType, Logs: []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, LogsResponse{Type: logType, Logs: logs})
}

func readSystemLogs(logType string, limit int) ([]string, error) {
	if runtime.GOOS != "linux" {
		return []string{"system logs are only available on Linux hosts"}, nil
	}

	args := []string{"-n", strconv.Itoa(limit), "--no-pager", "--output=short-iso"}
	if logType == "kernel" {
		args = append([]string{"-k"}, args...)
	}

	cmd := exec.Command("journalctl", args...)
	cmd.Env = append(os.Environ(), "SYSTEMD_COLORS=0")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}

	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines, nil
}
