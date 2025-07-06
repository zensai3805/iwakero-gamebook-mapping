# Pull Request管理

## PR作成の基本方針

- **PRでは必ずCLAUDE.mdを最新化する**
- **必ずAI最適化PRテンプレート使用**（プルリクエスト作成時に自動適用）
- **ユーザーマージ**: PRのマージは必ずユーザーが実行

## 開発フロー

1. **Issue確認**: SPECIFICATION.md・DEVELOPMENT.mdを確認してから作業開始
2. **featureブランチ作成**: AI管理スクリプト使用
3. **TDD実行**: RED → GREEN → REFACTOR
4. **品質チェック**: AI管理スクリプト使用
5. **PR作成**: 以下の手順で実行
6. **ユーザーマージ**: PRのマージは必ずユーザーが実行

## PR作成手順

### 1. コミット前必須チェック
- **ローカルLint**: `golangci-lint run` でエラーなし
- **ローカルテスト**: `go test ./...` で全テスト通過
- **commit・push前に必ず実行**

### 2. PR作成前チェック
- 関数の長さ50行以内
- 複雑度20以下
- 適切なエラーハンドリング
- 日本語コメントの追加
- 共有状態の適切な同期
- テストカバレッジの確認
- **CLAUDE.mdの更新**（必須）

### 3. PR作成コマンド

```bash
# PR作成（テンプレート自動適用）
gh pr create

# またはスクリプト使用
./scripts/claude-project-manager.sh pr {issue-number} "PRタイトル" "PR説明"
```

## PRテンプレート

自動適用されるPRテンプレートには以下が含まれます：
- 簡潔な品質チェックリスト
- AI実装に必要十分な内容
- CLAUDE.md更新確認

## 品質チェック

### 自動化コマンド
```bash
# 品質チェック実行
./scripts/claude-project-manager.sh quality

# または手動実行
go test ./...
golangci-lint run
gofmt -s -w .
goimports -w .
```

## PR承認とマージ

- **レビュー**: 品質チェックリストに基づく確認
- **マージ**: 必ずユーザーが実行
- **Issue Close**: PRマージによる自動クローズ