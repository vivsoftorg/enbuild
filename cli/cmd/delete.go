package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources from Enbuild",
	Long:  `Delete various resources from Enbuild such as stacks, etc.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("you must specify the type of resource to delete")
		}

		for _, subCmd := range cmd.Commands() {
			if subCmd.Name() == args[0] {
				subCmd.SetArgs(args[1:])
				return subCmd.Execute()
			}
		}

		return fmt.Errorf("unknown resource type: %s", args[0])
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}