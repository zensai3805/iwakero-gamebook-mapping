package domain

// Logger はドメイン層のログインターフェース
type Logger interface {
	// Debug はデバッグレベルのログを記録する
	Debug(msg string, fields ...Field)

	// Info は情報レベルのログを記録する
	Info(msg string, fields ...Field)

	// Warn は警告レベルのログを記録する
	Warn(msg string, fields ...Field)

	// Error はエラーレベルのログを記録する
	Error(msg string, fields ...Field)

	// Fatal は致命的エラーレベルのログを記録する
	Fatal(msg string, fields ...Field)

	// WithContext はコンテキストフィールドを追加した新しいロガーを返す
	WithContext(fields ...Field) Logger
}
