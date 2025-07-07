# テスト駆動開発（TDD）方針

## 最重要原則

### 🚨 テスト駆動開発（TDD）の徹底
**Kent Beck正典準拠のTDDサイクル必須**

## TDD実践の5つのステップ

### 1. テストリスト作成
**期待する振る舞いをリスト化**

```markdown
# AddParagraph機能のテストリスト
- [ ] 新しいパラグラフを正常に追加できる
- [ ] 重複するパラグラフ番号でエラーが発生する
- [ ] 空のDescriptionでエラーが発生する
- [ ] 負の番号でエラーが発生する
```

**ルール**: 振る舞いのみ記載、実装詳細は記載しない

### 2. 1つのテストを書く
**1つのテストのみ作成**

```go
func TestAddParagraph_WhenValidInput_ReturnsNoError(t *testing.T) {
    gamebook := domain.NewGamebook("test")
    paragraph := &domain.Paragraph{Number: 1, Description: "テスト"}
    
    err := gamebook.AddParagraph(paragraph)
    assert.NoError(t, err)
}
```

### 3. テストを通す
**そのテスト1つのみを通す最小限の実装**

- アサーション削除・変更は禁止
- そのテストのみ考慮、他のテストは無視
- 過度な機能追加禁止

### 4. リファクタリング
**必要に応じて改善**

- 既存テストが通る範囲でのみ実施
- 過度なリファクタリング禁止

### 5. 繰り返し
**テストリストが空になるまで2-4を繰り返す**

## 小さなステップの厳格ルール

### 定義
**1つのテストのみに焦点を絞って実装**

- 1回のサイクル = 1つのテストのみ
- 1つのテストが通るまで次に進まない
- そのテスト1つを通すコードのみ実装
- 「ついでに」実装は禁止

### 実践手順

#### 1. テストリストから1つ選択
```markdown
- [x] 新しいパラグラフを正常に追加できる  ← 完了
- [ ] 重複するパラグラフ番号でエラーが発生する  ← 次
- [ ] 空のDescriptionでエラーが発生する
```

#### 2. 1つのテストのみ実装
```go
// ✅ 正しい: 1つの振る舞いのみ検証
func TestAddParagraph_WhenDuplicate_ReturnsError(t *testing.T) {
    gamebook := domain.NewGamebook("test")
    paragraph := &domain.Paragraph{Number: 1, Description: "テスト"}
    
    gamebook.AddParagraph(paragraph)
    duplicateErr := gamebook.AddParagraph(paragraph)
    
    assert.Error(t, duplicateErr)
}

// ❌ 禁止: 複数の振る舞いを検証
func TestAddParagraph_MultipleScenarios(t *testing.T) {
    // 複数のテストケースを1つのテストで実装するのは禁止
}
```

#### 3. そのテスト1つを通す最小限実装
```go
// ✅ 正しい: そのテストのみ考慮
func (g *Gamebook) AddParagraph(p *Paragraph) error {
    for _, existing := range g.paragraphs {
        if existing.Number == p.Number {
            return fmt.Errorf("重複")
        }
    }
    g.paragraphs = append(g.paragraphs, p)
    return nil
}

// ❌ 禁止: テストしていない機能まで実装
func (g *Gamebook) AddParagraph(p *Paragraph) error {
    // 重複チェック（テスト済み）
    // 空チェック（未テスト） ← 禁止
    // 負の番号チェック（未テスト） ← 禁止
}
```

### 過渡的コード管理
```go
// ✅ 正しい: TODOで明示
func (g *Gamebook) AddParagraph(p *Paragraph) error {
    // TODO: 重複チェック実装予定
    g.paragraphs = append(g.paragraphs, p)
    return nil
}

// ❌ 禁止: 明示なし
func (g *Gamebook) AddParagraph(p *Paragraph) error {
    g.paragraphs = append(g.paragraphs, p)
    return nil
}
```

## テスト品質基準

### Logger活用
```go
func TestAddParagraph_WhenDuplicate_ReturnsError(t *testing.T) {
    os.Setenv("GAMEBOOK_AI_DEV", "true")
    os.Setenv("LOG_LEVEL", "DEBUG")
    defer os.Unsetenv("GAMEBOOK_AI_DEV")
    defer os.Unsetenv("LOG_LEVEL")
    
    gamebook := domain.NewGamebook("test")
    paragraph := &domain.Paragraph{Number: 1, Description: "テスト"}
    
    gamebook.AddParagraph(paragraph)
    duplicateErr := gamebook.AddParagraph(paragraph)
    
    assert.Error(t, duplicateErr)
}
```

### 命名規則
```go
// Test対象関数名_条件_期待結果
func TestAddParagraph_WhenDuplicate_ReturnsError(t *testing.T)
```

### テスト独立性
- 各テスト独立実行可能
- 順序依存なし
- 共有状態なし

## TDD違反パターンと対策

### 違反1: テストリスト省略
```
❌ いきなりテスト作成
✅ テストリスト作成 → テスト → 実装
```

### 違反2: 複数テスト同時実装
```
❌ 複数テスト一度に作成
✅ 1つのテスト → 実装 → 次のテスト

❌ 1つのテストで複数振る舞い検証
✅ 1つのテスト = 1つの振る舞い

❌ 未テスト機能を「ついでに」実装
✅ そのテスト1つのみ通す実装
```

### 違反3: テスト後付け
```
❌ 実装 → テスト作成
✅ テスト作成 → 実装
```

### 違反4: テスト修正でつじつま合わせ
```
❌ assert削除・変更でテスト通過
✅ 実装修正でテスト通過
```

## Logger検証義務

### 必須確認
```bash
# ログファイル確認
cat ./logs/interactive.log

# ログレベル確認
export LOG_LEVEL=DEBUG
./gamebook command

# 設定変更確認
export GAMEBOOK_AI_DEV=true
./gamebook
```

### 検証失敗時対応
- 根本原因特定まで追跡
- 修正後は再度検証実行