/*
Copyright © 2026 zigorou
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <category>",
	Short: "Decrypt a vault file, open it in an editor, then re-encrypt on save",
	Long: `Decrypt a vault file, open it in the default editor (respecting $EDITOR),
and re-encrypt it when you save and exit.

Example:
  ark edit obsidian
  ark edit aws`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
