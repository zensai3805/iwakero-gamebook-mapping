package presenter

import (
	"fmt"
	"strings"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/usecase/interfaces"
)

// NavigationPresenter はNavigationPresenterインターフェースの実装
type NavigationPresenter struct {
}

// NewNavigationPresenter は新しいNavigationPresenterを作成する
func NewNavigationPresenter() interfaces.NavigationPresenter {
	return &NavigationPresenter{}
}

// FormatNavigationHistory は移動履歴を表示用にフォーマットする
func (p *NavigationPresenter) FormatNavigationHistory(title string, history []domain.NavigationStep) (string, error) {
	var builder strings.Builder

	// タイトル部分
	builder.WriteString(fmt.Sprintf("# %s の移動履歴\n\n", title))

	// 履歴が空の場合
	if len(history) == 0 {
		builder.WriteString("移動履歴はありません。\n")
		return builder.String(), nil
	}

	// 各ステップをフォーマット
	for i, step := range history {
		builder.WriteString(fmt.Sprintf("%d. %d → %d\n", i+1, step.From, step.To))
	}

	return builder.String(), nil
}

// FormatCurrentPath は現在のパスを表示用にフォーマットする
func (p *NavigationPresenter) FormatCurrentPath(history []domain.NavigationStep) (string, error) {
	// TODO: 実装予定
	return "", nil
}
