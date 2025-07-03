package domain

// SessionRepository は現在のセッション状態を管理する
type SessionRepository interface {
	// SaveCurrentGame は現在のゲームブックタイトルを保存する
	SaveCurrentGame(title string) error
	
	// GetCurrentGame は現在のゲームブックタイトルを取得する
	GetCurrentGame() (string, error)
	
	// ClearCurrentGame は現在のゲーム設定をクリアする
	ClearCurrentGame() error
}