package interfaces

import "github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"

// NavigationRepository は移動履歴の永続化を担当するリポジトリインターフェース
type NavigationRepository interface {
	// SaveNavigationHistory は指定されたゲームブックの移動履歴を保存する
	SaveNavigationHistory(gamebookTitle string, history []domain.NavigationStep) error

	// LoadNavigationHistory は指定されたゲームブックの移動履歴を読み込む
	LoadNavigationHistory(gamebookTitle string) ([]domain.NavigationStep, error)
}
