# Issue管理

## 基本方針

- **すべての作業はGitHub Issueとして管理**
- **Issueの冒頭にSPECIFICATION.md・DEVELOPMENT.mdの確認を必須**
- **IssueのCloseはPRで行う**

## Issue作成

### テンプレート使用必須

1. **Feature実装**: `gh issue create --template ai_feature.md`
   - 新機能開発用
   - TDD方針、技術仕様、AI実装指示を含む

2. **Sub-Issue**: `gh issue create --template ai_sub_issue.md`
   - 複雑なFeatureの分割用
   - 親Issue連携、実装範囲明確化

3. **Bug Report**: `gh issue create --template ai_bug_report.md`
   - バグ修正用
   - 再現手順、AI修正指示を含む

### テンプレート遵守の重要性

- **AI（Claude Code）がプロジェクト内容を効率的に理解**
- **実装品質の統一性確保**
- **レビュー効率の向上**
- **プロジェクト管理の自動化促進**

## GitHub公式Sub-Issue管理

### 重要: gh CLIの現状
- **gh CLIはsub-issuesをネイティブサポートしていません**（2025年7月時点）
- GitHub CLI Issue #10298で機能追加が要求されている状況
- **現在はGraphQL API経由での操作が必要**

### GraphQL API操作方法

```bash
# Step 1: 親・子IssueのGraphQL IDを取得
gh api graphql -f query='
{
  repository(owner: "owner", name: "repo") {
    issue(number: 親Issue番号) {
      title
      id
    }
  }
}' --header "GraphQL-Features: sub_issues"

# Step 2: Sub-Issue追加
gh api graphql -f query='
mutation {
  addSubIssue(input: {
    issueId: "親IssueのGraphQL_ID"
    subIssueId: "子IssueのGraphQL_ID"
  }) {
    issue {
      title
    }
  }
}' --header "GraphQL-Features: sub_issues"

# Step 3: Sub-Issues一覧確認（REST API）
gh api -X GET /repos/{owner}/{repo}/issues/{issue_number}/sub_issues
```

### 実際の操作例（動作確認済み）

```bash
# Issue #50に Issue #51をsub-issueとして追加する例

# 1. GraphQL IDを取得
gh api graphql -f query='
{
  repository(owner: "zensai3805", name: "iwakero-gamebook-mapping") {
    issue(number: 50) {
      id
    }
  }
}' --header "GraphQL-Features: sub_issues"
# 結果: "I_kwDOPA7qR86_EALR"

# 2. Sub-Issue追加
gh api graphql -f query='
mutation {
  addSubIssue(input: {
    issueId: "I_kwDOPA7qR86_EALR"
    subIssueId: "I_kwDOPA7qR86_EAMX"
  }) {
    issue {
      title
    }
  }
}' --header "GraphQL-Features: sub_issues"
```

## ラベル体系

- **priority**: high/medium/low
- **type**: epic（マイルストーン追跡）
- **area**: cli/domain/repository（技術領域）