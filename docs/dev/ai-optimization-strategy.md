# AI開発最適化：スコープ制御戦略

## 課題：コード全体改修によるトークン消費問題

- **現状の問題**: 実装時にコード全体が改修範囲に含まれる
- **影響**: トークン消費量過多、AIの関心範囲の拡散
- **解決策**: cmd/internal分離によるSub-Issue戦略

## クリーンアーキテクチャ4層分離

### 1. 機能実装時の必須分離

**すべての機能実装を以下の4つのSub-Issueに分離**

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

### 2. 実装順序の原則

**必ずEntities → Usecase → Interface Adapters → Frameworks の順序で実装**

## トークン消費量最適化

### 読み込み対象の限定

| Sub-Issue | 読み込み対象 | 避けるべきディレクトリ |
|-----------|-------------|-------------------|
| A (Entities) | `internal/domain/` のみ | usecase/, infrastructure/, cmd/ |
| B (Usecase) | `internal/usecase/` のみ | infrastructure/, cmd/ の詳細 |
| C (Interface Adapters) | `internal/infrastructure/`, `internal/interface/` | cmd/ の詳細 |
| D (Frameworks) | `cmd/gamebook/` のみ | infrastructure/ の実装詳細 |

### AIプロンプトの最適化

- **Sub-Issue A**: 「Domainディレクトリのみに集中、純粋ビジネスロジック、外部依存なし」
- **Sub-Issue B**: 「Usecaseディレクトリのみに集中、アプリケーションルール、インターフェース定義」
- **Sub-Issue C**: 「Infrastructure・Interfaceディレクトリのみに集中、Usecaseインターフェース実装」
- **Sub-Issue D**: 「Frameworksディレクトリのみに集中、UI・依存性注入、Interface Adapters使用」

## Sub-Issue作成テンプレート

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

## スコープ制御の具体例

```
親Issue: 「入力支援機能実装」
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
└── Sub-Issue D: 「InputHelper Frameworks実装」
    ├── cmd/gamebook/interactive_pterm.go（修正）
    └── 統合テスト
```

## 開発効率化のメリット

1. **トークン消費量の削減**: 対象ディレクトリのみ読み込み
2. **AIの集中力向上**: 特定レイヤーへの集中
3. **実装品質の向上**: 各層の責任範囲内での最適化
4. **並行開発の可能性**: インターフェースで明確に分離