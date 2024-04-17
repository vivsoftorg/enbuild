package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	targetDirectory      = "target"
	valuesDirectoryName  = "bb_values"
	secretsDirectoryName = "bb_secrets"
	repositoryKeys       = "domain offline helmRepositories registryCredentials openshift git sso flux networkPolicies imagePullPolicy packages wrapper"
	secretsContent       = `domain: ""`
)

var createTemplateCmd = &cobra.Command{
	Use:   "create-template",
	Short: "Create a BigBang ENBUILD Catalog template for given version",
	Long:  "Create a BigBang ENBUILD Catalog template for given version",
	RunE:  createTemplate,
}

func init() {
	bigbangCmd.AddCommand(createTemplateCmd)
	createTemplateCmd.Flags().StringP("bb-version", "v", "", "Specify the BigBang version (required)")
	if err := createTemplateCmd.MarkFlagRequired("bb-version"); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting bb-version flag as required: %v\n", err)
		os.Exit(1)
	}
}

func createTemplate(cmd *cobra.Command, args []string) error {
	bbVersion, err := cmd.Flags().GetString("bb-version")
	if err != nil {
		return fmt.Errorf("failed to read 'bb-version' flag: %w", err)
	}
	fmt.Printf("Creating BigBang ENBUILD Catalog template for version %s\n", bbVersion)

	if err := os.MkdirAll(targetDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDirectory, err)
	}

	valuesDirectory := fmt.Sprintf("%s/%s", targetDirectory, valuesDirectoryName)
	secretsDirectory := fmt.Sprintf("%s/%s", targetDirectory, secretsDirectoryName)

	if err := os.MkdirAll(valuesDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", valuesDirectory, err)
	}

	if err := os.MkdirAll(secretsDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", secretsDirectory, err)
	}

	bbValuesFile := fmt.Sprintf("%s/values_%s.yaml", targetDirectory, bbVersion)

	if err := ensureBBValues(bbValuesFile, bbVersion); err != nil {
		return err
	}

	if err := splitBBValues(bbValuesFile, valuesDirectory, secretsDirectory); err != nil {
		return fmt.Errorf("failed to split BB values: %w", err)
	}

	return nil
}

func ensureBBValues(bbValuesFile, bbVersion string) error {
	_, err := os.Stat(bbValuesFile)
	if os.IsNotExist(err) {
		downloadURL := fmt.Sprintf("https://raw.githubusercontent.com/DoD-Platform-One/bigbang/%s/chart/values.yaml", bbVersion)
		return downloadAndSaveFile(downloadURL, bbValuesFile)
	} else if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", bbValuesFile, err)
	}
	fmt.Printf("File %s already exists. Skipping download.\n", bbValuesFile)
	return nil
}

func splitBBValues(bbValuesFile string, valuesDirectory string, secretsDirectory string) error {
	keys := strings.Split(repositoryKeys, " ")

	content, err := os.ReadFile(bbValuesFile)
	if err != nil {
		return fmt.Errorf("failed to read values file: %w", err)
	}

	values := make(map[string]interface{})
	if err := yaml.Unmarshal(content, &values); err != nil {
		return fmt.Errorf("failed to unmarshal values file: %w", err)
	}

	repositoryValues := make(map[string]interface{})
	addonsValues := make(map[string]interface{})
	for key, value := range values {
		if key == "addons" {
			addonValues, ok := value.(map[string]interface{})
			if !ok {
				return fmt.Errorf("type assertion failed for addons value, expected map[string]interface{}, got %T", value)
			}
			addonsValues = addonValues
		} else if contains(keys, key) {
			repositoryValues[key] = value
		} else {
			if err := writeValuesYAMLToFile(valuesDirectory, strings.ToLower(key), value); err != nil {
				return fmt.Errorf("failed to write values file for %s: %w", key, err)
			}

			if err := createBBSecretFiles(secretsDirectory, strings.ToLower(key)); err != nil {
				return fmt.Errorf("failed to write secret file for %s: %w", key, err)
			}
		}
	}

	for key, value := range addonsValues {
		addonsContent := map[string]interface{}{
			"addons": map[string]interface{}{
				key: value,
			},
		}

		if err := writeValuesYAMLToFile(valuesDirectory, strings.ToLower(key), addonsContent); err != nil {
			return fmt.Errorf("failed to write values file for %s: %w", key, err)
		}
		if err := createBBSecretFiles(secretsDirectory, strings.ToLower(key)); err != nil {
			return fmt.Errorf("failed to write secret file for %s: %w", key, err)
		}
	}

	if err := writeValuesYAMLToFile(valuesDirectory, "repo", repositoryValues); err != nil {
		return fmt.Errorf("failed to write repository.yaml file: %w", err)
	}

	if err := createBBSecretFiles(secretsDirectory, "repo"); err != nil {
		return fmt.Errorf("failed to write secret file for repo: %w", err)
	}

	if err := createBBSecretFiles(secretsDirectory, "repository-credentials"); err != nil {
		return fmt.Errorf("failed to write secret file for repository-credentials: %w", err)
	}

	return nil
}

// contains checks if a string is inside a slice.
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func writeValuesYAMLToFile(dir string, filename string, content interface{}) error {
	filePath := fmt.Sprintf("%s/%s.yaml", dir, filename)
	yamlData, err := yaml.Marshal(content)
	if err != nil {
		return fmt.Errorf("failed to marshal content: %w", err)
	}

	err = os.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	// // Save the marshaled content to a temporary file.
	// tmpFilePath := filePath + ".tmp"
	// if err = os.WriteFile(tmpFilePath, yamlData, 0644); err != nil {
	// 	return fmt.Errorf("failed to write temporary YAML file: %w", err)
	// }

	// // Use yq to process the temporary file and output to the final file path.
	// cmd := exec.Command("yq", "e", ".", "-i", tmpFilePath)
	// if err = cmd.Run(); err != nil {
	// 	return fmt.Errorf("failed to run yq: %w", err)
	// }

	// // Rename the temporary file to the final file name.
	// if err = os.Rename(tmpFilePath, filePath); err != nil {
	// 	return fmt.Errorf("failed to rename temporary YAML file to final path: %w", err)
	// }

	log.Printf("Created the BB Values File %s", filePath)
	return nil
}

// writeSecretsYAMLToFile writes the secrets content to the specified file as kubernetes secrets yaml file
func createBBSecretFiles(secretsDirectory string, key string) error {
	const tmpDir = "/tmp"
	component := strings.ToLower(key)
	switch key {
	case "repository-credentials":
		cmd := exec.Command("kubectl", "create", "secret", "generic", "repository-credentials", "--from-literal=username=dummy", "--from-literal=password=topsecret", "--dry-run=client", "-o", "yaml")
		output, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		secretFile := secretsDirectory + "/repository-credentials.enc.yaml"
		ioutil.WriteFile(secretFile, output, 0644)
		log.Printf("Created the BB Secret file: %s", secretFile)
		return nil
	default:
		secretInputFile := tmpDir + "/" + key + "-secret.yaml"
		ioutil.WriteFile(secretInputFile, []byte("{}"), 0644)
	}

	secretInputFile := tmpDir + "/" + key + "-secret.yaml"
	cmd := exec.Command("kubectl", "create", "secret", "generic", component+"-values", "--dry-run=client", "-o", "yaml", "--from-file", "values.yaml="+secretInputFile)
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	secretFile := secretsDirectory + "/" + component + ".enc.yaml"
	ioutil.WriteFile(secretFile, output, 0644)
	log.Printf("Created the BB Secret file: %s", secretFile)

	if err := os.Remove(secretInputFile); err != nil {
		return fmt.Errorf("failed to remove temporary secret file %s: %w", secretInputFile, err)
	}
	return nil
}
