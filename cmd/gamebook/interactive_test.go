package main

import (
	"fmt"
	"testing"
)

// MockExecutor はテスト用のモックエグゼキューター
type MockExecutor struct {
	NewCalled    bool
	NewTitle     string
	LoadCalled   bool
	LoadTitle    string
	AddCalled    bool
	AddNumber    int
	AddDesc      string
	ChoiceCalled bool
	SelectCalled bool
	ShowCalled   bool
	ShouldError  bool
}

func (m *MockExecutor) ExecuteNewCommand(title string) error {
	m.NewCalled = true
	m.NewTitle = title
	if m.ShouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (m *MockExecutor) ExecuteLoadCommand(title string) error {
	m.LoadCalled = true
	m.LoadTitle = title
	if m.ShouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (m *MockExecutor) ExecuteAddCommand(number int, description string) error {
	m.AddCalled = true
	m.AddNumber = number
	m.AddDesc = description
	if m.ShouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (m *MockExecutor) ExecuteChoiceCommand(paragraphNum int, description string, targetNum int) error {
	m.ChoiceCalled = true
	if m.ShouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (m *MockExecutor) ExecuteSelectCommand(paragraphNum int, choiceIndex int) error {
	m.SelectCalled = true
	if m.ShouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (m *MockExecutor) ExecuteShowCommand() error {
	m.ShowCalled = true
	if m.ShouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func TestInteractiveShell_New(t *testing.T) {
	// Setup: モックエグゼキューターを作成
	mockExecutor := &MockExecutor{}
	shell := &InteractiveShell{
		executor: mockExecutor,
	}

	// Execute: newコマンドを実行
	shell.handleNew([]string{"TestGame"})

	// Verify: エグゼキューターが呼び出されたことを確認
	if !mockExecutor.NewCalled {
		t.Error("ExecuteNewCommand should be called")
	}
	if mockExecutor.NewTitle != "TestGame" {
		t.Errorf("Expected title 'TestGame', got '%s'", mockExecutor.NewTitle)
	}
}

func TestInteractiveShell_NewWithInvalidArgs(t *testing.T) {
	// Setup
	mockExecutor := &MockExecutor{}
	shell := &InteractiveShell{
		executor: mockExecutor,
	}

	// Execute: 引数なしでnewコマンドを実行
	shell.handleNew([]string{})

	// Verify: エグゼキューターが呼び出されていないことを確認
	if mockExecutor.NewCalled {
		t.Error("ExecuteNewCommand should not be called with invalid args")
	}
}

func TestInteractiveShell_Load(t *testing.T) {
	// Setup
	mockExecutor := &MockExecutor{}
	shell := &InteractiveShell{
		executor: mockExecutor,
	}

	// Execute: loadコマンドを実行
	shell.handleLoad([]string{"TestGame"})

	// Verify: エグゼキューターが呼び出されたことを確認
	if !mockExecutor.LoadCalled {
		t.Error("ExecuteLoadCommand should be called")
	}
	if mockExecutor.LoadTitle != "TestGame" {
		t.Errorf("Expected title 'TestGame', got '%s'", mockExecutor.LoadTitle)
	}
}

func TestInteractiveShell_AddParagraph(t *testing.T) {
	// Setup
	mockExecutor := &MockExecutor{}
	shell := &InteractiveShell{
		executor: mockExecutor,
	}

	// Execute: addコマンドを実行
	shell.handleAdd([]string{"1", "Test", "paragraph", "description"})

	// Verify: エグゼキューターが呼び出されたことを確認
	if !mockExecutor.AddCalled {
		t.Error("ExecuteAddCommand should be called")
	}
	if mockExecutor.AddNumber != 1 {
		t.Errorf("Expected number 1, got %d", mockExecutor.AddNumber)
	}
	expectedDesc := "Test paragraph description"
	if mockExecutor.AddDesc != expectedDesc {
		t.Errorf("Expected description '%s', got '%s'", expectedDesc, mockExecutor.AddDesc)
	}
}

func TestInteractiveShell_AddChoice(t *testing.T) {
	// Setup
	mockExecutor := &MockExecutor{}
	shell := &InteractiveShell{
		executor: mockExecutor,
	}

	// Execute: choiceコマンドを実行
	shell.handleChoice([]string{"1", "Go north", "2"})

	// Verify: エグゼキューターが呼び出されたことを確認
	if !mockExecutor.ChoiceCalled {
		t.Error("ExecuteChoiceCommand should be called")
	}
}

func TestInteractiveShell_ExecuteCommand(t *testing.T) {
	// Setup
	mockExecutor := &MockExecutor{}
	shell := &InteractiveShell{
		executor: mockExecutor,
	}

	tests := []struct {
		name     string
		input    string
		expected bool // 終了すべきかどうか
	}{
		{"help command", "help", false},
		{"exit command", "exit", true},
		{"quit command", "quit", true},
		{"invalid command", "invalid", false},
		{"empty command", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shell.executeCommand(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for command '%s'", tt.expected, result, tt.input)
			}
		})
	}

	// エラーハンドリングテスト
	t.Run("error handling", func(t *testing.T) {
		mockExecutor.ShouldError = true
		result := shell.executeCommand("new test")
		if result {
			t.Error("Should not exit when error occurs")
		}
		if !mockExecutor.NewCalled {
			t.Error("NewCommand should be called even on error")
		}
	})
}

func TestInteractiveShell_GetCurrentPrompt(t *testing.T) {
	// Setup
	mockExecutor := &MockExecutor{}
	shell := &InteractiveShell{
		executor: mockExecutor,
	}

	// Test without current game
	currentGame = nil
	prompt := shell.getCurrentPrompt()
	expected := "> "
	if prompt != expected {
		t.Errorf("Expected prompt '%s', got '%s'", expected, prompt)
	}

	// 港理のため、グローバルな currentGame はテスト終了時にリセット
	defer func() { currentGame = nil }()
}

func TestNewInteractiveShell(t *testing.T) {
	// Execute
	shell := NewInteractiveShell()

	// Verify
	if shell == nil {
		t.Fatal("NewInteractiveShell should return non-nil shell")
	}
	if shell.executor == nil {
		t.Error("Shell should have executor")
	}
}
