package cmd

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
	"os"
)

func GetCurrentVersion() (*semver.Version, error) {
	version, err := semver.NewVersion(Version)
	if err != nil {
		return nil, fmt.Errorf("error trying to get semver from raw string `%s`, error: `%w`", Version, err)
	}

	return version, nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print installed version of the Enbuild CLI",
	Run: func(cmd *cobra.Command, args []string) {
		currentVersion, err := GetCurrentVersion()
		if err != nil {
			fmt.Printf("Error Can not get version %s\n", err)
			os.Exit(1)
			panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
		}
		fmt.Printf("enbuild version %v\n", currentVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
