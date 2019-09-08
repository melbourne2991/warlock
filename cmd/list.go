package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists vaults",
	Run: func(cmd *cobra.Command, args []string) {

		vaults := viper.GetStringMapStringSlice(configVaultKey)

		for vaultName, paths := range vaults {
			fmt.Println("\n" + cyanHighlight(underline(vaultName)))

			for i := 0; i < len(paths); i++ {
				pathValue := paths[i]
				fmt.Println(pathValue)
			}
		}

		fmt.Print("\n")
	},
}
