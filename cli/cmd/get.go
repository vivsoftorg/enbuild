package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources from Enbuild",
	Long:  `Get various resources from Enbuild such as catalogs, manifests, etc.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("you must specify the type of resource to get")
		}

		// Find the appropriate subcommand
		for _, subCmd := range cmd.Commands() {
			if subCmd.Name() == args[0] {
				// Run the subcommand with the remaining args
				subCmd.SetArgs(args[1:])
				return subCmd.Execute()
			}
		}

		return fmt.Errorf("unknown resource type: %s", args[0])
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
