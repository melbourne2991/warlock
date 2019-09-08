package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var password string

func init() {
	lockCmd.Flags().StringVarP(&password, "password", "p", "", "password")
	rootCmd.AddCommand(lockCmd)
}

var lockCmd = &cobra.Command{
	Use:   "lock [OPTIONS]",
	Short: "Locks one or more vaults",
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

		if password == "" {
			password = passwordPrompt()
		}

		fmt.Println(targetVault)

		for i := 0; i < len(targetVault); i++ {
			pathValue := targetVault[i]

			fmt.Println(fmt.Sprintf("locking %s...", pathValue))

			if err := pathLocker.LockPath(pathValue, password); err != nil {
				fmt.Println("Error: Could not lock path -", err)
				return
			}

			fmt.Println("locked.")
		}
	},
}
