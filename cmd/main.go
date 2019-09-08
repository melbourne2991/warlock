package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/melbourne2991/warlock/lib"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pathLocker *lib.PathLocker

var appDir string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&appDir, "dir", "", "Path to warlock store and config directory (default is $HOME/.warlock)")

	viper.SetDefault(configVaultKey, make(map[string][]string))
}

// Execute executes the command line app
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	home, err := homedir.Dir()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	dirMode := os.FileMode(int(0755))

	if appDir == "" {
		appDir = path.Join(home, ".warlock")
	}

	configFilePath := path.Join(appDir, "config.toml")
	storeDir := path.Join(appDir, "store")

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(appDir)

	createDirIfNotExist(appDir, dirMode)
	createDirIfNotExist(storeDir, dirMode)

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		fmt.Println("Creating config file:", configFilePath)

		if _, err := os.Create(configFilePath); err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println("Writing to config...")
		viper.WriteConfig()
	} else {
		viper.ReadInConfig()
	}

	pathLocker = lib.NewPathLocker(storeDir)
}

func createDirIfNotExist(dirPath string, dirMode os.FileMode) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("%s does not exist, initializing...", dirPath))

		if err := os.Mkdir(dirPath, dirMode); err != nil {
			log.Fatal(err)
			return
		}
	}
}
