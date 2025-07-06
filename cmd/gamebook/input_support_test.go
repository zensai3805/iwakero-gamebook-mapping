package main

import (
	"testing"

	"github.com/zensai3805/iwakero-gamebook-mapping/internal/domain"
)

// TestSuggestionEngine_GenerateParagraphSuggestions テスト
func TestSuggestionEngine_GenerateParagraphSuggestions(t *testing.T) {
	// テストデータ準備
	gamebook := domain.NewGamebook("TestGame")

	// 既存パラグラフを追加
	_ = gamebook.AddParagraph(&domain.Paragraph{
		Number:      1,
		Description: "開始",
		Choices:     []domain.Choice{},
	})
	_ = gamebook.AddParagraph(&domain.Paragraph{
		Number:      5,
		Description: "森",
		Choices:     []domain.Choice{},
	})
	_ = gamebook.AddParagraph(&domain.Paragraph{
		Number:      10,
		Description: "城",
		Choices:     []domain.Choice{},
	})

	// 未定義参照を追加（選択肢で参照されているが未定義）
	paragraph1, _ := gamebook.GetParagraph(1)
	paragraph1.Choices = []domain.Choice{
		{Description: "森へ", TargetNumber: 5, Selected: true},
		{Description: "山へ", TargetNumber: 20, Selected: false}, // 未定義
	}

	// 候補生成エンジンを作成
	engine := NewSuggestionEngine(gamebook)

	// 候補生成を実行
	suggestions := engine.GenerateParagraphSuggestions()

	// 既存パラグラフが含まれているか確認
	if len(suggestions.Existing) != 3 {
		t.Errorf("期待される既存パラグラフ数: 3, 実際: %d", len(suggestions.Existing))
	}

	// 未定義パラグラフが含まれているか確認
	if len(suggestions.Undefined) != 1 {
		t.Errorf("期待される未定義パラグラフ数: 1, 実際: %d", len(suggestions.Undefined))
	}

	if suggestions.Undefined[0] != 20 {
		t.Errorf("期待される未定義パラグラフ: 20, 実際: %d", suggestions.Undefined[0])
	}

	// 推奨パラグラフが生成されているか確認
	if len(suggestions.Recommended) == 0 {
		t.Error("推奨パラグラフが生成されていません")
	}

	// 次の論理的番号が適切か確認
	if suggestions.Next != 11 {
		t.Errorf("期待される次の論理的番号: 11, 実際: %d", suggestions.Next)
	}
}

// TestTabCompleter_Complete テスト
func TestTabCompleter_Complete(t *testing.T) {
	// テストデータ準備
	gamebook := domain.NewGamebook("TestGame")
	_ = gamebook.AddParagraph(&domain.Paragraph{Number: 1, Description: "開始"})
	_ = gamebook.AddParagraph(&domain.Paragraph{Number: 15, Description: "中間"})
	_ = gamebook.AddParagraph(&domain.Paragraph{Number: 150, Description: "終了"})

	engine := NewSuggestionEngine(gamebook)
	completer := NewTabCompleter(engine)

	// "1" で始まる入力での候補取得
	candidates := completer.Complete("1")

	// 候補の数を確認（推奨候補も含まれるため6個）
	if len(candidates) < 3 {
		t.Errorf("期待される最小候補数: 3, 実際: %d", len(candidates))
	}

	// 基本的な候補の値を確認
	foundValues := make(map[string]bool)
	for _, candidate := range candidates {
		foundValues[candidate.Value] = true
	}

	expectedValues := []string{"1", "15", "150"}
	for _, expected := range expectedValues {
		if !foundValues[expected] {
			t.Errorf("期待される候補値が見つかりません: %s", expected)
		}
	}

	// 候補の型を確認（最初の候補が既存パラグラフであることを確認）
	hasExistingType := false
	for _, candidate := range candidates {
		if candidate.Type == CandidateTypeExisting {
			hasExistingType = true
			break
		}
	}
	if !hasExistingType {
		t.Error("既存パラグラフ型の候補が見つかりません")
	}
}

// TestEnhancedInput_ValidateInput テスト
func TestEnhancedInput_ValidateInput(t *testing.T) {
	// テストデータ準備
	gamebook := domain.NewGamebook("TestGame")
	_ = gamebook.AddParagraph(&domain.Paragraph{Number: 1, Description: "開始"})

	validator := NewInputValidator(gamebook)

	// 既存パラグラフの重複入力
	result := validator.ValidateParagraphNumber(1)
	if result.Type != ValidationTypeDuplicate {
		t.Errorf("期待される検証結果: %v, 実際: %v", ValidationTypeDuplicate, result.Type)
	}

	// 新規パラグラフの入力
	result = validator.ValidateParagraphNumber(2)
	if result.Type != ValidationTypeValid {
		t.Errorf("期待される検証結果: %v, 実際: %v", ValidationTypeValid, result.Type)
	}

	// 無効な入力
	result = validator.ValidateParagraphNumber(0)
	if result.Type != ValidationTypeInvalid {
		t.Errorf("期待される検証結果: %v, 実際: %v", ValidationTypeInvalid, result.Type)
	}
}

// TestHintRenderer_RenderHints テスト
func TestHintRenderer_RenderHints(t *testing.T) {
	// テストデータ準備
	gamebook := domain.NewGamebook("TestGame")
	_ = gamebook.AddParagraph(&domain.Paragraph{Number: 1, Description: "開始"})

	engine := NewSuggestionEngine(gamebook)
	suggestions := engine.GenerateParagraphSuggestions()

	renderer := NewHintRenderer(suggestions)

	// ヒント文字列を生成
	hintStr := renderer.RenderHints("1")

	// ヒント文字列が空でないことを確認
	if hintStr == "" {
		t.Error("ヒント文字列が空です")
	}

	// PTerm要素が含まれていることを確認（簡易チェック）
	if len(hintStr) < 10 {
		t.Error("ヒント文字列が短すぎます")
	}
}

// TestInputSupport_Integration 統合テスト
func TestInputSupport_Integration(t *testing.T) {
	// テストデータ準備
	gamebook := domain.NewGamebook("TestGame")
	_ = gamebook.AddParagraph(&domain.Paragraph{Number: 1, Description: "開始"})

	// 拡張入力コンポーネントを作成
	input := NewEnhancedInput(gamebook)

	// 入力支援機能が正しく初期化されているか確認
	if input.suggestionEngine == nil {
		t.Error("候補生成エンジンが初期化されていません")
	}

	if input.tabCompleter == nil {
		t.Error("Tab補完機能が初期化されていません")
	}

	if input.hintRenderer == nil {
		t.Error("ヒント表示機能が初期化されていません")
	}

	if input.validator == nil {
		t.Error("入力検証機能が初期化されていません")
	}
}

// TestEnhancedInput_ShowWithSuggestions ShowWithSuggestionsのテスト
func TestEnhancedInput_ShowWithSuggestions(t *testing.T) {
	// テストデータ準備
	gamebook := domain.NewGamebook("TestGame")
	_ = gamebook.AddParagraph(&domain.Paragraph{Number: 1, Description: "開始"})
	
	// 拡張入力コンポーネントを作成
	input := NewEnhancedInput(gamebook)
	
	// ヒント生成が正常に動作するか確認
	hints := input.hintRenderer.RenderHints("")
	if hints == "" {
		t.Error("ヒント文字列が生成されていません")
	}
	
	// 候補生成が正常に動作するか確認
	suggestions := input.suggestionEngine.GenerateParagraphSuggestions()
	if len(suggestions.Existing) != 1 {
		t.Errorf("期待される既存パラグラフ数: 1, 実際: %d", len(suggestions.Existing))
	}
	
	// 入力検証が正常に動作するか確認
	validationResult := input.validator.ValidateParagraphNumber(1)
	if validationResult.IsValid {
		t.Error("重複パラグラフが有効として判定されています")
	}
	
	// 新規パラグラフの検証
	validationResult = input.validator.ValidateParagraphNumber(2)
	if !validationResult.IsValid {
		t.Error("新規パラグラフが無効として判定されています")
	}
}
