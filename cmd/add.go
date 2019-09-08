package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add [vault] [path]",
	Short: "Adds a path to a vault",
	Long:  "Adds a path to a vault, if the vault does not yet exist one will be created automatically",
	Args: func(cmd *cobra.Command, args []string) error {
		errors.New("requires vault name and path")

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		vaultName := args[0]
		pathValue := args[1]

		absPath, err := filepath.Abs(pathValue)

		if err != nil {
			fmt.Println("Problem with provided path:", err)
			os.Exit(1)
			return
		}

		vaults := viper.GetStringMapStringSlice(configVaultKey)

		vaults[vaultName] = append(vaults[vaultName], absPath)

		viper.Set(configVaultKey, vaults)
		viper.WriteConfig()

		fmt.Println(fmt.Sprintf("\nAdded %s to %s\n", cyanHighlight(absPath), cyanHighlight(vaultName)))
	},
}
