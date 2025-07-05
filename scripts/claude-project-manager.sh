#!/bin/bash

# Claude Code専用プロジェクト管理スクリプト
# AI開発効率化のためのGitHub CLI活用ツール

set -euo pipefail

# 設定値
PROJECT_NAME="iwakero-gamebook-mapping"
REPO_OWNER="iwapc"
DEFAULT_BRANCH="main"

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ログ関数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# GitHub CLI が利用可能かチェック
check_gh_cli() {
    if ! command -v gh &> /dev/null; then
        log_error "GitHub CLI (gh) がインストールされていません"
        log_info "インストール方法: https://cli.github.com/"
        exit 1
    fi
    
    if ! gh auth status &> /dev/null; then
        log_error "GitHub CLI認証が必要です"
        log_info "認証コマンド: gh auth login"
        exit 1
    fi
}

# Claude専用: Feature Issue作成
claude_create_feature_issue() {
    local title="$1"
    local description="$2"
    local milestone="${3:-}"
    
    log_info "Feature Issue作成中: $title"
    
    # Issue本文テンプレート
    local body=$(cat << EOF
## 概要
$description

## 主要タスク
- [ ] SPECIFICATION.md確認
- [ ] DEVELOPMENT.md確認
- [ ] 技術設計
- [ ] 実装
- [ ] テスト作成
- [ ] テスト実行
- [ ] CLAUDE.md更新

## 技術詳細
<!-- 実装に必要な技術詳細を記載 -->

## 依存関係
<!-- 他のIssueやタスクとの依存関係を記載 -->

## 完了条件
- [ ] 全てのテストが通過
- [ ] リンターエラーなし
- [ ] カバレッジ70%以上維持
- [ ] CLAUDE.md更新済み

## AI実装指示
- TDD厳守 (RED→GREEN→REFACTOR)
- エラーハンドリング適切実装
- 日本語コメント推奨
- 絵文字使用禁止
EOF
    )
    
    local issue_args=("issue" "create" "--title" "$title" "--body" "$body" "--label" "type:feature,priority:medium")
    
    if [ -n "$milestone" ]; then
        issue_args+=("--milestone" "$milestone")
    fi
    
    local issue_number=$(gh "${issue_args[@]}" --json number --jq '.number')
    
    log_success "Issue #$issue_number を作成しました"
    echo "$issue_number"
}

# Claude専用: Sub-Issue作成
claude_create_sub_issue() {
    local parent_issue="$1"
    local title="$2"
    local description="$3"
    
    log_info "Sub-Issue作成中: $title (親: #$parent_issue)"
    
    local body=$(cat << EOF
## 概要
$description

## 親Issue
関連: #$parent_issue

## 実装タスク
- [ ] 実装
- [ ] テスト作成
- [ ] テスト実行

## 完了条件
- [ ] テストが通過
- [ ] リンターエラーなし
- [ ] 親Issueのタスクリスト更新

## AI実装指示
- 単一責任の原則遵守
- 適切なエラーハンドリング
- 日本語コメント
EOF
    )
    
    local issue_number=$(gh issue create \
        --title "$title" \
        --body "$body" \
        --label "type:sub-issue,priority:medium" \
        --json number --jq '.number')
    
    # 親IssueにSub-Issueを追加
    gh issue comment "$parent_issue" --body "Sub-Issue作成: #$issue_number"
    
    log_success "Sub-Issue #$issue_number を作成しました"
    echo "$issue_number"
}

# Claude専用: Feature Branch作成
claude_create_feature_branch() {
    local issue_number="$1"
    local branch_name="feature/issue-${issue_number}"
    
    log_info "Feature Branch作成中: $branch_name"
    
    # メインブランチから最新取得
    git checkout "$DEFAULT_BRANCH"
    git pull origin "$DEFAULT_BRANCH"
    
    # Feature Branch作成
    git checkout -b "$branch_name"
    
    log_success "Branch $branch_name を作成しました"
    
    # Issue進捗更新
    claude_update_issue_progress "$issue_number" "🤖 Claude: 実装開始"
    
    echo "$branch_name"
}

# Claude専用: 進捗更新
claude_update_issue_progress() {
    local issue_number="$1"
    local status="$2"
    
    log_info "Issue #$issue_number の進捗更新中"
    
    gh issue comment "$issue_number" --body "$status"
    
    log_success "進捗を更新しました"
}

# Claude専用: PR作成
claude_create_pr() {
    local issue_number="$1"
    local title="$2"
    local description="$3"
    
    log_info "PR作成中: $title"
    
    local body=$(cat << EOF
## Summary
$description

## Changes
<!-- 変更内容の詳細 -->

## Test Plan
- [ ] 単体テスト実行
- [ ] 統合テスト実行
- [ ] 手動テスト実行

## Checklist
- [ ] テストが通過
- [ ] リンターエラーなし
- [ ] CLAUDE.md更新済み
- [ ] カバレッジ70%以上維持

Closes #$issue_number

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
EOF
    )
    
    local pr_number=$(gh pr create \
        --title "$title" \
        --body "$body" \
        --assignee "@me" \
        --json number --jq '.number')
    
    log_success "PR #$pr_number を作成しました"
    echo "$pr_number"
}

# Claude専用: 品質チェック実行
claude_quality_check() {
    log_info "品質チェック実行中"
    
    # テスト実行
    log_info "テスト実行中..."
    if go test ./...; then
        log_success "テスト通過"
    else
        log_error "テスト失敗"
        return 1
    fi
    
    # リンターチェック
    log_info "リンターチェック中..."
    if command -v golangci-lint &> /dev/null; then
        if golangci-lint run; then
            log_success "リンターチェック通過"
        else
            log_error "リンターエラー"
            return 1
        fi
    else
        log_warn "golangci-lint がインストールされていません"
    fi
    
    # カバレッジチェック
    log_info "カバレッジチェック中..."
    if go test -cover ./... | grep -E "coverage: [0-9]+"; then
        log_success "カバレッジチェック完了"
    else
        log_warn "カバレッジ情報取得に失敗"
    fi
    
    log_success "品質チェック完了"
}

# Claude専用: マイルストーン進捗確認
claude_check_milestone_progress() {
    local milestone="$1"
    
    log_info "マイルストーン進捗確認中: $milestone"
    
    # マイルストーンのIssue一覧取得
    local issues=$(gh issue list --milestone "$milestone" --json number,title,state)
    
    if [ "$issues" = "[]" ]; then
        log_warn "マイルストーン '$milestone' にIssueがありません"
        return 1
    fi
    
    echo "$issues" | jq -r '.[] | "Issue #\(.number): \(.title) (\(.state))"'
    
    local total=$(echo "$issues" | jq length)
    local closed=$(echo "$issues" | jq '[.[] | select(.state == "closed")] | length')
    local progress=$((closed * 100 / total))
    
    log_info "進捗: $closed/$total Issues完了 ($progress%)"
}

# Claude専用: 自動リリース準備
claude_prepare_release() {
    local version="$1"
    local milestone="$2"
    
    log_info "リリース準備中: $version"
    
    # マイルストーンのIssue確認
    claude_check_milestone_progress "$milestone"
    
    # 変更ログ生成
    log_info "変更ログ生成中..."
    local changelog=$(gh pr list --state merged --milestone "$milestone" --json title,number \
        | jq -r '.[] | "- \(.title) (#\(.number))"')
    
    if [ -n "$changelog" ]; then
        log_success "変更ログ:"
        echo "$changelog"
    else
        log_warn "変更ログが空です"
    fi
    
    # タグ作成確認
    read -p "リリース $version を作成しますか? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git tag -a "$version" -m "Release $version"
        git push origin "$version"
        log_success "タグ $version を作成しました"
    fi
}

# メイン関数
main() {
    check_gh_cli
    
    case "${1:-help}" in
        "feature")
            claude_create_feature_issue "$2" "$3" "${4:-}"
            ;;
        "sub-issue")
            claude_create_sub_issue "$2" "$3" "$4"
            ;;
        "branch")
            claude_create_feature_branch "$2"
            ;;
        "progress")
            claude_update_issue_progress "$2" "$3"
            ;;
        "pr")
            claude_create_pr "$2" "$3" "$4"
            ;;
        "quality")
            claude_quality_check
            ;;
        "milestone")
            claude_check_milestone_progress "$2"
            ;;
        "release")
            claude_prepare_release "$2" "$3"
            ;;
        "help"|*)
            echo "Claude Code専用プロジェクト管理スクリプト"
            echo ""
            echo "使用方法:"
            echo "  $0 feature <title> <description> [milestone]  # Feature Issue作成"
            echo "  $0 sub-issue <parent> <title> <description>   # Sub-Issue作成"
            echo "  $0 branch <issue_number>                      # Feature Branch作成"
            echo "  $0 progress <issue_number> <status>           # 進捗更新"
            echo "  $0 pr <issue_number> <title> <description>    # PR作成"
            echo "  $0 quality                                    # 品質チェック"
            echo "  $0 milestone <milestone_name>                 # マイルストーン進捗確認"
            echo "  $0 release <version> <milestone>              # リリース準備"
            echo "  $0 help                                       # このヘルプを表示"
            echo ""
            echo "例:"
            echo "  $0 feature 'ユーザー認証機能' 'ログイン・ログアウト機能の実装' 'v1.0.0'"
            echo "  $0 branch 42"
            echo "  $0 progress 42 '🤖 Claude: 実装完了'"
            echo "  $0 quality"
            ;;
    esac
}

main "$@"