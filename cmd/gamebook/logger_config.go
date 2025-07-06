package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
	"gopkg.in/yaml.v3"
)

// 設定定数
const (
	FormatText    = "text"
	FormatJSON    = "json"
	OutputConsole = "console"
	OutputFile    = "file"
	LevelINFO     = "INFO"
)

// LoggerConfig はロガーの設定を表す構造体
type LoggerConfig struct {
	Level             domain.LogLevel `yaml:"level"`
	Format            string          `yaml:"format"`
	OutputType        string          `yaml:"output_type"`
	FilePath          string          `yaml:"file_path"`
	AIDevelopmentMode bool            `yaml:"ai_development_mode"`
}

// LoadLoggerConfig はデフォルト設定またはファイルから設定を読み込む
func LoadLoggerConfig() (*LoggerConfig, error) {
	configPath := GetDefaultConfigPath()

	// 設定ファイルが存在するかチェック
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 設定ファイルが存在しない場合はデフォルト設定を返す
		return getDefaultConfig(), nil
	}

	// 設定ファイルから読み込む
	return LoadLoggerConfigFromFile(configPath)
}

// LoadLoggerConfigFromFile は指定されたファイルから設定を読み込む
func LoadLoggerConfigFromFile(configPath string) (*LoggerConfig, error) {
	// ファイルを読み込む
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("設定ファイル読み込みに失敗: %w", err)
	}

	// YAML形式でパース
	var yamlConfig struct {
		Level             string `yaml:"level"`
		Format            string `yaml:"format"`
		OutputType        string `yaml:"output_type"`
		FilePath          string `yaml:"file_path"`
		AIDevelopmentMode bool   `yaml:"ai_development_mode"`
	}

	if parseErr := yaml.Unmarshal(data, &yamlConfig); parseErr != nil {
		return nil, fmt.Errorf("設定ファイルのパースに失敗: %w", parseErr)
	}

	// ログレベルを変換
	logLevel, err := parseLogLevel(yamlConfig.Level)
	if err != nil {
		return nil, fmt.Errorf("ログレベルの変換に失敗: %w", err)
	}

	// 設定を作成
	config := &LoggerConfig{
		Level:             logLevel,
		Format:            yamlConfig.Format,
		OutputType:        yamlConfig.OutputType,
		FilePath:          yamlConfig.FilePath,
		AIDevelopmentMode: yamlConfig.AIDevelopmentMode,
	}

	// 設定を検証
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("設定検証に失敗: %w", err)
	}

	return config, nil
}

// SaveLoggerConfig は設定をファイルに保存する
func SaveLoggerConfig(configPath string, config *LoggerConfig) error {
	// 設定を検証
	if err := validateConfig(config); err != nil {
		return fmt.Errorf("設定検証に失敗: %w", err)
	}

	// ディレクトリを作成
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("設定ディレクトリの作成に失敗: %w", err)
	}

	// YAML形式に変換
	yamlConfig := struct {
		Level             string `yaml:"level"`
		Format            string `yaml:"format"`
		OutputType        string `yaml:"output_type"`
		FilePath          string `yaml:"file_path"`
		AIDevelopmentMode bool   `yaml:"ai_development_mode"`
	}{
		Level:             config.Level.String(),
		Format:            config.Format,
		OutputType:        config.OutputType,
		FilePath:          config.FilePath,
		AIDevelopmentMode: config.AIDevelopmentMode,
	}

	data, err := yaml.Marshal(&yamlConfig)
	if err != nil {
		return fmt.Errorf("YAML変換に失敗: %w", err)
	}

	// ファイルに書き込む
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("設定ファイルの書き込みに失敗: %w", err)
	}

	return nil
}

// GetDefaultConfigPath はデフォルトの設定ファイルパスを取得する
func GetDefaultConfigPath() string {
	// ホームディレクトリを取得
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// エラーの場合はカレントディレクトリを使用
		return ".gamebook/logger.yaml"
	}

	return filepath.Join(homeDir, ".gamebook", "logger.yaml")
}

// ToggleAIDevelopmentMode はAI開発モードを切り替える
func ToggleAIDevelopmentMode(configPath string) error {
	// 現在の設定を読み込む
	config, err := LoadLoggerConfigFromFile(configPath)
	if err != nil {
		// 設定ファイルが存在しない場合はデフォルト設定を作成
		config = getDefaultConfig()
	}

	// AI開発モードを切り替える
	config.AIDevelopmentMode = !config.AIDevelopmentMode

	// AI開発モードが有効な場合は詳細ログを有効にする
	if config.AIDevelopmentMode {
		config.Level = domain.LogLevelDebug
		config.Format = FormatJSON
	}

	// 設定を保存
	return SaveLoggerConfig(configPath, config)
}

// getDefaultConfig はデフォルト設定を返す
func getDefaultConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:             domain.LogLevelInfo,
		Format:            FormatText,
		OutputType:        OutputConsole,
		FilePath:          "",
		AIDevelopmentMode: false,
	}
}

// validateConfig は設定を検証する
func validateConfig(config *LoggerConfig) error {
	if config == nil {
		return fmt.Errorf("設定がnullです")
	}

	// フォーマットの検証
	if config.Format != FormatText && config.Format != FormatJSON {
		return fmt.Errorf("無効なフォーマット: %s", config.Format)
	}

	// 出力タイプの検証
	if config.OutputType != OutputConsole && config.OutputType != OutputFile {
		return fmt.Errorf("無効な出力タイプ: %s", config.OutputType)
	}

	// ファイル出力の場合はファイルパスが必要
	if config.OutputType == OutputFile && config.FilePath == "" {
		return fmt.Errorf("ファイル出力にはファイルパスが必要です")
	}

	return nil
}
