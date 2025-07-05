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
2. **featureブランチ作成**: `git checkout -b feature/issue-{number}`
3. **TDD実行**: RED → GREEN → REFACTOR
4. **コミット**: Issue単位で適切な粒度
5. **PR作成**: ユーザー確認が必要な段階で作成
6. **ユーザーマージ**: PRのマージは必ずユーザーが実行

### 必須チェック
- **ローカルLint**: `golangci-lint run` でエラーなし
- **ローカルテスト**: `go test ./...` で全テスト通過
- **commit・push前に必ず実行**

## GitHub公式Sub-Issue管理

### 使用方法
- **マイルストーン開始時にSub-Issue洗い出しを必須実施**
- REST API: `POST /repos/{owner}/{repo}/issues/{issue_number}/sub_issues`
- 各Sub-IssueでもSPECIFICATION.md・DEVELOPMENT.md確認必須

### 洗い出しフロー
1. **SPECIFICATION.md・DEVELOPMENT.md確認**: 要件と開発方針の詳細確認
2. **現状分析**: 実装済み/未実装の棚卸し
3. **Sub-Issue作成**: 具体的タスクごとに作成
4. **依存関係整理**: 実装順序決定
5. **実装開始**: 最優先Sub-Issueから

## 重要資料の位置づけ

- **SPECIFICATION.md**: ツール仕様の最重要資料
- **DEVELOPMENT.md**: 開発方針の最重要資料
- Issue作業前に必ず両方確認
- 更新が必要な場合はユーザーに確認を求める

## 開発コマンド

```bash
# テスト実行
go test ./...

# Lint実行
golangci-lint run

# ビルド
go build -o gamebook ./cmd/gamebook

# 実行
./gamebook                      # インタラクティブモード
./gamebook new "GameTitle"      # 新規作成
./gamebook add 1 "Description"  # パラグラフ追加
./gamebook choice 1 "Go north" 2  # 選択肢追加
./gamebook select 1 1           # 選択実行
./gamebook show                 # 現在状態表示
```

## ラベル体系

- **priority**: high/medium/low
- **type**: epic（マイルストーン追跡）
- **area**: cli/domain/repository（技術領域）