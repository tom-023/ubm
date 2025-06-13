package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DefaultBrowser string `yaml:"default_browser"`
	Editor         string `yaml:"editor"`
	DisplayFormat  string `yaml:"display_format"`
	AutoBackup     bool   `yaml:"auto_backup"`
	MaxBackups     int    `yaml:"max_backups"`
}

var defaultConfig = Config{
	DefaultBrowser: "",
	Editor:         "",
	DisplayFormat:  "tree",
	AutoBackup:     true,
	MaxBackups:     5,
}

func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "ubm"), nil
}

func Load() (*Config, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "config.yaml")

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := defaultConfig
		return &cfg, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse config
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults for missing values
	if cfg.DisplayFormat == "" {
		cfg.DisplayFormat = defaultConfig.DisplayFormat
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = defaultConfig.MaxBackups
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")

	// Marshal config
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func Init() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Create config directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		return nil // Config already exists
	}

	// Save default config
	cfg := defaultConfig
	return Save(&cfg)
}