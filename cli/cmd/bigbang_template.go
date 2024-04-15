package cmd

import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"gopkg.in/yaml.v3"

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
	valuesDirectory := fmt.Sprintf("%s/bb_values", targetDirectory)
	secretsDirectory := fmt.Sprintf("%s/bb_secrets", targetDirectory)
	bbvaluesFile := fmt.Sprintf("%s/values_%s.yaml", targetDirectory, bbVersion)

	if err := os.MkdirAll(valuesDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDirectory, err)
	}

	if err := os.MkdirAll(secretsDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDirectory, err)
	}

	if err := downloadAndSaveFile(fmt.Sprintf("https://raw.githubusercontent.com/DoD-Platform-One/bigbang/%s/chart/values.yaml", bbVersion), bbvaluesFile); err != nil {
		return err
	}

	splitBBValues(bbvaluesFile, valuesDirectory)

	fmt.Printf("Created BigBang ENBUILD Catalog template for version %s at %s\n", bbVersion, valuesDirectory)
	return nil
}



func splitBBValues(bbvaluesFile string, valuesDirectory string) error {
    // Define the repository keys
	REPOSITORY_KEYS := "domain offline helmRepositories registryCredentials openshift git sso flux networkPolicies imagePullPolicy"
    repositoryKeys := strings.Split(REPOSITORY_KEYS, " ")

    // Read the input values file
    content, err := ioutil.ReadFile(bbvaluesFile)
    if err != nil {
        return fmt.Errorf("failed to read values file: %w", err)
    }

    // Split the content into multiple files based on the key
    values := make(map[string]interface{})
    if err := yaml.Unmarshal(content, &values); err != nil {
        return fmt.Errorf("failed to unmarshal values file: %w", err)
    }

    repositoryValues := make(map[string]interface{})
    for key, value := range values {
        if contains(repositoryKeys, key) {
            // Add the value to the repository values
            repositoryValues[key] = value
        } else {
            // Convert the value back to YAML
            yamlContent, err := yaml.Marshal(value)
            if err != nil {
                return fmt.Errorf("failed to marshal value to YAML: %w", err)
            }

            // Write the YAML content to a file
            outputFile := fmt.Sprintf("%s/%s.yaml", valuesDirectory, key)
            if err := ioutil.WriteFile(outputFile, yamlContent, 0644); err != nil {
                return fmt.Errorf("failed to write output file: %w", err)
            }
        }
    }

    // Write the repository values to repository.yaml
    repositoryYamlContent, err := yaml.Marshal(repositoryValues)
    if err != nil {
        return fmt.Errorf("failed to marshal repository values to YAML: %w", err)
    }
    repositoryFile := fmt.Sprintf("%s/repository.yaml", valuesDirectory)
    if err := ioutil.WriteFile(repositoryFile, repositoryYamlContent, 0644); err != nil {
        return fmt.Errorf("failed to write repository file: %w", err)
    }

    return nil
}

// Helper function to check if a slice contains a string
func contains(slice []string, str string) bool {
    for _, v := range slice {
        if v == str {
            return true
        }
    }
    return false
}
