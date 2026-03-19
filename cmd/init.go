/*
Copyright © 2026 zigorou
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init --repo <repository>",
	Short: "Initialize ark by linking a GitHub private repository as the vault",
	Long: `Initialize ark by linking a GitHub private repository as the encrypted
vault. This sets up the local configuration and clones the repository if needed.

Example:
  ark init --repo github.com/zigorou/my-vault`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}

func init() {
	initCmd.Flags().StringP("repo", "r", "", "GitHub repository to use as vault (e.g. github.com/user/my-vault)")
	initCmd.MarkFlagRequired("repo")
	rootCmd.AddCommand(initCmd)
}
