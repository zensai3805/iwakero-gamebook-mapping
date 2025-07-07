# テスト駆動開発（TDD）方針

## 最重要原則

### 🚨 テスト駆動開発（TDD）の徹底

## TDD実践：RED/GREEN/REFACTORサイクル

### 0. テストリスト作成（事前準備）
**期待する振る舞いをリスト化**

**テストリスト詳細管理** → `docs/dev/test-list-management.md` 参照

**基本原則**: 
- 振る舞いのみ記載、実装詳細は記載しない
- 「変更後の振る舞いが満たすべき様々な動作を網羅的に考える」
- 小さな問題から始めて段階的に進める
- 「もしも〜だったら」の観点で網羅的に分析

### 🔴 RED: 失敗するテストを書く
**1つのテストのみ作成し、失敗を確認**

#### REDフェーズの重要ポイント
- **失敗の確認**: テストが確実に失敗することを確認
- **最小限のテスト**: そのテスト1つのみに集中
- **明確な失敗理由**: なぜ失敗するかを理解

```go
func TestAddParagraph_WhenValidInput_ReturnsNoError(t *testing.T) {
    gamebook := domain.NewGamebook("test")
    paragraph := &domain.Paragraph{Number: 1, Description: "テスト"}
    
    err := gamebook.AddParagraph(paragraph)
    assert.NoError(t, err)  // この時点で失敗する（実装未完成）
}
```

### 🟢 GREEN: テストを通す
**そのテスト1つのみを通す最小限の実装**

#### GREENフェーズの重要ポイント
- **最小限の実装**: そのテストのみを通すコード
- **アサーション保護**: テストの削除・変更は禁止
- **他テスト無視**: 未実装テストは考慮しない
- **重複コード許容**: この段階では重複を恐れない

```go
// ✅ 正しい: そのテストのみ考慮
func (g *Gamebook) AddParagraph(p *Paragraph) error {
    g.paragraphs = append(g.paragraphs, p)
    return nil  // 最小限でテストを通す
}

// ❌ 禁止: 未テスト機能まで実装
func (g *Gamebook) AddParagraph(p *Paragraph) error {
    if p.Number < 0 { return errors.New("invalid") }  // 未テスト
    g.paragraphs = append(g.paragraphs, p)
    return nil
}
```

### 🔵 REFACTOR: 重複を除去する
**テストを保護しながら設計を改善**

#### REFACTORフェーズの重要ポイント
- **テスト保護**: 既存テストは必ず通り続ける
- **重複除去**: コードの重複を取り除く
- **設計改善**: 可読性・保守性の向上
- **小さな変更**: 一度に大きく変更しない

### 🔄 サイクル継続
**テストリストが空になるまでRED→GREEN→REFACTORを繰り返す**

## RED/GREEN/REFACTORサイクルの厳格ルール

### 小さなステップの定義
**1回のサイクル = 1つのテストのみに焦点**

- 🔴 RED: 1つのテストを書いて失敗させる
- 🟢 GREEN: そのテスト1つのみを通す最小実装
- 🔵 REFACTOR: テストを保護しながら重複除去
- 🔄 次のテストに進む前に必ずサイクル完了

### サイクル実践例

#### 🔴 REDフェーズ実例
```markdown
テストリストから選択:
- [x] 新しいパラグラフを追加できる  ← 完了
- [ ] 重複するパラグラフ番号でエラーが発生する  ← 🔴 RED
- [ ] 空のDescriptionでエラーが発生する
```

```go
func TestAddParagraph_WhenDuplicate_ReturnsError(t *testing.T) {
    gamebook := domain.NewGamebook("test")
    paragraph := &domain.Paragraph{Number: 1, Description: "テスト"}
    
    gamebook.AddParagraph(paragraph)
    duplicateErr := gamebook.AddParagraph(paragraph)  // 🔴 失敗する
    
    assert.Error(t, duplicateErr)
}
```

#### 🟢 GREENフェーズ実例
```go
func (g *Gamebook) AddParagraph(p *Paragraph) error {
    // 🟢 このテストのみを通す最小実装
    for _, existing := range g.paragraphs {
        if existing.Number == p.Number {
            return fmt.Errorf("重複")  // 最小限のエラー
        }
    }
    g.paragraphs = append(g.paragraphs, p)
    return nil
}
```

#### 🔵 REFACTORフェーズ実例
```go
// リファクタリング: エラーメッセージを改善
func (g *Gamebook) AddParagraph(p *Paragraph) error {
    for _, existing := range g.paragraphs {
        if existing.Number == p.Number {
            return fmt.Errorf("paragraph %d already exists", p.Number)
        }
    }
    g.paragraphs = append(g.paragraphs, p)
    return nil
}
// ✅ 全テストが通ることを確認してからサイクル完了
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

## RED/GREEN/REFACTORサイクル違反パターンと対策

### 🔴 RED違反パターン
```
❌ テストリスト省略していきなりテスト作成
✅ テストリスト作成 → RED → GREEN → REFACTOR

❌ 複数テストを一度に作成
✅ 1つのテストのみ作成して失敗確認

❌ 成功するテストを書く
✅ 確実に失敗するテストを書く
```

### 🟢 GREEN違反パターン
```
❌ 未テスト機能を「ついでに」実装
✅ そのテスト1つのみを通す最小実装

❌ アサーション削除・変更でテスト通過
✅ 実装修正でテスト通過

❌ 複数テストを同時に通そうとする
✅ 現在のテスト1つのみに集中
```

### 🔵 REFACTOR違反パターン
```
❌ テストを修正してリファクタリング
✅ テストを保護してコードをリファクタリング

❌ 機能追加をリファクタリングと混同
✅ 重複除去・設計改善のみに集中

❌ 大規模な設計変更
✅ 小さな改善の積み重ね
```

### 🔄 サイクル違反パターン
```
❌ RED → GREEN → 次のRED（REFACTOR省略）
✅ RED → GREEN → REFACTOR → 次のRED

❌ 実装 → テスト作成（サイクル逆転）
✅ テスト作成 → 実装（正しいサイクル）
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