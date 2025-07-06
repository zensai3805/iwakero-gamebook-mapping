package main

import (
	"fmt"
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

// TestInputHelper_GetDefaultText tests the default text generation for current location
func TestInputHelper_GetDefaultText(t *testing.T) {
	tests := []struct {
		name        string
		currentGame *domain.Gamebook
		want        string
	}{
		{
			name:        "現在地が設定されている場合",
			currentGame: createGamebookWithCurrentParagraph(15),
			want:        "15",
		},
		{
			name:        "ゲームが読み込まれていない場合",
			currentGame: nil,
			want:        "",
		},
		{
			name:        "現在地が設定されていない場合",
			currentGame: createGamebookWithoutCurrent(),
			want:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := NewInputHelper(tt.currentGame)
			got := helper.GetDefaultText()
			if got != tt.want {
				t.Errorf("GetDefaultText() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestInputHelper_GetCLIPrompt tests CLI prompt generation with current location
func TestInputHelper_GetCLIPrompt(t *testing.T) {
	tests := []struct {
		name        string
		prompt      string
		currentGame *domain.Gamebook
		want        string
	}{
		{
			name:        "現在地がある場合はプロンプトに含める",
			prompt:      "パラグラフ番号",
			currentGame: createGamebookWithCurrentParagraph(15),
			want:        "パラグラフ番号 [現在地: 15]: ",
		},
		{
			name:        "現在地がない場合は通常のプロンプト",
			prompt:      "パラグラフ番号",
			currentGame: nil,
			want:        "パラグラフ番号: ",
		},
		{
			name:        "ゲームはあるが現在地が未設定",
			prompt:      "パラグラフ番号",
			currentGame: createGamebookWithoutCurrent(),
			want:        "パラグラフ番号: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := NewInputHelper(tt.currentGame)
			got := helper.GetCLIPrompt(tt.prompt)
			if got != tt.want {
				t.Errorf("GetCLIPrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestInputHelper_GetSuggestions tests suggestion generation
func TestInputHelper_GetSuggestions(t *testing.T) {
	tests := []struct {
		name        string
		currentGame *domain.Gamebook
		want        []string
	}{
		{
			name:        "現在地がある場合は候補として返す",
			currentGame: createGamebookWithCurrentParagraph(15),
			want:        []string{"15"},
		},
		{
			name:        "現在地がない場合は空の候補",
			currentGame: nil,
			want:        []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := NewInputHelper(tt.currentGame)
			got := helper.GetSuggestions()
			if len(got) != len(tt.want) {
				t.Errorf("GetSuggestions() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("GetSuggestions()[%d] = %v, want %v", i, v, tt.want[i])
				}
			}
		})
	}
}

// TestInputHelper_ProcessEmptyInput tests empty input processing
func TestInputHelper_ProcessEmptyInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		currentGame *domain.Gamebook
		want        string
	}{
		{
			name:        "空入力で現在地がある場合は自動入力",
			input:       "",
			currentGame: createGamebookWithCurrentParagraph(15),
			want:        "15",
		},
		{
			name:        "値が入力されている場合はそのまま",
			input:       "20",
			currentGame: createGamebookWithCurrentParagraph(25),
			want:        "20",
		},
		{
			name:        "空入力で現在地がない場合はそのまま",
			input:       "",
			currentGame: nil,
			want:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := NewInputHelper(tt.currentGame)
			got := helper.ProcessEmptyInput(tt.input)
			if got != tt.want {
				t.Errorf("ProcessEmptyInput(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// createGamebookWithCurrentParagraph creates a gamebook with a specified current paragraph
func createGamebookWithCurrentParagraph(number int) *domain.Gamebook {
	gb := domain.NewGamebook("test")
	p := domain.NewParagraph(number, fmt.Sprintf("Paragraph %d", number))
	_ = gb.AddParagraph(p)
	_ = gb.MoveTo(number)
	return gb
}

// createGamebookWithoutCurrent creates a gamebook without current paragraph
func createGamebookWithoutCurrent() *domain.Gamebook {
	gb := domain.NewGamebook("test")
	gb.Current = nil // テスト用に現在地をクリア
	return gb
}
