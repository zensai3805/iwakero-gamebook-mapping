package domain

// GamebookRepository はゲームブックの永続化を担当するインターフェース
type GamebookRepository interface {
	// Save はゲームブックを保存する
	Save(gamebook *Gamebook) error
	
	// Load は指定されたタイトルのゲームブックを読み込む
	Load(title string) (*Gamebook, error)
	
	// List は保存されているゲームブックのタイトル一覧を返す
	List() ([]string, error)
	
	// Delete は指定されたタイトルのゲームブックを削除する
	Delete(title string) error
}