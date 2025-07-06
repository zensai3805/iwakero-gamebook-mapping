package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/infrastructure/logger"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/usecase"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/usecase/interfaces"
)

// LoggingConfig はログ設定を表す構造体
type LoggingConfig struct {
	Level      domain.LogLevel
	Format     string    // "text" or "json"
	OutputType string    // "console" or "file"
	FilePath   string    // OutputType が "file" の場合に必要
	Writer     io.Writer // テスト用のライター（オプション）
	SyncMode   bool      // テスト用の同期モード（オプション）
}

// LoggingController はログ設定を管理するコントローラー
type LoggingController struct {
	config         LoggingConfig
	loggingService *usecase.LoggingService
	mu             sync.RWMutex
}

// NewLoggingController は新しいLoggingControllerを生成する
func NewLoggingController(config LoggingConfig) (*LoggingController, error) {
	// フォーマッターの作成
	formatter, err := createFormatter(config.Format)
	if err != nil {
		return nil, fmt.Errorf("フォーマッター作成に失敗: %w", err)
	}

	// ライターの作成
	writer, err := createWriter(config, formatter)
	if err != nil {
		return nil, fmt.Errorf("ライター作成に失敗: %w", err)
	}

	// LoggingServiceの作成
	var loggingService *usecase.LoggingService
	if config.SyncMode {
		loggingService = usecase.NewLoggingService(writer, formatter, usecase.WithSyncMode())
	} else {
		loggingService = usecase.NewLoggingService(writer, formatter)
	}

	return &LoggingController{
		config:         config,
		loggingService: loggingService,
	}, nil
}

// GetLogger はロガーを取得する
func (c *LoggingController) GetLogger() domain.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.loggingService
}

// SetLevel はログレベルを設定する
func (c *LoggingController) SetLevel(level domain.LogLevel) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.config.Level = level
}

// GetLevel はログレベルを取得する
func (c *LoggingController) GetLevel() domain.LogLevel {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.config.Level
}

// Close はリソースをクリーンアップする
func (c *LoggingController) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.loggingService != nil {
		if closeErr := c.loggingService.Close(); closeErr != nil {
			return fmt.Errorf("ロギングサービスのクローズに失敗: %w", closeErr)
		}
	}

	return nil
}

// createFormatter はフォーマッターを作成する
func createFormatter(format string) (interfaces.LogFormatter, error) {
	switch format {
	case "text":
		return logger.NewTextFormatter(), nil
	case "json":
		return logger.NewJSONFormatter(), nil
	default:
		return nil, fmt.Errorf("不正なフォーマット: %s", format)
	}
}

// createWriter はライターを作成する
func createWriter(config LoggingConfig, formatter interfaces.LogFormatter) (interfaces.LogWriter, error) {
	switch config.OutputType {
	case "console":
		var writer io.Writer = os.Stdout
		if config.Writer != nil {
			writer = config.Writer
		}
		return logger.NewConsoleWriterWithLevelFilter(writer, formatter, config.Level), nil
	case "file":
		if config.FilePath == "" {
			return nil, fmt.Errorf("ファイル出力にはファイルパスが必要です")
		}
		return logger.NewFileWriter(config.FilePath, formatter)
	default:
		return nil, fmt.Errorf("不正な出力タイプ: %s", config.OutputType)
	}
}
