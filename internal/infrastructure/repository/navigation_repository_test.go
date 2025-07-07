package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestNewNavigationRepository_WhenValidDataDir_ReturnsRepository(t *testing.T) {
	// Arrange
	dataDir := "./test_data"
	defer os.RemoveAll(dataDir)

	// Act
	repo := NewNavigationRepository(dataDir)

	// Assert
	assert.NotNil(t, repo)
}

func TestSaveNavigationHistory_WhenValidInput_SavesFile(t *testing.T) {
	// Arrange
	dataDir := "./test_data"
	defer os.RemoveAll(dataDir)
	
	repo := NewNavigationRepository(dataDir)
	gamebookTitle := "test_game"
	history := []domain.NavigationStep{
		*domain.NewNavigationStep(1, 2, []int{}),
	}

	// Act
	err := repo.SaveNavigationHistory(gamebookTitle, history)

	// Assert
	assert.NoError(t, err)
	// ファイルが実際に作成されているかを確認
	filePath := "./test_data/test_game_history.md"
	_, statErr := os.Stat(filePath)
	assert.NoError(t, statErr, "履歴ファイルが作成されていない")
}

func TestLoadNavigationHistory_WhenSavedData_ReturnsCorrectHistory(t *testing.T) {
	// Arrange
	dataDir := "./test_data"
	defer os.RemoveAll(dataDir)
	
	repo := NewNavigationRepository(dataDir)
	gamebookTitle := "test_game"
	originalHistory := []domain.NavigationStep{
		*domain.NewNavigationStep(1, 2, []int{}),
		*domain.NewNavigationStep(2, 5, []int{}),
	}

	// 履歴を保存
	saveErr := repo.SaveNavigationHistory(gamebookTitle, originalHistory)
	assert.NoError(t, saveErr)

	// Act
	loadedHistory, err := repo.LoadNavigationHistory(gamebookTitle)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, len(originalHistory), len(loadedHistory))
	assert.Equal(t, originalHistory[0].From, loadedHistory[0].From)
	assert.Equal(t, originalHistory[0].To, loadedHistory[0].To)
	assert.Equal(t, originalHistory[1].From, loadedHistory[1].From)
	assert.Equal(t, originalHistory[1].To, loadedHistory[1].To)
}