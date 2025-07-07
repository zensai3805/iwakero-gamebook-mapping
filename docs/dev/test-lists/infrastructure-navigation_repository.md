# navigation_repository.go テストリスト

## 対象ファイル
- internal/infrastructure/repository/navigation_repository.go
- internal/infrastructure/repository/navigation_repository_test.go

## TODOリストとしての期待する振る舞い

### NavigationRepository実装の基本動作
- [x] NavigationRepositoryインターフェースを実装する構造体が定義される
- [x] NewNavigationRepositoryコンストラクタが適切に動作する
- [ ] 必要な依存関係（dataDir等）が適切に管理される
- [ ] 実装がインターフェースを満たしていることが確認される

### SaveNavigationHistory メソッドの基本動作
- [x] 有効なgamebookTitleと履歴でSaveNavigationHistoryが成功する
- [ ] 保存されたデータが適切なMarkdown形式である
- [ ] 既存の履歴データが上書きされる
- [ ] 空の履歴でもSaveNavigationHistoryが成功する
- [ ] 複数のNavigationStepを含む履歴が正しく保存される

### LoadNavigationHistory メソッドの基本動作
- [x] 保存された履歴がLoadNavigationHistoryで正しく読み込める
- [x] 存在しないゲームブックでLoadNavigationHistoryが空配列を返す
- [ ] 空の履歴ファイルでLoadNavigationHistoryが空配列を返す
- [ ] 複数のNavigationStepを含む履歴が正しく読み込める
- [ ] 読み込んだデータがdomain.NavigationStep型に正しく変換される

### エラー処理
- [x] ファイル書き込み時のI/Oエラーが適切に処理される
- [x] ループ内の書き込みエラーが適切に伝播される
- [ ] 空文字列のgamebookTitleでSaveNavigationHistoryがエラーを返す
- [ ] nilの履歴配列でSaveNavigationHistoryがエラーを返す
- [ ] 空文字列のgamebookTitleでLoadNavigationHistoryがエラーを返す
- [ ] ファイル作成権限がない場合にSaveNavigationHistoryがエラーを返す
- [ ] 破損したファイルでLoadNavigationHistoryがエラーを返す
- [ ] 無効なMarkdown形式でLoadNavigationHistoryがエラーを返す

### ファイルシステム操作
- [ ] 適切なディレクトリが存在しない場合に自動作成される
- [ ] ファイル名が適切にエスケープされる
- [ ] 特殊文字を含むタイトルでも正しく動作する
- [ ] パスの区切り文字が適切に処理される
- [ ] ファイルロックが適切に処理される

### Markdown形式の永続化
- [ ] NavigationStepがMarkdown形式で適切に保存される
- [ ] From、To、ViaPathsが正しく保存される
- [ ] 空のViaPathsが適切に処理される
- [ ] ViaPathsに複数の値がある場合に正しく保存される
- [ ] 保存されたMarkdownが人間が読める形式である

### データ整合性
- [ ] 保存と読み込みのデータが完全に一致する
- [ ] NavigationStepの各フィールドが正確に復元される
- [ ] 履歴の順序が保持される
- [ ] 大量のNavigationStepでもデータ整合性が保たれる
- [ ] エンコーディングの問題が発生しない

### 境界条件
- [ ] 最大値のパラグラフ番号でも正しく動作する
- [ ] 最小値（1）のパラグラフ番号でも正しく動作する
- [ ] 非常に長いゲームブックタイトルでも動作する
- [ ] 1つだけのNavigationStepでも正しく動作する
- [ ] 大量のNavigationStep（1000個以上）でも動作する

### 「もしも」の状況
- [ ] ディスク容量不足時にSaveNavigationHistoryが適切にエラーを返す
- [ ] ファイルが他のプロセスによって使用中の場合に適切に処理される
- [ ] システム再起動後でも保存されたデータが読み込める
- [ ] ファイルの一部が破損している場合に適切にエラーを返す
- [ ] 権限変更によりファイルアクセスできない場合にエラーを返す

### 既存機能との関係
- [ ] 既存のMarkdownRepositoryパターンと整合性が保たれる
- [ ] 既存のデータディレクトリ構造と競合しない
- [ ] 既存のゲームブックファイルに影響を与えない
- [ ] 既存のエラーハンドリングパターンと統一される
- [ ] 既存のLogger活用パターンに従っている

### パフォーマンス
- [ ] 大量データの保存が合理的な時間で完了する
- [ ] 大量データの読み込みが合理的な時間で完了する
- [ ] メモリ使用量が適切に管理される
- [ ] ファイルハンドルが適切に管理される
- [ ] 並行アクセス時のパフォーマンス低下が最小限である

### ログ出力（Logger活用）
- [ ] ファイル保存操作がDEBUGレベルでログ出力される
- [ ] ファイル読み込み操作がDEBUGレベルでログ出力される
- [ ] エラー発生時にERRORレベルでログ出力される
- [ ] 成功操作がINFOレベルでログ出力される
- [ ] AI開発モード時に詳細なコンテキスト情報が出力される

## 実装制約
- Infrastructure層として外部依存（ファイルシステム）を扱う
- NavigationRepositoryインターフェースの完全実装
- 既存のMarkdownRepositoryパターンに準拠
- Logger活用による適切なログ出力
- エラーラッピング（%w）の使用
- テスト可能な設計
- 振る舞いのみ記載、実装詳細は記載しない
- 1つのテストは1つの振る舞いのみ検証

## TDD進行管理
- [x] 全振る舞いのテスト実装完了（基本機能）
- [x] 全テストが通過
- [x] リファクタリング完了
- [ ] Lintエラーなし