package domain

import (
	"reflect"
	"testing"
)

func TestUndefinedAnalyzer_AnalyzeUndefinedParagraphs(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *Gamebook
		currentPos  int
		expected    *UndefinedAnalysis
		description string
	}{
		{
			name: "接続された未定義パラグラフのみ検出",
			setup: func() *Gamebook {
				gb := NewGamebook("Test")
				// パラグラフ1を追加（現在位置）
				p1 := NewParagraph(1, "冒険の始まり")
				_ = gb.AddParagraph(p1)
				// パラグラフ1から未定義パラグラフ10への選択肢
				_ = gb.AddChoiceToParagraph(1, "北へ進む", 10)

				// パラグラフ5を追加（接続なし）
				p5 := NewParagraph(5, "別の場所")
				_ = gb.AddParagraph(p5)
				// パラグラフ5から未定義パラグラフ20への選択肢
				_ = gb.AddChoiceToParagraph(5, "東へ進む", 20)

				return gb
			},
			currentPos: 1,
			expected: &UndefinedAnalysis{
				Connected: []int{10},
				Orphaned:  []int{20},
			},
			description: "現在位置1から接続された未定義パラグラフ10のみ検出、20は孤立",
		},
		{
			name: "複数の接続された未定義パラグラフ",
			setup: func() *Gamebook {
				gb := NewGamebook("Test")
				p1 := NewParagraph(1, "冒険の始まり")
				_ = gb.AddParagraph(p1)
				_ = gb.AddChoiceToParagraph(1, "北へ進む", 10)
				_ = gb.AddChoiceToParagraph(1, "南へ進む", 15)
				_ = gb.AddChoiceToParagraph(1, "東へ進む", 20)

				return gb
			},
			currentPos: 1,
			expected: &UndefinedAnalysis{
				Connected: []int{10, 15, 20},
				Orphaned:  []int{},
			},
			description: "現在位置から複数の未定義パラグラフが接続",
		},
		{
			name: "未定義パラグラフが存在しない場合",
			setup: func() *Gamebook {
				gb := NewGamebook("Test")
				p1 := NewParagraph(1, "冒険の始まり")
				p2 := NewParagraph(2, "次の場所")
				_ = gb.AddParagraph(p1)
				_ = gb.AddParagraph(p2)
				_ = gb.AddChoiceToParagraph(1, "進む", 2)

				return gb
			},
			currentPos: 1,
			expected: &UndefinedAnalysis{
				Connected: []int{},
				Orphaned:  []int{},
			},
			description: "全てのパラグラフが定義済みの場合",
		},
		{
			name: "現在位置がnilの場合のゼロ処理",
			setup: func() *Gamebook {
				gb := NewGamebook("Test")
				p1 := NewParagraph(1, "冒険の始まり")
				_ = gb.AddParagraph(p1)
				_ = gb.AddChoiceToParagraph(1, "北へ進む", 10)

				return gb
			},
			currentPos: 0, // 現在位置なし
			expected: &UndefinedAnalysis{
				Connected: []int{},
				Orphaned:  []int{10},
			},
			description: "現在位置がない場合、全ての未定義パラグラフが孤立",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gb := tt.setup()
			analyzer := NewUndefinedAnalyzer(gb)

			result := analyzer.AnalyzeUndefinedParagraphs(tt.currentPos)

			if !reflect.DeepEqual(result.Connected, tt.expected.Connected) {
				t.Errorf("Connected mismatch.\nExpected: %v\nGot: %v", tt.expected.Connected, result.Connected)
			}

			if !reflect.DeepEqual(result.Orphaned, tt.expected.Orphaned) {
				t.Errorf("Orphaned mismatch.\nExpected: %v\nGot: %v", tt.expected.Orphaned, result.Orphaned)
			}
		})
	}
}

func TestUndefinedAnalyzer_IsConnectedToUndefined(t *testing.T) {
	tests := []struct {
		name         string
		setup        func() *Gamebook
		currentPos   int
		targetNumber int
		expected     bool
		description  string
	}{
		{
			name: "直接接続された未定義パラグラフ",
			setup: func() *Gamebook {
				gb := NewGamebook("Test")
				p1 := NewParagraph(1, "冒険の始まり")
				_ = gb.AddParagraph(p1)
				_ = gb.AddChoiceToParagraph(1, "北へ進む", 10)

				return gb
			},
			currentPos:   1,
			targetNumber: 10,
			expected:     true,
			description:  "現在位置から直接接続された未定義パラグラフ",
		},
		{
			name: "接続されていない未定義パラグラフ",
			setup: func() *Gamebook {
				gb := NewGamebook("Test")
				p1 := NewParagraph(1, "冒険の始まり")
				p5 := NewParagraph(5, "別の場所")
				_ = gb.AddParagraph(p1)
				_ = gb.AddParagraph(p5)
				_ = gb.AddChoiceToParagraph(5, "東へ進む", 20)

				return gb
			},
			currentPos:   1,
			targetNumber: 20,
			expected:     false,
			description:  "現在位置から接続されていない未定義パラグラフ",
		},
		{
			name: "定義済みパラグラフ",
			setup: func() *Gamebook {
				gb := NewGamebook("Test")
				p1 := NewParagraph(1, "冒険の始まり")
				p2 := NewParagraph(2, "次の場所")
				_ = gb.AddParagraph(p1)
				_ = gb.AddParagraph(p2)
				_ = gb.AddChoiceToParagraph(1, "進む", 2)

				return gb
			},
			currentPos:   1,
			targetNumber: 2,
			expected:     false,
			description:  "定義済みパラグラフは未定義ではない",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gb := tt.setup()
			analyzer := NewUndefinedAnalyzer(gb)

			result := analyzer.IsConnectedToUndefined(tt.currentPos, tt.targetNumber)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v for %s", tt.expected, result, tt.description)
			}
		})
	}
}
