package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/infrastructure/presenter"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/infrastructure/repository"
)

func TestNewNavigationController_WhenValidDependencies_ReturnsController(t *testing.T) {
	// Arrange
	repo := repository.NewNavigationRepository("./test_data")
	pres := presenter.NewNavigationPresenter()

	// Act
	controller := NewNavigationController(repo, pres)

	// Assert
	assert.NotNil(t, controller)
}