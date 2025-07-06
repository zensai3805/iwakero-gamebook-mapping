package main

import (
	"fmt"
	"strconv"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

// InputHelper は入力支援機能を提供する
type InputHelper struct {
	currentGame *domain.Gamebook
}

// NewInputHelper は新しいInputHelperを作成する
func NewInputHelper(currentGame *domain.Gamebook) *InputHelper {
	return &InputHelper{
		currentGame: currentGame,
	}
}

// GetDefaultText は現在地のデフォルトテキストを返す
func (h *InputHelper) GetDefaultText() string {
	if h.currentGame == nil || h.currentGame.Current == nil {
		return ""
	}
	return strconv.Itoa(h.currentGame.Current.Number)
}

// GetDefaultTextForAdd はパラグラフ追加時のデフォルトテキストを返す（未定義の場合のみ）
func (h *InputHelper) GetDefaultTextForAdd() string {
	if h.currentGame == nil || h.currentGame.Current == nil {
		return ""
	}
	// 現在地が未定義パラグラフの場合のみデフォルト値を提供
	if h.currentGame.Current.Description == "(未定義)" {
		return strconv.Itoa(h.currentGame.Current.Number)
	}
	return ""
}

// GetCLIPrompt はCLIモード用のプロンプトを生成する
func (h *InputHelper) GetCLIPrompt(prompt string) string {
	if h.currentGame != nil && h.currentGame.Current != nil {
		return fmt.Sprintf("%s [現在地: %d]: ", prompt, h.currentGame.Current.Number)
	}
	return prompt + ": "
}

// GetSuggestions は補完候補を返す
func (h *InputHelper) GetSuggestions() []string {
	if h.currentGame == nil || h.currentGame.Current == nil {
		return []string{}
	}
	return []string{strconv.Itoa(h.currentGame.Current.Number)}
}

// ProcessEmptyInput は空入力を処理する
func (h *InputHelper) ProcessEmptyInput(input string) string {
	if input == "" && h.currentGame != nil && h.currentGame.Current != nil {
		return strconv.Itoa(h.currentGame.Current.Number)
	}
	return input
}
