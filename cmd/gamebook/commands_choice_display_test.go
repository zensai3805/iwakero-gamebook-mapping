package main

import (
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

// TestFormatChoicesDisplay_ImprovedFormat 選択肢番号表示の改善をテスト
func TestFormatChoicesDisplay_ImprovedFormat(t *testing.T) {
	// Arrange: テスト用のゲームブック作成
	gamebook := &domain.Gamebook{
		Title: "Test Game",
		Paragraphs: map[int]*domain.Paragraph{
			1: {
				Number:      1,
				Description: "森の入り口",
				Visited:     true,
				Choices: []domain.Choice{
					{Description: "北へ進む", TargetNumber: 5, Selected: true},
					{Description: "南へ進む", TargetNumber: 10, Selected: false},
					{Description: "東へ進む", TargetNumber: 15, Selected: false},
				},
			},
		},
		Current: &domain.Paragraph{
			Number:      1,
			Description: "森の入り口",
			Visited:     true,
			Choices: []domain.Choice{
				{Description: "北へ進む", TargetNumber: 5, Selected: true},
				{Description: "南へ進む", TargetNumber: 10, Selected: false},
				{Description: "東へ進む", TargetNumber: 15, Selected: false},
			},
		},
	}

	// Act: 選択肢表示フォーマットを取得
	choicesFormat := formatChoicesDisplay(gamebook.Current.Choices)

	// Assert: 期待するフォーマットで表示されることを確認
	expected := `🎯 選択肢:
  [1] ✅ 北へ進む → #5
  [2] ⚪ 南へ進む → #10
  [3] ⚪ 東へ進む → #15`

	if choicesFormat != expected {
		t.Errorf("期待される選択肢フォーマットと異なります。\n期待値:\n%s\n\n実際の値:\n%s", expected, choicesFormat)
	}
}

// TestFormatChoicesDisplay_NoChoices 選択肢がない場合のテスト
func TestFormatChoicesDisplay_NoChoices(t *testing.T) {
	// Arrange: 選択肢がない場合
	choices := []domain.Choice{}

	// Act: 選択肢表示フォーマットを取得
	choicesFormat := formatChoicesDisplay(choices)

	// Assert: 空文字列が返されることを確認
	expected := ""
	if choicesFormat != expected {
		t.Errorf("選択肢がない場合は空文字列が返されるべきです。実際の値: %s", choicesFormat)
	}
}

// TestFormatChoicesDisplay_SingleChoice 単一の選択肢の場合のテスト
func TestFormatChoicesDisplay_SingleChoice(t *testing.T) {
	// Arrange: 単一の選択肢
	choices := []domain.Choice{
		{Description: "扉を開ける", TargetNumber: 100, Selected: false},
	}

	// Act: 選択肢表示フォーマットを取得
	choicesFormat := formatChoicesDisplay(choices)

	// Assert: 期待するフォーマットで表示されることを確認
	expected := `🎯 選択肢:
  [1] ⚪ 扉を開ける → #100`

	if choicesFormat != expected {
		t.Errorf("単一選択肢のフォーマットが期待値と異なります。\n期待値:\n%s\n\n実際の値:\n%s", expected, choicesFormat)
	}
}

// TestShowCommand_ChoiceDisplayIntegration 実際のshowコマンドでの選択肢表示統合テスト
func TestShowCommand_ChoiceDisplayIntegration(t *testing.T) {
	// テスト用のゲームブック作成 (NewGamebook を使用)
	gamebook := domain.NewGamebook("Test Game")

	// 既存のパラグラフ1を更新
	gamebook.Paragraphs[1].Description = "森の入り口"
	gamebook.Paragraphs[1].Visited = true
	gamebook.Paragraphs[1].Choices = []domain.Choice{
		{Description: "北へ進む", TargetNumber: 5, Selected: true},
		{Description: "南へ進む", TargetNumber: 10, Selected: false},
	}

	// 現在位置を設定
	gamebook.Current = gamebook.Paragraphs[1]

	// グローバル変数を設定
	oldCurrentGame := currentGame
	currentGame = gamebook

	// Act: 選択肢表示部分のみをテスト
	choicesDisplay := formatChoicesDisplay(currentGame.Current.Choices)

	// Assert: 期待するフォーマットが生成されることを確認
	expectedChoiceFormat := `🎯 選択肢:
  [1] ✅ 北へ進む → #5
  [2] ⚪ 南へ進む → #10`

	if choicesDisplay != expectedChoiceFormat {
		t.Errorf("showコマンドの選択肢表示が期待値と異なります。\n期待値:\n%s\n\n実際の値:\n%s", expectedChoiceFormat, choicesDisplay)
	}

	// クリーンアップ
	currentGame = oldCurrentGame
}
