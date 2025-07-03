# iwakero-gamebook-mapping

🎮 ゲームブック（選択型冒険小説）のプレイ支援ツール。訪れたパラグラフと選択を記録し、ゲームブックの全体構造を可視化します。

## ✨ 機能

### v0.1.5 - リッチ対話モード 🎯
- **🎨 美しいターミナルUI**: PTerm使用でカラフルなインターフェース
- **📊 自動状態表示**: 現在のゲーム状況を即座に確認
- **⚡ 最適化されたメニュー**: 利用頻度に基づく効率的な操作
- **🔄 リアルタイム更新**: パラグラフ追加・選択が即座に反映

### 基本機能
- 📝 パラグラフの訪問記録・管理
- 🎯 選択肢の追跡（選択済み・未選択の状態管理）
- 🗺️ フロー図の自動生成（Mermaid形式）
- 📈 プレイ統計情報の表示

## 🚀 使用方法

```bash
# 🎮 対話モード（推奨）
./gamebook

# 📋 CLIコマンド
./gamebook new "ゲーム名"          # 新規ゲーム作成
./gamebook add 1 "冒険の始まり"     # パラグラフ追加
./gamebook choice 1 "北へ進む" 5   # 選択肢追加
./gamebook select 1 1             # 選択肢選択・移動
./gamebook show                   # 現在状態表示
```

## 開発方法

このプロジェクトはテスト駆動開発（TDD）で開発されています。

### 💻 必要な環境

- Go 1.23以上 (PTerm要件)

### 🔧 クイックスタート

```bash
# ビルド
go build -o gamebook ./cmd/gamebook

# 🎮 対話モード開始
./gamebook
```

### 🧪 開発・テスト

```bash
# テスト実行
go test ./...

# Lint実行
golangci-lint run

# ローカル開発時の必須チェック
go test ./... && golangci-lint run
```

## 🏗️ プロジェクト構造

```
.
├── cmd/gamebook/              # CLIアプリケーション
│   ├── main.go               # エントリーポイント（対話モード・CLI分岐）
│   ├── commands.go           # CLIコマンド実装（CommandExecutor）
│   ├── interactive_pterm.go  # PTerm対話モード（リッチUI）
│   ├── interactive.go        # 基本対話モード
│   └── choice_commands.go    # 選択肢関連コマンド
├── internal/
│   ├── domain/              # ドメインモデル（パラグラフ、ゲームブック等）
│   └── infrastructure/     # インフラ層
│       └── repository/      # データ永続化（Markdown形式）
└── test/                    # 統合テスト
```

## 🎯 ロードマップ

- ✅ **v0.1.0**: 基本機能（記録、フロー図、エラー検出）
- ✅ **v0.1.5**: PTerm対話モード（リッチUI、CLIラッパー）
- 🔄 **v0.2.0**: PTerm可視化機能（動的フロー、インタラクティブ2Dマップ）
- 🚀 **v1.0.0**: 音声入力機能

## ライセンス

MIT