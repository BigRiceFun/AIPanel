package system

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/aipanel/aipanel/server/config"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type Handler struct {
	cfg *config.Config
}

type StatusResponse struct {
	CPU      int    `json:"cpu"`
	Memory   int    `json:"memory"`
	Disk     int    `json:"disk"`
	Hostname string `json:"hostname"`
	Uptime   string `json:"uptime"`
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) Status(c *gin.Context) {
	cpuPercents, err := cpu.Percent(200*time.Millisecond, false)
	if err != nil || len(cpuPercents) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read cpu usage"})
		return
	}

	memory, err := mem.VirtualMemory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read memory usage"})
		return
	}

	diskUsage, err := disk.Usage(h.cfg.System.DiskPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read disk usage"})
		return
	}

	info, err := host.Info()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read host info"})
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		CPU:      roundPercent(cpuPercents[0]),
		Memory:   roundPercent(memory.UsedPercent),
		Disk:     roundPercent(diskUsage.UsedPercent),
		Hostname: info.Hostname,
		Uptime:   formatUptime(info.Uptime),
	})
}

func roundPercent(value float64) int {
	return int(math.Round(value))
}

func formatUptime(seconds uint64) string {
	duration := time.Duration(seconds) * time.Second
	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	if days > 0 {
		return fmt.Sprintf("%d days %d hours", days, hours)
	}
	return fmt.Sprintf("%d hours", hours)
}
