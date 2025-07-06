# ログシステム使用ガイド

## 概要

iwakero-gamebook-mappingのログシステムは、AI開発効率化を目的とした構造化ログシステムです。
Clean Architecture 4層（Entities, Usecase, Interface Adapters, Frameworks）全体で統一されたログ出力を提供します。

## 🎯 目的

- **AI開発効率化**: 構造化されたログで開発状況を把握
- **デバッグ支援**: エラーの詳細なコンテキスト情報
- **操作追跡**: ユーザー操作の記録と分析
- **品質向上**: バリデーション結果と設定変更の記録

## 🏗️ アーキテクチャ

### 4層設計
```
┌─────────────────────────────────────┐
│ Frameworks Layer (cmd/gamebook/)    │ ← ログシステム統合
├─────────────────────────────────────┤
│ Interface Adapters Layer            │ ← Repository/Presenter/Controller
├─────────────────────────────────────┤
│ Usecase Layer                       │ ← アプリケーションロジック
├─────────────────────────────────────┤
│ Entities Layer (internal/domain/)   │ ← Loggerインターフェース
└─────────────────────────────────────┘
```

### 主要コンポーネント
- **domain.Logger**: ログインターフェース定義
- **LoggingController**: 設定管理とログライター制御
- **StructuredLogger**: AI開発最適化されたログフォーマット
- **OperationLogger**: 高レベル操作ログ関数

## ✅ 実装状況

### Issue #92完了: CLI統一実装
**2025-07-07時点**: すべてのCLIワンショットコマンドでログ出力が正常動作します。

#### 対応済みコマンド
- `./gamebook new "タイトル"` - ゲームブック作成
- `./gamebook add 1 "説明"` - パラグラフ追加
- `./gamebook choice 1 "選択肢" 2` - 選択肢追加
- `./gamebook select 1 1` - 選択肢選択
- `./gamebook move 2` - パラグラフ移動
- `./gamebook load "タイトル"` - ゲームブック読み込み
- `./gamebook list` - ゲームブック一覧
- `./gamebook show` - 状態表示

#### CLIログ出力確認
```bash
export LOG_LEVEL=DEBUG
export LOG_OUTPUT=stderr  
export GAMEBOOK_AI_DEV=true
./gamebook new "テストゲーム" 2>&1 | grep -E "(DEBUG|INFO|WARN|ERROR)"
```

期待される出力例:
```
2025-07-07T04:44:25+09:00 [INFO] [operation.new_gamebook] ゲームブック作成 | title="テストゲーム" | commands.go:85
```

## 🚀 基本使用方法

### 1. 環境変数設定

```bash
# ログレベル設定
export LOG_LEVEL=DEBUG    # DEBUG, INFO, WARN, ERROR, FATAL

# 出力先設定
export LOG_OUTPUT=console # console, stderr, file

# フォーマット設定
export LOG_FORMAT=text    # text, json

# AI開発モード（自動DEBUGレベル、詳細出力）
export GAMEBOOK_AI_DEV=true

# ファイル出力時のパス
export LOG_FILE_PATH=/path/to/logfile.log
```

### 2. YAML設定ファイル

`~/.gamebook/logger.yaml`:
```yaml
level: DEBUG
format: text
output: console
file_path: ""
```

### 3. コード内での使用

#### 基本ログ出力
```go
// Frameworks Layer
logger := GetGlobalLogger()
logger.Info("ゲームブック作成", domain.Field{Key: "title", Value: "冒険の書"})
logger.Error("ファイル読み込みエラー", domain.Field{Key: "path", Value: "/data/game.md"})
```

#### 操作ログ（推奨）
```go
// ユーザー操作記録
LogUserOperation("new_gamebook", map[string]interface{}{
    "title": "新しいゲーム",
    "title_length": 12,
})

// バリデーションエラー記録
LogValidationError("title", title, "空のタイトル", map[string]interface{}{
    "command": "new",
})

// コマンド結果記録
LogCommandResult("new", true, "ゲームブック作成成功", map[string]interface{}{
    "file_path": "./data/game.md",
})

// UI操作記録（AI開発モード時のみ）
LogUIInteraction("menu_selection", map[string]interface{}{
    "selected_option": "新しいゲーム",
    "selection_time_ms": 245.67,
})
```

## 🎨 ログフォーマット

### テキストフォーマット（AI開発最適化）
```
2025-07-06T19:30:15+09:00 [INFO] [operation.new_gamebook] ゲームブック作成 | title="冒険の書" duration=15.23ms | gamebook.go:45 | session: current_game=true paragraphs=0
```

### JSONフォーマット
```json
{
  "timestamp": "2025-07-06T19:30:15+09:00",
  "level": "INFO",
  "category": "operation",
  "action": "new_gamebook",
  "component": "gamebook",
  "message": "ゲームブック作成",
  "context": {"title": "冒険の書"},
  "performance": {"duration_ms": 15.23, "success": true},
  "location": {"file": "gamebook.go", "line": 45, "function": "ExecuteNewCommand"},
  "session": {"has_current_game": true, "total_paragraphs": 0}
}
```

## 🤖 AI開発モード

### 有効化
```bash
export GAMEBOOK_AI_DEV=true
```

### 特徴
- **自動DEBUGレベル**: `LOG_LEVEL`を自動的にDEBUGに昇格
- **詳細出力**: パフォーマンス情報、ソースコード位置、セッション状態
- **UI操作ログ**: メニュー選択、入力操作の詳細記録
- **構造化フォーマット**: AI解析に最適化された出力

### AI開発モード時の追加情報
- **パフォーマンス**: 処理時間、成功/失敗ステータス
- **ソースコード位置**: ファイル名、行番号、関数名
- **セッション状態**: 現在のゲーム状態、パラグラフ数
- **操作コンテキスト**: 入力値、選択肢、タイミング情報

## 📂 設定ファイル優先順位

1. **環境変数** (最優先)
2. **YAML設定ファイル** (`~/.gamebook/logger.yaml`)
3. **デフォルト設定** (INFO, text, console)

## 🔧 トラブルシューティング

### ログが出力されない
```bash
# 設定確認
env | grep -E "(LOG_|GAMEBOOK_)"

# デバッグモード有効化
export LOG_LEVEL=DEBUG
export GAMEBOOK_AI_DEV=true

# stderr出力に変更
export LOG_OUTPUT=stderr
```

### パフォーマンス問題
```bash
# ログレベルを上げる
export LOG_LEVEL=WARN

# AI開発モードを無効化
unset GAMEBOOK_AI_DEV
```

## 🎯 ベストプラクティス

### 1. 適切なレベル選択
- **DEBUG**: 開発時の詳細追跡
- **INFO**: 通常操作の記録
- **WARN**: 注意が必要な状況
- **ERROR**: エラー発生時
- **FATAL**: アプリケーション終了レベル

### 2. 構造化データ活用
```go
// 良い例: 構造化されたフィールド
LogUserOperation("add_paragraph", map[string]interface{}{
    "paragraph_id": 42,
    "description_length": len(description),
    "has_choices": hasChoices,
})

// 悪い例: 文字列内に情報を埋め込み
logger.Info(fmt.Sprintf("パラグラフ%d追加: %s", id, description))
```

### 3. AI開発モード活用
- 開発中: `GAMEBOOK_AI_DEV=true`で詳細ログ
- 本番: `GAMEBOOK_AI_DEV=false`で最小限ログ

## 🔗 関連ファイル

- `cmd/gamebook/logger_setup.go` - 初期化
- `cmd/gamebook/logger_config.go` - 設定管理  
- `cmd/gamebook/structured_logger.go` - フォーマット
- `cmd/gamebook/operation_logger.go` - 操作ログ関数
- `internal/domain/logger.go` - インターフェース定義

## 📋 設定例

### 開発環境
```bash
export LOG_LEVEL=DEBUG
export LOG_OUTPUT=stderr
export LOG_FORMAT=text
export GAMEBOOK_AI_DEV=true
```

### 本番環境
```bash
export LOG_LEVEL=WARN
export LOG_OUTPUT=file
export LOG_FORMAT=json
export LOG_FILE_PATH=/var/log/gamebook.log
```

### テスト環境
```bash
export LOG_LEVEL=ERROR
export LOG_OUTPUT=console
export LOG_FORMAT=text
```