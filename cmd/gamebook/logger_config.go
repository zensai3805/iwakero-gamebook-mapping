package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

// LoggerConfigYAML はYAMLファイルから読み込むログ設定
type LoggerConfigYAML struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	Output   string `yaml:"output"`
	FilePath string `yaml:"file_path"`
}

// LoadLoggerConfig は設定ファイルからログ設定を読み込む
func LoadLoggerConfig(configPath string) (LoggingConfig, error) {
	// デフォルト設定
	defaultConfig := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "text",
		OutputType: "console",
		FilePath:   "",
		SyncMode:   false,
	}

	// 設定ファイルパスの決定
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	// ファイルが存在しない場合はデフォルト設定を返す
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return defaultConfig, nil
	}

	// 設定ファイルの読み込み
	yamlContent, err := os.ReadFile(configPath)
	if err != nil {
		return defaultConfig, fmt.Errorf("設定ファイル読み込みに失敗: %w", err)
	}

	// YAMLをパース
	var yamlConfig LoggerConfigYAML
	if err := yaml.Unmarshal(yamlContent, &yamlConfig); err != nil {
		return defaultConfig, fmt.Errorf("YAML解析に失敗: %w", err)
	}

	// LoggingConfigに変換
	config := LoggingConfig{
		Level:      parseLogLevel(yamlConfig.Level),
		Format:     parseFormat(yamlConfig.Format),
		OutputType: parseOutputType(yamlConfig.Output),
		FilePath:   yamlConfig.FilePath,
		SyncMode:   false,
	}

	return config, nil
}

// SetDynamicLogLevel は動的にログレベルを変更する
func SetDynamicLogLevel(controller *LoggingController, level domain.LogLevel) error {
	if controller == nil {
		return fmt.Errorf("LoggingController が nil です")
	}

	controller.SetLevel(level)
	return nil
}

// IsAIDevelopmentMode はAI開発モードかどうかを判定する
func IsAIDevelopmentMode() bool {
	aiDevMode := os.Getenv("GAMEBOOK_AI_DEV")
	switch strings.ToLower(aiDevMode) {
	case "true", "1":
		return true
	default:
		return false
	}
}

// getDefaultConfigPath はデフォルトの設定ファイルパスを取得する
func getDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// フォールバック: 環境変数HOMEを使用
		homeDir = os.Getenv("HOME")
	}

	if homeDir == "" {
		// さらなるフォールバック: カレントディレクトリ
		return ".gamebook/logger.yaml"
	}

	return filepath.Join(homeDir, ".gamebook", "logger.yaml")
}

// parseLogLevel は文字列からLogLevelに変換する
func parseLogLevel(levelStr string) domain.LogLevel {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return domain.LogLevelDebug
	case "INFO":
		return domain.LogLevelInfo
	case "WARN":
		return domain.LogLevelWarn
	case "ERROR":
		return domain.LogLevelError
	case "FATAL":
		return domain.LogLevelFatal
	default:
		return domain.LogLevelInfo // デフォルト
	}
}

// parseFormat は文字列からフォーマットに変換する
func parseFormat(formatStr string) string {
	switch strings.ToLower(formatStr) {
	case "text", "json":
		return strings.ToLower(formatStr)
	default:
		return "text" // デフォルト
	}
}

// parseOutputType は文字列から出力タイプに変換する
func parseOutputType(outputStr string) string {
	switch strings.ToLower(outputStr) {
	case "console", "file":
		return strings.ToLower(outputStr)
	default:
		return "console" // デフォルト
	}
}
