package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// Version of the application, should be set during build time.
var Version = "v0.0.0" // Replace with current version number as needed.

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "enbuild",
	Short: "enbuild cli",
	Long:  `enbuild is a CLI to work with ENBUILD`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the version flag was passed
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			currentVersion, err := GetCurrentVersion()
			if err != nil {
				fmt.Printf("Error Can not get version %s\n", err)
				os.Exit(1)
				panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
			}
			fmt.Printf("enbuild version %v\n", currentVersion)
			os.Exit(0) // Exit after printing the version
		}
		// If no flag, print default help
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.Flags().BoolP("version", "v", false, "Print the version and exit")
	rootCmd.Flags().StringVar(&token, "token", "", "API token (or set env variable ENBUILD_API_TOKEN)")
	rootCmd.Flags().StringVar(&baseURL, "base-url", "https://enbuild-dev.vivplatform.io/enbuild-bk/", "API base URL for ENBUILD (or set env variable ENBUILD_BASE_URL)")
	rootCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
}
