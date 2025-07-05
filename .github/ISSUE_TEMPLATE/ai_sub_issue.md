---
name: AI Sub-Issue Implementation
about: AI最適化されたSub-Issue実装用
title: '[SUB-ISSUE] '
labels: 'type:sub-issue, priority:medium'
assignees: ''
---

## 親のIssue
<!-- このSub-Issueの親のIssue番号 -->
#123

## 概要
<!-- このSub-Issueで実装する具体的な機能・タスク -->
親のIssueの中で、このSub-Issueが担当する具体的な範囲を記載

## 実装範囲
<!-- このSub-Issueで実装する具体的な範囲 -->

### 実装対象ファイル
- `cmd/gamebook/example.go`
- `internal/domain/example.go`
- `internal/domain/example_test.go`

### 実装対象機能
- 関数: `ExampleFunction()`
- 構造体: `ExampleStruct`
- インターフェース: `ExampleInterface`

## 実装タスク
- [ ] テストケース設計
- [ ] テストコード実装（RED）
- [ ] 最小限の実装（GREEN）
- [ ] リファクタリング（REFACTOR）
- [ ] エラーハンドリング実装
- [ ] 日本語コメント追加
- [ ] 統合テスト実行
- [ ] 品質チェック実行

## AI実装指示
<!-- Claude Code向けの具体的な実装指示 -->

### 実装手順
1. テストファイル作成（`*_test.go`）
2. 失敗するテストケース実装（RED）
3. 最小限の実装で通るようにする（GREEN）
4. コードをリファクタリング（REFACTOR）
5. エラーハンドリング追加
6. 日本語コメント追加
7. 統合テスト実行

### コーディング規約
- 日本語コメント必須
- エラーは`fmt.Errorf("説明: %w", err)`でラップ
- 変数shadowing回避のため明示的命名
- 関数は50行以内、複雑度20以下

### 品質基準
- `go test ./...` でエラーなし
- `golangci-lint run` でエラーなし
- カバレッジ100%（新規コード）