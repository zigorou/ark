/*
Copyright © 2026 zigorou
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var injectCmd = &cobra.Command{
	Use:   "inject -i <template> -o <output>",
	Short: "Process a template file, substituting ark:// URIs with secret values",
	Long: `Process a template file by resolving all ark:// URIs and writing the
result to an output file. Equivalent to 1Password's "op inject".

Example:
  ark inject -i config.tpl -o config.json
  ark inject -i .env.tpl -o .env`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}

func init() {
	injectCmd.Flags().StringP("in", "i", "", "input template file (required)")
	injectCmd.Flags().StringP("out", "o", "", "output file (required)")
	injectCmd.MarkFlagRequired("in")
	injectCmd.MarkFlagRequired("out")
	rootCmd.AddCommand(injectCmd)
}
