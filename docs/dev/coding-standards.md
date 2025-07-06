# コーディング規約

## 基本方針

- **絵文字使用禁止**: コード・コメント・コミットメッセージで絵文字を使用しない
- **日本語コメント推奨**: 仕様書が日本語のため、コメントも日本語で記載
- 可読性とプロフェッショナル性を重視
- 例外: ユーザーから明示的に要求された場合のみ

## Error Handling（エラーハンドリング）

- **エラーラッピング**: 必ず `%w` を使用してエラーをラップする
- **エラーメッセージ**: 日本語で統一する
- **エラーチェック**: すべてのエラーを適切に処理する

```go
// ✅ 推奨
func (r *MarkdownRepository) Load(title string) (*domain.Gamebook, error) {
    file, err := os.Open(filepath.Join(r.dataDir, title+".md"))
    if err != nil {
        return nil, fmt.Errorf("ファイルオープンに失敗: %w", err)
    }
    defer func() {
        if closeErr := file.Close(); closeErr != nil {
            log.Printf("ファイルクローズエラー: %v", closeErr)
        }
    }()
    // ...
}

// ❌ 非推奨
func (r *MarkdownRepository) Load(title string) (*domain.Gamebook, error) {
    file, err := os.Open(filepath.Join(r.dataDir, title+".md"))
    if err != nil {
        return nil, fmt.Errorf("error: %v", err)  // %vは非推奨
    }
    defer file.Close()  // エラーハンドリングなし
    // ...
}
```

## Variable Naming（変数命名）

- **Variable Shadowing回避**: 同名変数の再宣言を避ける
- **明確な変数名**: 目的が分かる変数名を使用する
- **一時変数**: 明示的な名前を付ける

```go
// ✅ 推奨
func ExecuteShowCommand() error {
    if initErr := ui.Initialize(data); initErr != nil {
        return fmt.Errorf("UI初期化エラー: %w", initErr)
    }
    
    if layoutErr := ui.SetupLayout(); layoutErr != nil {
        return fmt.Errorf("レイアウト設定エラー: %w", layoutErr)
    }
    
    if renderErr := ui.Render(); renderErr != nil {
        return fmt.Errorf("描画エラー: %w", renderErr)
    }
    
    return nil
}

// ❌ 非推奨（shadowエラー）
func ExecuteShowCommand() error {
    if err := ui.Initialize(data); err != nil {
        return fmt.Errorf("UI初期化エラー: %w", err)
    }
    
    if err := ui.SetupLayout(); err != nil {  // shadowエラー
        return fmt.Errorf("レイアウト設定エラー: %w", err)
    }
    
    return nil
}
```

## Function Design（関数設計）

- **関数の長さ**: 50行以内を目標とする
- **複雑度**: 20以下を維持する
- **単一責任の原則**: 1つの関数は1つの責任のみ
- **適切な分割**: 大きな関数は小さなヘルパー関数に分割

```go
// ✅ 推奨
func (r *MarkdownRepository) Load(title string) (*domain.Gamebook, error) {
    file, err := r.openGamebookFile(title)
    if err != nil {
        return nil, fmt.Errorf("ファイルオープンに失敗: %w", err)
    }
    defer r.closeFile(file)
    
    gamebook, err := r.parseGamebook(file, title)
    if err != nil {
        return nil, fmt.Errorf("パース処理に失敗: %w", err)
    }
    
    return gamebook, nil
}
```

## Constants and Magic Numbers（定数とマジックナンバー）

- **マジックナンバー禁止**: 意味のある数値は定数として定義
- **定数名**: 用途が分かる名前を付ける
- **グループ化**: 関連する定数はまとめて定義

```go
// ✅ 推奨
const (
    // パラグラフ番号の上限
    MaxParagraphNumber = 1000
    
    // デフォルトのグリッドサイズ
    DefaultGridSize = 3
    
    // ファイル権限
    FilePermission = 0644
    DirPermission  = 0755
)

// ❌ 非推奨
for i := 0; i < 1000; i++ {  // 1000はマジックナンバー
    // ...
}
```

## Concurrency Safety（並行性安全）

- **共有状態**: 適切な同期機構を使用
- **Mutex使用**: 読み書きが混在する場合は`sync.RWMutex`を使用
- **Goroutine管理**: 適切なライフサイクル管理

```go
// ✅ 推奨
type SafeGamebook struct {
    mu   sync.RWMutex
    data *domain.Gamebook
}

func (sg *SafeGamebook) GetParagraph(num int) (*domain.Paragraph, error) {
    sg.mu.RLock()
    defer sg.mu.RUnlock()
    
    return sg.data.GetParagraph(num)
}

func (sg *SafeGamebook) AddParagraph(p *domain.Paragraph) error {
    sg.mu.Lock()
    defer sg.mu.Unlock()
    
    return sg.data.AddParagraph(p)
}
```

## コードレビュー観点

### 必須チェック項目

- [ ] **Variable Shadowing**: 特に`err`変数の再宣言
- [ ] **Error Wrapping**: `fmt.Errorf("message: %w", err)`の使用
- [ ] **Resource Management**: `defer`を使った適切なリソース解放
- [ ] **Function Length**: 50行を超える関数の分割
- [ ] **Magic Numbers**: 定数化されているか
- [ ] **Concurrency**: 共有状態の適切な保護
- [ ] **Test Coverage**: 新機能のテストが追加されているか

### よくある問題と解決策

#### Variable Shadowing

```go
// 問題: shadowエラー
func example() error {
    if err := step1(); err != nil {
        return err
    }
    if err := step2(); err != nil {  // shadowエラー
        return err
    }
    return nil
}

// 解決: 明示的な変数名
func example() error {
    if step1Err := step1(); step1Err != nil {
        return fmt.Errorf("step1エラー: %w", step1Err)
    }
    if step2Err := step2(); step2Err != nil {
        return fmt.Errorf("step2エラー: %w", step2Err)
    }
    return nil
}
```

#### Error Wrapping

```go
// 問題: エラーラッピング不適切
return fmt.Errorf("error: %v", err)

// 解決: 適切なエラーラッピング
return fmt.Errorf("処理に失敗: %w", err)
```

#### Resource Management

```go
// 問題: エラーハンドリングなし
defer file.Close()

// 解決: 適切なエラーハンドリング
defer func() {
    if closeErr := file.Close(); closeErr != nil {
        log.Printf("ファイルクローズエラー: %v", closeErr)
    }
}()
```