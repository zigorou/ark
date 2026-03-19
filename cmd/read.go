/*
Copyright © 2026 zigorou
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read <uri>",
	Short: "Read a secret value from the vault",
	Long: `Read a secret value from the vault by its URI and print it to stdout.

URI format: ark://<category>/<item>/<field>

Example:
  ark read ark://aws/prod/access_key_id
  ark read ark://obsidian/work/api_key`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
