package usecase

import "github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"

type NavigationManager struct {
	repo      NavigationRepository
	presenter NavigationPresenter
	logger    domain.Logger
}

type NavigationRepository interface {
	SaveNavigationStep(step domain.NavigationStep) error
	GetNavigationHistory() ([]domain.NavigationStep, error)
	GetCurrentGamebook() (*domain.Gamebook, error)
}

type NavigationPresenter interface {
	FormatNavigationHistory(history []domain.NavigationStep) (interface{}, error)
	FormatCurrentPath(path []int) (interface{}, error)
}

func NewNavigationManager(repo NavigationRepository, presenter NavigationPresenter, logger domain.Logger) *NavigationManager {
	return &NavigationManager{
		repo:      repo,
		presenter: presenter,
		logger:    logger,
	}
}

func (nm *NavigationManager) RecordChoiceMove(from, to int) error {
	return nm.recordNavigationStep(from, to, []int{})
}

func (nm *NavigationManager) RecordJumpMove(from, to int, viaPaths []int) error {
	return nm.recordNavigationStep(from, to, viaPaths)
}

func (nm *NavigationManager) recordNavigationStep(from, to int, viaPaths []int) error {
	step := domain.NavigationStep{
		From:     from,
		To:       to,
		ViaPaths: viaPaths,
	}
	return nm.repo.SaveNavigationStep(step)
}

func (nm *NavigationManager) GetNavigationHistory() (interface{}, error) {
	history, err := nm.repo.GetNavigationHistory()
	if err != nil {
		return nil, err
	}
	return nm.presenter.FormatNavigationHistory(history)
}

func (nm *NavigationManager) GetCurrentPath() (interface{}, error) {
	history, err := nm.repo.GetNavigationHistory()
	if err != nil {
		return nil, err
	}

	path := nm.calculateCurrentPath(history)
	return nm.presenter.FormatCurrentPath(path)
}

func (nm *NavigationManager) calculateCurrentPath(history []domain.NavigationStep) []int {
	if len(history) == 0 {
		return []int{}
	}

	var path []int
	for i, step := range history {
		if i == 0 {
			path = append(path, step.From)
		}
		path = append(path, step.To)
	}

	return path
}
