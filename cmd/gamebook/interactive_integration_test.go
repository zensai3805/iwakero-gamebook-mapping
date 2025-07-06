package main

import (
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

// テスト用にcurrentGameを別のローカル変数として使用
var testCurrentGame *domain.Gamebook

// TestInteractiveInputIntegration 対話式入力の統合テスト
func TestInteractiveInputIntegration(t *testing.T) {
	// テストデータ準備
	gamebook := domain.NewGamebook("TestGame")
	_ = gamebook.AddParagraph(&domain.Paragraph{
		Number:      1,
		Description: "開始",
		Choices:     []domain.Choice{},
	})

	// テスト用変数を設定
	testCurrentGame = gamebook

	// 入力支援機能が利用可能か確認
	enhancedInput := NewEnhancedInput(testCurrentGame)

	if enhancedInput == nil {
		t.Fatal("拡張入力コンポーネントの作成に失敗")
	}

	// 候補生成が正常に動作するか確認
	suggestions := enhancedInput.suggestionEngine.GenerateParagraphSuggestions()

	if len(suggestions.Existing) != 1 {
		t.Errorf("期待される既存パラグラフ数: 1, 実際: %d", len(suggestions.Existing))
	}

	if suggestions.Existing[0] != 1 {
		t.Errorf("期待される既存パラグラフ: 1, 実際: %d", suggestions.Existing[0])
	}

	// 次の番号提案機能は削除された
	if suggestions.Next != 0 {
		t.Errorf("次の番号提案は使用しません。実際: %d", suggestions.Next)
	}

	// ヒント表示が機能するか確認
	hints := enhancedInput.hintRenderer.RenderHints("")
	if hints == "" {
		t.Error("ヒント文字列が生成されていません")
	}

	// 入力検証が機能するか確認（既存パラグラフは有効）
	validationResult := enhancedInput.validator.ValidateParagraphNumber(1)
	if !validationResult.IsValid {
		t.Error("既存パラグラフが無効として判定されています")
	}

	if validationResult.Type != ValidationTypeValid {
		t.Errorf("期待される検証結果: %v, 実際: %v", ValidationTypeValid, validationResult.Type)
	}

	// 新規パラグラフの検証
	validationResult = enhancedInput.validator.ValidateParagraphNumber(2)
	if !validationResult.IsValid {
		t.Error("新規パラグラフが無効として判定されています")
	}

	// クリーンアップ
	testCurrentGame = nil
}

// TestInputSupportPerformance 入力支援機能のパフォーマンステスト
func TestInputSupportPerformance(t *testing.T) {
	// 大きなゲームブックでのパフォーマンステスト
	gamebook := domain.NewGamebook("LargeTestGame")

	// 100個のパラグラフを追加
	for i := 1; i <= 100; i++ {
		_ = gamebook.AddParagraph(&domain.Paragraph{
			Number:      i,
			Description: "テストパラグラフ",
			Choices:     []domain.Choice{},
		})
	}

	// 入力支援機能を作成
	enhancedInput := NewEnhancedInput(gamebook)

	// 候補生成のパフォーマンステスト
	suggestions := enhancedInput.suggestionEngine.GenerateParagraphSuggestions()

	if len(suggestions.Existing) != 100 {
		t.Errorf("期待される既存パラグラフ数: 100, 実際: %d", len(suggestions.Existing))
	}

	// Tab補完のパフォーマンステスト
	candidates := enhancedInput.tabCompleter.Complete("1")

	// "1"で始まるパラグラフは1, 10-19, 100なので最低11個
	if len(candidates) < 11 {
		t.Errorf("期待される最小候補数: 11, 実際: %d", len(candidates))
	}

	// ヒント表示のパフォーマンステスト
	hints := enhancedInput.hintRenderer.RenderHints("")
	if hints == "" {
		t.Error("ヒント文字列が生成されていません")
	}
}
