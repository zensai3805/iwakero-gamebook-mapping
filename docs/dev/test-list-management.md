# テストリスト管理方針

## 概要
TDD開発におけるテストリストの作成・管理・保存方針を定義します。

## 基本方針

### 1. テストリスト管理の目的
- **TDD第1ステップの支援**: テストリスト作成・保存・更新の体系化
- **プロジェクト全体での一貫性**: .goファイル単位での統一的管理
- **AI開発最適化**: feature-name揺らぎ防止、永続的価値の確保
- **進捗追跡**: テスト実装状況の可視化

**TDD実装方針の詳細** → `docs/dev/tdd-policy.md` 参照

### 2. 配置場所

#### メインディレクトリ
```
docs/dev/test-lists/
```

#### ファイル命名規則
```
docs/dev/test-lists/
├── entities-{go-filename}.md                 # 例: entities-navigation.md
├── usecase-{go-filename}.md                  # 例: usecase-navigation_manager.md
├── infrastructure-{go-filename}.md           # 例: infrastructure-current_state.md
└── frameworks-{go-filename}.md               # 例: frameworks-tree_printer.md
```

**重要**: 
- .goファイル名と完全一致させることで、AI開発でのfeature-name揺らぎを防止
- Issue番号ベースのファイルは作成しない（Issue close後に不要になるため）

### 3. テストリストの種類

#### .goファイル別テストリスト  
- **ファイル名**: `{layer}-{go-filename}.md`
- **目的**: 特定.goファイルの詳細テストケース + 関連するエンドツーエンドテスト
- **作成タイミング**: Sub-Issue 0で.goファイル名確定後、各Sub-Issue開始時
- **管理者**: 各レイヤー実装担当者
- **重要**: .goファイル名が変更される場合、テストリストも同時に変更

#### 統合テストの管理方針
- Issue別統合テストリストは作成しない
- エンドツーエンドテストは、主要な.goファイルのテストリストに含める
- 複数ファイルに跨る統合テストは、最も中心的な.goファイルのテストリストで管理

## テストリスト作成プロセス

### 1. Sub-Issue 0（.goファイル設計とテストリスト準備）

#### 作成手順
1. **仕様書確認**: SPECIFICATION.md + 関連仕様書
2. **期待振る舞い分析**: ユーザー視点での機能要件
3. **.goファイル名確定**: 各レイヤーで作成する.goファイル名を決定
4. **レイヤー別責務分割**: 各レイヤーの責務に応じた分割
5. **インターフェース設計**: API設計の明確化
6. **テストリストファイル作成**: 各.goファイル用のテストリストファイルを作成

#### .goファイル設計記録（一時的ドキュメント）
Sub-Issue 0では、実装する.goファイル名を確定し、以下の形式で記録します：

```markdown
# {機能名} .goファイル設計

## 概要
- **機能**: {機能名}
- **関連Issue**: #{issue-number}（参考）
- **作成日**: {YYYY-MM-DD}

## .goファイル設計
### Entities Layer
- `internal/domain/{filename}.go` → テストリスト: `entities-{filename}.md`

### Usecase Layer  
- `internal/usecase/{filename}.go` → テストリスト: `usecase-{filename}.md`

### Interface Adapters Layer
- `internal/infrastructure/{filename}.go` → テストリスト: `infrastructure-{filename}.md`

### Frameworks & Drivers Layer
- `cmd/gamebook/{filename}.go` → テストリスト: `frameworks-{filename}.md`

## インターフェース設計
```go
// 主要インターフェース定義
type {InterfaceName} interface {
    {Method}({params}) {return}
}
```

## 実装制約
- {制約事項1}
- {制約事項2}
```

**注意**: このドキュメントは.goファイル名確定後に削除（一時的な設計記録のみ）

### 2. .goファイル別テストリスト作成

#### 作成タイミング
各Sub-Issue A～D開始時に、.goファイルごとの詳細テストリスト作成

#### テストリスト作成の要点
**TDD原則の詳細** → `docs/dev/tdd-policy.md` 参照

**管理上の要点**:
1. **.goファイル名との一致**: テストリストファイル名と実装ファイル名の対応
2. **段階的な詳細化**: Sub-Issue 0で概要、各Sub-Issueで詳細化
3. **進捗管理**: チェックボックスによる実装状況の追跡
4. **永続的保存**: Issue close後も価値を持つ管理

#### テンプレート
```markdown
# {go-filename}.go テストリスト

## 対象ファイル
- {layer-path}/{go-filename}.go
- {layer-path}/{go-filename}_test.go

## TODOリストとしての期待する振る舞い

### 基本的な動作
- [ ] {機能}が期待通りに動作する
- [ ] {機能}で適切な戻り値が返される  
- [ ] {機能}で状態が適切に更新される
- [ ] 複数回の操作で{機能}が動作する

### エラー処理
- [ ] 無効な入力で{機能}でエラーが発生する
- [ ] nil入力で{機能}でエラーが発生する
- [ ] 空の入力で{機能}でエラーが発生する
- [ ] 重複データで{機能}でエラーが発生する

### 境界条件
- [ ] 最小値で{機能}が動作する
- [ ] 最大値で{機能}が動作する
- [ ] 最小値-1で{機能}でエラーが発生する
- [ ] 最大値+1で{機能}でエラーが発生する
- [ ] nullや未定義値で{機能}の動作が適切である

### 「もしも」の状況
- [ ] サービスタイムアウト時に{機能}でエラーが発生する
- [ ] データベース接続失敗時に{機能}でエラーが発生する
- [ ] ネットワーク障害時に{機能}でエラーが発生する

### 既存機能との関係
- [ ] 新機能が既存の{関連機能}を破壊しない
- [ ] 既存の{API}契約が維持される
- [ ] 既存の{データ整合性}が保たれる

### 統合動作（該当する場合のみ）
- [ ] 統合環境で{機能}が動作する
- [ ] 他のコンポーネントとの連携で{機能}が動作する

## 実装制約
- {このファイル固有の制約事項}
- 振る舞いのみ記載、実装詳細は記載しない
- 1つのテストは1つの振る舞いのみ検証

## TDD進行管理
- [ ] 全振る舞いのテスト実装完了
- [ ] 全テストが通過
- [ ] リファクタリング完了
- [ ] Lintエラーなし
```

## ライフサイクル管理

### 1. 作成フェーズ
```bash
# ディレクトリ作成
mkdir -p docs/dev/test-lists

# .goファイル別テストリスト作成（各Sub-Issue開始時）
# .goファイル名ベースで作成
touch docs/dev/test-lists/entities-navigation.md                    # navigation.go
touch docs/dev/test-lists/usecase-navigation_manager.md             # navigation_manager.go  
touch docs/dev/test-lists/infrastructure-navigation_repository.md   # navigation_repository.go
touch docs/dev/test-lists/frameworks-tree_printer.md                # tree_printer.go修正
```

### 2. 運用フェーズ
- **更新**: 実装進捗に応じてチェック状態を更新
- **追加**: 新たに発見されたテストケースを追加
- **参照**: 次のテスト実装時の参照資料として活用

#### .goファイル追加時の対応
```bash
# 新しい.goファイルが必要になった場合
# 1. テストリスト作成
touch docs/dev/test-lists/{layer}-{new_filename}.md

# 2. 統合テストリストに追加
# issue-{number}-overview.md の「.goファイル設計」セクションを更新

# 3. 既存テストリストからの移行
# 関連テストケースを新しいテストリストに移動または複製
```

#### .goファイル名変更時の対応
```bash
# ファイル名変更の場合
# 1. テストリストもリネーム
mv docs/dev/test-lists/{layer}-{old_name}.md docs/dev/test-lists/{layer}-{new_name}.md

# 2. git履歴保持
git mv docs/dev/test-lists/{layer}-{old_name}.md docs/dev/test-lists/{layer}-{new_name}.md
```

### 3. 完了フェーズ
- **保存**: 機能実装完了後も残存（将来の機能拡張時の参考）
- **アーカイブ**: 古いバージョンは `docs/dev/test-lists/archive/` に移動

## Git管理方針

### コミット対象
- **.goファイル別テストリスト**: 各Sub-Issueブランチにコミット
- **完了後**: メインブランチにマージ

### .gitignore例外
```gitignore
# テストリストは必ずコミット対象
!docs/dev/test-lists/**/*.md
```

## AI最適化との連携

### Sub-Issue 0での活用
- **トークン削減**: 実装コードを読まずに.goファイル設計と振る舞い分析
- **スコープ制御**: .goファイル別責務の明確化
- **並行開発**: インターフェース確定による独立開発

### 各Sub-Issueでの活用  
- **集中実装**: 該当.goファイルのテストリストのみ参照
- **品質確保**: テスト先行による実装品質向上
- **進捗追跡**: チェック状態による可視化
- **永続性**: Issue close後も.goファイルと共に残存

## ツール支援

### テストリスト生成スクリプト（将来実装予定）
```bash
# .goファイル別テストリスト生成
./scripts/generate-test-list.sh entities navigation
./scripts/generate-test-list.sh usecase navigation_manager

# 複数ファイル一括生成
./scripts/generate-test-lists.sh entities-navigation usecase-navigation_manager
```

### GitHub Integration（将来実装予定）
- PR作成時のテストリスト完了確認
- .goファイル変更時の対応テストリスト自動更新
- テストカバレッジと連携した進捗確認

## 品質基準（TDD方針準拠）

### 管理品質基準

#### ファイル管理要件
- [ ] .goファイル名とテストリストファイル名が一致
- [ ] テストリストが該当する.goファイルの責務に適合
- [ ] Issue close後も価値を持つ内容で構成
- [ ] 重複ファイルの回避（feature-name揺らぎ防止）

#### 内容品質要件
**TDD品質基準の詳細** → `docs/dev/tdd-policy.md` 参照

- [ ] 振る舞いベースの記述（実装詳細回避）
- [ ] 適切な粒度での項目分割
- [ ] 進捗追跡可能なチェックボックス構成

### 推奨要件
- [ ] テストリスト項目の優先度設定
- [ ] 関連する.goファイルとの依存関係明記
- [ ] Sub-Issue間での整合性確保
- [ ] 定期的な見直しと更新

## テストリスト更新フロー

### 実装済みテストの更新手順

#### 1. 実装進捗に応じた更新
```markdown
# テストリスト例
- [x] 基本的な操作が期待通りに動作する ← ✅ 実装完了時にチェック
- [ ] エラー処理が適切に行われる      ← 次の実装対象
```

#### 2. 更新タイミング
- **RED→GREEN→REFACTORサイクル完了後**: 該当テストにチェックマーク
- **新たなテストケース発見時**: テストリストに追加
- **実装制約変更時**: 関連テストケースの見直し

#### 3. 更新責任
- **実装者**: 各TDDサイクル完了時に即座に更新
- **レビュアー**: PR時に更新状況を確認

#### 4. 更新確認コマンド
```bash
# テストリスト更新状況確認
grep -r "\- \[x\]" docs/dev/test-lists/
grep -r "\- \[ \]" docs/dev/test-lists/
```

### ドキュメント更新フロー統合

このテストリスト更新フローは、以下の開発ガイドと連携します：

- **TDD実装時**: `docs/dev/tdd-policy.md` のRED→GREEN→REFACTORサイクル完了後に更新
- **品質チェック時**: `docs/dev/lint-configuration.md` のチェック項目に含める
- **PR作成時**: `docs/dev/pull-request-management.md` の確認項目として含める

---

**更新履歴**
- 2025-07-07: テストリスト更新フローを追加（実装済みテスト管理プロセス標準化）
- 2025-07-07: tdd-policy.mdとの役割分離（管理特化、TDD実装詳細は分離）
- 2025-07-07: TDD正典準拠用語に修正（「正常系・異常系」を「基本動作・エラー処理」等に変更）
- 2025-07-07: Kent Beck TDD正典準拠に強化（「もしも」観点、TODOリストとして整理）
- 2025-07-07: TDD方針準拠に修正（振る舞い網羅、振る舞い記述）
- 2025-07-07: .goファイル名ベース管理方針に変更（Issue番号依存を削除）
- 2025-07-07: 初版作成（Issue #114対応）