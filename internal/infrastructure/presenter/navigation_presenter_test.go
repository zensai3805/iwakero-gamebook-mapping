package presenter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestNewNavigationPresenter_WhenCalled_ReturnsPresenter(t *testing.T) {
	// Act
	presenter := NewNavigationPresenter()

	// Assert
	assert.NotNil(t, presenter)
}

func TestFormatNavigationHistory_WhenValidInput_ReturnsFormattedString(t *testing.T) {
	// Arrange
	presenter := NewNavigationPresenter()
	title := "テストゲーム"
	history := []domain.NavigationStep{
		*domain.NewNavigationStep(1, 2, []int{}),
		*domain.NewNavigationStep(2, 5, []int{}),
	}

	// Act
	result, err := presenter.FormatNavigationHistory(title, history)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "テストゲーム")
	assert.Contains(t, result, "1")
	assert.Contains(t, result, "2")
	assert.Contains(t, result, "5")
}
