# navigation_controller.go テストリスト

## 対象ファイル
- internal/interface/controllers/navigation_controller.go
- internal/interface/controllers/navigation_controller_test.go

## TODOリストとしての期待する振る舞い

### NavigationController実装の基本動作
- [x] NavigationControllerクラスが適切に定義される
- [x] NewNavigationControllerコンストラクタが適切に動作する
- [x] 必要な依存関係（repository, presenter）が適切に管理される
- [x] 基本的なメソッドが実装される

### SaveHistory メソッドの基本動作
- [x] 有効なgamebookTitleと履歴でSaveHistoryが成功する
- [ ] RepositoryのSaveNavigationHistoryが適切に呼び出される
- [ ] エラーが発生した場合に適切にエラーが返される

### LoadHistory メソッドの基本動作
- [x] 有効なgamebookTitleでLoadHistoryが成功する
- [ ] RepositoryのLoadNavigationHistoryが適切に呼び出される
- [ ] 読み込んだ履歴が正しく返される
- [ ] エラーが発生した場合に適切にエラーが返される

### FormatHistory メソッドの基本動作
- [x] 有効なgamebookTitleでFormatHistoryが成功する
- [ ] Repository → Presenterの連携が適切に動作する
- [ ] フォーマット結果が正しく返される
- [ ] エラーが発生した場合に適切にエラーが返される

## 実装制約
- Interface Adaptersレイヤーのコントローラー機能
- Repository と Presenter の組み合わせ使用
- エラーラッピング（%w）の使用
- テスト可能な設計
- 振る舞いのみ記載、実装詳細は記載しない
- 1つのテストは1つの振る舞いのみ検証

## TDD進行管理
- [x] 全振る舞いのテスト実装完了（基本機能）
- [x] 全テストが通過
- [x] リファクタリング完了
- [ ] Lintエラーなし