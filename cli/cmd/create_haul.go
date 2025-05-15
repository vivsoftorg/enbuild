package cmd

import (
	"github.com/spf13/cobra"
)

// createHaulCmd represents the haul command under create
var createHaulCmd = &cobra.Command{
	Use:   "haul",
	Short: "Create haul manifests",
	Long:  `Create haul manifests for various components.`,
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
	createCmd.AddCommand(createHaulCmd)
}
