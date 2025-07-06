package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/infrastructure/logger"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/usecase"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/usecase/interfaces"
)

var (
	globalLogger      domain.Logger
	globalLoggerMu    sync.RWMutex
	loggingController *LoggingController
)

// InitializeLogger はアプリケーション起動時にログシステムを初期化する
func InitializeLogger() (domain.Logger, error) {
	// 環境変数から設定を読み取り
	config := LoggingConfig{
		Level:      getLogLevelFromEnv(),
		Format:     getLogFormatFromEnv(),
		OutputType: getLogOutputFromEnv(),
		FilePath:   getLogFilePathFromEnv(),
		SyncMode:   false, // 本番環境では非同期モード
	}

	// LoggingControllerを作成
	controller, err := NewLoggingController(config)
	if err != nil {
		return nil, fmt.Errorf("LoggingController作成に失敗: %w", err)
	}

	// グローバル変数に保存
	loggingController = controller
	logger := controller.GetLogger()

	// グローバルロガーを設定
	SetGlobalLogger(logger)

	return logger, nil
}

// SetGlobalLogger はグローバルロガーを設定する
func SetGlobalLogger(logger domain.Logger) {
	globalLoggerMu.Lock()
	defer globalLoggerMu.Unlock()
	globalLogger = logger
}

// GetGlobalLogger はグローバルロガーを取得する
func GetGlobalLogger() domain.Logger {
	globalLoggerMu.RLock()
	defer globalLoggerMu.RUnlock()
	return globalLogger
}

// CleanupLogger はアプリケーション終了時にリソースをクリーンアップする
func CleanupLogger() error {
	globalLoggerMu.Lock()
	defer globalLoggerMu.Unlock()

	var cleanupErr error

	// グローバルロガーがClosableInterface（CloseMethodを持つ）の場合はClose
	if globalLogger != nil {
		if closableLogger, ok := globalLogger.(interface{ Close() error }); ok {
			if closeErr := closableLogger.Close(); closeErr != nil {
				cleanupErr = fmt.Errorf("グローバルロガーのクローズに失敗: %w", closeErr)
			}
		}
	}

	if loggingController != nil {
		if closeErr := loggingController.Close(); closeErr != nil {
			if cleanupErr != nil {
				cleanupErr = fmt.Errorf("複数のクローズエラー: %w, %v", cleanupErr, closeErr)
			} else {
				cleanupErr = fmt.Errorf("LoggingController クローズに失敗: %w", closeErr)
			}
		}
		loggingController = nil
	}

	globalLogger = nil
	return cleanupErr
}

// getLogLevelFromEnv は環境変数からログレベルを取得する
func getLogLevelFromEnv() domain.LogLevel {
	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		return domain.LogLevelInfo // デフォルト
	}

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
		return domain.LogLevelInfo
	}
}

// getLogFormatFromEnv は環境変数からログフォーマットを取得する
func getLogFormatFromEnv() string {
	format := os.Getenv("LOG_FORMAT")
	if format == "" {
		return LogFormatText // デフォルト
	}
	return format
}

// getLogOutputFromEnv は環境変数からログ出力タイプを取得する
func getLogOutputFromEnv() string {
	output := os.Getenv("LOG_OUTPUT")
	if output == "" {
		// AI開発モード時はデフォルトでstderrに出力
		if IsAIDevelopmentMode() {
			return "stderr"
		}
		return LogOutputConsole // デフォルト
	}
	return output
}

// getLogFilePathFromEnv は環境変数からログファイルパスを取得する
func getLogFilePathFromEnv() string {
	return os.Getenv("LOG_FILE_PATH")
}

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
	case LogFormatText:
		return logger.NewTextFormatter(), nil
	case LogFormatJSON:
		return logger.NewJSONFormatter(), nil
	default:
		return nil, fmt.Errorf("不正なフォーマット: %s", format)
	}
}

// createWriter はライターを作成する
func createWriter(config LoggingConfig, formatter interfaces.LogFormatter) (interfaces.LogWriter, error) {
	switch config.OutputType {
	case LogOutputConsole:
		var writer io.Writer = os.Stdout
		if config.Writer != nil {
			writer = config.Writer
		}
		return logger.NewConsoleWriterWithLevelFilter(writer, formatter, config.Level), nil
	case LogOutputStderr:
		var writer io.Writer = os.Stderr
		if config.Writer != nil {
			writer = config.Writer
		}
		return logger.NewConsoleWriterWithLevelFilter(writer, formatter, config.Level), nil
	case LogOutputFile:
		if config.FilePath == "" {
			return nil, fmt.Errorf("ファイル出力にはファイルパスが必要です")
		}
		return logger.NewFileWriter(config.FilePath, formatter)
	default:
		return nil, fmt.Errorf("不正な出力タイプ: %s", config.OutputType)
	}
}
