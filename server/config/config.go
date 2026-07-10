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

	applyDefaults(&cfg)
	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
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
