package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resources for Enbuild",
	Long:  `Create various resources for Enbuild such as haul manifests.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		// Find the appropriate subcommand
		for _, subCmd := range cmd.Commands() {
			if subCmd.Name() == args[0] {
				// Run the subcommand with the remaining args
				subCmd.SetArgs(args[1:])
				return subCmd.Execute()
			}
		}

		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
