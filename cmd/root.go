/*
Copyright © 2026 zigorou
*/
package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/zigorou/ark/internal/vault"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ark",
	Short: "A secrets manager CLI using GitHub private repos as vault storage",
	Long: `ark is an open-source secrets manager CLI that uses GitHub private
repositories as encrypted vault storage (via SOPS + age).

It provides a 1Password-like UX — URI scheme, run, inject, and edit —
without the cost or centralized server dependency.

URI scheme: ark://<category>/<item>/<field>

Example:
  ark read ark://aws/prod/access_key_id
  ark run --env API_KEY=ark://obsidian/work/api_key -- some-command
  ark inject -i config.tpl -o config.json`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		dir, err := vault.ResolveDir(cmd.Flags().Lookup("vault").Value.String())
		if err != nil {
			cmd.PrintErrln("ark: " + err.Error())
			cmd.SilenceUsage = true
			return err
		}
		if err := vault.CheckInitialized(dir); err != nil {
			if errors.Is(err, vault.ErrVaultNotInitialized) {
				cmd.PrintErrln("ark: vault is not initialized.")
				cmd.PrintErrln("Run 'ark init --repo <repository>' to set up a new vault,")
				cmd.PrintErrln("or 'ark init' to initialize a local vault without a remote.")
			} else {
				cmd.PrintErrln("ark: " + err.Error())
			}
			cmd.SilenceUsage = true
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("vault", "v", "", "vault repository path (overrides config)")
}
