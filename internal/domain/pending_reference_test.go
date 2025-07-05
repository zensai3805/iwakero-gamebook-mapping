package domain

import (
	"testing"
)

// TestPendingReferenceManager_AddReference 保留参照の追加テスト
func TestPendingReferenceManager_AddReference(t *testing.T) {
	manager := NewPendingReferenceManager()

	// 未定義段落への参照を追加
	pendingErr := manager.AddReference(1, "洞窟に入る", 23)
	if pendingErr != nil {
		t.Errorf("保留参照の追加に失敗: %v", pendingErr)
	}

	// 保留参照が正しく追加されていることを確認
	references := manager.GetReferences(23)
	if len(references) != 1 {
		t.Errorf("保留参照数が不正: expected 1, got %d", len(references))
	}

	ref := references[0]
	if ref.FromParagraph != 1 {
		t.Errorf("参照元段落が不正: expected 1, got %d", ref.FromParagraph)
	}
	if ref.ChoiceDescription != "洞窟に入る" {
		t.Errorf("選択肢説明が不正: expected '洞窟に入る', got '%s'", ref.ChoiceDescription)
	}
	if ref.TargetParagraph != 23 {
		t.Errorf("参照先段落が不正: expected 23, got %d", ref.TargetParagraph)
	}
}

// TestPendingReferenceManager_ResolveReference 保留参照の解決テスト
func TestPendingReferenceManager_ResolveReference(t *testing.T) {
	manager := NewPendingReferenceManager()

	// 保留参照を追加
	addErr := manager.AddReference(1, "洞窟に入る", 23)
	if addErr != nil {
		t.Fatalf("保留参照の追加に失敗: %v", addErr)
	}

	// 参照を解決
	resolveErr := manager.ResolveReference(23)
	if resolveErr != nil {
		t.Errorf("保留参照の解決に失敗: %v", resolveErr)
	}

	// 解決後は保留参照が存在しないことを確認
	references := manager.GetReferences(23)
	if len(references) != 0 {
		t.Errorf("解決後に保留参照が残存: expected 0, got %d", len(references))
	}
}

// TestPendingReferenceManager_GetAllPendingTargets 全保留対象の取得テスト
func TestPendingReferenceManager_GetAllPendingTargets(t *testing.T) {
	manager := NewPendingReferenceManager()

	// 複数の保留参照を追加
	addErr1 := manager.AddReference(1, "洞窟に入る", 23)
	if addErr1 != nil {
		t.Fatalf("保留参照1の追加に失敗: %v", addErr1)
	}

	addErr2 := manager.AddReference(2, "森に進む", 45)
	if addErr2 != nil {
		t.Fatalf("保留参照2の追加に失敗: %v", addErr2)
	}

	addErr3 := manager.AddReference(3, "街に戻る", 23) // 同じ対象への複数参照
	if addErr3 != nil {
		t.Fatalf("保留参照3の追加に失敗: %v", addErr3)
	}

	// 全保留対象を取得
	targets := manager.GetAllPendingTargets()
	if len(targets) != 2 {
		t.Errorf("保留対象数が不正: expected 2, got %d", len(targets))
	}

	// 対象に23と45が含まれることを確認
	targetMap := make(map[int]bool)
	for _, target := range targets {
		targetMap[target] = true
	}

	if !targetMap[23] || !targetMap[45] {
		t.Errorf("保留対象が不正: expected [23, 45], got %v", targets)
	}
}
