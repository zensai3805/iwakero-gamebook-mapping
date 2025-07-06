# ゲームブック記録支援ツール仕様書

## プロジェクト概要
ファイティングファンタジーなどのゲームブックプレイ時に、訪れたパラグラフの記録と選択の追跡を支援し、ゲームブックの全貌を解明することを目的とした可視化ツール。

## 仕様書構成（クリーンアーキテクチャ4層）

### 🎯 Entities Layer（ビジネスロジック）
→ `docs/spec/entities-spec.md`
- エンティティ定義（Gamebook, Paragraph, Choice, Session）
- ビジネスルール（パラグラフ管理、選択肢管理、移動ルール）
- データ整合性

### 📋 Usecase Layer（アプリケーションロジック）
→ `docs/spec/usecase-spec.md`
- ユースケース（ゲームブック管理、パラグラフ管理、選択肢管理）
- ワークフロー（プレイフロー、データ入力フロー）
- エラーハンドリング

### 🔌 Interface Adapters Layer（外部接続）
→ `docs/spec/interface-adapters-spec.md`
- Repository（Markdown永続化、セッション管理）
- Presenter（TreePrinter, AreaPrinter, Mermaid）
- Controller（CLI, Interactive）
- 入力支援（InputHelper）

### 🖥️ Frameworks & Drivers Layer（フレームワーク）
→ `docs/spec/frameworks-drivers-spec.md`
- CLI Framework（cobra）
- UI Framework（PTerm）
- 外部ツール（Git, Make）
- プラットフォーム対応

## 開発状況
→ `CHANGELOG.md` - バージョン履歴と開発状況

## 技術スタック
- **言語**: Go (Golang)
- **UI**: PTerm (リッチターミナルUI)
- **データ**: Markdown形式
- **プラットフォーム**: Windows/macOS/Linux

## 制約事項
- Gitリポジトリを通じた共有のみ（クラウド同期なし）
- テキストベースのみ（画像生成なし）