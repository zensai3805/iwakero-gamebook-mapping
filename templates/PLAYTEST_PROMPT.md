# プレイテスト実行プロンプト集

## 🎮 基本プロンプト

```
iwakero-gamebook-mappingツールの実際のプレイテストを依頼します。

templates/README.mdの指針に従い、自然なゲームブックプレイを体験して、templates/PLAY_FEEDBACK_TEMPLATE.mdの形式でフィードバックを作成してください。

**出力**: docs/feedback/PLAY_FEEDBACK_YYYY-MM-DD.md として保存

重要: 単なる機能テストではなく、具体的な冒険ストーリーを作りながら、想定外の使用パターンも積極的に探求してください。
```

## 🎯 用途別バリエーション

### 新機能フォーカス
```
iwakero-gamebook-mappingツールの[新機能名]に特化したプレイテストを依頼します。

[新機能名]を重点的に検証しながら、templates/の指針に従って自然なゲームブックプレイを体験し、フィードバックを作成してください。

**出力**: docs/feedback/PLAY_FEEDBACK_v[バージョン]_[新機能名].md として保存
```

### 比較検証
```
iwakero-gamebook-mappingツールと従来の記録方法（紙・メモアプリ等）との比較を含むプレイテストを依頼します。

相対的な価値を評価しながら、templates/の指針に従ってフィードバックを作成してください。

**出力**: docs/feedback/PLAY_FEEDBACK_YYYY-MM-DD_comparison.md として保存
```

### 初心者体験
```
iwakero-gamebook-mappingツールを初めて使用するプレイヤーの体験に特化したプレイテストを依頼します。

初心者視点（理解しやすさ、学習コスト等）を重視して、templates/の指針に従ってフィードバックを作成してください。

**出力**: docs/feedback/PLAY_FEEDBACK_YYYY-MM-DD_beginner.md として保存
```

---

## 📋 使用方法

1. **適切なプロンプトを選択**してコピー
2. **AI（Claude Code等）に送信**
3. **AIが自動的に `docs/feedback/` に適切なファイル名で保存**
4. **結果確認**: 作成されたフィードバックファイルをレビュー

**注意**: プロンプトには保存先が明記されているため、AIは指定されたディレクトリとファイル名で出力します。

詳細はtemplates/README.mdを参照。