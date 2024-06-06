package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	targetDirectory      = "target"
	valuesDirectoryName  = "bb_values"
	secretsDirectoryName = "bb_secrets"
	repositoryKeys       = "domain offline helmRepositories registryCredentials openshift git sso flux networkPolicies imagePullPolicy wrapper packages"
	sourceType           = "helmRepo" // Default sourceType is "git" in BigBang , but we want helmrepo
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
	repo_keys := strings.Split(repositoryKeys, " ")
	values_kustomizationFile := fmt.Sprintf("%s/kustomization.yaml", valuesDirectory)
	secrets_kustomizationFile := fmt.Sprintf("%s/kustomization.yaml", secretsDirectory)
	var kustomization_keys []string

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
		} else if contains(repo_keys, key) {
			repositoryValues[key] = value
		} else {
			// valueContent := map[string]interface{}{
			// 	key: value,
			// }
			// if err := writeValuesYAMLToFile(valuesDirectory, strings.ToLower(key+"_nocomments"), valueContent); err != nil {
			// 	return fmt.Errorf("failed to write values file for %s: %w", key, err)
			// }
			filePath := fmt.Sprintf("%s/%s.yaml", valuesDirectory, strings.ToLower(key))
			c := fmt.Sprintf(
				"yq 'with(.%s.sourceType; . = \"%s\" | . style=\"double\") | .%s | {\"%s\" : . }' %s > %s",
				key, sourceType, key, key, bbValuesFile, filePath,
			)
			cmd := exec.Command("sh", "-c", c)
			if err := cmd.Run(); err != nil {
				log.Fatalf("Failed to run yq command: %v", err)
			}

			log.Printf("Created the BB Values File %s", filePath)

			if err := createBBSecretFiles(secretsDirectory, strings.ToLower(key)); err != nil {
				return fmt.Errorf("failed to write secret file for %s: %w", key, err)
			}

			kustomization_keys = append(kustomization_keys, key)
		}
	}

	// for addon_key, value := range addonsValues {
	for addon_key := range addonsValues {
		// addonsContent := map[string]interface{}{
		// 	"addons": map[string]interface{}{
		// 		addon_key: value,
		// 	},
		// }
		// if err := writeValuesYAMLToFile(valuesDirectory, strings.ToLower(addon_key+"_nocomments"), addonsContent); err != nil {
		// 	return fmt.Errorf("failed to write values file for %s: %w", addon_key, err)
		// }

		filePath := fmt.Sprintf("%s/%s.yaml", valuesDirectory, strings.ToLower(addon_key))
		c := fmt.Sprintf(
			"yq 'with(.addons.%s.sourceType; . = \"%s\" | . style=\"double\") | .addons.%s | {\"addons\": {\"%s\" : . }}' %s > %s",
			addon_key, sourceType, addon_key, addon_key, bbValuesFile, filePath,
		)
		cmd := exec.Command("sh", "-c", c)
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to run yq command: %v", err)
		}

		log.Printf("Created the BB Values File %s", filePath)

		if err := createBBSecretFiles(secretsDirectory, strings.ToLower(addon_key)); err != nil {
			return fmt.Errorf("failed to write secret file for %s: %w", addon_key, err)
		}
		kustomization_keys = append(kustomization_keys, addon_key)
	}
	// -------------------------------------------------------------------------------------------------------------------------------------------------
	// Create the repo.yaml file

	// if err := writeValuesYAMLToFile(valuesDirectory, "repo_nocomments", repositoryValues); err != nil {
	// 	return fmt.Errorf("failed to write repo.yaml file: %w", err)
	// }

	keys := strings.Join(repo_keys, `","`)
	filePath := fmt.Sprintf("%s/%s.yaml", valuesDirectory, "repo")
	c := fmt.Sprintf(`yq '. |= pick(["%s"])' %s > %s`, keys, bbValuesFile, filePath)
	cmd := exec.Command("sh", "-c", c)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run yq command: %v", err)
	}

	// // update the .helmRepositories key with default values
	// c = fmt.Sprintf("yq eval '.helmRepositories += [{\"name\": \"registry1\", \"repository\": \"oci://registry1.dso.mil/bigbang\", \"existingSecret\": \"private-registry\", \"type\": \"oci\"}]' -i %s", filePath)
	// cmd = exec.Command("sh", "-c", c)
	// if err := cmd.Run(); err != nil {
	// 	log.Fatalf("Failed to run yq command: %v", err)
	// }

	log.Printf("Created the BB Values File %s", filePath)

	if err := createBBSecretFiles(secretsDirectory, "repo"); err != nil {
		return fmt.Errorf("failed to write secret file for repo: %w", err)
	}

	kustomization_keys = append(kustomization_keys, "repo")
	// -------------------------------------------------------------------------------------------------------------------------------------------------
	if err := createBBSecretFiles(secretsDirectory, "repository-credentials"); err != nil {
		return fmt.Errorf("failed to write secret file for repository-credentials: %w", err)
	}
	// -------------------------------------------------------------------------------------------------------------------------------------------------
	// log.Printf("Creating values kustomization file %s\n", values_kustomizationFile)
	file, err := os.Create(values_kustomizationFile)
	if err != nil {
		return fmt.Errorf("Error creating kustomization file for values at %s: %w", values_kustomizationFile, err)
	}
	defer file.Close() // Ensure the file is closed after writing

	// Write the fixed part of the YAML content
	header := `---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: bigbang
generatorOptions:
  disableNameSuffixHash: true
  labels:
    app.kubernetes.io/part-of: bigbang
configMapGenerator:
- name: bb-helm-values
  files:
`

	// Write the header to the file
	_, err = file.WriteString(header)
	if err != nil {
		return fmt.Errorf("failed to write kustomization file for values at %s: %w", values_kustomizationFile, err)
	}

	// Write each sorted item to the file
	sort.Strings(kustomization_keys)
	for _, item := range kustomization_keys {
		line := fmt.Sprintf("    - %s.yaml\n", strings.ToLower(item))
		_, err = file.WriteString(line)
		if err != nil {
			return fmt.Errorf("failed to write kustomization file for values at %s: %w", values_kustomizationFile, err)
		}
	}

	log.Printf("Created values kustomization file %s successfully.\n", values_kustomizationFile)
	// -------------------------------------------------------------------------------------------------------------------------------------------------
	// fmt.Printf("Creating secrets kustomization file %s\n", secrets_kustomizationFile)
	file, err = os.Create(secrets_kustomizationFile)
	if err != nil {
		return fmt.Errorf("Error creating kustomization file for secrets at %s: %w", secrets_kustomizationFile, err)
	}
	defer file.Close() // Ensure the file is closed after writing

	// Write the fixed part of the YAML content
	secrets_header := `---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: bigbang
resources:
`

	// Write the header to the file
	_, err = file.WriteString(secrets_header)
	if err != nil {
		return fmt.Errorf("failed to write kustomization file for secrets at %s: %w", secrets_kustomizationFile, err)
	}

	// Write each sorted item to the file
	sort.Strings(kustomization_keys)
	for _, item := range kustomization_keys {
		line := fmt.Sprintf("  - %s.enc.yaml\n", strings.ToLower(item))
		_, err = file.WriteString(line)
		if err != nil {
			return fmt.Errorf("failed to write kustomization file for secrets at %s: %w", secrets_kustomizationFile, err)
		}
	}

	log.Printf("Created secrets kustomization file %s successfully.\n", secrets_kustomizationFile)
	// -------------------------------------------------------------------------------------------------------------------------------------------------
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
	contentMap, ok := content.(map[string]interface{})
	if !ok {
		return fmt.Errorf("content is not of type map[string]interface{}")
	}

	updateSourceType(contentMap)

	filePath := fmt.Sprintf("%s/%s.yaml", dir, filename)
	yamlData, err := yaml.Marshal(contentMap)
	if err != nil {
		return fmt.Errorf("failed to marshal content: %w", err)
	}

	//use yq to format the yamlData
	cmd := exec.Command("yq", "eval", ".", "-")
	cmd.Stdin = strings.NewReader(string(yamlData))
	yamlData, err = cmd.Output()

	err = os.WriteFile(filePath, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	log.Printf("Created the BB Values File %s", filePath)
	return nil
}

// updateSourceType recursively searches for the key "sourceType" and updates its value.
func updateSourceType(data map[string]interface{}) {
	for key, val := range data {
		if key == "sourceType" {
			data[key] = sourceType
			// } else if key == "helmRepositories" {
			// 	repositories, ok := val.([]interface{})
			// 	if ok {
			// 		repository := map[string]interface{}{
			// 			"name": "registry1",
			// 			"repository": "oci://registry1.dso.mil/bigbang",
			// 			"type": "oci",
			// 			"username": "demo",
			// 			"password": "demo",
			// 		}
			// 		repositories = append(repositories, repository)
			// 		data[key] = repositories
			// 	}
		} else if subMap, ok := val.(map[string]interface{}); ok {
			updateSourceType(subMap) // Recurse into nested maps
		}
	}
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
