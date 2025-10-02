package main

import (
	"bufio"
	"errors"
	"strings"
)

type YAMLFrontmatterExtractor struct{}

func NewYAMLFrontmatterExtractor() *YAMLFrontmatterExtractor {
	return &YAMLFrontmatterExtractor{}
}

type ExtractionResult struct {
	Frontmatter     string
	Body            string
	HasFrontmatter  bool
	IsValid         bool
	ValidationError error
}

var (
	ErrInvalidFrontmatter = errors.New("invalid frontmatter: missing closing delimiter")
	ErrScalarFrontmatter  = errors.New("invalid frontmatter: must be a YAML object (key-value pairs), not a scalar value")
)

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

	result := &ExtractionResult{
		Frontmatter:    frontmatter.String(),
		Body:           body.String(),
		HasFrontmatter: delimiterCount >= 2,
		IsValid:        true,
	}

	// Validate frontmatter structure
	if delimiterCount == 1 {
		result.IsValid = false
		result.ValidationError = ErrInvalidFrontmatter
		return result, nil
	}

	// Check if frontmatter is a scalar value (not a key-value structure)
	if delimiterCount >= 2 {
		fm := strings.TrimSpace(result.Frontmatter)
		if fm != "" && !strings.Contains(fm, ":") && fm != "{}" {
			result.IsValid = false
			result.ValidationError = ErrScalarFrontmatter
			return result, nil
		}
	}

	return result, nil
}
