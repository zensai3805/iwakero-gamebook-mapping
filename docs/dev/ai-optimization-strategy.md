# AI開発最適化：スコープ制御戦略

## 課題：コード全体改修によるトークン消費問題

- **現状の問題**: 実装時にコード全体が改修範囲に含まれる
- **影響**: トークン消費量過多、AIの関心範囲の拡散
- **解決策**: cmd/internal分離によるSub-Issue戦略

## クリーンアーキテクチャ4層分離

### 1. 機能実装時の必須分離

**すべての機能実装を以下の6つのSub-Issueに分離**

#### Sub-Issue 0: テストリスト作成（TDD第1ステップ）
- テストリスト.md - 期待する振る舞いの網羅的分析
- テストケース設計 - 各レイヤーのテスト要件定義
- インターフェース設計 - 利用者視点でのAPI設計
- 実装設計分離 - 実装詳細は後回し

#### Sub-Issue A: Entities Layer（エンティティ層）
- `internal/domain/` - エンティティ・値オブジェクト・ドメインサービス
- 純粋ビジネスロジック・ドメインルール
- 外部依存なし（Go標準ライブラリのみ）
- 単体テスト（モック不要）

#### Sub-Issue B: Usecase Layer（ユースケース層）
- `internal/usecase/` - アプリケーションサービス・ユースケース
- ビジネスワークフロー・アプリケーションルール
- Repository/Presenterインターフェース定義
- 単体テスト（Repository/Presenterモック使用）

#### Sub-Issue C: Interface Adapters Layer（インターフェースアダプター層）
- `internal/infrastructure/repository/` - データアクセス実装
- `internal/infrastructure/presenter/` - 外部出力実装
- `internal/interface/` - 入出力境界・コントローラー
- 結合テスト（Usecaseモック使用）

#### Sub-Issue D: Frameworks & Drivers Layer（フレームワーク・ドライバー層）
- `cmd/gamebook/` - フレームワーク・UI・エントリーポイント
- CLI・PTerm・可視化・依存性注入
- 統合テスト（エンドツーエンド）

#### Sub-Issue E: 統合検証・リファクタリング（TDD第5ステップ）
- 全レイヤー統合テスト実行
- レイヤー間の依存関係検証
- コードの重複除去・リファクタリング
- 過渡的コードの整理・最適化

### 2. 実装順序の原則

**必ずテストリスト作成 → Entities → Usecase → Interface Adapters → Frameworks → 統合検証の順序で実装**

1. **Sub-Issue 0**: テストリスト作成（振る舞い分析・インターフェース設計）
2. **Sub-Issue A**: Entities Layer（純粋ビジネスロジック）
3. **Sub-Issue B**: Usecase Layer（アプリケーションルール）
4. **Sub-Issue C**: Interface Adapters Layer（外部連携）
5. **Sub-Issue D**: Frameworks & Drivers Layer（UI・フレームワーク）
6. **Sub-Issue E**: 統合検証・リファクタリング（品質向上）

## トークン消費量最適化

### 読み込み対象の限定

| Sub-Issue | 読み込み対象 | 避けるべきディレクトリ |
|-----------|-------------|-------------------|
| 0 (テストリスト) | 仕様書・既存テストのみ | 実装コード全般 |
| A (Entities) | `internal/domain/` のみ | usecase/, infrastructure/, cmd/ |
| B (Usecase) | `internal/usecase/` のみ | infrastructure/, cmd/ の詳細 |
| C (Interface Adapters) | `internal/infrastructure/`, `internal/interface/` | cmd/ の詳細 |
| D (Frameworks) | `cmd/gamebook/` のみ | infrastructure/ の実装詳細 |
| E (統合検証) | 全レイヤー | なし（全体確認が必要） |

### AIプロンプトの最適化

- **Sub-Issue 0**: 「仕様書のみに集中、テストリスト作成、インターフェース設計、実装詳細は後回し」
- **Sub-Issue A**: 「Domainディレクトリのみに集中、純粋ビジネスロジック、外部依存なし」
- **Sub-Issue B**: 「Usecaseディレクトリのみに集中、アプリケーションルール、インターフェース定義」
- **Sub-Issue C**: 「Infrastructure・Interfaceディレクトリのみに集中、Usecaseインターフェース実装」
- **Sub-Issue D**: 「Frameworksディレクトリのみに集中、UI・依存性注入、Interface Adapters使用」
- **Sub-Issue E**: 「全レイヤー統合検証、リファクタリング、過渡的コード整理」

## Sub-Issue作成テンプレート

### Sub-Issue 0（テストリスト作成）

```markdown
## 目的
[機能名]のテストリスト作成・テスト設計

## 実装対象
- テストリスト.md（一時的なドキュメント）
- テストケース設計書
- インターフェース設計書

## 実装制約
- 実装コードは一切読まない
- 仕様書・既存テストのみ参照
- 振る舞いに焦点、実装詳細は後回し
- インターフェース設計を優先

## 完了条件
- [ ] 期待する振る舞いの網羅的分析完了
- [ ] 各レイヤーのテストケース明確化
- [ ] Sub-Issue A～D のテスト要件確定
- [ ] インターフェース設計完了
- [ ] 実装設計と分離済み
```

### Sub-Issue A（Entities Layer）

```markdown
## 目的
[機能名]のエンティティ・ドメインロジック実装

## 実装対象
- internal/domain/[エンティティ名].go
- internal/domain/[エンティティ名]_test.go

## 実装制約
- 外部依存禁止（Go標準ライブラリのみ）
- 他層への依存禁止
- 純粋関数・不変オブジェクト中心
- 100%のテストカバレッジ（モック不要）

## 完了条件
- [ ] 全単体テストが通過
- [ ] 純粋ビジネスルールが完全実装
- [ ] エンティティ・値オブジェクトが確定
- [ ] 外部依存なし
- [ ] Lintエラーなし
```

### Sub-Issue B（Usecase Layer）

```markdown
## 目的
[機能名]のユースケース・アプリケーションサービス実装

## 実装対象
- internal/usecase/[ユースケース名].go
- internal/usecase/[ユースケース名]_test.go
- internal/usecase/interfaces/[インターフェース名].go

## 実装制約
- Entitiesのみに依存（Go標準ライブラリ含む）
- Repository/Presenterはインターフェースのみ定義
- Infrastructure/Frameworksへの依存禁止
- 単体テスト（Repository/Presenterモック使用）

## 完了条件
- [ ] 全単体テストが通過
- [ ] アプリケーションルールが完全実装
- [ ] インターフェースが確定
- [ ] ビジネスワークフローが完全実装
- [ ] Lintエラーなし
```

### Sub-Issue C（Interface Adapters Layer）

```markdown
## 目的
[機能名]のインターフェースアダプター実装

## 実装対象
- internal/infrastructure/repository/[リポジトリ名].go
- internal/infrastructure/presenter/[プレゼンター名].go
- internal/interface/controllers/[コントローラー名].go

## 実装制約
- Usecaseのインターフェースを実装
- Usecase Layer実装詳細への依存禁止
- Frameworks Layerへの依存禁止
- 外部ライブラリ使用OK

## 完了条件
- [ ] 全結合テストが通過
- [ ] Repository実装が完全動作
- [ ] Presenter実装が完全動作
- [ ] Controller実装が完全動作
- [ ] Lintエラーなし
```

### Sub-Issue D（Frameworks & Drivers Layer）

```markdown
## 目的
[機能名]のフレームワーク・UI実装

## 実装対象
- cmd/gamebook/[フレームワーク名].go
- cmd/gamebook/[UI名].go
- cmd/gamebook/main.go（DI設定）

## 実装制約
- Interface Adaptersのみに依存
- Infrastructure Layer実装詳細への依存禁止
- PTerm UI一貫性維持
- DIP（依存性逆転の原則）厳守

## 完了条件
- [ ] 全統合テストが通過
- [ ] CLIフレームワークが正常動作
- [ ] PTerm UIが正常動作
- [ ] 依存性注入が正常動作
- [ ] Lintエラーなし
```

### Sub-Issue E（統合検証・リファクタリング）

```markdown
## 目的
[機能名]の全レイヤー統合検証・リファクタリング

## 実装対象
- 全レイヤー統合テスト実行
- レイヤー間の依存関係検証
- コードの重複除去・リファクタリング
- 過渡的コード整理・最適化

## 実装制約
- 全レイヤーのテストが通過した状態で開始
- 破壊的変更は行わない
- 既存テストの修正は最小限に留める
- 「身軽であることが備え」を維持

## 完了条件
- [ ] 全レイヤー統合テストが通過
- [ ] レイヤー間の依存関係が適切
- [ ] コードの重複が除去済み
- [ ] 過渡的コードが整理済み
- [ ] リファクタリングが完了
- [ ] Lintエラーなし
```

## スコープ制御の具体例

```
親Issue: 「入力支援機能実装」
├── Sub-Issue 0: 「InputHelper テストリスト作成」
│   ├── テストリスト.md
│   ├── テストケース設計書
│   └── インターフェース設計書
├── Sub-Issue A: 「InputHelper Entities実装」
│   ├── internal/domain/input_helper.go
│   ├── internal/domain/input_helper_test.go
│   └── 純粋ビジネスロジック
├── Sub-Issue B: 「InputHelper Usecase実装」
│   ├── internal/usecase/input_assistance.go
│   ├── internal/usecase/input_assistance_test.go
│   └── アプリケーションルール
├── Sub-Issue C: 「InputHelper Interface Adapters実装」
│   ├── internal/infrastructure/repository/current_state.go
│   ├── internal/infrastructure/presenter/input_formatter.go
│   └── 結合テスト
├── Sub-Issue D: 「InputHelper Frameworks実装」
│   ├── cmd/gamebook/interactive_pterm.go（修正）
│   └── 統合テスト
└── Sub-Issue E: 「InputHelper 統合検証・リファクタリング」
    ├── 全レイヤー統合テスト
    ├── 依存関係検証
    └── コード最適化
```

## 開発効率化のメリット

1. **トークン消費量の削減**: 対象ディレクトリのみ読み込み
2. **AIの集中力向上**: 特定レイヤーへの集中
3. **実装品質の向上**: 各層の責任範囲内での最適化
4. **並行開発の可能性**: インターフェースで明確に分離