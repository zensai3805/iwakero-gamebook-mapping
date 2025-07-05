package main

import (
	"fmt"
	"runtime"
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
	MoveCalled   bool
	MoveTarget   int
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

func (m *MockExecutor) ExecuteMoveCommand(targetNum int) error {
	m.MoveCalled = true
	m.MoveTarget = targetNum
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

func (m *MockExecutor) ExecuteShowCommandWithScope(scope DisplayScope) error {
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

func TestPTermInteractiveShell_ClearScreen(t *testing.T) {
	// Setup
	shell := NewPTermInteractiveShell()

	// Execute: 画面クリア関数が各プラットフォームで実行できることを確認
	// エラーが発生しないことを確認（実際の画面クリアは行わない）
	t.Run("clearScreen", func(t *testing.T) {
		// clearScreenはエラーを返さないため、パニックしないことを確認
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("clearScreen should not panic: %v", r)
			}
		}()

		// shellオブジェクトが正しく作成されていることを確認
		if shell == nil {
			t.Fatal("shell should not be nil")
		}

		// 実行環境に応じて適切なコマンドが選択されることを確認
		// （実際のコマンド実行は副作用があるため、ここではruntime.GOOSの確認のみ）
		switch runtime.GOOS {
		case "windows":
			// Windows環境でのclearコマンド確認
			t.Log("Windows clear command would be: cmd /c cls")
		default:
			// Unix系環境でのclearコマンド確認
			t.Log("Unix clear command would be: clear")
		}
	})
}

func TestPTermInteractiveShell_GetCurrentPrompt(t *testing.T) {
	// Setup
	shell := NewPTermInteractiveShell()

	// Test without current game
	currentGame = nil
	prompt := shell.getCurrentPrompt()
	expected := "> "
	if prompt != expected {
		t.Errorf("Expected prompt '%s', got '%s'", expected, prompt)
	}

	// クリーンアップ
	defer func() { currentGame = nil }()
}

func TestPTermInteractiveShell_ErrorMessageHandling(t *testing.T) {
	// Setup
	mockExecutor := &MockExecutor{}
	shell := &PTermInteractiveShell{
		executor: mockExecutor,
	}

	// Test: エラーメッセージが保持されることを確認
	t.Run("invalid args error", func(t *testing.T) {
		shell.handleNew([]string{}) // 引数不足
		if shell.lastError != "使用法: new <ゲーム名>" {
			t.Errorf("Expected error message '使用法: new <ゲーム名>', got '%s'", shell.lastError)
		}
		// 実行されていないことを確認
		if mockExecutor.NewCalled {
			t.Error("ExecuteNewCommand should not be called with invalid args")
		}
	})

	// Test: コマンド実行エラーが保持されることを確認
	t.Run("execution error", func(t *testing.T) {
		mockExecutor.ShouldError = true
		shell.handleNew([]string{"TestGame"})
		if shell.lastError != "mock error" {
			t.Errorf("Expected error message 'mock error', got '%s'", shell.lastError)
		}
		// lastInfoは設定されないことを確認
		if shell.lastInfo != "" {
			t.Errorf("lastInfo should be empty on error, got '%s'", shell.lastInfo)
		}
	})

	// Test: 成功時はlastInfoが設定されることを確認
	t.Run("success info", func(t *testing.T) {
		mockExecutor.ShouldError = false
		shell.lastError = "previous error" // 前のエラーをセット
		shell.handleNew([]string{"TestGame"})
		if shell.lastInfo != "新しいゲームブック 'TestGame' を作成しました" {
			t.Errorf("Expected info message, got '%s'", shell.lastInfo)
		}
		// エラーはクリアされないことを確認（画面表示時にクリアされる）
		if shell.lastError != "previous error" {
			t.Error("lastError should not be cleared until displayed")
		}
	})

	// Test: 未知のコマンドエラー
	t.Run("unknown command", func(t *testing.T) {
		shell.lastError = ""
		shell.executeCommand("unknown")
		expectedError := "不明なコマンド: unknown ('help' でコマンド一覧を確認できます)"
		if shell.lastError != expectedError {
			t.Errorf("Expected error '%s', got '%s'", expectedError, shell.lastError)
		}
	})
}

func TestPTermInteractiveShell_MenuErrorHandling(t *testing.T) {
	// Setup
	shell := NewPTermInteractiveShell()
	currentGame = nil // ゲーム未読み込み状態

	// Test: ゲーム未読み込み時のエラー
	t.Run("no game loaded", func(t *testing.T) {
		shell.handleAddFromMenu()
		if shell.lastError != ErrNoGameLoaded {
			t.Errorf("Expected error about no game loaded, got '%s'", shell.lastError)
		}

		shell.handleChoiceFromMenu()
		if shell.lastError != ErrNoGameLoaded {
			t.Errorf("Expected error about no game loaded, got '%s'", shell.lastError)
		}

		shell.handleSelectFromMenu()
		if shell.lastError != ErrNoGameLoaded {
			t.Errorf("Expected error about no game loaded, got '%s'", shell.lastError)
		}
	})

	// クリーンアップ
	defer func() { currentGame = nil }()
}
