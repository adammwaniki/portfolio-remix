package content

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

var md = goldmark.New()

// LoadCardFromMarkdown reads a markdown file and returns a Card.
// Frontmatter is parsed as simple "key: value" pairs between --- delimiters.
func LoadCardFromMarkdown(path string) (Card, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Card{}, err
	}

	front, body := parseFrontmatter(string(data))

	var buf bytes.Buffer
	if err := md.Convert([]byte(body), &buf); err != nil {
		return Card{}, err
	}

	id := strings.TrimSuffix(filepath.Base(path), ".md")

	tags := parseTags(front["tags"])

	return Card{
		ID:          id,
		Title:       front["title"],
		Subtitle:    front["tags"], // Keep original "Go · Architecture" style for display
		Tags:        tags,
		Description: front["description"],
		CardIcon:    front["icon"],
		ReadingTime: front["reading_time"],
		Date:        front["date"],
		Updated:     front["updated"],
		Detail:      buf.String(),
		DemoURL:     front["demo_url"],
	}, nil
}

// LoadCardsFromDir reads all .md files from a directory in the given order.
// If order is nil, files are loaded alphabetically.
func LoadCardsFromDir(dir string, order []string) ([]Card, error) {
	if order != nil {
		cards := make([]Card, 0, len(order))
		for _, id := range order {
			path := filepath.Join(dir, id+".md")
			card, err := LoadCardFromMarkdown(path)
			if err != nil {
				return nil, err
			}
			cards = append(cards, card)
		}
		return cards, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var cards []Card
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		card, err := LoadCardFromMarkdown(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

// parseFrontmatter splits --- delimited frontmatter from the body.
func parseFrontmatter(content string) (map[string]string, string) {
	front := make(map[string]string)

	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "---") {
		return front, content
	}

	rest := content[3:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return front, content
	}

	frontBlock := rest[:idx]
	body := strings.TrimSpace(rest[idx+4:])

	for _, line := range strings.Split(frontBlock, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		colon := strings.Index(line, ":")
		if colon < 0 {
			continue
		}
		key := strings.TrimSpace(line[:colon])
		value := strings.TrimSpace(line[colon+1:])
		// Strip surrounding quotes
		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}
		front[key] = value
	}

	return front, body
}

// parseTags splits a comma-or-dot-separated tag string into a slice.
func parseTags(s string) []string {
	// Handle both "Go, Architecture" and "Go · Architecture" formats
	s = strings.ReplaceAll(s, " · ", ", ")
	parts := strings.Split(s, ",")
	var tags []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			tags = append(tags, p)
		}
	}
	return tags
}
