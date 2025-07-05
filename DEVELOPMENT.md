# 開発方針とプロジェクト管理

## 最重要原則

### 🚨 テスト駆動開発（TDD）の徹底
**t_wadaメソッドによるテストファーストを必須とする**

1. **RED**: 失敗するテストを最初に作成
2. **GREEN**: テストが通る最低限のコードを実装
3. **REFACTOR**: コードを改善・整理

- **全機能はテストファーストで実装**
- テストコードは実装コードと同じ重要度で扱う
- TDDを守らない実装は認めない

## Issue駆動開発

### 基本方針
- **すべての作業はGitHub Issueとして管理**
- **Issueの冒頭にSPECIFICATION.md・DEVELOPMENT.mdの確認を必須**
- **IssueのCloseはPRで行う**
- **PRでは必ずCLAUDE.mdを最新化する**

### 開発フロー
1. **Issue確認**: SPECIFICATION.md・DEVELOPMENT.mdを確認してから作業開始
2. **Issue作成**: **必ずAI最適化テンプレート使用**（`gh issue create --template ai_feature.md`、`--template ai_sub_issue.md`、`--template ai_bug_report.md`）
3. **featureブランチ作成**: AI管理スクリプト使用（`./scripts/claude-project-manager.sh branch {issue-number}`）
4. **TDD実行**: RED → GREEN → REFACTOR
5. **品質チェック**: AI管理スクリプト使用（`./scripts/claude-project-manager.sh quality`）
6. **PR作成**: **必ずAI最適化PRテンプレート使用**（プルリクエスト作成時に自動適用）
7. **ユーザーマージ**: PRのマージは必ずユーザーが実行

## GitHub Template使用方針

### Issue・PRテンプレート利用義務
- **全Issue・PRでAI最適化テンプレート使用必須**
- テンプレートはClaude Code最適化済み（AIが効率的に作業可能）
- カスタマイズ不要（必要十分な分量に調整済み）

### テンプレート種別
1. **Feature実装**: `gh issue create --template ai_feature.md`
   - 新機能開発用
   - TDD方針、技術仕様、AI実装指示を含む

2. **Sub-Issue**: `gh issue create --template ai_sub_issue.md`
   - 複雑なFeatureの分割用
   - 親Issue連携、実装範囲明確化

3. **Bug Report**: `gh issue create --template ai_bug_report.md`
   - バグ修正用
   - 再現手順、AI修正指示を含む

4. **Pull Request**: 自動適用されるPRテンプレート
   - 簡潔な品質チェックリスト
   - AI実装に必要十分な内容

### テンプレート遵守の重要性
- **AI（Claude Code）がプロジェクト内容を効率的に理解**
- **実装品質の統一性確保**
- **レビュー効率の向上**
- **プロジェクト管理の自動化促進**

### テンプレート利用コマンド
```bash
# Feature実装Issue作成
gh issue create --template ai_feature.md

# Sub-Issue作成  
gh issue create --template ai_sub_issue.md

# Bug Report作成
gh issue create --template ai_bug_report.md

# PR作成（テンプレート自動適用）
gh pr create
```

### 必須チェック
- **ローカルLint**: `golangci-lint run` でエラーなし
- **ローカルテスト**: `go test ./...` で全テスト通過
- **commit・push前に必ず実行**

## 開発者チェックリスト

### コミット前必須チェック
- [ ] `golangci-lint run` でエラーなし
- [ ] `go test ./...` で全テスト通過
- [ ] `gofmt -s -w .` でフォーマット済み
- [ ] エラーメッセージは日本語で統一
- [ ] Variable shadowingなし（特に`err`変数）
- [ ] エラーラッピングに`%w`を使用
- [ ] マジックナンバーは定数化

### PR作成前必須チェック
- [ ] 関数の長さ50行以内
- [ ] 複雑度20以下
- [ ] 適切なエラーハンドリング
- [ ] 日本語コメントの追加
- [ ] 共有状態の適切な同期
- [ ] テストカバレッジの確認
- [ ] CLAUDE.mdの更新

### コードレビュー観点
- [ ] **Variable Shadowing**: 特に`err`変数の再宣言
- [ ] **Error Wrapping**: `fmt.Errorf("message: %w", err)`の使用
- [ ] **Resource Management**: `defer`を使った適切なリソース解放
- [ ] **Function Length**: 50行を超える関数の分割
- [ ] **Magic Numbers**: 定数化されているか
- [ ] **Concurrency**: 共有状態の適切な保護
- [ ] **Test Coverage**: 新機能のテストが追加されているか

### よくある問題と解決策
#### Variable Shadowing
```go
// 問題: shadowエラー
func example() error {
    if err := step1(); err != nil {
        return err
    }
    if err := step2(); err != nil {  // shadowエラー
        return err
    }
    return nil
}

// 解決: 明示的な変数名
func example() error {
    if step1Err := step1(); step1Err != nil {
        return fmt.Errorf("step1エラー: %w", step1Err)
    }
    if step2Err := step2(); step2Err != nil {
        return fmt.Errorf("step2エラー: %w", step2Err)
    }
    return nil
}
```

#### Error Wrapping
```go
// 問題: エラーラッピング不適切
return fmt.Errorf("error: %v", err)

// 解決: 適切なエラーラッピング
return fmt.Errorf("処理に失敗: %w", err)
```

#### Resource Management
```go
// 問題: エラーハンドリングなし
defer file.Close()

// 解決: 適切なエラーハンドリング
defer func() {
    if closeErr := file.Close(); closeErr != nil {
        log.Printf("ファイルクローズエラー: %v", closeErr)
    }
}()
```

## GitHub公式Sub-Issue管理

### 重要: gh CLIの現状
- **gh CLIはsub-issuesをネイティブサポートしていません**（2025年7月時点）
- GitHub CLI Issue #10298で機能追加が要求されている状況
- **現在はGraphQL API経由での操作が必要**

### Beta機能の利用条件
- **組織レベルでのサインアップが必要**: https://github.com/features/issues/signup
- 個人リポジトリでも利用可能（要申請）
- Evolving GitHub Issues (Public Beta) への参加が前提

### 使用方法
- **マイルストーン開始時にSub-Issue洗い出しを必須実施**
- GraphQL API経由でのみ操作可能
- 各Sub-IssueでもSPECIFICATION.md・DEVELOPMENT.md確認必須

### GraphQL API操作方法
```bash
# Step 1: 親・子IssueのGraphQL IDを取得
gh api graphql -f query='
{
  repository(owner: "owner", name: "repo") {
    issue(number: 親Issue番号) {
      title
      id
    }
  }
}' --header "GraphQL-Features: sub_issues"

gh api graphql -f query='
{
  repository(owner: "owner", name: "repo") {
    issue(number: 子Issue番号) {
      title
      id
    }
  }
}' --header "GraphQL-Features: sub_issues"

# Step 2: Sub-Issue追加
gh api graphql -f query='
mutation {
  addSubIssue(input: {
    issueId: "親IssueのGraphQL_ID"
    subIssueId: "子IssueのGraphQL_ID"
  }) {
    issue {
      title
    }
  }
}' --header "GraphQL-Features: sub_issues"

# Step 3: Sub-Issues一覧確認（REST API）
gh api -X GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues
```

### 実際の操作例（動作確認済み）
```bash
# Issue #50に Issue #51をsub-issueとして追加する例

# 1. GraphQL IDを取得
gh api graphql -f query='
{
  repository(owner: "zensai3805", name: "iwakero-gamebook-mapping") {
    issue(number: 50) {
      id
    }
  }
}' --header "GraphQL-Features: sub_issues"
# 結果: "I_kwDOPA7qR86_EALR"

gh api graphql -f query='
{
  repository(owner: "zensai3805", name: "iwakero-gamebook-mapping") {
    issue(number: 51) {
      id
    }
  }
}' --header "GraphQL-Features: sub_issues"
# 結果: "I_kwDOPA7qR86_EAMX"

# 2. Sub-Issue追加
gh api graphql -f query='
mutation {
  addSubIssue(input: {
    issueId: "I_kwDOPA7qR86_EALR"
    subIssueId: "I_kwDOPA7qR86_EAMX"
  }) {
    issue {
      title
    }
  }
}' --header "GraphQL-Features: sub_issues"
# 結果: {"data":{"addSubIssue":{"issue":{"title":"テスト: 親Issue"}}}}

# 3. 確認
gh api -X GET /repos/zensai3805/iwakero-gamebook-mapping/issues/50/sub_issues
# 結果: Issue #51が子Issueとして表示される
```

### 洗い出しフロー
1. **SPECIFICATION.md・DEVELOPMENT.md確認**: 要件と開発方針の詳細確認
2. **現状分析**: 実装済み/未実装の棚卸し
3. **Sub-Issue作成**: 具体的タスクごとに作成
4. **GraphQL API呼び出し**: 親子関係を設定
5. **依存関係整理**: 実装順序決定
6. **実装開始**: 最優先Sub-Issueから

## 重要資料の位置づけ

- **SPECIFICATION.md**: ツール仕様の最重要資料
- **DEVELOPMENT.md**: 開発方針の最重要資料
- Issue作業前に必ず両方確認
- 更新が必要な場合はユーザーに確認を求める

## コーディング規約

### 基本方針
- **絵文字使用禁止**: コード・コメント・コミットメッセージで絵文字を使用しない
- **日本語コメント推奨**: 仕様書が日本語のため、コメントも日本語で記載
- 可読性とプロフェッショナル性を重視
- 例外: ユーザーから明示的に要求された場合のみ

### Error Handling（エラーハンドリング）
- **エラーラッピング**: 必ず `%w` を使用してエラーをラップする
- **エラーメッセージ**: 日本語で統一する
- **エラーチェック**: すべてのエラーを適切に処理する

```go
// ✅ 推奨
func (r *MarkdownRepository) Load(title string) (*domain.Gamebook, error) {
    file, err := os.Open(filepath.Join(r.dataDir, title+".md"))
    if err != nil {
        return nil, fmt.Errorf("ファイルオープンに失敗: %w", err)
    }
    defer func() {
        if closeErr := file.Close(); closeErr != nil {
            log.Printf("ファイルクローズエラー: %v", closeErr)
        }
    }()
    // ...
}

// ❌ 非推奨
func (r *MarkdownRepository) Load(title string) (*domain.Gamebook, error) {
    file, err := os.Open(filepath.Join(r.dataDir, title+".md"))
    if err != nil {
        return nil, fmt.Errorf("error: %v", err)  // %vは非推奨
    }
    defer file.Close()  // エラーハンドリングなし
    // ...
}
```

### Variable Naming（変数命名）
- **Variable Shadowing回避**: 同名変数の再宣言を避ける
- **明確な変数名**: 目的が分かる変数名を使用する
- **一時変数**: 明示的な名前を付ける

```go
// ✅ 推奨
func ExecuteShowCommand() error {
    if initErr := ui.Initialize(data); initErr != nil {
        return fmt.Errorf("UI初期化エラー: %w", initErr)
    }
    
    if layoutErr := ui.SetupLayout(); layoutErr != nil {
        return fmt.Errorf("レイアウト設定エラー: %w", layoutErr)
    }
    
    if renderErr := ui.Render(); renderErr != nil {
        return fmt.Errorf("描画エラー: %w", renderErr)
    }
    
    return nil
}

// ❌ 非推奨（shadowエラー）
func ExecuteShowCommand() error {
    if err := ui.Initialize(data); err != nil {
        return fmt.Errorf("UI初期化エラー: %w", err)
    }
    
    if err := ui.SetupLayout(); err != nil {  // shadowエラー
        return fmt.Errorf("レイアウト設定エラー: %w", err)
    }
    
    return nil
}
```

### Function Design（関数設計）
- **関数の長さ**: 50行以内を目標とする
- **複雑度**: 20以下を維持する
- **単一責任の原則**: 1つの関数は1つの責任のみ
- **適切な分割**: 大きな関数は小さなヘルパー関数に分割

```go
// ✅ 推奨
func (r *MarkdownRepository) Load(title string) (*domain.Gamebook, error) {
    file, err := r.openGamebookFile(title)
    if err != nil {
        return nil, fmt.Errorf("ファイルオープンに失敗: %w", err)
    }
    defer r.closeFile(file)
    
    gamebook, err := r.parseGamebook(file, title)
    if err != nil {
        return nil, fmt.Errorf("パース処理に失敗: %w", err)
    }
    
    return gamebook, nil
}

func (r *MarkdownRepository) openGamebookFile(title string) (*os.File, error) {
    // ファイルオープン処理
}

func (r *MarkdownRepository) parseGamebook(file *os.File, title string) (*domain.Gamebook, error) {
    // パース処理
}
```

### Constants and Magic Numbers（定数とマジックナンバー）
- **マジックナンバー禁止**: 意味のある数値は定数として定義
- **定数名**: 用途が分かる名前を付ける
- **グループ化**: 関連する定数はまとめて定義

```go
// ✅ 推奨
const (
    // パラグラフ番号の上限
    MaxParagraphNumber = 1000
    
    // デフォルトのグリッドサイズ
    DefaultGridSize = 3
    
    // ファイル権限
    FilePermission = 0644
    DirPermission  = 0755
)

// ❌ 非推奨
for i := 0; i < 1000; i++ {  // 1000はマジックナンバー
    // ...
}
```

### Concurrency Safety（並行性安全）
- **共有状態**: 適切な同期機構を使用
- **Mutex使用**: 読み書きが混在する場合は`sync.RWMutex`を使用
- **Goroutine管理**: 適切なライフサイクル管理

```go
// ✅ 推奨
type SafeGamebook struct {
    mu   sync.RWMutex
    data *domain.Gamebook
}

func (sg *SafeGamebook) GetParagraph(num int) (*domain.Paragraph, error) {
    sg.mu.RLock()
    defer sg.mu.RUnlock()
    
    return sg.data.GetParagraph(num)
}

func (sg *SafeGamebook) AddParagraph(p *domain.Paragraph) error {
    sg.mu.Lock()
    defer sg.mu.Unlock()
    
    return sg.data.AddParagraph(p)
}
```

## Lint設定の詳細

### 有効化しているLinter（理由付き）
- **govet**: Go公式の静的解析ツール
  - `shadow`: Variable shadowingを検出（最重要）
  - 型チェック、未使用変数、構造体タグの検証
- **staticcheck**: 高品質な静的解析
  - パフォーマンス、正確性、スタイルの問題を検出
  - SA4011など詳細な分析
- **stylecheck**: Go標準のスタイル
  - 命名規則、コメント規則の検証
  - ST1000系のルール適用
- **errcheck**: エラーチェック漏れ検出
  - 戻り値のerrorを無視していないかチェック
  - 重要なエラーハンドリング漏れを防止
- **goimports**: import文の自動整理
  - 不要なimportの削除
  - 標準・外部・内部パッケージの適切な順序
- **ineffassign**: 無効な代入検出
  - 使用されない変数への代入を検出
- **unused**: 未使用コード検出
  - 未使用の変数、関数、型、定数を検出
- **unconvert**: 不要な型変換検出
  - 同じ型への無意味な変換を検出

### 意図的に無効化しているLinter（理由付き）
- **gochecknoglobals**: プロジェクト特性上、一部グローバル変数を使用
  - `main.go`での`repo`, `sessionRepo`, `currentGame`
  - CLIツールとして状態管理に必要
- **funlen**: 関数長制限は独自基準で管理
  - 50行以内の基準を設定済み
  - Markdownパーサーなど特定の処理で長くなりがち
- **lll**: 行長制限を緩和
  - 日本語コメントが含まれるため
  - 120文字制限では日本語で適切な説明が困難
- **gosec**: セキュリティチェックを一部緩和
  - G301, G302, G307（ファイル権限）が厳しすぎる
  - 開発環境での利便性を優先
- **gocyclo**: 複雑度チェックを独自基準で管理
  - 20以下の基準を設定済み
  - Markdownパーサーは構造上複雑になりがち

### Linter設定の詳細
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
```

### 除外ルール
- **テストファイル**: `_test.go`では一部のlintを無効化
  - `gomnd`: テストでのマジックナンバーを許可
  - `dupl`: テストの重複コードを許可
  - `goconst`: テストでの定数化を強制しない
- **統合テスト**: `test/`ディレクトリでは緩和
  - `ineffassign`, `errcheck`, `gomnd`を無効化
  - 統合テストの特性を考慮
- **main関数**: `cmd/`ディレクトリでは一部緩和
  - `gochecknoinits`: 初期化処理を許可

## 開発コマンド

### AI最適化コマンド（推奨）
```bash
# Claude専用プロジェクト管理
./scripts/claude-project-manager.sh feature "機能名" "説明" "v1.0.0"  # Feature Issue作成
./scripts/claude-project-manager.sh sub-issue 42 "Sub機能名" "説明"     # Sub-Issue作成
./scripts/claude-project-manager.sh branch 42                        # Feature Branch作成
./scripts/claude-project-manager.sh progress 42 "進捗メッセージ"        # 進捗更新
./scripts/claude-project-manager.sh quality                          # 品質チェック実行
./scripts/claude-project-manager.sh pr 42 "PRタイトル" "PR説明"       # PR作成

# GitHub CLI活用
gh issue create --template ai_feature.md      # AI最適化Feature Issue
gh issue create --template ai_sub_issue.md    # AI最適化Sub-Issue
gh issue create --template ai_bug_report.md   # AI最適化Bug Report
```

### 基本コマンド
```bash
# テスト実行
go test ./...

# Lint実行
golangci-lint run

# フォーマット
gofmt -s -w .

# import整理
goimports -w .

# ビルド
go build -o gamebook ./cmd/gamebook
```

### 品質チェック（推奨）
```bash
# 完全な品質チェック
make check

# または手動で実行
go test ./...
golangci-lint run
gofmt -s -w .
goimports -w .
```

### アプリケーション実行
```bash
# インタラクティブモード（推奨）
./gamebook

# CLIコマンド
./gamebook new "GameTitle"      # 新規作成
./gamebook add 1 "Description"  # パラグラフ追加
./gamebook choice 1 "Go north" 2  # 選択肢追加
./gamebook select 1 1           # 選択実行
./gamebook show                 # 現在状態表示
```

### 開発支援コマンド
```bash
# リアルタイムテスト（開発中）
go test ./... -watch

# Lintの詳細確認
golangci-lint run --verbose

# 特定のLinterのみ実行
golangci-lint run --enable-only=govet

# カバレッジ計測
go test ./... -cover
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## ラベル体系

- **priority**: high/medium/low
- **type**: epic（マイルストーン追跡）
- **area**: cli/domain/repository（技術領域）