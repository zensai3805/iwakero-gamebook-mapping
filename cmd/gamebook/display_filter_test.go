package main

import (
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

func TestDisplayScope_String(t *testing.T) {
	tests := []struct {
		name     string
		scope    DisplayScope
		expected string
	}{
		{
			name:     "Connected scope",
			scope:    ScopeConnected,
			expected: "connected",
		},
		{
			name:     "All scope",
			scope:    ScopeAll,
			expected: "all",
		},
		{
			name:     "None scope",
			scope:    ScopeNone,
			expected: "none",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := string(tt.scope)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestDisplayFilter_ShouldDisplayUndefined(t *testing.T) {
	tests := []struct {
		name            string
		scope           DisplayScope
		setup           func() (*domain.Gamebook, *domain.UndefinedAnalysis)
		targetParagraph int
		expected        bool
		description     string
	}{
		{
			name:  "Connected scope - 接続された未定義パラグラフを表示",
			scope: ScopeConnected,
			setup: func() (*domain.Gamebook, *domain.UndefinedAnalysis) {
				gb := domain.NewGamebook("Test")
				p1 := domain.NewParagraph(1, "開始")
				_ = gb.AddParagraph(p1)
				_ = gb.AddChoiceToParagraph(1, "進む", 10)

				analyzer := domain.NewUndefinedAnalyzer(gb)
				analysis := analyzer.AnalyzeUndefinedParagraphs(1)
				return gb, analysis
			},
			targetParagraph: 10,
			expected:        true,
			description:     "接続された未定義パラグラフは表示する",
		},
		{
			name:  "Connected scope - 孤立した未定義パラグラフを非表示",
			scope: ScopeConnected,
			setup: func() (*domain.Gamebook, *domain.UndefinedAnalysis) {
				gb := domain.NewGamebook("Test")
				p1 := domain.NewParagraph(1, "開始")
				p5 := domain.NewParagraph(5, "別の場所")
				_ = gb.AddParagraph(p1)
				_ = gb.AddParagraph(p5)
				_ = gb.AddChoiceToParagraph(5, "進む", 20)

				analyzer := domain.NewUndefinedAnalyzer(gb)
				analysis := analyzer.AnalyzeUndefinedParagraphs(1)
				return gb, analysis
			},
			targetParagraph: 20,
			expected:        false,
			description:     "孤立した未定義パラグラフは非表示にする",
		},
		{
			name:  "All scope - 全ての未定義パラグラフを表示",
			scope: ScopeAll,
			setup: func() (*domain.Gamebook, *domain.UndefinedAnalysis) {
				gb := domain.NewGamebook("Test")
				p1 := domain.NewParagraph(1, "開始")
				p5 := domain.NewParagraph(5, "別の場所")
				_ = gb.AddParagraph(p1)
				_ = gb.AddParagraph(p5)
				_ = gb.AddChoiceToParagraph(5, "進む", 20)

				analyzer := domain.NewUndefinedAnalyzer(gb)
				analysis := analyzer.AnalyzeUndefinedParagraphs(1)
				return gb, analysis
			},
			targetParagraph: 20,
			expected:        true,
			description:     "All scopeでは全ての未定義パラグラフを表示",
		},
		{
			name:  "None scope - 全ての未定義パラグラフを非表示",
			scope: ScopeNone,
			setup: func() (*domain.Gamebook, *domain.UndefinedAnalysis) {
				gb := domain.NewGamebook("Test")
				p1 := domain.NewParagraph(1, "開始")
				_ = gb.AddParagraph(p1)
				_ = gb.AddChoiceToParagraph(1, "進む", 10)

				analyzer := domain.NewUndefinedAnalyzer(gb)
				analysis := analyzer.AnalyzeUndefinedParagraphs(1)
				return gb, analysis
			},
			targetParagraph: 10,
			expected:        false,
			description:     "None scopeでは全ての未定義パラグラフを非表示",
		},
		{
			name:  "定義済みパラグラフは常に表示",
			scope: ScopeNone,
			setup: func() (*domain.Gamebook, *domain.UndefinedAnalysis) {
				gb := domain.NewGamebook("Test")
				p1 := domain.NewParagraph(1, "開始")
				p2 := domain.NewParagraph(2, "次")
				_ = gb.AddParagraph(p1)
				_ = gb.AddParagraph(p2)

				analyzer := domain.NewUndefinedAnalyzer(gb)
				analysis := analyzer.AnalyzeUndefinedParagraphs(1)
				return gb, analysis
			},
			targetParagraph: 2,
			expected:        true,
			description:     "定義済みパラグラフはスコープに関係なく表示",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gb, analysis := tt.setup()
			filter := NewDisplayFilter(tt.scope, analysis)

			result := filter.ShouldDisplayUndefined(gb, tt.targetParagraph)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v for %s", tt.expected, result, tt.description)
			}
		})
	}
}

func TestDisplayFilter_FilterChoices(t *testing.T) {
	tests := []struct {
		name        string
		scope       DisplayScope
		setup       func() (*domain.Gamebook, *domain.UndefinedAnalysis, []domain.Choice)
		expected    []domain.Choice
		description string
	}{
		{
			name:  "Connected scope - 接続された未定義選択肢のみ表示",
			scope: ScopeConnected,
			setup: func() (*domain.Gamebook, *domain.UndefinedAnalysis, []domain.Choice) {
				gb := domain.NewGamebook("Test")
				p1 := domain.NewParagraph(1, "開始")
				p2 := domain.NewParagraph(2, "定義済み")
				_ = gb.AddParagraph(p1)
				_ = gb.AddParagraph(p2)
				_ = gb.AddChoiceToParagraph(1, "接続未定義", 10)
				_ = gb.AddChoiceToParagraph(1, "定義済み", 2)

				// 孤立した未定義パラグラフ
				p5 := domain.NewParagraph(5, "別の場所")
				_ = gb.AddParagraph(p5)
				_ = gb.AddChoiceToParagraph(5, "孤立未定義", 20)

				analyzer := domain.NewUndefinedAnalyzer(gb)
				analysis := analyzer.AnalyzeUndefinedParagraphs(1)

				choices := []domain.Choice{
					{Description: "接続未定義", TargetNumber: 10, Selected: false},
					{Description: "定義済み", TargetNumber: 2, Selected: false},
					{Description: "孤立未定義", TargetNumber: 20, Selected: false},
				}

				return gb, analysis, choices
			},
			expected: []domain.Choice{
				{Description: "接続未定義", TargetNumber: 10, Selected: false},
				{Description: "定義済み", TargetNumber: 2, Selected: false},
			},
			description: "接続された未定義選択肢と定義済み選択肢のみ表示",
		},
		{
			name:  "None scope - 未定義選択肢を全て非表示",
			scope: ScopeNone,
			setup: func() (*domain.Gamebook, *domain.UndefinedAnalysis, []domain.Choice) {
				gb := domain.NewGamebook("Test")
				p1 := domain.NewParagraph(1, "開始")
				p2 := domain.NewParagraph(2, "定義済み")
				_ = gb.AddParagraph(p1)
				_ = gb.AddParagraph(p2)
				_ = gb.AddChoiceToParagraph(1, "未定義", 10)
				_ = gb.AddChoiceToParagraph(1, "定義済み", 2)

				analyzer := domain.NewUndefinedAnalyzer(gb)
				analysis := analyzer.AnalyzeUndefinedParagraphs(1)

				choices := []domain.Choice{
					{Description: "未定義", TargetNumber: 10, Selected: false},
					{Description: "定義済み", TargetNumber: 2, Selected: false},
				}

				return gb, analysis, choices
			},
			expected: []domain.Choice{
				{Description: "定義済み", TargetNumber: 2, Selected: false},
			},
			description: "未定義選択肢を全て非表示",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gb, analysis, choices := tt.setup()
			filter := NewDisplayFilter(tt.scope, analysis)

			result := filter.FilterChoices(gb, choices)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d choices, got %d for %s", len(tt.expected), len(result), tt.description)
				return
			}

			for i, expectedChoice := range tt.expected {
				if result[i].Description != expectedChoice.Description ||
					result[i].TargetNumber != expectedChoice.TargetNumber ||
					result[i].Selected != expectedChoice.Selected {
					t.Errorf("Choice %d mismatch. Expected %+v, got %+v", i, expectedChoice, result[i])
				}
			}
		})
	}
}
