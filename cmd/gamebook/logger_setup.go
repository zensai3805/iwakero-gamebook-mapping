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

// LoggingConfig はログ設定を表す構造体
type LoggingConfig struct {
	Level      domain.LogLevel
	Format     string    // "text" or "json"
	OutputType string    // "console" or "file"
	FilePath   string    // OutputType が "file" の場合に必要
	Writer     io.Writer // テスト用のライター（オプション）
	SyncMode   bool      // テスト用の同期モード（オプション）
}

// LoggingController はFrameworksレイヤー用の軽量ログコントローラー
type LoggingController struct {
	config         LoggingConfig
	loggingService *usecase.LoggingService
	writer         interfaces.LogWriter
	mu             sync.RWMutex
}

// グローバルロガーインスタンス
var (
	globalLogger domain.Logger
	globalMu     sync.RWMutex
)

// SetupLogger はアプリケーション起動時のログシステムを初期化する
func SetupLogger() (domain.Logger, func(), error) {
	// 環境変数からログ設定を読み込む
	logLevel, err := parseLogLevel(getEnv("LOG_LEVEL", "INFO"))
	if err != nil {
		return nil, nil, fmt.Errorf("ログレベルの解析に失敗: %w", err)
	}

	outputType := getEnv("LOG_OUTPUT", "console")
	format := getEnv("LOG_FORMAT", "text")
	filePath := getEnv("LOG_FILE_PATH", "")

	// ログ設定を作成
	config := LoggingConfig{
		Level:      logLevel,
		Format:     format,
		OutputType: outputType,
		FilePath:   filePath,
	}

	// ログコントローラーを作成
	controller, err := NewLoggingController(config)
	if err != nil {
		return nil, nil, fmt.Errorf("ログコントローラーの作成に失敗: %w", err)
	}

	// ロガーを取得
	logger := controller.GetLogger()

	// グローバルロガーに設定
	setGlobalLogger(logger)

	// クリーンアップ関数を作成
	cleanupFunc := func() {
		if closeErr := controller.Close(); closeErr != nil {
			// エラー時はstderrに出力
			fmt.Fprintf(os.Stderr, "ログシステムのクローズに失敗: %v\n", closeErr)
		}
	}

	return logger, cleanupFunc, nil
}

// GetGlobalLogger はグローバルロガーインスタンスを取得する
func GetGlobalLogger() domain.Logger {
	globalMu.RLock()
	defer globalMu.RUnlock()

	if globalLogger == nil {
		// デフォルトのロガーを作成
		defaultLogger, cleanupFunc, err := createDefaultLogger()
		if err != nil {
			// エラー時はstderrに出力して、nil safe なロガーを作成
			fmt.Fprintf(os.Stderr, "デフォルトロガーの作成に失敗: %v\n", err)
			return createNilSafeLogger()
		}

		// デフォルトのクリーンアップ関数を設定（この場合は何もしない）
		_ = cleanupFunc

		globalLogger = defaultLogger
	}

	return globalLogger
}

// setGlobalLogger はグローバルロガーを設定する
func setGlobalLogger(logger domain.Logger) {
	globalMu.Lock()
	defer globalMu.Unlock()

	globalLogger = logger
}

// parseLogLevel は文字列をLogLevelに変換する
func parseLogLevel(levelStr string) (domain.LogLevel, error) {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return domain.LogLevelDebug, nil
	case "INFO":
		return domain.LogLevelInfo, nil
	case "WARN":
		return domain.LogLevelWarn, nil
	case "ERROR":
		return domain.LogLevelError, nil
	case "FATAL":
		return domain.LogLevelFatal, nil
	default:
		return domain.LogLevelInfo, fmt.Errorf("無効なログレベル: %s", levelStr)
	}
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// createDefaultLogger はデフォルトのロガーを作成する
func createDefaultLogger() (domain.Logger, func(), error) {
	config := LoggingConfig{
		Level:      domain.LogLevelInfo,
		Format:     "text",
		OutputType: "console",
	}

	controller, err := NewLoggingController(config)
	if err != nil {
		return nil, nil, fmt.Errorf("デフォルトロガーの作成に失敗: %w", err)
	}

	cleanupFunc := func() {
		if closeErr := controller.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "デフォルトログシステムのクローズに失敗: %v\n", closeErr)
		}
	}

	return controller.GetLogger(), cleanupFunc, nil
}

// createNilSafeLogger はnilセーフなロガーを作成する
func createNilSafeLogger() domain.Logger {
	return &nilSafeLogger{}
}

// nilSafeLogger は最低限の機能を持つロガー実装
type nilSafeLogger struct{}

func (n *nilSafeLogger) Debug(msg string, fields ...domain.Field) {
	// 何もしない
}

func (n *nilSafeLogger) Info(msg string, fields ...domain.Field) {
	fmt.Println(msg)
}

func (n *nilSafeLogger) Warn(msg string, fields ...domain.Field) {
	fmt.Printf("WARNING: %s\n", msg)
}

func (n *nilSafeLogger) Error(msg string, fields ...domain.Field) {
	fmt.Printf("ERROR: %s\n", msg)
}

func (n *nilSafeLogger) Fatal(msg string, fields ...domain.Field) {
	fmt.Printf("FATAL: %s\n", msg)
	os.Exit(1)
}

func (n *nilSafeLogger) WithContext(fields ...domain.Field) domain.Logger {
	return n
}

// NewLoggingController は新しいLoggingControllerを生成する
// NOTE: この実装はFrameworksレイヤー専用の軽量実装です
func NewLoggingController(config LoggingConfig) (*LoggingController, error) {
	// フォーマッターの作成
	formatter, formatErr := createSimpleFormatter(config.Format)
	if formatErr != nil {
		return nil, fmt.Errorf("フォーマッター作成に失敗: %w", formatErr)
	}

	// ライターの作成
	writer, writerErr := createSimpleWriter(config, formatter)
	if writerErr != nil {
		return nil, fmt.Errorf("ライター作成に失敗: %w", writerErr)
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
		writer:         writer,
	}, nil
}

// GetLogger はロガーを取得する
func (c *LoggingController) GetLogger() domain.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.loggingService
}

// SetLevel はログレベルを設定する
func (c *LoggingController) SetLevel(level domain.LogLevel) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.config.Level = level

	// ConsoleWriterWithLevelFilter の場合はレベルを更新
	if levelWriter, ok := c.writer.(interface{ SetLevel(domain.LogLevel) }); ok {
		levelWriter.SetLevel(level)
	}

	return nil
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

// createSimpleFormatter はフォーマッターを作成する (Frameworks Layer用)
func createSimpleFormatter(format string) (interfaces.LogFormatter, error) {
	switch format {
	case FormatText:
		return logger.NewTextFormatter(), nil
	case FormatJSON:
		return logger.NewJSONFormatter(), nil
	default:
		return nil, fmt.Errorf("不正なフォーマット: %s", format)
	}
}

// createSimpleWriter はライターを作成する (Frameworks Layer用)
func createSimpleWriter(config LoggingConfig, formatter interfaces.LogFormatter) (interfaces.LogWriter, error) {
	switch config.OutputType {
	case OutputConsole:
		var writer io.Writer = os.Stdout
		if config.Writer != nil {
			writer = config.Writer
		}
		return logger.NewConsoleWriterWithLevelFilter(writer, formatter, config.Level), nil
	case OutputFile:
		if config.FilePath == "" {
			return nil, fmt.Errorf("ファイル出力にはファイルパスが必要です")
		}
		return logger.NewFileWriter(config.FilePath, formatter)
	default:
		return nil, fmt.Errorf("不正な出力タイプ: %s", config.OutputType)
	}
}
