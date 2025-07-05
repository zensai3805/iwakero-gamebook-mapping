# Claude Code AI開発ガイド

Claude Code専用の最適化された開発ガイドです。AI主導開発における効率的なプロジェクト管理とコード実装を支援します。

## 🎯 **AI開発の基本方針**

### **1. 必須事前確認**
```bash
# 作業開始前に必ず実行
1. SPECIFICATION.md を読了
2. DEVELOPMENT.md を読了
3. CLAUDE.md を読了
4. Issue内容を完全理解
```

### **2. TDD厳守**
```bash
# 実装フロー
RED    → 失敗するテスト実装
GREEN  → 最小限の実装で通す
REFACTOR → コードをリファクタリング
```

### **3. 品質基準**
```bash
# コミット前必須チェック
go test ./...         # 全テスト通過
golangci-lint run     # リンターエラーなし
gofmt -s -w .        # フォーマット済み
goimports -w .       # import整理済み
```

## 🛠️ **AI専用開発ツール**

### **Claude Project Manager**
```bash
# 利用可能なコマンド
./scripts/claude-project-manager.sh feature "機能名" "説明" "v1.0.0"
./scripts/claude-project-manager.sh branch 42
./scripts/claude-project-manager.sh progress 42 "🤖 Claude: 実装完了"
./scripts/claude-project-manager.sh pr 42 "PR Title" "PR Description"
./scripts/claude-project-manager.sh quality
```

### **GitHub CLI活用**
```bash
# Issue管理
gh issue create --template ai_feature
gh issue create --template ai_sub_issue
gh issue create --template ai_bug_report

# PR管理
gh pr create --fill-first
gh pr merge --squash
gh pr view --web
```

## 📋 **AI実装チェックリスト**

### **Issue作成時**
- [ ] SPECIFICATION.md確認済み
- [ ] DEVELOPMENT.md確認済み
- [ ] 適切なテンプレート選択
- [ ] 完了条件明確定義
- [ ] AI実装指示記載

### **実装開始時**
- [ ] Feature Branch作成
- [ ] 進捗更新（開始）
- [ ] テストファイル作成
- [ ] 実装方針決定

### **実装中**
- [ ] TDDサイクル遵守
- [ ] 変数shadowing回避
- [ ] エラーハンドリング適切実装
- [ ] 日本語コメント追加
- [ ] 絵文字使用禁止

### **実装完了時**
- [ ] 品質チェック実行
- [ ] 統合テスト実行
- [ ] CLAUDE.md更新
- [ ] PR作成
- [ ] 進捗更新（完了）

## 🔧 **コーディング規約（AI特化）**

### **Variable Shadowing回避**
```go
// ❌ 問題のあるコード
func ExampleFunction() error {
    if err := step1(); err != nil {
        return err
    }
    if err := step2(); err != nil {  // shadowエラー
        return err
    }
    return nil
}

// ✅ 推奨コード
func ExampleFunction() error {
    if step1Err := step1(); step1Err != nil {
        return fmt.Errorf("step1処理エラー: %w", step1Err)
    }
    if step2Err := step2(); step2Err != nil {
        return fmt.Errorf("step2処理エラー: %w", step2Err)
    }
    return nil
}
```

### **Error Handling**
```go
// ❌ 問題のあるコード
return fmt.Errorf("error: %v", err)

// ✅ 推奨コード
return fmt.Errorf("処理に失敗しました: %w", err)
```

### **Function Design**
```go
// ✅ 推奨設計
func ProcessGamebook(gb *domain.Gamebook) error {
    // 50行以内、複雑度20以下
    if validateErr := validateGamebook(gb); validateErr != nil {
        return fmt.Errorf("バリデーションエラー: %w", validateErr)
    }
    
    if processErr := processGamebookData(gb); processErr != nil {
        return fmt.Errorf("データ処理エラー: %w", processErr)
    }
    
    return nil
}
```

### **Comments**
```go
// ✅ 推奨コメント
// ProcessGamebook はゲームブックの処理を実行します
// 入力: ゲームブックオブジェクト
// 出力: 処理結果のエラー
// 副作用: ファイルシステムへの書き込み
func ProcessGamebook(gb *domain.Gamebook) error {
    // バリデーション処理
    if validateErr := validateGamebook(gb); validateErr != nil {
        return fmt.Errorf("バリデーションエラー: %w", validateErr)
    }
    
    // データ処理
    if processErr := processGamebookData(gb); processErr != nil {
        return fmt.Errorf("データ処理エラー: %w", processErr)
    }
    
    return nil
}
```

## 🧪 **テスト実装ガイド**

### **TDD実装例**
```go
// 1. RED: 失敗するテストを実装
func TestGamebook_AddParagraph_Success(t *testing.T) {
    gb := &domain.Gamebook{
        Title:      "テストゲーム",
        Paragraphs: make(map[int]*domain.Paragraph),
    }
    
    paragraph := &domain.Paragraph{
        Number:      1,
        Description: "テストパラグラフ",
    }
    
    err := gb.AddParagraph(paragraph)
    assert.NoError(t, err)
    assert.Equal(t, paragraph, gb.Paragraphs[1])
}

// 2. GREEN: 最小限の実装
func (g *Gamebook) AddParagraph(p *Paragraph) error {
    g.Paragraphs[p.Number] = p
    return nil
}

// 3. REFACTOR: リファクタリング
func (g *Gamebook) AddParagraph(p *Paragraph) error {
    if p == nil {
        return fmt.Errorf("パラグラフがnilです")
    }
    if p.Number <= 0 {
        return fmt.Errorf("パラグラフ番号は正の整数である必要があります: %d", p.Number)
    }
    
    g.Paragraphs[p.Number] = p
    return nil
}
```

### **テストケース設計**
```go
func TestGamebook_AddParagraph(t *testing.T) {
    tests := []struct {
        name        string
        paragraph   *domain.Paragraph
        expectError bool
        errorMsg    string
    }{
        {
            name: "正常ケース",
            paragraph: &domain.Paragraph{
                Number:      1,
                Description: "テストパラグラフ",
            },
            expectError: false,
        },
        {
            name:        "nilパラグラフ",
            paragraph:   nil,
            expectError: true,
            errorMsg:    "パラグラフがnilです",
        },
        {
            name: "無効なパラグラフ番号",
            paragraph: &domain.Paragraph{
                Number:      -1,
                Description: "テストパラグラフ",
            },
            expectError: true,
            errorMsg:    "パラグラフ番号は正の整数である必要があります",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gb := &domain.Gamebook{
                Title:      "テストゲーム",
                Paragraphs: make(map[int]*domain.Paragraph),
            }
            
            err := gb.AddParagraph(tt.paragraph)
            
            if tt.expectError {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errorMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## 🔄 **AI実装フロー**

### **Feature実装フロー**
```bash
# 1. Issue確認・理解
./scripts/claude-project-manager.sh feature "機能名" "説明" "v1.0.0"

# 2. Branch作成
./scripts/claude-project-manager.sh branch 42

# 3. TDD実装
# RED → GREEN → REFACTOR を繰り返し

# 4. 品質チェック
./scripts/claude-project-manager.sh quality

# 5. PR作成
./scripts/claude-project-manager.sh pr 42 "PR Title" "PR Description"
```

### **Sub-Issue実装フロー**
```bash
# 1. 親Issue確認
gh issue view 123

# 2. Sub-Issue作成
gh issue create --template ai_sub_issue

# 3. 実装（TDD）
# 4. 品質チェック
# 5. 親Issue進捗更新
```

### **Bug修正フロー**
```bash
# 1. Bug Report作成
gh issue create --template ai_bug_report

# 2. 再現テスト作成
# 3. 修正実装
# 4. 回帰テスト実行
# 5. PR作成
```

## 📊 **品質管理**

### **品質チェックコマンド**
```bash
# 基本品質チェック
go test ./...
golangci-lint run
gofmt -s -w .
goimports -w .

# 詳細品質チェック
go test ./... -cover
go test ./... -race
go vet ./...
```

### **品質基準**
- **テストカバレッジ**: 70%以上
- **Lint**: エラーゼロ
- **関数長**: 50行以内
- **複雑度**: 20以下
- **Variable Shadowing**: なし

## 🚀 **自動化ツール**

### **GitHub Actions活用**
```yaml
# 品質チェック自動化
name: Quality Check
on: [push, pull_request]
jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: 1.23
      - run: go test ./...
      - run: golangci-lint run
```

### **Pre-commit Hook**
```bash
#!/bin/bash
# .git/hooks/pre-commit

# 品質チェック実行
go test ./...
golangci-lint run
gofmt -s -w .
goimports -w .
```

## 📚 **リファレンス**

### **重要ファイル**
- `SPECIFICATION.md`: 仕様書
- `DEVELOPMENT.md`: 開発方針
- `CLAUDE.md`: プロジェクト指示書
- `CLAUDE_AI_DEVELOPMENT_GUIDE.md`: このファイル（AI専用）

### **Templates**
- `.github/ISSUE_TEMPLATE/ai_feature.yml`: Feature Issue
- `.github/ISSUE_TEMPLATE/ai_sub_issue.yml`: Sub-Issue
- `.github/ISSUE_TEMPLATE/ai_bug_report.yml`: Bug Report
- `.github/PULL_REQUEST_TEMPLATE.md`: PR Template

### **Scripts**
- `scripts/claude-project-manager.sh`: AI専用プロジェクト管理

## 💡 **最適化Tips**

### **効率的なIssue管理**
1. Issue作成時に完了条件を明確に定義
2. Sub-Issueは技術境界で分割
3. 依存関係を明確に記載
4. 進捗は定期的に更新

### **効率的なコード実装**
1. TDDサイクルを厳守
2. 変数shadowing回避を最優先
3. エラーハンドリングを適切に実装
4. 日本語コメントで意図を明確化

### **効率的な品質管理**
1. コミット前に品質チェック実行
2. テストカバレッジを常に確認
3. リンターエラーを即座に修正
4. 定期的にコードレビューを実施

---

**このガイドは、Claude Code による AI 主導開発を最適化するために設計されています。**
**質問や改善提案がある場合は、Issue を作成してください。**