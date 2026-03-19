/*
Copyright © 2026 zigorou
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set <uri> <value>",
	Short: "Add or update a secret value in the vault",
	Long: `Add or update a secret in the vault identified by a URI.

URI format: ark://<category>/<item>/<field>

Example:
  ark set ark://obsidian/work/new_api_key "my-secret-value"
  ark set ark://aws/prod/secret_access_key "AKIAIOSFODNN7EXAMPLE"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
