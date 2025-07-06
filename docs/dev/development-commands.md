# 開発コマンド集

## AI最適化コマンド（推奨）

### Claude専用プロジェクト管理

```bash
# Feature Issue作成
./scripts/claude-project-manager.sh feature "機能名" "説明" "v1.0.0"

# Sub-Issue作成
./scripts/claude-project-manager.sh sub-issue 42 "Sub機能名" "説明"

# Feature Branch作成
./scripts/claude-project-manager.sh branch 42

# 進捗更新
./scripts/claude-project-manager.sh progress 42 "進捗メッセージ"

# 品質チェック実行
./scripts/claude-project-manager.sh quality

# PR作成
./scripts/claude-project-manager.sh pr 42 "PRタイトル" "PR説明"
```

### GitHub CLI活用

```bash
# AI最適化Feature Issue
gh issue create --template ai_feature.md

# AI最適化Sub-Issue
gh issue create --template ai_sub_issue.md

# AI最適化Bug Report
gh issue create --template ai_bug_report.md

# PR作成（テンプレート自動適用）
gh pr create
```

## 基本開発コマンド

### テスト関連

```bash
# 全テスト実行
go test ./...

# 特定パッケージのテスト
go test ./internal/domain/...

# テスト詳細表示
go test -v ./...

# カバレッジ計測
go test ./... -cover

# カバレッジレポート生成
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# 特定のテストのみ実行
go test -run TestAddParagraph ./internal/domain

# Logger活用テスト（AI開発モード）
export GAMEBOOK_AI_DEV=true
export LOG_LEVEL=DEBUG
export LOG_OUTPUT=stderr
go test -v ./... 2>&1 | grep -E "(DEBUG|INFO|WARN|ERROR|Test)"
```

### Lint・フォーマット

```bash
# Lint実行
golangci-lint run

# Lint詳細確認
golangci-lint run --verbose

# 特定のLinterのみ実行
golangci-lint run --enable-only=govet

# フォーマット
gofmt -s -w .

# import整理
goimports -w .

# 自動修正（可能な範囲）
golangci-lint run --fix
```

### ビルド・実行

```bash
# ビルド
go build -o gamebook ./cmd/gamebook

# ビルド（全OS向け）
GOOS=windows go build -o gamebook.exe ./cmd/gamebook
GOOS=darwin go build -o gamebook-mac ./cmd/gamebook
GOOS=linux go build -o gamebook-linux ./cmd/gamebook

# 実行（ビルド後）
./gamebook

# 直接実行
go run ./cmd/gamebook
```

## Logger設定確認

### 現在のLogger設定確認

```bash
# 環境変数確認
env | grep -E "(LOG_|GAMEBOOK_)"

# Logger設定ファイル確認
cat ~/.gamebook/logger.yaml

# Logger動作確認（CLI）
export LOG_LEVEL=DEBUG
export LOG_OUTPUT=stderr
export GAMEBOOK_AI_DEV=true
./gamebook list 2>&1 | head -5

# Logger動作確認（インタラクティブ）
./gamebook 2>&1 | grep -E "(DEBUG|INFO|WARN|ERROR)" | head -3
```

## 品質チェック（推奨）

### 完全な品質チェック

```bash
# Makefileを使用
make check

# 手動実行（Logger活用）
export GAMEBOOK_AI_DEV=true
export LOG_LEVEL=DEBUG
go test ./...
golangci-lint run
gofmt -s -w .
goimports -w .
```

### コミット前チェックリスト

```bash
# 1. テスト実行
go test ./...

# 2. Lint確認
golangci-lint run

# 3. フォーマット
gofmt -s -w .
goimports -w .

# 4. ビルド確認
go build -o gamebook ./cmd/gamebook
```

## アプリケーション操作

### インタラクティブモード（推奨）

```bash
# 起動
./gamebook

# メニュー操作
# - 矢印キー: 選択
# - Enter: 決定
# - Esc/Ctrl+C: 終了
```

### CLIコマンド

```bash
# 新規ゲームブック作成
./gamebook new "GameTitle"

# パラグラフ追加
./gamebook add 1 "冒険の始まり"

# 選択肢追加
./gamebook choice 1 "北へ進む" 2

# 選択実行
./gamebook select 1 1

# 現在状態表示
./gamebook show

# 直接移動
./gamebook move 5

# ゲーム切り替え
./gamebook switch "AnotherGame"

# ゲーム一覧
./gamebook list
```

## 開発支援ツール

### 依存関係管理

```bash
# 依存関係インストール
go mod download

# 依存関係更新
go get -u ./...

# 不要な依存関係削除
go mod tidy

# 依存関係グラフ確認
go mod graph
```

### デバッグ

```bash
# デバッグビルド
go build -gcflags="all=-N -l" -o gamebook-debug ./cmd/gamebook

# レースコンディション検出
go test -race ./...

# プロファイリング
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

### ドキュメント生成

```bash
# GoDoc起動
godoc -http=:6060

# 特定パッケージのドキュメント確認
go doc github.com/zensai3805/iwakero-gamebook-mapping/internal/domain
```

## Git操作

### ブランチ管理

```bash
# Feature Branch作成
git checkout -b feature/issue-42

# ブランチ一覧
git branch -a

# リモートブランチ更新
git fetch --prune
```

### コミット

```bash
# ステージング
git add .

# コミット（メッセージ付き）
git commit -m "機能: パラグラフ追加機能の実装"

# 修正コミット
git commit --amend
```

### プッシュ・PR

```bash
# プッシュ
git push origin feature/issue-42

# PR作成
gh pr create

# PR状態確認
gh pr status
```

## トラブルシューティング

### よくある問題

```bash
# go.modエラー
go mod tidy

# キャッシュクリア
go clean -modcache

# テストキャッシュクリア
go clean -testcache

# 完全クリーン
go clean -cache
```