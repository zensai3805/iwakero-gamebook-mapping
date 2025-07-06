package interfaces

import "github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"

// LogWriter はログエントリの書き込みを抽象化するインターフェース
type LogWriter interface {
	// Write は複数のログエントリをバッチで書き込む
	Write(entries []domain.LogEntry) error

	// Close はリソースをクリーンアップする
	Close() error
}
