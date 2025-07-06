# Frameworks & Drivers Layer 仕様書

## 概要
フレームワーク統合、UI実装、外部ツールとの連携仕様

## CLI Framework

### cobra
コマンドラインインターフェースフレームワーク

#### コマンド構造
```
gamebook
├── new <title>              # 新規作成
├── add <number> <desc>      # パラグラフ追加
├── choice <number> <desc> <target>  # 選択肢追加
├── select <para> <choice>   # 選択実行
├── move <number>            # 直接移動
├── show                     # 状態表示
├── list                     # 一覧表示
└── switch <title>           # 切り替え
```

## UI Framework

### PTerm
リッチターミナルUI

#### コンポーネント
- **InteractiveSelect**: メニュー選択
- **TextInput**: テキスト入力
- **Table**: 統計情報表示
- **Tree**: 階層表示
- **Panel**: 情報パネル
- **Spinner**: 処理中表示

#### 画面レイアウト
```
┌─────────────────────────────────┐
│      ゲームブック名              │
├─────────────────────────────────┤
│ 統計情報:                       │
│ - 総パラグラフ数: 10            │
│ - 訪問済み: 5                   │
├─────────────────────────────────┤
│ 現在位置: パラグラフ 5          │
│ 概要: 暗い森                    │
├─────────────────────────────────┤
│ メニュー:                       │
│ > 1. パラグラフを追加           │
│   2. 選択肢を追加               │
│   3. 選択肢を選択               │
└─────────────────────────────────┘
```

### 統合UI
3分割レイアウトによる可視化

#### レイアウト構成
- **左側**: 2Dマップ表示
- **右側**: フロー図表示
- **下部**: 操作メニュー

## 外部ツール連携

### Git
バージョン管理

#### 自動生成される.gitignore
```
# ビルド成果物
gamebook
gamebook.exe
*.exe

# テスト関連
*.test
*.out
coverage.html

# エディタ
.vscode/
.idea/
```

### Make
ビルドタスク管理

#### Makefile
```makefile
build:
    go build -o gamebook ./cmd/gamebook

test:
    go test ./...

lint:
    golangci-lint run

check: test lint
```

## プラットフォーム対応

### クロスプラットフォーム
- Windows: 画面クリアコマンド調整
- macOS: 標準対応
- Linux: 標準対応

### ターミナル互換性
- UTF-8サポート必須
- 256色対応推奨
- 最小80x24文字

## 依存性注入

### main関数の構成
```go
1. Repository初期化
2. Usecase初期化
3. Presenter初期化
4. Controller初期化
5. アプリケーション起動
```

### エラーハンドリング
- パニックリカバリー
- グレースフルシャットダウン
- エラーログ出力

## 将来の拡張

### WebUI（v0.4.0）
- HTTPサーバー
- WebSocket通信
- RESTful API

### 音声入力（v1.0.0）
- 音声認識エンジン統合
- リアルタイム変換
- コマンド認識