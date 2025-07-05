# Pull Request Template

## Summary
<!-- このPRで実装した内容を簡潔に説明 -->

## Related Issue
<!-- 関連するIssueを記載 -->
Closes #

## Changes Made
<!-- 実装した変更内容を詳細に記載 -->

### 🔧 Modified Files
- `path/to/file.go` - 変更内容の説明
- `path/to/test.go` - テストケースの追加/修正

### ➕ Added Files
- `path/to/new_file.go` - 新規追加ファイルの説明

### 🗑️ Removed Files
- `path/to/old_file.go` - 削除理由

## Technical Details
<!-- 技術的な詳細を記載 -->

### Architecture Changes
<!-- アーキテクチャの変更がある場合 -->
- 

### API Changes
<!-- APIの変更がある場合 -->
- 

### Database Changes
<!-- データベースの変更がある場合 -->
- 

## Test Plan
<!-- テスト計画を記載 -->

### Unit Tests
- [ ] 新機能の単体テスト追加
- [ ] 既存機能の単体テスト修正
- [ ] エラーケースのテスト追加

### Integration Tests
- [ ] 統合テストの実行
- [ ] 主要ユースケースの動作確認
- [ ] 後方互換性の確認

### Manual Testing
- [ ] インタラクティブモードでの動作確認
- [ ] CLIモードでの動作確認
- [ ] エラーハンドリングの確認

## Quality Checklist
<!-- 品質チェックリスト -->

### Code Quality
- [ ] **テスト通過**: `go test ./...` でエラーなし
- [ ] **Lint通過**: `golangci-lint run` でエラーなし
- [ ] **フォーマット**: `gofmt -s -w .` 実行済み
- [ ] **Import整理**: `goimports -w .` 実行済み
- [ ] **カバレッジ**: 新規コードのカバレッジ70%以上

### Code Standards
- [ ] **Variable Shadowing**: 変数shadowing回避（特にerr変数）
- [ ] **Error Handling**: エラーラッピングに`%w`を使用
- [ ] **Function Length**: 関数の長さ50行以内
- [ ] **Complexity**: 複雑度20以下
- [ ] **Comments**: 日本語コメント追加
- [ ] **No Emojis**: 絵文字使用禁止

### Documentation
- [ ] **CLAUDE.md更新**: 変更内容を反映
- [ ] **コメント更新**: 必要に応じてコメント追加/修正
- [ ] **README更新**: 必要に応じてREADME更新

## Breaking Changes
<!-- 破壊的変更がある場合 -->
- [ ] 破壊的変更なし
- [ ] 破壊的変更あり（下記に詳細記載）

### Breaking Change Details
<!-- 破壊的変更がある場合の詳細 -->

## Migration Guide
<!-- マイグレーションが必要な場合 -->

## Performance Impact
<!-- パフォーマンスへの影響 -->
- [ ] パフォーマンスへの影響なし
- [ ] パフォーマンス向上
- [ ] パフォーマンス低下（理由と対策を記載）

## Deployment Notes
<!-- デプロイ時の注意事項 -->

## Screenshots/Demos
<!-- スクリーンショットやデモがある場合 -->

## AI Implementation Notes
<!-- AI実装時の特記事項 -->

### TDD Process
- [ ] RED: 失敗するテストを実装
- [ ] GREEN: 最小限の実装で通すように修正
- [ ] REFACTOR: コードをリファクタリング

### AI Specific Considerations
- 実装時に注意した点
- 参考にした既存コード
- 今後の改善点

## Review Checklist for Reviewers
<!-- レビュアー向けチェックリスト -->

### Code Review
- [ ] コードの可読性・保守性
- [ ] エラーハンドリングの適切性
- [ ] テストの網羅性
- [ ] パフォーマンスの考慮
- [ ] セキュリティの考慮

### Functional Review
- [ ] 仕様通りの動作
- [ ] エッジケースの処理
- [ ] 既存機能への影響
- [ ] ユーザビリティ

## Additional Notes
<!-- その他の補足事項 -->

---

**🤖 Generated with [Claude Code](https://claude.ai/code)**

**Co-Authored-By: Claude <noreply@anthropic.com>**