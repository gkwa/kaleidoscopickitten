package main

import (
	"bufio"
	"strings"
)

type YAMLFrontmatterExtractor struct{}

func NewYAMLFrontmatterExtractor() *YAMLFrontmatterExtractor {
	return &YAMLFrontmatterExtractor{}
}

type ExtractionResult struct {
	Frontmatter    string
	Body           string
	HasFrontmatter bool
}

func (e *YAMLFrontmatterExtractor) Extract(content string) (*ExtractionResult, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var frontmatter strings.Builder
	var body strings.Builder
	var inFrontmatter bool
	var delimiterCount int

	for scanner.Scan() {
		line := scanner.Text()

		// Handle frontmatter delimiters
		if line == "---" {
			delimiterCount++
			if delimiterCount == 1 {
				inFrontmatter = true
				continue
			}
			if delimiterCount == 2 {
				inFrontmatter = false
				continue
			}
		}

		if inFrontmatter {
			frontmatter.WriteString(line)
			frontmatter.WriteString("\n")
			continue
		}

		body.WriteString(line)
		body.WriteString("\n")
	}

	return &ExtractionResult{
		Frontmatter:    frontmatter.String(),
		Body:           body.String(),
		HasFrontmatter: delimiterCount >= 2,
	}, nil
}
