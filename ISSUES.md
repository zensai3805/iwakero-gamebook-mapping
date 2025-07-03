# 現在のIssue管理

## Issue #1: 🐛 選択肢の読み込み処理にバグがある

**Status**: Open  
**Priority**: High  
**Assignee**: Claude  

### 問題の説明

現在、選択肢をMarkdownファイルに保存することはできているが、ファイルから読み込む際に選択肢情報が失われている。

### 現在の状況

#### ✅ 正常に動作している機能
- 選択肢の保存（Markdownファイル形式）
- Mermaidフロー図の自動生成
- パラグラフ間のネットワーク関係の記録

#### ❌ 問題のある機能
- 選択肢の読み込み（自動ロード時）
- `./gamebook show` で選択肢が表示されない

### 再現手順

1. `./gamebook new "テスト"`
2. `./gamebook add 1 "開始"`
3. `./gamebook add 2 "次"`
4. `./gamebook choice 1 "進む" 2`
5. `./gamebook show` → 選択肢が表示されない

### 期待される結果

`show` コマンドで選択肢情報が正しく表示される

### 技術的詳細

- 問題箇所: `internal/infrastructure/repository/markdown_repository.go` の Load メソッド
- 自動ロード機能により、過去の状態が上書きされている可能性

### タスク

- [ ] **必須**: SPECIFICATION.md を読み直して要件を確認
- [ ] 選択肢読み込み処理のデバッグ
- [ ] テストケースの追加
- [ ] 修正の実装
- [ ] 動作確認

### 関連ファイル

- `internal/infrastructure/repository/markdown_repository.go`
- `internal/infrastructure/repository/markdown_repository_test.go`
- `SPECIFICATION.md`（要件確認用）

---

## Issue #2: 📝 プロジェクト管理体制の確立

**Status**: In Progress  
**Priority**: Medium  
**Assignee**: Claude  

### 説明

GitHub Issue管理への移行とワークフロー確立

### タスク

- [x] PROJECT_MANAGEMENT.md作成
- [x] ISSUES.md作成（このファイル）
- [ ] **必須**: SPECIFICATION.md確認
- [ ] 今後の開発フロー確立
- [ ] ユーザーへの移行報告

### 関連ファイル

- `PROJECT_MANAGEMENT.md`
- `SPECIFICATION.md`

---

## 今後のIssue作成予定

### v0.1.0完成に向けて

1. **Issue #3**: 🧪 選択肢読み込みのテストケース強化
2. **Issue #4**: 🎤 音声入力機能の設計（低優先度）
3. **Issue #5**: 🗺️ マップ自動生成機能の設計（低優先度）

### 注意

- **すべてのIssueでSPECIFICATION.mdの確認を必須とする**
- 新しいタスクが発生した場合は即座にIssue追加
- SPECIFICATION.md更新が必要な場合はユーザーに確認