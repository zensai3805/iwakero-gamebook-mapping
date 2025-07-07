package controllers

import (
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/usecase/interfaces"
)

// NavigationController は移動履歴の操作を制御するコントローラー
type NavigationController struct {
	repository interfaces.NavigationRepository
	presenter  interfaces.NavigationPresenter
}

// NewNavigationController は新しいNavigationControllerを作成する
func NewNavigationController(
	repository interfaces.NavigationRepository,
	presenter interfaces.NavigationPresenter,
) *NavigationController {
	return &NavigationController{
		repository: repository,
		presenter:  presenter,
	}
}

// SaveHistory は移動履歴を保存する
func (c *NavigationController) SaveHistory(gamebookTitle string, history []domain.NavigationStep) error {
	return c.repository.SaveNavigationHistory(gamebookTitle, history)
}

// LoadHistory は移動履歴を読み込む
func (c *NavigationController) LoadHistory(gamebookTitle string) ([]domain.NavigationStep, error) {
	return c.repository.LoadNavigationHistory(gamebookTitle)
}

// FormatHistory は移動履歴をフォーマットして返す
func (c *NavigationController) FormatHistory(gamebookTitle string) (string, error) {
	history, err := c.repository.LoadNavigationHistory(gamebookTitle)
	if err != nil {
		return "", err
	}

	return c.presenter.FormatNavigationHistory(gamebookTitle, history)
}
