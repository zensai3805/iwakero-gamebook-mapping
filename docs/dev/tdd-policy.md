# テスト駆動開発（TDD）方針

## 最重要原則

### 🚨 テスト駆動開発（TDD）の徹底
**t_wadaメソッドによるテストファーストを必須とする**

1. **RED**: 失敗するテストを最初に作成
2. **GREEN**: テストが通る最低限のコードを実装
3. **REFACTOR**: コードを改善・整理

- **全機能はテストファーストで実装**
- テストコードは実装コードと同じ重要度で扱う
- TDDを守らない実装は認めない

## TDD実践の手順

### 1. テストファースト
```go
// ❌ 非推奨: 実装を先に書く
func NewFeature() {
    // 実装...
}

// ✅ 推奨: テストを先に書く
func TestNewFeature(t *testing.T) {
    // 期待する振る舞いを定義
    result := NewFeature()
    assert.Equal(t, expected, result)
}
```

### 2. 最小限の実装
- テストが通る最小限のコードのみ実装
- 過度な一般化は避ける
- YAGNI原則を守る

### 3. リファクタリング
- テストが通る状態を維持しながら改善
- コードの重複を削除
- 可読性の向上

## テストの種類と方針

### 単体テスト
- **モック不要の層**: Entities Layer
- **モック使用の層**: Usecase Layer以上
- **カバレッジ目標**: 80%以上

### 結合テスト
- Interface Adapters Layerで実施
- 外部システムとの連携を検証
- 実際のファイルシステムを使用

### 統合テスト
- Frameworks & Drivers Layerで実施
- エンドツーエンドのシナリオテスト
- ユーザーワークフローの検証

## テストコードの品質基準

### Logger活用のTDD方針
**実装中はLogger活用でデバッグ効率向上**

```go
// ✅ 推奨: テスト実行中のログ出力で問題箇所特定
func TestAddParagraph_WhenDuplicate_ReturnsError(t *testing.T) {
    // AI開発モード有効化（テスト専用）
    os.Setenv("GAMEBOOK_AI_DEV", "true")
    os.Setenv("LOG_LEVEL", "DEBUG")
    defer os.Unsetenv("GAMEBOOK_AI_DEV")
    defer os.Unsetenv("LOG_LEVEL")
    
    // Arrange
    gamebook := domain.NewGamebook("test")
    paragraph := &domain.Paragraph{Number: 1, Description: "テスト"}
    
    // 初回追加（成功ケース）
    err := gamebook.AddParagraph(paragraph)
    assert.NoError(t, err)
    
    // Act & Assert（重複追加でエラー）
    duplicateErr := gamebook.AddParagraph(paragraph)
    assert.Error(t, duplicateErr)
    assert.Contains(t, duplicateErr.Error(), "重複")
}
```

### 命名規則
```go
// テスト関数名: Test対象関数名_条件_期待結果
func TestAddParagraph_WhenDuplicate_ReturnsError(t *testing.T) {
    // ...
}
```

### AAA原則
```go
func TestExample(t *testing.T) {
    // Arrange（準備）
    gamebook := domain.NewGamebook("test")
    
    // Act（実行）
    err := gamebook.AddParagraph(paragraph)
    
    // Assert（検証）
    assert.NoError(t, err)
}
```

### テストの独立性
- 各テストは独立して実行可能
- 順序依存なし
- 共有状態なし

## TDD違反の例と対策

### 違反例1: 実装後のテスト追加
```
❌ 実装 → テスト作成
✅ テスト作成 → 実装 → リファクタリング
```

### 違反例2: テストなしのバグ修正
```
❌ バグ修正 → 動作確認
✅ バグ再現テスト作成 → 修正 → テスト通過確認
```

### 違反例3: カバレッジのためのテスト
```
❌ カバレッジ率を上げるためだけのテスト
✅ 振る舞いを検証する意味のあるテスト
```

### 違反例4: ログ出力の検証不足
```
❌ ログ機能実装 → 目視確認なしで「動作確認済み」
✅ ログ機能実装 → 実際にログファイル確認 → 期待するログ出力確認
```

## Logger活用における検証義務

### 必須検証項目
**ログ関連の実装では以下を必ず実行:**

1. **実際のログファイル確認**
```bash
# ファイル出力の場合
cat ./logs/interactive.log
ls -la ./logs/

# コンソール出力の場合  
./gamebook command 2>&1 | grep -E "(DEBUG|INFO|WARN|ERROR)"
```

2. **期待するログレベルでの出力確認**
```bash
# DEBUGレベル確認
export LOG_LEVEL=DEBUG
./gamebook command

# ログが期待通り出力されることを確認
tail -f ./logs/interactive.log
```

3. **ログ設定変更の動作確認**
```bash
# 設定変更前後でのログ出力差分確認
export GAMEBOOK_AI_DEV=true
./gamebook  # → ファイル出力確認

export GAMEBOOK_AI_DEV=false  
./gamebook  # → コンソール制限確認
```

### 検証失敗時の対応
- ❌ 「動作するはず」で完了しない
- ✅ 根本原因を特定して修正完了まで追跡
- ✅ 修正後は必ず再度検証を実行