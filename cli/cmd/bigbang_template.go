package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


// createTemplateCmd represents the createtemplate command
var createTemplateCmd = &cobra.Command{
	Use:   "create-template",
	Short: "Create a BigBang ENBUILD Catalog template for given version",
	Long:  "Create a BigBang ENBUILD Catalog template for given version",
	RunE:  createtemplate,
}

func init() {
	bigbangCmd.AddCommand(createTemplateCmd)
	createTemplateCmd.Flags().StringP("bb-version", "v", "", "Specify the BigBang version (required)")
	_ = createTemplateCmd.MarkFlagRequired("bb-version")
}

func createtemplate(cmd *cobra.Command, args []string) error {
	bbVersion, err := cmd.Flags().GetString("bb-version")
	if err != nil {
		return fmt.Errorf("failed to read 'bb-version' flag: %w", err)
	}
	fmt.Printf("Creating BigBang ENBUILD Catalog template for version %s\n", bbVersion)

	targetDirectory := "target"
	if err := os.MkdirAll(targetDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDirectory, err)
	}


	return createBigBangTemplate(bbVersion, targetDirectory)
}



func createBigBangTemplate( bbVersion string, targetDirectory string) error {
	valuesDirectory := targetDirectory + "/bb_values"
	// if err := os.MkdirAll(targetDirectory, 0755); err != nil {
	// 	return fmt.Errorf("failed to create directory %s: %w", targetDirectory, err)
	// }
	fmt.Printf("Created BigBang ENBUILD Catalog template for version %s at %s\n", bbVersion, valuesDirectory)
	return nil
}
