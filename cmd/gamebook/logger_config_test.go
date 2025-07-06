package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestLoadLoggerConfig_DefaultPath(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "gamebook_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// デフォルトパス作成
	defaultConfigDir := filepath.Join(tempDir, ".gamebook")
	err = os.MkdirAll(defaultConfigDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	configContent := `
level: DEBUG
format: json
output: console
file_path: ""
`
	configPath := filepath.Join(defaultConfigDir, "logger.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// 元のホームディレクトリを保存し、テンポラリディレクトリに変更
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Act
	config, err := LoadLoggerConfig("")

	// Assert
	if err != nil {
		t.Fatalf("LoadLoggerConfig failed: %v", err)
	}
	if config.Level != domain.LogLevelDebug {
		t.Errorf("Expected level DEBUG, got %v", config.Level)
	}
	if config.Format != "json" {
		t.Errorf("Expected format json, got %s", config.Format)
	}
	if config.OutputType != "console" {
		t.Errorf("Expected output console, got %s", config.OutputType)
	}
}

func TestLoadLoggerConfig_CustomPath(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "gamebook_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configContent := `
level: ERROR
format: text
output: file
file_path: /var/log/gamebook.log
`
	configPath := filepath.Join(tempDir, "custom_logger.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Act
	config, err := LoadLoggerConfig(configPath)

	// Assert
	if err != nil {
		t.Fatalf("LoadLoggerConfig failed: %v", err)
	}
	if config.Level != domain.LogLevelError {
		t.Errorf("Expected level ERROR, got %v", config.Level)
	}
	if config.Format != "text" {
		t.Errorf("Expected format text, got %s", config.Format)
	}
	if config.OutputType != "file" {
		t.Errorf("Expected output file, got %s", config.OutputType)
	}
	if config.FilePath != "/var/log/gamebook.log" {
		t.Errorf("Expected file_path /var/log/gamebook.log, got %s", config.FilePath)
	}
}

func TestLoadLoggerConfig_FileNotExists(t *testing.T) {
	// Arrange
	nonExistentPath := "/does/not/exist/logger.yaml"

	// Act
	config, err := LoadLoggerConfig(nonExistentPath)

	// Assert
	if err != nil {
		t.Fatalf("LoadLoggerConfig should not fail when file doesn't exist: %v", err)
	}
	// デフォルト設定が返されるはず
	if config.Level != domain.LogLevelInfo {
		t.Errorf("Expected default level INFO, got %v", config.Level)
	}
	if config.Format != "text" {
		t.Errorf("Expected default format text, got %s", config.Format)
	}
	if config.OutputType != "console" {
		t.Errorf("Expected default output console, got %s", config.OutputType)
	}
}

func TestSetDynamicLogLevel(t *testing.T) {
	// Arrange
	config := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "text",
		OutputType: "console",
	}

	controller, err := NewLoggingController(config)
	if err != nil {
		t.Fatalf("Failed to create LoggingController: %v", err)
	}

	// Act
	err = SetDynamicLogLevel(controller, domain.LogLevelDebug)

	// Assert
	if err != nil {
		t.Fatalf("SetDynamicLogLevel failed: %v", err)
	}
	if controller.GetLevel() != domain.LogLevelDebug {
		t.Errorf("Expected level DEBUG, got %v", controller.GetLevel())
	}
}

func TestIsAIDevelopmentMode(t *testing.T) {
	// Arrange
	originalAIMode := os.Getenv("GAMEBOOK_AI_DEV")
	defer os.Setenv("GAMEBOOK_AI_DEV", originalAIMode)

	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{"AIモード有効", "true", true},
		{"AIモード有効（大文字）", "TRUE", true},
		{"AIモード有効（1）", "1", true},
		{"AIモード無効", "false", false},
		{"AIモード無効（空文字）", "", false},
		{"AIモード無効（その他）", "other", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			os.Setenv("GAMEBOOK_AI_DEV", tt.envValue)

			// Act
			result := IsAIDevelopmentMode()

			// Assert
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
