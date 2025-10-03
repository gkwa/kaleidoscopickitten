package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "kaleidoscopickitten",
		Short: "Query and modify YAML frontmatter in markdown files",
		Long:  "Query and modify YAML frontmatter in markdown files using yq expressions.",
	}

	rootCmd.AddCommand(NewFrontmatterCmd())

	return rootCmd
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
