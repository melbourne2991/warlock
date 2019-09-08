package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(unlockCmd)
}

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlocks one or more vaults",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires vault name")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		vaultName := args[0]

		vaults := viper.GetStringMapStringSlice(configVaultKey)
		targetVault := vaults[vaultName]

		if targetVault == nil {
			fmt.Println(fmt.Sprintf("Error: \"%s\" vault does not exist", vaultName))
			return
		}

		password := passwordPrompt()

		fmt.Println(targetVault)

		for i := 0; i < len(targetVault); i++ {
			pathValue := targetVault[i]
			fmt.Println(fmt.Sprintf("unlocking %s...", pathValue))

			if err := pathLocker.UnlockPath(pathValue, password); err != nil {
				fmt.Println("Error: Could not unlock path -", err)
				return
			}

			fmt.Println("unlocked.")
		}
	},
}
