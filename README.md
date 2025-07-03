# 岩ケロのゲームブック・マッピングツール

[![CI](https://github.com/zensai3805/iwakero-gamebook-mapping/actions/workflows/ci.yml/badge.svg)](https://github.com/zensai3805/iwakero-gamebook-mapping/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/zensai3805/iwakero-gamebook-mapping)](https://goreportcard.com/report/github.com/zensai3805/iwakero-gamebook-mapping)

ゲームブック（選択型冒険小説）のプレイ支援ツール。訪れたパラグラフと選択を記録し、ゲームブックの全体構造を可視化します。

## 機能

### v0.1.0 基本機能
- ✅ **ゲームブック管理**: 新規作成、読み込み、一覧表示、削除
- ✅ **パラグラフ管理**: 追加、削除、概要編集
- ✅ **選択肢管理**: 追加、削除、選択・移動
- ✅ **フロー図生成**: Mermaid形式でのリアルタイム可視化
- ✅ **セッション管理**: 現在のゲーム状態の永続化
- ✅ **エラー検出**: パラグラフ重複検出、未定義遷移先警告
- ✅ **データ整合性**: 原子的保存、包括的検証、破損データ処理

## インストール

```bash
# リポジトリをクローン
git clone https://github.com/zensai3805/iwakero-gamebook-mapping.git
cd iwakero-gamebook-mapping

# ビルド
go build -o gamebook ./cmd/gamebook

# 実行
./gamebook --help
```

## 基本的な使い方

```bash
# 新しいゲームブックを作成
./gamebook new "冒険の書"

# パラグラフを追加
./gamebook add 1 "物語の始まり"

# 選択肢を追加
./gamebook choice 1 "北へ進む" 2
./gamebook choice 1 "南へ進む" 3

# 選択肢を選択して移動
./gamebook select 1 1  # パラグラフ1の1番目の選択肢を選択

# 現在の状態を確認
./gamebook show

# ゲームブック一覧
./gamebook list
```

## アーキテクチャ

- **Clean Architecture**: ドメイン駆動設計による保守性の高い設計
- **Repository Pattern**: データ永続化の抽象化
- **CLI Interface**: Cobraライブラリによる使いやすいコマンドライン
- **Markdown Storage**: 人間が読めるMarkdown形式でのデータ保存

## 開発

### 前提条件
- Go 1.21以降
- Git

### テスト実行
```bash
# 単体テスト
go test -v ./internal/...

# カバレッジ付き
go test -v -race -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out
```

### 開発フロー
1. Issue作成・SPECIFICATION.md確認
2. Feature branchで開発
3. テスト追加・実装
4. PR作成・レビュー
5. mainブランチにマージ

## プロジェクト構造

```
.
├── cmd/gamebook/          # CLIアプリケーションのエントリーポイント
├── internal/
│   ├── domain/           # ドメインモデル（パラグラフ、ゲームブック等）
│   ├── usecase/          # ユースケース層
│   ├── infrastructure/   # インフラ層
│   │   ├── repository/   # データ永続化
│   │   └── presenter/    # 出力フォーマット（Mermaid等）
│   └── interface/        # インターフェース層
│       └── cli/          # CLIコマンド
└── test/                 # 統合テスト
```

## ライセンス

MIT License

## 作者

- **zensai3805** - *Initial work* - [GitHub](https://github.com/zensai3805)

## Contributing

1. このリポジトリをFork
2. Feature branchを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をCommit (`git commit -m 'Add amazing feature'`)
4. BranchをPush (`git push origin feature/amazing-feature`)
5. Pull Requestを作成

詳細は[SPECIFICATION.md](SPECIFICATION.md)を参照してください。