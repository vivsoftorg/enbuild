package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var Version = "v0.0.0" // Set during build time.

var rootCmd = &cobra.Command{
	Use:   "enbuild",
	Short: "enbuild cli",
	Long:  `enbuild is a CLI to work with ENBUILD`,
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			if err := printVersion(); err != nil {
				exitWithError(err, "Cannot get version")
			}
			os.Exit(0)
		}
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithError(err, "")
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolP("version", "v", false, "Print the version and exit")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Username for ENBUILD (or set env variable ENBUILD_USERNAME)")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Password for ENBUILD (or set env variable ENBUILD_PASSWORD)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "API base URL for ENBUILD (or set env variable ENBUILD_BASE_URL)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")
}

func initConfig() {
	//fmt.Println("Loading environment variables...")
}

func printVersion() error {
	currentVersion, err := GetCurrentVersion()
	if err != nil {
		return err
	}
	fmt.Printf("enbuild version %v\n", currentVersion)
	return nil
}

func exitWithError(err error, msg string) {
	if msg != "" {
		fmt.Printf("Error: %s - %s\n", msg, err)
	} else {
		fmt.Printf("Error: %s\n", err)
	}
	os.Exit(1)
}
