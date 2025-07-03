# iwakero-gamebook-mapping

ゲームブック（選択型冒険小説）のプレイ支援ツール。訪れたパラグラフと選択を記録し、ゲームブックの全体構造を可視化します。

## 機能

- パラグラフの訪問記録
- 選択肢の追跡（選択済み・未選択）
- フロー図の自動生成（Mermaid形式）
- テキストベースマップの生成（今後実装予定）

## 開発方法

このプロジェクトはテスト駆動開発（TDD）で開発されています。

### 必要な環境

- Go 1.21以上

### セットアップ

```bash
# 依存関係のインストール
make deps
```

### テストの実行

```bash
# テスト実行
make test

# 詳細なテスト結果を表示
make test-v

# カバレッジ付きテスト
make test-cover

# カバレッジレポートをブラウザで表示
make test-cover-html
```

### ビルド

```bash
make build
```

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

MIT