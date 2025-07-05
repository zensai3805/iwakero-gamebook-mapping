package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/iwapc/iwakero-gamebook-mapping/internal/domain"
)

// MarkdownRepository はMarkdownファイルでゲームブックを永続化する
type MarkdownRepository struct {
	baseDir string
}

// NewMarkdownRepository は新しいMarkdownRepositoryを作成する
func NewMarkdownRepository(baseDir string) *MarkdownRepository {
	return &MarkdownRepository{
		baseDir: baseDir,
	}
}

// Save はゲームブックをMarkdownファイルとして保存する
func (r *MarkdownRepository) Save(gamebook *domain.Gamebook) error {
	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(r.baseDir, 0755); err != nil {
		return err
	}

	filename := filepath.Join(r.baseDir, gamebook.Title+".md")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	// ヘッダーを書き込み
	_, _ = fmt.Fprintf(file, "# %s\n\n", gamebook.Title)

	// 現在地情報を保存
	if gamebook.Current != nil {
		_, _ = fmt.Fprintf(file, "## 現在地\n")
		_, _ = fmt.Fprintf(file, "- 現在のパラグラフ: %d\n\n", gamebook.Current.Number)
	}

	// 各パラグラフを書き込み（ソートされた順序で）
	keys := make([]int, 0, len(gamebook.Paragraphs))
	for number := range gamebook.Paragraphs {
		keys = append(keys, number)
	}
	sort.Ints(keys)

	for _, number := range keys {
		p := gamebook.Paragraphs[number]

		_, _ = fmt.Fprintf(file, "## パラグラフ %d\n", p.Number)
		_, _ = fmt.Fprintf(file, "- 概要：%s\n", p.Description)

		if len(p.Choices) > 0 {
			_, _ = fmt.Fprintf(file, "- 選択肢：\n")
			for _, choice := range p.Choices {
				selected := " "
				if choice.Selected {
					selected = "x"
				}
				_, _ = fmt.Fprintf(file, "  - [%s] %s → %d\n", selected, choice.Description, choice.TargetNumber)
			}
		}

		if p.Visited {
			_, _ = fmt.Fprintf(file, "- 訪問済み：はい\n")
		}

		_, _ = fmt.Fprintln(file)
	}

	// Mermaid形式のフロー図を追加
	_, _ = fmt.Fprintln(file, "## フロー図")
	_, _ = fmt.Fprintln(file, "```mermaid")
	_, _ = fmt.Fprintln(file, "graph TD")

	// ノードの定義
	for _, number := range keys {
		p := gamebook.Paragraphs[number]
		nodeStyle := ""
		if p.Visited {
			nodeStyle = ":::"
		}
		_, _ = fmt.Fprintf(file, "    %d[%d: %s]%s\n", p.Number, p.Number, truncate(p.Description, 20), nodeStyle)
	}

	_, _ = fmt.Fprintln(file)

	// エッジの定義
	for _, number := range keys {
		p := gamebook.Paragraphs[number]

		for _, choice := range p.Choices {
			arrow := "-.->|%s|"
			if choice.Selected {
				arrow = "-->|%s|"
			}
			_, _ = fmt.Fprintf(file, "    %d %s %d\n", p.Number, fmt.Sprintf(arrow, truncate(choice.Description, 15)), choice.TargetNumber)
		}
	}

	_, _ = fmt.Fprintln(file, "```")

	return nil
}

// Load は指定されたタイトルのゲームブックを読み込む
func (r *MarkdownRepository) Load(title string) (*domain.Gamebook, error) {
	filename := filepath.Join(r.baseDir, title+".md")
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	gamebook := domain.NewGamebook(title)
	lines := strings.Split(string(content), "\n")

	var currentParagraph *domain.Paragraph
	var currentLocationNumber int
	inChoices := false
	inCurrentSection := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 現在地セクションの開始
		if trimmedLine == "## 現在地" {
			inCurrentSection = true
			continue
		}

		// 他のセクションの開始（現在地セクション終了）
		if strings.HasPrefix(trimmedLine, "## ") && trimmedLine != "## 現在地" {
			inCurrentSection = false
		}

		// 現在地情報を読み込み（現在地セクション内のみ）
		if inCurrentSection && strings.HasPrefix(trimmedLine, "- 現在のパラグラフ: ") {
			if n, err := fmt.Sscanf(trimmedLine, "- 現在のパラグラフ: %d", &currentLocationNumber); err != nil || n != 1 {
				currentLocationNumber = 0
			}
			continue
		}

		// パラグラフの開始
		if strings.HasPrefix(trimmedLine, "## パラグラフ ") {
			if currentParagraph != nil {
				_ = gamebook.AddParagraph(currentParagraph)
			}

			var number int
			if n, err := fmt.Sscanf(trimmedLine, "## パラグラフ %d", &number); err != nil || n != 1 {
				continue // パラグラフ番号が読み取れない場合はスキップ
			}
			currentParagraph = domain.NewParagraph(number, "")
			inChoices = false
			continue
		}

		if currentParagraph == nil {
			continue
		}

		// 概要
		if desc, found := strings.CutPrefix(trimmedLine, "- 概要："); found {
			currentParagraph.Description = desc
			continue
		}

		// 選択肢セクション
		if trimmedLine == "- 選択肢：" {
			inChoices = true
			continue
		}

		// 選択肢の内容（元の行でチェック）
		if inChoices && strings.HasPrefix(line, "  - [") {
			var description string
			var targetNumber int
			var selected bool

			// 選択状態を確認
			if strings.Contains(line, "[x]") {
				selected = true
				// "  - [x] " を除去して矢印で分割
				parts := strings.Split(strings.TrimPrefix(line, "  - [x] "), " → ")
				if len(parts) == 2 {
					description = parts[0]
					if n, err := fmt.Sscanf(parts[1], "%d", &targetNumber); err != nil || n != 1 {
						targetNumber = 0 // パース失敗時は0に設定
					}
				}
			} else {
				selected = false
				// "  - [ ] " を除去して矢印で分割
				parts := strings.Split(strings.TrimPrefix(line, "  - [ ] "), " → ")
				if len(parts) == 2 {
					description = parts[0]
					if n, err := fmt.Sscanf(parts[1], "%d", &targetNumber); err != nil || n != 1 {
						targetNumber = 0 // パース失敗時は0に設定
					}
				}
			}

			if description != "" && targetNumber > 0 {
				currentParagraph.AddChoice(description, targetNumber)
				if selected {
					_ = currentParagraph.SelectChoice(len(currentParagraph.Choices) - 1)
				}
			}
			continue
		}

		// 訪問済み
		if strings.HasPrefix(trimmedLine, "- 訪問済み：はい") {
			currentParagraph.Visited = true
			inChoices = false // 選択肢セクション終了
			continue
		}

		// 他の "- " で始まる行も選択肢セクション終了の合図
		if inChoices && strings.HasPrefix(trimmedLine, "- ") && !strings.HasPrefix(line, "  - [") {
			inChoices = false
		}

		// フロー図セクションに達したら終了
		if trimmedLine == "## フロー図" {
			break
		}
	}

	// 最後のパラグラフを追加
	if currentParagraph != nil {
		_ = gamebook.AddParagraph(currentParagraph)
	}

	// 現在地を設定
	if currentLocationNumber > 0 {
		_ = gamebook.MoveTo(currentLocationNumber)
		// 現在地が存在しない場合は無視して続行
	}

	return gamebook, nil
}

// List は保存されているゲームブックのタイトル一覧を返す
func (r *MarkdownRepository) List() ([]string, error) {
	files, err := os.ReadDir(r.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var titles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			title := strings.TrimSuffix(file.Name(), ".md")
			titles = append(titles, title)
		}
	}

	return titles, nil
}

// Delete は指定されたタイトルのゲームブックを削除する
func (r *MarkdownRepository) Delete(title string) error {
	filename := filepath.Join(r.baseDir, title+".md")
	return os.Remove(filename)
}

// truncate は文字列を指定された長さに切り詰める
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}
