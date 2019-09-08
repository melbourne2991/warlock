package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "warlock [subcommand]",
	Short: "warlock provides simple file/directory encryption",
	Long: `
warlock provides simple file/directory encryption by allowing different "vaults" to be defined,
a vault is effectively an identifier for a set of paths to be encrypted and decrypted together`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
