package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/usecase/interfaces"
)

// FileWriter はファイルにログを書き込むWriter実装
type FileWriter struct {
	file      *os.File
	formatter interfaces.LogFormatter
	mu        sync.Mutex
	closed    bool
}

// NewFileWriter は新しいFileWriterを生成する
func NewFileWriter(filePath string, formatter interfaces.LogFormatter) (*FileWriter, error) {
	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(filePath)
	if createErr := os.MkdirAll(dir, 0755); createErr != nil {
		return nil, fmt.Errorf("ディレクトリ作成に失敗: %w", createErr)
	}

	// ファイルを追記モードで開く
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("ファイルオープンに失敗: %w", err)
	}

	return &FileWriter{
		file:      file,
		formatter: formatter,
	}, nil
}

// Write は複数のログエントリをファイルに書き込む
func (w *FileWriter) Write(entries []domain.LogEntry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return fmt.Errorf("ファイルライターは既にクローズされています")
	}

	for _, entry := range entries {
		// フォーマットして書き込み
		formatted, formatErr := w.formatter.Format(entry)
		if formatErr != nil {
			return fmt.Errorf("ログエントリのフォーマットに失敗: %w", formatErr)
		}

		// 改行を追加
		_, writeErr := w.file.Write(append(formatted, '\n'))
		if writeErr != nil {
			return fmt.Errorf("ファイル書き込みに失敗: %w", writeErr)
		}
	}

	// ファイルをフラッシュ
	if syncErr := w.file.Sync(); syncErr != nil {
		return fmt.Errorf("ファイル同期に失敗: %w", syncErr)
	}

	return nil
}

// Close はファイルをクローズする
func (w *FileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return nil
	}

	w.closed = true

	if closeErr := w.file.Close(); closeErr != nil {
		return fmt.Errorf("ファイルクローズに失敗: %w", closeErr)
	}

	return nil
}
