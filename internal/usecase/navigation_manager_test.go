package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestNewNavigationManager_WhenValidDependencies_ReturnsManager(t *testing.T) {
	manager := NewNavigationManager(nil, nil, nil)
	assert.NotNil(t, manager)
}

type MockNavigationRepository struct {
	mock.Mock
}

func (m *MockNavigationRepository) SaveNavigationStep(step domain.NavigationStep) error {
	args := m.Called(step)
	return args.Error(0)
}

func (m *MockNavigationRepository) GetNavigationHistory() ([]domain.NavigationStep, error) {
	args := m.Called()
	return args.Get(0).([]domain.NavigationStep), args.Error(1)
}

func (m *MockNavigationRepository) GetCurrentGamebook() (*domain.Gamebook, error) {
	args := m.Called()
	return args.Get(0).(*domain.Gamebook), args.Error(1)
}

type MockNavigationPresenter struct {
	mock.Mock
}

func (m *MockNavigationPresenter) FormatNavigationHistory(history []domain.NavigationStep) (interface{}, error) {
	args := m.Called(history)
	return args.Get(0), args.Error(1)
}

func (m *MockNavigationPresenter) FormatCurrentPath(path []int) (interface{}, error) {
	args := m.Called(path)
	return args.Get(0), args.Error(1)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(message string, fields ...domain.Field) {
	m.Called(message, fields)
}

func (m *MockLogger) Info(message string, fields ...domain.Field) {
	m.Called(message, fields)
}

func (m *MockLogger) Warn(message string, fields ...domain.Field) {
	m.Called(message, fields)
}

func (m *MockLogger) Error(message string, fields ...domain.Field) {
	m.Called(message, fields)
}

func (m *MockLogger) Fatal(message string, fields ...domain.Field) {
	m.Called(message, fields)
}

func (m *MockLogger) WithContext(fields ...domain.Field) domain.Logger {
	args := m.Called(fields)
	return args.Get(0).(domain.Logger)
}

func (m *MockLogger) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestRecordChoiceMove_WhenValidMove_SavesStep(t *testing.T) {
	mockRepo := &MockNavigationRepository{}
	mockPresenter := &MockNavigationPresenter{}
	mockLogger := &MockLogger{}

	manager := NewNavigationManager(mockRepo, mockPresenter, mockLogger)

	expectedStep := domain.NavigationStep{From: 1, To: 2, ViaPaths: []int{}}
	mockRepo.On("SaveNavigationStep", expectedStep).Return(nil)

	err := manager.RecordChoiceMove(1, 2)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRecordJumpMove_WhenValidMove_SavesStepWithViaPaths(t *testing.T) {
	mockRepo := &MockNavigationRepository{}
	mockPresenter := &MockNavigationPresenter{}
	mockLogger := &MockLogger{}

	manager := NewNavigationManager(mockRepo, mockPresenter, mockLogger)

	viaPaths := []int{10, 11}
	expectedStep := domain.NavigationStep{From: 1, To: 12, ViaPaths: viaPaths}
	mockRepo.On("SaveNavigationStep", expectedStep).Return(nil)

	err := manager.RecordJumpMove(1, 12, viaPaths)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetNavigationHistory_WhenCalled_ReturnsFormattedHistory(t *testing.T) {
	mockRepo := &MockNavigationRepository{}
	mockPresenter := &MockNavigationPresenter{}
	mockLogger := &MockLogger{}

	manager := NewNavigationManager(mockRepo, mockPresenter, mockLogger)

	history := []domain.NavigationStep{
		{From: 1, To: 2, ViaPaths: []int{}},
	}
	mockRepo.On("GetNavigationHistory").Return(history, nil)
	mockPresenter.On("FormatNavigationHistory", history).Return("formatted", nil)

	result, err := manager.GetNavigationHistory()

	assert.NoError(t, err)
	assert.Equal(t, "formatted", result)
	mockRepo.AssertExpectations(t)
	mockPresenter.AssertExpectations(t)
}

func TestGetCurrentPath_WhenCalled_ReturnsFormattedPath(t *testing.T) {
	mockRepo := &MockNavigationRepository{}
	mockPresenter := &MockNavigationPresenter{}
	mockLogger := &MockLogger{}

	manager := NewNavigationManager(mockRepo, mockPresenter, mockLogger)

	history := []domain.NavigationStep{
		{From: 1, To: 2, ViaPaths: []int{}},
		{From: 2, To: 3, ViaPaths: []int{}},
	}
	mockRepo.On("GetNavigationHistory").Return(history, nil)
	mockPresenter.On("FormatCurrentPath", []int{1, 2, 3}).Return("path", nil)

	result, err := manager.GetCurrentPath()

	assert.NoError(t, err)
	assert.Equal(t, "path", result)
	mockRepo.AssertExpectations(t)
	mockPresenter.AssertExpectations(t)
}
