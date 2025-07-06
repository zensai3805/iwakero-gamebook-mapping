package interfaces

import "github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"

// LogFormatter はログエントリのフォーマットを抽象化するインターフェース
type LogFormatter interface {
	// Format はログエントリを指定された形式のバイト列に変換する
	Format(entry domain.LogEntry) ([]byte, error)
}
