# Lint設定の詳細

## 有効化しているLinter（理由付き）

### govet
- **Go公式の静的解析ツール**
- `shadow`: Variable shadowingを検出（最重要）
- 型チェック、未使用変数、構造体タグの検証

### staticcheck
- **高品質な静的解析**
- パフォーマンス、正確性、スタイルの問題を検出
- SA4011など詳細な分析

### stylecheck
- **Go標準のスタイル**
- 命名規則、コメント規則の検証
- ST1000系のルール適用

### errcheck
- **エラーチェック漏れ検出**
- 戻り値のerrorを無視していないかチェック
- 重要なエラーハンドリング漏れを防止

### goimports
- **import文の自動整理**
- 不要なimportの削除
- 標準・外部・内部パッケージの適切な順序

### ineffassign
- **無効な代入検出**
- 使用されない変数への代入を検出

### unused
- **未使用コード検出**
- 未使用の変数、関数、型、定数を検出

### unconvert
- **不要な型変換検出**
- 同じ型への無意味な変換を検出

## 意図的に無効化しているLinter（理由付き）

### gochecknoglobals
- **プロジェクト特性上、一部グローバル変数を使用**
- `main.go`での`repo`, `sessionRepo`, `currentGame`
- CLIツールとして状態管理に必要

### funlen
- **関数長制限は独自基準で管理**
- 50行以内の基準を設定済み
- Markdownパーサーなど特定の処理で長くなりがち

### lll
- **行長制限を緩和**
- 日本語コメントが含まれるため
- 120文字制限では日本語で適切な説明が困難

### gosec
- **セキュリティチェックを一部緩和**
- G301, G302, G307（ファイル権限）が厳しすぎる
- 開発環境での利便性を優先

### gocyclo
- **複雑度チェックを独自基準で管理**
- 20以下の基準を設定済み
- Markdownパーサーは構造上複雑になりがち

## Linter設定ファイル (.golangci.yml)

```yaml
linters-settings:
  govet:
    enable:
      - shadow  # Variable shadowingを検出
  gocyclo:
    min-complexity: 15  # 複雑度の閾値
  misspell:
    locale: US  # 英語のスペルチェック
  unused:
    check-exported: false  # 外部パッケージ向けは除外
  unparam:
    check-exported: false  # 外部パッケージ向けは除外

linters:
  enable:
    - govet
    - staticcheck
    - stylecheck
    - errcheck
    - goimports
    - ineffassign
    - unused
    - unconvert
```

## 除外ルール

### テストファイル
- `_test.go`では一部のlintを無効化
- `gomnd`: テストでのマジックナンバーを許可
- `dupl`: テストの重複コードを許可
- `goconst`: テストでの定数化を強制しない

### 統合テスト
- `test/`ディレクトリでは緩和
- `ineffassign`, `errcheck`, `gomnd`を無効化
- 統合テストの特性を考慮

### main関数
- `cmd/`ディレクトリでは一部緩和
- `gochecknoinits`: 初期化処理を許可

## Lintエラーへの対処法

### shadow エラー

```bash
# エラー例
internal/domain/gamebook.go:50:4: declaration of "err" shadows declaration at line 45

# 対処法
# 異なる変数名を使用
if loadErr := load(); loadErr != nil {
    return fmt.Errorf("読み込みエラー: %w", loadErr)
}
```

### errcheck エラー

```bash
# エラー例
internal/repository/markdown_repository.go:100:12: Error return value is not checked

# 対処法
# エラーを適切に処理
if closeErr := file.Close(); closeErr != nil {
    log.Printf("ファイルクローズエラー: %v", closeErr)
}
```

### ineffassign エラー

```bash
# エラー例
internal/domain/paragraph.go:25:2: ineffectual assignment to err

# 対処法
# 未使用の代入を削除するか、適切に使用
if err != nil {
    return err
}
```

## Lint実行コマンド

```bash
# 通常実行
golangci-lint run

# 詳細表示
golangci-lint run --verbose

# 特定のLinterのみ実行
golangci-lint run --enable-only=govet

# 新規導入Linterのテスト
golangci-lint run --enable=gochecknoglobals --new-from-rev=HEAD~1

# 自動修正（対応しているLinterのみ）
golangci-lint run --fix
```