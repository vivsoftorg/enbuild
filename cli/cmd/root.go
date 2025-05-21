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
	cobra.OnInitialize(func() {
		// Set username from flag or ENV
		if username == "" {
			if envUser := os.Getenv("ENBUILD_USERNAME"); envUser != "" {
				username = envUser
			}
		}
		// Set password from flag or ENV
		if password == "" {
			if envPass := os.Getenv("ENBUILD_PASSWORD"); envPass != "" {
				password = envPass
			}
		}
		// Set baseURL from flag or ENV or default
		if baseURL == "" {
			if envBase := os.Getenv("ENBUILD_BASE_URL"); envBase != "" {
				baseURL = envBase
			}
		}
		// set default baseURL if not set
		if baseURL == "" {
			baseURL = "https://enbuild.vivplatform.io"
		}
		// Set debug from flag or ENV
		if !debug {
			if envDebug := os.Getenv("ENBUILD_DEBUG"); envDebug != "" {
				if envDebug == "1" || envDebug == "true" || envDebug == "TRUE" {
					debug = true
				}
			}
		}
	})
	rootCmd.Flags().BoolP("version", "v", false, "Print the version and exit")

	// Make these flags persistent so they apply to all subcommands
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Username for ENBUILD (or set env variable ENBUILD_USERNAME)")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Password for ENBUILD (or set env variable ENBUILD_PASSWORD)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "API base URL for ENBUILD default (or set env variable ENBUILD_BASE_URL)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output(or set env variable ENBUILD_DEBUG=1)")
}
