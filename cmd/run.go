/*
Copyright © 2026 zigorou
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [--env <NAME=ark://...>]... -- <command> [args...]",
	Short: "Run a command with secrets injected as environment variables",
	Long: `Run a command with secrets resolved from the vault and injected as
environment variables. Equivalent to 1Password's "op run".

Example:
  ark run \
    --env API_KEY=ark://obsidian/work/api_key \
    --env AWS_ACCESS_KEY_ID=ark://aws/prod/access_key_id \
    -- some-command --some-flag`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}

func init() {
	runCmd.Flags().StringArrayP("env", "e", nil, "environment variable mapping: NAME=ark://category/item/field")
	rootCmd.AddCommand(runCmd)
}
