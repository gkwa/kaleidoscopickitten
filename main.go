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
		return "", fmt.Errorf("yaml processing error: %w", err)
	}

	return result, nil
}

func reconstructFile(processedFrontmatter, body string) string {
	var output strings.Builder
	output.WriteString("---\n")
	output.WriteString(strings.TrimSpace(processedFrontmatter) + "\n")
	output.WriteString("---\n")
	output.WriteString(body)
	return output.String()
}

func processFrontmatter(extractor FrontmatterExtractor, expression, mdContent string) (string, error) {
	result, err := extractor.Extract(mdContent)
	if err != nil {
		return "", err
	}

	// Check for validation errors
	if !result.IsValid {
		return "", result.ValidationError
	}

	processed, err := processYAML(result.Frontmatter, expression)
	if err != nil {
		return "", err
	}

	return processed, nil
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

func createPrintCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "print [expression] <file.md>",
		Aliases: []string{"p"},
		Short:   "Print frontmatter (read-only)",
		Long:    "Query and display YAML frontmatter without modifying the file. If no expression is provided, defaults to '.' (show all).",
		Args:    cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			var expression, filename string

			if len(args) == 1 {
				expression = "."
				filename = args[0]
			} else {
				expression = args[0]
				filename = args[1]
			}

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
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "kaleidoscopickitten",
		Short: "Query and modify YAML frontmatter in markdown files",
		Long:  "Query and modify YAML frontmatter in markdown files using yq expressions.",
	}

	frontmatterCmd := &cobra.Command{
		Use:     "frontmatter [expression] <file.md>",
		Aliases: []string{"fm"},
		Short:   "Work with YAML frontmatter",
		Long:    "Query and modify YAML frontmatter in markdown files. Defaults to print mode if no subcommand specified.",
		Args:    cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			// Default to print behavior when no subcommand is specified
			var expression, filename string

			if len(args) == 1 {
				expression = "."
				filename = args[0]
			} else {
				expression = args[0]
				filename = args[1]
			}

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

	printCmd := createPrintCommand()

	editCmd := &cobra.Command{
		Use:     "edit <expression> <file.md>",
		Aliases: []string{"e"},
		Short:   "Edit frontmatter (in-place)",
		Long:    "Modify YAML frontmatter and write changes back to the file.",
		Args:    cobra.ExactArgs(2),
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

	frontmatterCmd.AddCommand(printCmd)
	frontmatterCmd.AddCommand(editCmd)
	rootCmd.AddCommand(frontmatterCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
