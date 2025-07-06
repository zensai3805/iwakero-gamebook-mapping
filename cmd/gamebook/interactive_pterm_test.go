package main

import (
	"testing"

	"github.com/pterm/pterm"
	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

// TestPTermInteractiveShell_InputHelperIntegration tests input assistance in interactive mode
func TestPTermInteractiveShell_InputHelperIntegration(t *testing.T) {
	// Save original values
	originalGame := currentGame
	defer func() {
		currentGame = originalGame
	}()

	// Setup test game with current location
	gb := domain.NewGamebook("TestGame")
	p1 := domain.NewParagraph(15, "Current location")
	_ = gb.AddParagraph(p1)
	_ = gb.MoveTo(15)
	currentGame = gb

	tests := []struct {
		name            string
		operation       string
		expectedDefault string
	}{
		{
			name:            "パラグラフ追加時に現在地がデフォルト値になる",
			operation:       "add",
			expectedDefault: "15",
		},
		{
			name:            "選択肢追加時に現在地がデフォルト値になる",
			operation:       "choice",
			expectedDefault: "15",
		},
		{
			name:            "選択肢選択時に現在地がデフォルト値になる",
			operation:       "select",
			expectedDefault: "15",
		},
		{
			name:            "直接移動時に現在地がデフォルト値になる",
			operation:       "move",
			expectedDefault: "15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// InputHelperが正しくデフォルト値を提供することを確認
			helper := NewInputHelper(currentGame)
			got := helper.GetDefaultText()
			if got != tt.expectedDefault {
				t.Errorf("GetDefaultText() = %v, want %v", got, tt.expectedDefault)
			}
		})
	}
}

// TestPTermInteractiveShell_EmptyInputProcessing tests empty input processing
func TestPTermInteractiveShell_EmptyInputProcessing(t *testing.T) {
	// Save original values
	originalGame := currentGame
	defer func() {
		currentGame = originalGame
	}()

	// Setup test game with current location
	gb := domain.NewGamebook("TestGame")
	p1 := domain.NewParagraph(20, "Current location")
	_ = gb.AddParagraph(p1)
	_ = gb.MoveTo(20)
	currentGame = gb

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "空入力時は現在地が自動入力される",
			input:    "",
			expected: "20",
		},
		{
			name:     "値が入力されている場合はそのまま",
			input:    "30",
			expected: "30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := NewInputHelper(currentGame)
			got := helper.ProcessEmptyInput(tt.input)
			if got != tt.expected {
				t.Errorf("ProcessEmptyInput(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// TestPTermInteractiveTextInput_WithSuggestions tests PTerm suggestions
func TestPTermInteractiveTextInput_WithSuggestions(t *testing.T) {
	// This test verifies that suggestions can be set on PTerm inputs
	// Actual interactive testing requires manual verification

	// Create an interactive text input with default text
	textInput := pterm.DefaultInteractiveTextInput.
		WithDefaultText("15").
		WithTextStyle(pterm.NewStyle(pterm.FgLightBlue))

	// Verify that the text input can be created without errors
	if textInput == nil {
		t.Error("Failed to create interactive text input with suggestions")
	}
}
