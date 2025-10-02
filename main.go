package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/spf13/cobra"
	logging "gopkg.in/op/go-logging.v1"
)

func init() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backendLeveled)
}

type FrontmatterExtractor interface {
	Extract(content string) (*ExtractionResult, error)
}

func processYAML(yamlString, expression string) (string, error) {
	yqlib.InitExpressionParser()

	decoder := yqlib.NewYamlDecoder(yqlib.ConfiguredYamlPreferences)
	encoder := yqlib.NewYamlEncoder(yqlib.ConfiguredYamlPreferences)

	input := yamlString
	if strings.TrimSpace(yamlString) == "" {
		input = "{}"
	}

	stringEval := yqlib.NewStringEvaluator()
	result, err := stringEval.Evaluate(expression, input, encoder, decoder)
	if err != nil {
		return "", err
	}

	return result, nil
}

func reconstructFile(processedFrontmatter, body string) string {
	var output strings.Builder
	output.WriteString("---\n")
	output.WriteString(strings.TrimSpace(processedFrontmatter))
	output.WriteString("\n---\n")
	output.WriteString(body)
	return output.String()
}

func processFrontmatter(extractor FrontmatterExtractor, expression, mdContent string) (string, error) {
	result, err := extractor.Extract(mdContent)
	if err != nil {
		return "", err
	}

	return processYAML(result.Frontmatter, expression)
}

func run(extractor FrontmatterExtractor, expression, mdContent string, fullFile bool) (string, error) {
	processedFrontmatter, err := processFrontmatter(extractor, expression, mdContent)
	if err != nil {
		return "", err
	}

	if !fullFile {
		return processedFrontmatter, nil
	}

	result, err := extractor.Extract(mdContent)
	if err != nil {
		return "", err
	}

	return reconstructFile(processedFrontmatter, result.Body), nil
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "kaleidoscopickitten",
		Short: "Query and modify YAML frontmatter in markdown files",
		Long:  "Query and modify YAML frontmatter in markdown files using yq expressions.",
	}

	frontmatterCmd := &cobra.Command{
		Use:   "frontmatter",
		Short: "Work with YAML frontmatter",
		Long:  "Query and modify YAML frontmatter in markdown files.",
	}

	viewCmd := &cobra.Command{
		Use:   "view <expression> <file.md>",
		Short: "View frontmatter (read-only)",
		Long:  "Query and display YAML frontmatter without modifying the file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			expression := args[0]
			filename := args[1]

			mdBytes, err := os.ReadFile(filename)
			if err != nil {
				log.Fatalf("failed to read file %s: %v", filename, err)
			}
			mdContent := string(mdBytes)

			extractor := NewYAMLFrontmatterExtractor()
			result, err := run(extractor, expression, mdContent, false)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(result)
		},
	}

	editCmd := &cobra.Command{
		Use:   "edit <expression> <file.md>",
		Short: "Edit frontmatter (in-place)",
		Long:  "Modify YAML frontmatter and write changes back to the file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			expression := args[0]
			filename := args[1]

			mdBytes, err := os.ReadFile(filename)
			if err != nil {
				log.Fatalf("failed to read file %s: %v", filename, err)
			}
			mdContent := string(mdBytes)

			extractor := NewYAMLFrontmatterExtractor()
			result, err := run(extractor, expression, mdContent, true)
			if err != nil {
				log.Fatal(err)
			}

			if result == mdContent {
				return
			}

			if err := os.WriteFile(filename, []byte(result), 0o644); err != nil {
				log.Fatalf("failed to write file %s: %v", filename, err)
			}
		},
	}

	frontmatterCmd.AddCommand(viewCmd)
	frontmatterCmd.AddCommand(editCmd)
	rootCmd.AddCommand(frontmatterCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
