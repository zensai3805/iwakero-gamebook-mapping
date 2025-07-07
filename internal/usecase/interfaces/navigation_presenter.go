package interfaces

import "github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"

// NavigationPresenter は移動履歴の表示を担当するプレゼンターインターフェース
type NavigationPresenter interface {
	// FormatNavigationHistory は移動履歴を表示用にフォーマットする
	FormatNavigationHistory(title string, history []domain.NavigationStep) (string, error)

	// FormatCurrentPath は現在のパスを表示用にフォーマットする
	FormatCurrentPath(history []domain.NavigationStep) (string, error)
}
