# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Important Project Documentation

- **Tool Specification**: @SPECIFICATION.md - Complete functional and technical specifications
- **Development Guidelines**: @DEVELOPMENT.md - Development methodology, TDD approach, and coding standards

These documents are the most important references for understanding the project requirements and development practices.

## Project Overview

This is a gamebook (choose-your-own-adventure) playing assistant tool called "iwakero-gamebook-mapping". The project aims to help gamebook players track their progress through paragraphs and visualize their journey.

## Project Goals

### Primary Goals
- Create a tool to record paragraph summaries and player choices during gamebook sessions
- Use paragraph numbers as primary keys for all records
- Dynamically generate dungeon maps or flow diagrams based on recorded data

### Technical Preferences
- **Language**: Go (Golang) preferred
- **Platform**: Cross-platform support for Windows, macOS, and Linux
- **Data Format**: Prefer common formats like Markdown (potentially with Mermaid syntax for diagrams)
- **Input Method**: Voice input support desired for easier recording during gameplay

## Current Status

- **v0.2.5 In Progress**: 入力支援機能実装
- ✅ **v0.1.0**: Core functionality (recording, flow diagrams, error detection)
- ✅ **v0.1.5**: Interactive mode implementation complete
  - PTerm-based rich terminal UI with colors, tables, and menus
  - CLI wrapper architecture avoiding code duplication
  - Auto-display of current game status
  - Optimized menu ordering based on usage frequency
  - Test-driven development (TDD) workflow restored
- ✅ **v0.2.0**: PTerm visualization features complete
  - Tree Printer: Hierarchical flow visualization with PTerm styling
  - Area Printer: 2D ASCII grid-based maps with borders and positioning
  - Integrated UI: 3-split layout management with synchronized components
  - Heatmap: Visit frequency tracking with color gradients
  - Full TDD implementation across all visualization components
- 🔄 **v0.2.5**: 入力支援機能実装 (Issue #59)
  - Tab補完・候補表示・半透明ヒント
  - インテリジェント候補生成エンジン
  - 高度な入力バリデーション機能
  - 現代的なツールとしての入力体験向上

## Development Guidelines

When implementing this project:
1. **Always check SPECIFICATION.md and DEVELOPMENT.md before starting any work** - these are the most important documents
2. Follow Issue-driven development using GitHub Issues and PRs
3. Update CLAUDE.md in every PR
4. Use Test-Driven Development (TDD) approach - this is the highest priority
5. Ensure cross-platform compatibility from the start
6. Design the data structure to efficiently use paragraph numbers as primary keys

## AI Development Optimization

### Claude Code Specific Tools
- **AI Development Guide**: `CLAUDE_AI_DEVELOPMENT_GUIDE.md` - Claude Code最適化指示書
- **Project Manager Script**: `scripts/claude-project-manager.sh` - AI実行可能な管理コマンド
- **Issue Templates**: AI最適化されたGitHub Issue Templates
  - `ai_feature.yml` - Feature実装専用
  - `ai_sub_issue.yml` - Sub-Issue管理専用
  - `ai_bug_report.yml` - Bug Report専用
- **PR Template**: `.github/PULL_REQUEST_TEMPLATE.md` - 包括的な品質チェックリスト

### AI管理コマンド
```bash
# Claude専用プロジェクト管理
./scripts/claude-project-manager.sh feature "機能名" "説明" "v1.0.0"
./scripts/claude-project-manager.sh branch 42
./scripts/claude-project-manager.sh progress 42 "🤖 Claude: 実装完了"
./scripts/claude-project-manager.sh pr 42 "PR Title" "PR Description"
./scripts/claude-project-manager.sh quality

# GitHub CLI活用
gh issue create --template ai_feature
gh issue create --template ai_sub_issue
gh issue create --template ai_bug_report
```

## Project Management

### Issue-Driven Development
- All work is managed through GitHub Issues
- Issues must start with SPECIFICATION.md confirmation
- Issues are closed via PRs that update CLAUDE.md
- Use Sub-Issues for complex features

### Milestones
- **v0.1.0**: Basic functionality (recording, flow diagrams, error detection) ✅ **Complete**
- **v0.1.5**: Interactive mode with PTerm rich UI ✅ **Complete**
- **v0.2.0**: PTerm visualization features (dynamic flow, interactive 2D maps) ✅ **Complete**
- **v0.2.1**: 段落順序依存問題の解消 - プレイテスターフィードバック対応 (Issue #43, #46) ✅ **Complete**
- **v0.2.2**: フロー図の視認性と情報表示の改善 - 追加プレイテスターフィードバック対応 (Issue #53) ✅ **Complete**
- **v0.2.3-hotfix**: 現在地表示機能の緊急実装 - 致命的な基本機能欠如の修正 (Issue #67) 🔄 **最高優先**
- **v0.2.4**: 段落順序依存問題の解決 - 未定義移動先選択肢の暫定対応 (Issue #57)
- **v0.2.5**: フロー図スケーラビリティ対応 - 大規模ゲームブック対応 (Issue #58)
- **v0.2.6**: 入力支援機能実装 - Tab補完・候補表示・半透明ヒント (Issue #59)
- **v0.2.7**: 直接移動パラグラフ可視化対応 - moveコマンド対応 (Issue #60)
- **v0.3.0**: 操作性の大幅改善 - コマンド短縮とショートカット強化 (Issue #44)
- **v0.3.1**: プレイヤーシート機能実装 - 体力・技量・運・金貨・アイテム管理 (Issue #61)
- **v0.4.0**: スマートフォン連携 - WebUI実装によるハンズフリー操作 (Issue #45)
- **v0.5.0**: 東西南北対応マップ機能の再設計 - 方向概念とマップ座標系の実装 (Issue #47)
- **v1.0.0**: Voice input functionality

### Development Commands
```bash
# Test execution
go test ./...

# Build CLI
go build -o gamebook ./cmd/gamebook

# Run CLI
./gamebook [command]

# Interactive mode (no arguments):
./gamebook                      # Launch rich PTerm interactive mode

# CLI commands:
./gamebook new "GameTitle"      # Create new gamebook
./gamebook add 1 "Description"  # Add paragraph
./gamebook choice 1 "Go north" 2  # Add choice
./gamebook select 1 1           # Select choice and move
./gamebook show                 # Display current state
```


## Document Maintenance Strategy

This CLAUDE.md file should be updated:
- After major feature implementations
- When project structure changes significantly
- When new technical decisions are made
- Monthly review (use `/init` command if comprehensive update needed)

### Update History
- 2025-07-05: **v0.2.5 入力支援機能実装完了** - Issue #59対応
  - **現代的な入力支援機能の実装**: Tab補完・候補表示・ヒント機能
  - **インテリジェント候補生成エンジン**: 既存・未定義・推奨パラグラフの自動分類
  - **高度な入力バリデーション**: 重複検出・代替案提案・リアルタイム警告
  - **interactive_pterm.go統合**: パラグラフ追加・選択肢追加・移動処理に入力支援機能を適用
  - **TDD徹底**: RED→GREEN→REFACTORサイクルでの実装
  - **品質保証**: 全テスト通過・lint通過・ビルド成功
- 2025-07-05: **現在地表示機能の致命的な問題発見** - Issue #67を最高優先で緊急対応
  - **新たな最高優先問題発見**: 現在地表示機能の完全な欠如（⭐⭐⭐⭐⭐）
  - **技術的根本原因**: UI層での表示処理と永続化機能の未実装による基本機能の破綻
  - **開発者観点での解決策**: commands.goへの現在地表示処理追加とMarkdownRepository永続化実装
  - Issue #67を緊急作成（現在地表示機能の実装）
  - v0.2.3-hotfixとして緊急修正を予定
  - プレイヤー体験の根本的な問題として最優先対応
- 2025-07-05: **プレイテスターフィードバック対応（フロー図汚染問題発見）** - Issue #57を最高優先で再設計
  - **新たな最高優先問題発見**: 未定義パラグラフによるフロー図汚染問題（⭐⭐⭐⭐⭐）
  - **技術的根本原因**: TreePrinterの無差別表示による実用性の完全な破綻
  - **開発者観点での解決策**: グラフ接続性分析とフィルタリング表示システム
  - Issue #57を技術的に再設計（グラフ理論に基づく接続性分析実装）
  - 重複Issue（#62-65）のクリーンアップ完了
  - 段階的表示システム：接続のみ（デフォルト）→ 全て（管理用）→ 非表示（クリーン）
- 2025-07-05: **v0.2.2 Issue #53完了** - フロー図の視認性と情報表示の改善
  - TreePrinterで選択肢情報の完全表示を実装
  - TDD手法でRED→GREEN→REFACTORサイクルを徹底
  - 選択済み/未選択の視覚的区別（[✓]/[ ]シンボル）
  - 選択肢の階層表示により分岐構造が明確化
  - プレイヤビリティが大幅改善（何を選択するか一目で分かる）
- 2025-07-05: **v0.2.1.hotfix完了** - Issue #46マップ機能一時除去対応
  - showコマンドでマップ表示を一時的に無効化、フロー図のみ表示に変更
  - commands.go の showVisualization() を統合UIからTreePrinter直接使用に変更
  - TDD手法により回帰テスト防止（commands_hotfix_test.go追加）
  - 統合UIシステムは将来のマップ機能再実装に備えて保持
  - プレイテスター体験改善の第一歩として、複雑なマップ表示を簡素化
- 2025-07-05: **GitHub sub-issues管理方法の正確化** - GraphQL API使用方法を文書化
  - gh CLIがsub-issuesをネイティブサポートしていないことを確認・記録
  - 過去のClaude Codeが使用していたGraphQL API方法を再現・文書化
  - 実際にテスト実行してDEVELOPMENT.mdに動作確認済み手順を追加
  - IssueテンプレートをYMLからMarkdownに移行（.md拡張子で統一）
  - PRテンプレートにGitHub Copilot向けプロジェクトコンテキスト情報を追加
- 2025-07-05: **プレイテスターフィードバック対応** - 新バージョン計画策定
  - PLAY_FEEDBACK.md分析により致命的な不便さを特定
  - v0.2.1: 段落順序依存問題の解消 (Issue #43)
  - v0.3.0: 操作性の大幅改善 (Issue #44)
  - v0.4.0: スマートフォン連携WebUI実装 (Issue #45)
  - 実際の使用体験に基づく段階的改善計画
- 2025-07-05: **AI-Scrum システム分離完了** - 独立リポジトリへの移行
  - ai-scrum-system/ ディレクトリを独立したプロジェクトとして分離
  - iwakero-claude-scrum-agent リポジトリへの移行準備完了
  - Go TUI + Shell Script のハイブリッド設計で真のエージェント自律性実現
  - プロジェクト独立性により汎用的なスクラム開発環境を提供
  - AI開発方針（CLAUDE_AI_DEVELOPMENT_GUIDE.md等）を新リポジトリに継承
- 2025-07-05: AI開発最適化ツール実装 - Claude Code専用管理スクリプト、Issue Templates、PR Template、開発ガイド追加
- 2025-07-05: v0.2.0 completed - PTerm visualization features with comprehensive TDD implementation
- 2025-07-04: Initial version created
- 2025-07-03: Issue #1 resolved - Fixed critical choice loading bug in markdown_repository.go
- 2025-07-03: v0.1.5 completed - PTerm interactive mode implementation with CLI wrapper architecture

## Next Steps

1. Initialize a Go module for the project
2. Create the basic project structure
3. Design the data model for storing paragraph information
4. Implement core functionality for recording paragraph visits and choices
5. Add visualization capabilities for maps and flow diagrams

## Technology Stack

- **Backend**: Go (Golang)
- **Data Storage**: Markdown files with structured format
- **UI**: PTerm (modern Go terminal UI library)
- **Visualization**: Mermaid.js for diagram generation + PTerm for interactive displays
- **Cross-platform**: Standard Go build for multiple OS
- **Interface**: CLI + Rich interactive mode
- **Future**: Voice input integration planned