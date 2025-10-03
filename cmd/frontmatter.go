package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/gkwa/kaleidoscopickitten/frontmatter"
	"github.com/spf13/cobra"
)

func NewFrontmatterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "frontmatter [expression] <file.md>",
		Aliases: []string{"fm"},
		Short:   "Work with YAML frontmatter",
		Long:    "Query and modify YAML frontmatter in markdown files. Defaults to print mode if no subcommand specified.",
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

			extractor := frontmatter.NewYAMLFrontmatterExtractor()
			result, err := frontmatter.Run(extractor, expression, mdContent, false)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(result)
		},
	}

	cmd.AddCommand(NewPrintCmd())
	cmd.AddCommand(NewEditCmd())

	return cmd
}

func NewPrintCmd() *cobra.Command {
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

			extractor := frontmatter.NewYAMLFrontmatterExtractor()
			result, err := frontmatter.Run(extractor, expression, mdContent, false)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(result)
		},
	}
}

func NewEditCmd() *cobra.Command {
	return &cobra.Command{
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

			extractor := frontmatter.NewYAMLFrontmatterExtractor()
			result, err := frontmatter.Run(extractor, expression, mdContent, true)
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
}
