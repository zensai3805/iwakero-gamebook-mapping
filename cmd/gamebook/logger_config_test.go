package main

import (
	"os"
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestLoadLoggerConfig_DefaultConfig(t *testing.T) {
	// Arrange
	// 設定ファイルが存在しない場合

	// Act
	config, err := LoadLoggerConfig()

	// Assert
	if err != nil {
		t.Fatalf("LoadLoggerConfig()でエラーが発生: %v", err)
	}

	if config.Level != domain.LogLevelInfo {
		t.Fatalf("期待値: %v, 実際の値: %v", domain.LogLevelInfo, config.Level)
	}

	if config.Format != "text" {
		t.Fatalf("期待値: text, 実際の値: %s", config.Format)
	}

	if config.OutputType != "console" {
		t.Fatalf("期待値: console, 実際の値: %s", config.OutputType)
	}
}

func TestLoadLoggerConfig_WithConfigFile(t *testing.T) {
	// Arrange
	configDir := "/tmp/test_gamebook"
	configFile := configDir + "/logger.yaml"

	// テスト用のディレクトリを作成
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("ディレクトリ作成に失敗: %v", err)
	}
	defer os.RemoveAll(configDir)

	// テスト用の設定ファイルを作成
	configContent := `level: DEBUG
format: json
output_type: file
file_path: /tmp/test.log
ai_development_mode: true
`
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("設定ファイル作成に失敗: %v", err)
	}

	// Act
	config, err := LoadLoggerConfigFromFile(configFile)

	// Assert
	if err != nil {
		t.Fatalf("LoadLoggerConfigFromFile()でエラーが発生: %v", err)
	}

	if config.Level != domain.LogLevelDebug {
		t.Fatalf("期待値: %v, 実際の値: %v", domain.LogLevelDebug, config.Level)
	}

	if config.Format != "json" {
		t.Fatalf("期待値: json, 実際の値: %s", config.Format)
	}

	if config.OutputType != "file" {
		t.Fatalf("期待値: file, 実際の値: %s", config.OutputType)
	}

	if config.FilePath != "/tmp/test.log" {
		t.Fatalf("期待値: /tmp/test.log, 実際の値: %s", config.FilePath)
	}
}

func TestLoadLoggerConfig_WithInvalidConfigFile(t *testing.T) {
	// Arrange
	configDir := "/tmp/test_gamebook"
	configFile := configDir + "/logger.yaml"

	// テスト用のディレクトリを作成
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("ディレクトリ作成に失敗: %v", err)
	}
	defer os.RemoveAll(configDir)

	// 不正な設定ファイルを作成
	configContent := `level: INVALID_LEVEL
format: json
output_type: console
`
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("設定ファイル作成に失敗: %v", err)
	}

	// Act
	config, err := LoadLoggerConfigFromFile(configFile)

	// Assert
	if err == nil {
		t.Fatal("無効な設定ファイルでエラーが発生するべきです")
	}

	if config != nil {
		t.Fatal("エラー時にconfigはnilであるべきです")
	}
}

func TestSaveLoggerConfig(t *testing.T) {
	// Arrange
	configDir := "/tmp/test_gamebook"
	configFile := configDir + "/logger.yaml"

	// テスト用のディレクトリを作成
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("ディレクトリ作成に失敗: %v", err)
	}
	defer os.RemoveAll(configDir)

	config := &LoggerConfig{
		Level:             domain.LogLevelDebug,
		Format:            "json",
		OutputType:        "file",
		FilePath:          "/tmp/test.log",
		AIDevelopmentMode: true,
	}

	// Act
	err := SaveLoggerConfig(configFile, config)

	// Assert
	if err != nil {
		t.Fatalf("SaveLoggerConfig()でエラーが発生: %v", err)
	}

	// ファイルが作成されているか確認
	if _, statErr := os.Stat(configFile); os.IsNotExist(statErr) {
		t.Fatal("設定ファイルが作成されていません")
	}

	// 設定を読み込み直して確認
	loadedConfig, err := LoadLoggerConfigFromFile(configFile)
	if err != nil {
		t.Fatalf("設定ファイル読み込みに失敗: %v", err)
	}

	if loadedConfig.Level != config.Level {
		t.Fatalf("期待値: %v, 実際の値: %v", config.Level, loadedConfig.Level)
	}

	if loadedConfig.Format != config.Format {
		t.Fatalf("期待値: %s, 実際の値: %s", config.Format, loadedConfig.Format)
	}
}

func TestGetDefaultConfigPath(t *testing.T) {
	// Act
	configPath := GetDefaultConfigPath()

	// Assert
	if configPath == "" {
		t.Fatal("設定ファイルパスが空です")
	}

	// パスが".gamebook/logger.yaml"で終わることを確認
	expectedSuffix := ".gamebook/logger.yaml"
	if len(configPath) < len(expectedSuffix) || configPath[len(configPath)-len(expectedSuffix):] != expectedSuffix {
		t.Fatalf("期待されるパス末尾: %s, 実際のパス: %s", expectedSuffix, configPath)
	}
}

func TestToggleAIDevelopmentMode(t *testing.T) {
	// Arrange
	configDir := "/tmp/test_gamebook"
	configFile := configDir + "/logger.yaml"

	// テスト用のディレクトリを作成
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("ディレクトリ作成に失敗: %v", err)
	}
	defer os.RemoveAll(configDir)

	// 初期設定を作成
	initialConfig := &LoggerConfig{
		Level:             domain.LogLevelInfo,
		Format:            "text",
		OutputType:        "console",
		AIDevelopmentMode: false,
	}

	if err := SaveLoggerConfig(configFile, initialConfig); err != nil {
		t.Fatalf("初期設定の保存に失敗: %v", err)
	}

	// Act
	err := ToggleAIDevelopmentMode(configFile)

	// Assert
	if err != nil {
		t.Fatalf("ToggleAIDevelopmentMode()でエラーが発生: %v", err)
	}

	// 設定を読み込み直して確認
	updatedConfig, err := LoadLoggerConfigFromFile(configFile)
	if err != nil {
		t.Fatalf("設定ファイル読み込みに失敗: %v", err)
	}

	if !updatedConfig.AIDevelopmentMode {
		t.Fatal("AI開発モードが有効になっていません")
	}

	// もう一度実行して無効にする
	err = ToggleAIDevelopmentMode(configFile)
	if err != nil {
		t.Fatalf("ToggleAIDevelopmentMode()でエラーが発生: %v", err)
	}

	// 設定を読み込み直して確認
	updatedConfig, err = LoadLoggerConfigFromFile(configFile)
	if err != nil {
		t.Fatalf("設定ファイル読み込みに失敗: %v", err)
	}

	if updatedConfig.AIDevelopmentMode {
		t.Fatal("AI開発モードが無効になっていません")
	}
}
