# Issue #107 移動履歴管理機能 - .goファイル設計

## 概要
- **機能**: 移動履歴管理機能（Issue #107バグ修正対応）
- **関連Issue**: #107（参考）、#114（Sub-Issue 0）
- **作成日**: 2025-07-07

## .goファイル設計

### Entities Layer
- `internal/domain/navigation.go` → テストリスト: `entities-navigation.md`
  - NavigationStepエンティティ定義
  - 移動履歴の基本構造と純粋なビジネスロジック

- `internal/domain/gamebook.go`（既存ファイル拡張） → テストリスト: `entities-gamebook.md`
  - NavigationHistoryフィールド追加
  - 移動履歴管理メソッド追加

### Usecase Layer  
- `internal/usecase/navigation_manager.go` → テストリスト: `usecase-navigation_manager.md`
  - 移動履歴記録ユースケース
  - 選択肢移動とジャンプ移動の処理ロジック

### Interface Adapters Layer
- `internal/usecase/interfaces/navigation_repository.go`（新規インターフェース定義）
  - NavigationRepository インターフェース定義
  - 移動履歴の永続化インターフェース

### Frameworks & Drivers Layer
- `cmd/gamebook/tree_printer.go`（既存ファイル修正） → テストリスト: `frameworks-tree_printer.md`
  - 経路表示機能の拡張
  - 現在地ハイライト修正

- `cmd/gamebook/data_converter.go`（既存ファイル修正）
  - 経路連続性を考慮したフローデータ変換

## インターフェース設計

### NavigationStep エンティティ
```go
// NavigationStep は移動履歴の1ステップを表現する不変オブジェクト
type NavigationStep struct {
    From int // 移動元パラグラフ番号（1以上）
    To   int // 移動先パラグラフ番号（1以上、From≠To）
}

// NewNavigationStep は新しいNavigationStepを作成する
func NewNavigationStep(from, to int) (*NavigationStep, error)

// String は文字列表現を返す
func (ns *NavigationStep) String() string

// Equals は等価性を判定する
func (ns *NavigationStep) Equals(other *NavigationStep) bool
```

### Gamebook 拡張
```go
// Gamebook に追加されるフィールドとメソッド
type Gamebook struct {
    // 既存フィールド
    Title      string
    Paragraphs map[int]*Paragraph
    Current    *Paragraph
    
    // 新規追加
    NavigationHistory []NavigationStep // 移動履歴（不変配列）
}

// AddNavigationStep は移動履歴を追加する
func (g *Gamebook) AddNavigationStep(step NavigationStep) error

// GetNavigationHistory は移動履歴を取得する（読み取り専用）
func (g *Gamebook) GetNavigationHistory() []NavigationStep
```

### NavigationManager ユースケース
```go
// NavigationManager は移動履歴管理のユースケースを定義する
type NavigationManager interface {
    // RecordChoiceMove は選択肢による移動を記録する
    RecordChoiceMove(from, to int) error
    
    // RecordJumpMove はジャンプによる移動を記録する
    RecordJumpMove(from, to int) error
    
    // GetNavigationHistory は全移動履歴を取得する
    GetNavigationHistory() ([]NavigationStep, error)
    
    // GetCurrentPath は現在の連続経路を取得する
    GetCurrentPath() ([]int, error)
}
```

### NavigationRepository インターフェース
```go
// NavigationRepository は移動履歴の永続化インターフェース
type NavigationRepository interface {
    // SaveNavigationHistory は移動履歴を保存する
    SaveNavigationHistory(gamebookTitle string, history []NavigationStep) error
    
    // LoadNavigationHistory は移動履歴を読み込む
    LoadNavigationHistory(gamebookTitle string) ([]NavigationStep, error)
}
```

### TreePrinter 拡張
```go
// TreePrinter に追加されるメソッド（既存構造体の拡張）
type TreePrinter struct {
    // 既存フィールド
    // ...
    
    // 新規追加
    navigationHistory []NavigationStep // 表示用移動履歴
}

// PrintTreeWithPath は経路を考慮したツリー表示を行う
func (tp *TreePrinter) PrintTreeWithPath(gamebook *domain.Gamebook) error

// highlightCurrentPath は現在の経路をハイライトする
func (tp *TreePrinter) highlightCurrentPath(history []NavigationStep) error
```

## 実装制約
- NavigationStepは不変オブジェクトとして実装
- 移動履歴は追加のみ可能、削除・変更は禁止
- Entitiesレイヤーは外部依存を持たない純粋なビジネスロジック
- Usecaseレイヤーはビジネスルールに基づいた移動判定を実装
- 既存機能への非破壊的な拡張のみ実施
- 全レイヤーでLogger活用による詳細なデバッグ情報出力
- TDD厳格遵守による品質確保

## レイヤー責務分離
- **Entities**: 移動履歴の純粋なデータ構造とビジネスルール
- **Usecase**: 移動パターンの判定と履歴記録のアプリケーションロジック
- **Interface Adapters**: データ永続化インターフェースの定義
- **Frameworks**: 経路表示の技術的実装とUI制御

## データフロー設計
1. ユーザー操作（選択肢選択/ジャンプ移動）
2. NavigationManager による移動記録
3. Gamebook エンティティへの履歴追加
4. Repository による永続化
5. TreePrinter による経路表示更新

## 完了条件
- [ ] 全インターフェースの設計完了
- [ ] 各レイヤーの責務明確化完了
- [ ] テストリスト作成完了
- [ ] 実装制約の明確化完了
- [ ] Sub-Issue A～E の準備完了