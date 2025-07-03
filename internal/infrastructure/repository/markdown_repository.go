package repository

import (
	"fmt"
	"os"
	"path/filepath"
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
	defer file.Close()

	// ヘッダーを書き込み
	fmt.Fprintf(file, "# %s\n\n", gamebook.Title)

	// 各パラグラフを書き込み
	for number := 1; number <= 1000; number++ { // 適当な上限
		p, exists := gamebook.Paragraphs[number]
		if !exists {
			continue
		}

		fmt.Fprintf(file, "## パラグラフ %d\n", p.Number)
		fmt.Fprintf(file, "- 概要：%s\n", p.Description)

		if len(p.Choices) > 0 {
			fmt.Fprintf(file, "- 選択肢：\n")
			for _, choice := range p.Choices {
				selected := " "
				if choice.Selected {
					selected = "x"
				}
				fmt.Fprintf(file, "  - [%s] %s → %d\n", selected, choice.Description, choice.TargetNumber)
			}
		}

		if p.Visited {
			fmt.Fprintf(file, "- 訪問済み：はい\n")
		}

		fmt.Fprintln(file)
	}

	// Mermaid形式のフロー図を追加
	fmt.Fprintln(file, "## フロー図")
	fmt.Fprintln(file, "```mermaid")
	fmt.Fprintln(file, "graph TD")

	// ノードの定義
	for number := 1; number <= 1000; number++ {
		p, exists := gamebook.Paragraphs[number]
		if !exists {
			continue
		}

		nodeStyle := ""
		if p.Visited {
			nodeStyle = ":::"
		}
		fmt.Fprintf(file, "    %d[%d: %s]%s\n", p.Number, p.Number, truncate(p.Description, 20), nodeStyle)
	}

	fmt.Fprintln(file)

	// エッジの定義
	for number := 1; number <= 1000; number++ {
		p, exists := gamebook.Paragraphs[number]
		if !exists {
			continue
		}

		for _, choice := range p.Choices {
			arrow := "-.->|%s|"
			if choice.Selected {
				arrow = "-->|%s|"
			}
			fmt.Fprintf(file, "    %d %s %d\n", p.Number, fmt.Sprintf(arrow, truncate(choice.Description, 15)), choice.TargetNumber)
		}
	}

	fmt.Fprintln(file, "```")

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
	inChoices := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// パラグラフの開始
		if strings.HasPrefix(trimmedLine, "## パラグラフ ") {
			if currentParagraph != nil {
				_ = gamebook.AddParagraph(currentParagraph)
			}

			var number int
			if _, err := fmt.Sscanf(trimmedLine, "## パラグラフ %d", &number); err != nil {
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
		if strings.HasPrefix(trimmedLine, "- 概要：") {
			currentParagraph.Description = strings.TrimPrefix(trimmedLine, "- 概要：")
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
					_, _ = fmt.Sscanf(parts[1], "%d", &targetNumber)
				}
			} else {
				selected = false
				// "  - [ ] " を除去して矢印で分割
				parts := strings.Split(strings.TrimPrefix(line, "  - [ ] "), " → ")
				if len(parts) == 2 {
					description = parts[0]
					_, _ = fmt.Sscanf(parts[1], "%d", &targetNumber)
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
