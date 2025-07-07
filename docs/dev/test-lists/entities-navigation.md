# navigation.go テストリスト

## 対象ファイル
- internal/domain/navigation.go
- internal/domain/navigation_test.go

## TODOリストとしての期待する振る舞い

### NavigationStep基本動作
- [ ] NavigationStepが正常に作成できる
- [ ] From、Toフィールドが適切に設定される
- [ ] 作成後にFrom、Toフィールドが正常に取得できる
- [ ] NavigationStepの複数作成が正常に動作する

### NavigationStepのバリデーション
- [ ] 負の番号でFromが設定された場合エラーが発生する
- [ ] 負の番号でToが設定された場合エラーが発生する
- [ ] 0でFromが設定された場合エラーが発生する
- [ ] 0でToが設定された場合エラーが発生する
- [ ] FromとToが同じ値の場合エラーが発生する

### NavigationStepの等価性比較
- [ ] 同じFrom、Toを持つNavigationStepは等価と判定される
- [ ] Fromが異なるNavigationStepは等価でないと判定される
- [ ] Toが異なるNavigationStepは等価でないと判定される
- [ ] nilとの比較で適切に処理される

### NavigationStepの文字列化
- [ ] NavigationStepが適切な文字列形式で表現される
- [ ] From、To情報が文字列に含まれる
- [ ] 複数のNavigationStepで一意な文字列表現になる

### 境界条件
- [ ] 最小値（1）でFrom、Toが正常に動作する
- [ ] 最大値（int最大値）でFrom、Toが正常に動作する
- [ ] 最小値-1（0）でエラーが発生する
- [ ] 大きな値での処理が正常に動作する

### 「もしも」の状況
- [ ] メモリ不足状況でのNavigationStep作成が適切に処理される
- [ ] 非常に大量のNavigationStep作成時のパフォーマンスが適切である
- [ ] Garbage Collection発生時の動作が安定している
- [ ] 並行アクセス時の安全性が確保される

### 既存機能との関係
- [ ] NavigationStepが既存のParagraphエンティティを破壊しない
- [ ] NavigationStepが既存のChoiceエンティティを破壊しない
- [ ] NavigationStepが既存のGamebookエンティティと適切に連携する
- [ ] 既存のバリデーションルールが維持される
- [ ] 既存のデータ整合性が保たれる

## 実装制約
- Entitiesレイヤーでは外部依存を持たない純粋なビジネスロジックのみ実装
- 振る舞いのみ記載、実装詳細は記載しない
- 1つのテストは1つの振る舞いのみ検証
- NavigationStepは不変オブジェクトとして設計
- ドメインモデルの純粋性を保持

## TDD進行管理
- [ ] 全振る舞いのテスト実装完了
- [ ] 全テストが通過
- [ ] リファクタリング完了
- [ ] Lintエラーなし