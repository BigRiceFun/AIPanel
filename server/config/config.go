package config

import (
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	JWT      JWTConfig      `yaml:"jwt"`
	Database DatabaseConfig `yaml:"database"`
	Security SecurityConfig `yaml:"security"`
	System   SystemConfig   `yaml:"system"`
	Terminal TerminalConfig `yaml:"terminal"`
	AI       AIConfig       `yaml:"ai"`
}

type ServerConfig struct {
	Port        int      `yaml:"port"`
	CorsOrigins []string `yaml:"cors_origins"`
}

type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpireHour int    `yaml:"expire_hour"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type SecurityConfig struct {
	InitialUsername string `yaml:"initial_username"`
	InitialPassword string `yaml:"initial_password"`
}

type SystemConfig struct {
	DiskPath string `yaml:"disk_path"`
}

type TerminalConfig struct {
	Enabled bool   `yaml:"enabled"`
	Shell   string `yaml:"shell"`
}

type AIConfig struct {
	Provider string `yaml:"provider"`
	BaseURL  string `yaml:"base_url"`
	APIKey   string `yaml:"api_key"`
	Model    string `yaml:"model"`
}

func Load() (*Config, error) {
	path := os.Getenv("AIPANEL_CONFIG")
	if path == "" {
		path = "config.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	applyEnvOverrides(&cfg)
	applyDefaults(&cfg)
	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	if origins := os.Getenv("AIPANEL_CORS_ORIGINS"); origins != "" {
		cfg.Server.CorsOrigins = splitCSV(origins)
	}
}

func applyDefaults(cfg *Config) {
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if len(cfg.Server.CorsOrigins) == 0 {
		cfg.Server.CorsOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	}
	if cfg.JWT.ExpireHour == 0 {
		cfg.JWT.ExpireHour = 24
	}
	if cfg.System.DiskPath == "" {
		cfg.System.DiskPath = "."
	}
	if cfg.Terminal.Shell == "" {
		cfg.Terminal.Shell = defaultShell()
	}
	if cfg.AI.Provider == "" {
		cfg.AI.Provider = "openai"
	}
}

func defaultShell() string {
	if os.Getenv("OS") == "Windows_NT" {
		return "powershell.exe"
	}
	return "/bin/bash"
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func validate(cfg *Config) error {
	if strings.TrimSpace(cfg.JWT.Secret) == "" {
		return errors.New("jwt.secret is required")
	}
	if strings.TrimSpace(cfg.Database.Path) == "" {
		return errors.New("database.path is required")
	}
	if strings.TrimSpace(cfg.Security.InitialUsername) == "" {
		return errors.New("security.initial_username is required")
	}
	if strings.TrimSpace(cfg.Security.InitialPassword) == "" {
		return errors.New("security.initial_password is required")
	}
	return nil
}
