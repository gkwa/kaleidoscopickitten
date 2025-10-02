package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
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

	// If frontmatter is empty, use an empty YAML map
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
	var inPlace bool
	flag.BoolVar(&inPlace, "w", false, "write result back to file (in-place)")
	flag.BoolVar(&inPlace, "write", false, "write result back to file (in-place)")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("usage: kaleidoscopickitten [-w] <expression> <file.md>")
	}

	expression := args[0]
	filename := args[1]

	mdBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read file %s: %v", filename, err)
	}
	mdContent := string(mdBytes)

	extractor := NewYAMLFrontmatterExtractor()
	result, err := run(extractor, expression, mdContent, inPlace)
	if err != nil {
		log.Fatal(err)
	}

	if !inPlace {
		fmt.Print(result)
		return
	}

	// Only write if content has changed
	if result == mdContent {
		return
	}

	if err := os.WriteFile(filename, []byte(result), 0o644); err != nil {
		log.Fatalf("failed to write file %s: %v", filename, err)
	}
}
