package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version of the application, should be set during build time.
var version = "v0.0.6" // Replace with current version number as needed.

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "enbuild",
	Short: "enbuild cli",
	Long:  `enbuild is a CLI to help generate the ENBUILD catalog templates`,
	// Uncomment the following line if your bare application has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	rootCmd.Version = version
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number of enbuild CLI")
	rootCmd.SetVersionTemplate("{{printf \"%s\\n\" .Version}}") // Custom version output template, change as needed
}
