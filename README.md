# iwakero-gamebook-mapping

🎮 ゲームブック（選択型冒険小説）のプレイ支援ツール。訪れたパラグラフと選択を記録し、ゲームブックの全体構造を可視化します。

## 🚨 重要: 開発者向け注意事項

### Claude Code 使用時の推奨設定

`/clear` 後の固定プロンプトとして以下を設定することを強く推奨します：

```
以後すべてのアクションは CLAUDE.md から辿るアクション別のドキュメントを読んでから実施する事を守ってください。
```

### 開発フロー
1. **必須**: `CLAUDE.md` → `DEVELOPMENT.md` → 対応する専門ドキュメントの順で確認
2. **効果**: トークン消費最適化、開発品質向上、一貫性確保
3. **詳細**: すべての開発ガイドラインは `CLAUDE.md` に記載済み

## ✨ 主な機能

### 📊 可視化機能 (v0.2.0)
- **フロー図**: 階層的な選択肢構造を視覚化
- **2Dマップ**: ASCIIグリッドによる空間的配置
- **ヒートマップ**: パラグラフの訪問頻度を色彩で表現
- **統計情報**: プレイ状況の詳細分析

### 🎨 リッチUI (v0.1.5)
- **PTerm対話モード**: カラフルなターミナルインターフェース
- **メニュー選択**: 矢印キーとEnterで直感的操作
- **リアルタイム状態表示**: 現在位置と選択肢を即座に表示
- **入力支援**: 自動補完とプレースホルダー表示

### 📝 基本機能
- パラグラフの訪問記録・管理
- 選択肢の追跡（選択済み・未選択の状態管理）
- 複数ゲームブック対応
- Markdown形式でのデータ保存

## 🚀 使用方法

### 推奨: 対話モード
```bash
./gamebook
```
- 美しいターミナルUIで直感的操作
- 現在の状況を常に表示
- メニューから機能を選択

### CLI コマンド
```bash
# 新規ゲームブック作成
./gamebook new "ゲーム名"

# パラグラフ追加
./gamebook add 1 "冒険の始まり"

# 選択肢追加
./gamebook choice 1 "北へ進む" 5

# 選択実行
./gamebook select 1 1

# 現在状態表示
./gamebook show
```

## 📁 プロジェクト構造

```
iwakero-gamebook-mapping/
├── cmd/gamebook/              # Frameworks Layer (UI, CLI)
│   ├── main.go               # エントリーポイント
│   ├── commands.go           # CLIコマンド実装
│   ├── interactive_pterm.go  # PTerm対話モード
│   └── interactive.go        # 基本対話モード
├── internal/
│   ├── domain/              # Entities Layer (純粋なビジネスロジック)
│   │   ├── gamebook.go      # ゲームブックエンティティ
│   │   ├── paragraph.go     # パラグラフエンティティ
│   │   └── choice.go        # 選択肢エンティティ
│   ├── usecase/             # Usecase Layer (アプリケーションロジック)
│   │   ├── gamebook_usecase.go
│   │   └── paragraph_usecase.go
│   ├── infrastructure/      # Interface Adapters (外部インターフェース)
│   │   ├── repository/      # データ永続化
│   │   └── presenter/       # データ表示
│   └── interface/           # Interface Adapters (コントローラー)
├── docs/                    # ドキュメント
│   ├── dev/                # 開発ガイド (8ファイル)
│   └── spec/               # 仕様書 (4層構造)
├── data/                   # ゲームデータ (Markdown)
└── test/                   # 統合テスト
```

## 🛠️ 開発環境

### 必要な環境
- **Go**: 1.23以上 (PTerm要件)
- **Git**: バージョン管理
- **GitHub CLI**: Issue/PR管理 (推奨)

### クイックスタート
```bash
# ビルド
go build -o gamebook ./cmd/gamebook

# 実行
./gamebook

# テスト
go test ./...

# Lint
golangci-lint run
```

### 開発原則
- **TDD (Test-Driven Development)**: テストファースト開発
- **Issue-Driven Development**: GitHub Issue による作業管理
- **Clean Architecture**: 4層アーキテクチャによる責任分離

## 📈 バージョン情報

- **現在**: v0.2.5 (入力支援機能実装)
- **前版**: v0.2.0 (PTerm可視化機能完成)
- **次版**: v0.2.6 (フロー図スケーラビリティ対応)

詳細は `CHANGELOG.md` を参照してください。

## 🎯 ロードマップ

- ✅ **v0.1.0**: 基本機能（記録、フロー図、エラー検出）
- ✅ **v0.1.5**: PTerm対話モード（リッチUI）
- ✅ **v0.2.0**: PTerm可視化機能（フロー図、2Dマップ）
- 🔄 **v0.3.0**: 操作性の大幅改善
- 🚀 **v0.4.0**: スマートフォン連携WebUI
- 🌟 **v1.0.0**: 音声入力機能

## 📚 技術スタック

- **言語**: Go (Golang)
- **UI**: PTerm (リッチターミナルUI)
- **データ**: Markdown形式
- **アーキテクチャ**: Clean Architecture (4層)
- **開発**: TDD, Issue-Driven Development
- **プラットフォーム**: Windows/macOS/Linux

## 🤝 貢献

このプロジェクトは厳格な開発プロセスを採用しています：

1. **Issue作成**: GitHub Issueで作業を管理
2. **ブランチ作成**: `feature/issue-XX` 形式
3. **TDD実装**: テストファースト開発
4. **品質チェック**: テスト・Lint・フォーマット
5. **PR作成**: レビュー後マージ

詳細は `CLAUDE.md` → `DEVELOPMENT.md` を参照してください。

## ライセンス

MIT License