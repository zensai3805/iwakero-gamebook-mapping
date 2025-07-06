package logger

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/usecase/interfaces"
)

// ConsoleWriter はコンソールにログを出力するWriter実装
type ConsoleWriter struct {
	writer      io.Writer
	formatter   interfaces.LogFormatter
	levelFilter domain.LogLevel
	mu          sync.Mutex
}

// NewConsoleWriter は新しいConsoleWriterを生成する
func NewConsoleWriter(writer io.Writer, formatter interfaces.LogFormatter) *ConsoleWriter {
	if writer == nil {
		writer = os.Stdout
	}

	return &ConsoleWriter{
		writer:      writer,
		formatter:   formatter,
		levelFilter: domain.LogLevelDebug, // デフォルトではすべてのレベルを出力
	}
}

// NewConsoleWriterWithLevelFilter はレベルフィルタ付きのConsoleWriterを生成する
func NewConsoleWriterWithLevelFilter(writer io.Writer, formatter interfaces.LogFormatter, minLevel domain.LogLevel) *ConsoleWriter {
	if writer == nil {
		writer = os.Stdout
	}

	return &ConsoleWriter{
		writer:      writer,
		formatter:   formatter,
		levelFilter: minLevel,
	}
}

// Write は複数のログエントリをコンソールに書き込む
func (w *ConsoleWriter) Write(entries []domain.LogEntry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, entry := range entries {
		// レベルフィルタリング
		if !entry.Level.IsHigherThan(w.levelFilter) && entry.Level != w.levelFilter {
			continue
		}

		// フォーマットして書き込み
		formatted, formatErr := w.formatter.Format(entry)
		if formatErr != nil {
			return fmt.Errorf("ログエントリのフォーマットに失敗: %w", formatErr)
		}

		// 改行を追加して書き込み
		_, writeErr := w.writer.Write(append(formatted, '\n'))
		if writeErr != nil {
			return fmt.Errorf("コンソール書き込みに失敗: %w", writeErr)
		}
	}

	return nil
}

// SetLevel はレベルフィルタを動的に変更する
func (w *ConsoleWriter) SetLevel(level domain.LogLevel) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.levelFilter = level
}

// Close はリソースをクリーンアップする
func (w *ConsoleWriter) Close() error {
	// コンソールライターはクローズ処理は不要
	return nil
}
