package log

import (
	"bufio"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
		return readLogFiles(logType, limit)
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

func readLogFiles(logType string, limit int) ([]string, error) {
	paths := candidateLogFiles(logType)
	for _, path := range paths {
		lines, err := tailFile(path, limit)
		if err == nil && len(lines) > 0 {
			return lines, nil
		}
	}
	return []string{
		"journalctl is not available and no readable host log files were found",
		"mount /var/log to /host/var/log or run the server directly on a Linux host",
	}, nil
}

func candidateLogFiles(logType string) []string {
	base := "/host/var/log"
	if logType == "kernel" {
		return []string{
			filepath.Join(base, "kern.log"),
			filepath.Join(base, "dmesg"),
			filepath.Join(base, "messages"),
		}
	}
	return []string{
		filepath.Join(base, "syslog"),
		filepath.Join(base, "messages"),
		filepath.Join(base, "auth.log"),
	}
}

func tailFile(path string, limit int) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, errors.New("path is a directory")
	}

	ring := make([]string, 0, limit)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(ring) < limit {
			ring = append(ring, line)
			continue
		}
		copy(ring, ring[1:])
		ring[len(ring)-1] = line
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ring, nil
}
