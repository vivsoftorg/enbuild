package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// Template for the YAML file
const yamlTemplate = `
apiVersion: content.hauler.cattle.io/v1alpha1
kind: Charts
metadata:
  name: enbuild-chart-hauler
spec:
  charts:
    - name: enbuild
      repoURL: https://vivsoftorg.github.io/enbuild
      version: {{.ChartVersion}}
---
apiVersion: content.hauler.cattle.io/v1alpha1
kind: Images
metadata:
  name: enbuild-images-hauler 
spec:
  images:
    {{range .Images}}- name: {{.}}
      platform: linux/amd64
    {{end}}
`

// ChartData holds the chart version and image list for the template
type ChartData struct {
	ChartVersion string
	Images       []string
}

// createHaulEnbuildCmd represents the enbuild command under create haul
var createHaulEnbuildCmd = &cobra.Command{
	Use:   "enbuild",
	Short: "Create a haul manifest file for the ENBUILD Helm Chart and images",
	Long:  "Create a haul manifest.yaml file given the ENBUILD Helm Chart version.",
	RunE:  runCreateENBUILDHaul,
}

func init() {
	createHaulCmd.AddCommand(createHaulEnbuildCmd)
	createHaulEnbuildCmd.Flags().StringP("helm-chart-version", "v", "", "Specify the ENBUILD Helm Chart version")
}

func getHelmChartVersion() (string, error) {
	if _, err := runCommand("helm", "repo", "add", "vivsoft", "https://vivsoftorg.github.io/enbuild"); err != nil {
		return "", fmt.Errorf("failed to add Helm repo: %v", err)
	}
	log.Println("ENBUILD Helm repo added successfully")

	if _, err := runCommand("helm", "repo", "update"); err != nil {
		return "", fmt.Errorf("failed to update Helm repo: %v", err)
	}
	log.Println("Helm repo updated successfully")

	versionJSON, err := runCommand("helm", "search", "repo", "vivsoft/enbuild", "-o", "json")
	if err != nil {
		return "", fmt.Errorf("failed to search Helm chart: %v", err)
	}

	log.Println("Extracting chart version from JSON using jq")
	chartVersion, err := runPipedCommand(fmt.Sprintf("echo '%s' | jq -r '.[0].version'", versionJSON))
	if err != nil {
		return "", fmt.Errorf("failed to extract chart version: %v", err)
	}
	log.Printf("Latest ENBUILD chart version is: %s", chartVersion)
	return chartVersion, nil
}

func runCreateENBUILDHaul(cmd *cobra.Command, args []string) error {
	// Get the version flag value
	chartVersion, err := cmd.Flags().GetString("helm-chart-version")
	if err != nil {
		return fmt.Errorf("failed to read 'helm-chart-version' flag: %w", err)
	}

	// If version is not provided, fetch the latest version
	if chartVersion == "" {
		chartVersion, err = getHelmChartVersion()
		if err != nil {
			return err
		}
	}

	if _, err := runCommand("helm", "repo", "update"); err != nil {
		return fmt.Errorf("failed to update Helm repo: %v", err)
	}
	log.Println("Helm repo updated successfully")
	// Fetch images from the Helm chart
	imagesOutput, err := runCommand("helm", "template", "vivsoft/enbuild", "--version", chartVersion)
	if err != nil {
		return fmt.Errorf("failed to template Helm chart %s: %v", chartVersion, err)
	}

	log.Println("Extracting images using yq and sorting them")
	images, err := runPipedCommand(fmt.Sprintf("echo '%s' | yq -N '..|.image? | select(.)' | sort -u", imagesOutput))
	if err != nil {
		return fmt.Errorf("failed to extract images: %v", err)
	}

	// Prepare data for the template
	data := ChartData{
		ChartVersion: chartVersion,
		Images:       strings.Split(images, "\n"),
	}

	// Render the YAML template
	tmpl, err := template.New("yaml").Parse(yamlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}
	log.Println("Template rendered successfully")

	// Create the target directory and write the output file
	targetDir := "target"
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	outputFile := fmt.Sprintf("%s/enbuild_%s_haul.yaml", targetDir, chartVersion)
	if err = os.WriteFile(outputFile, rendered.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write the YAML file: %w", err)
	}

	fmt.Printf("Haul file created successfully: %s\n", outputFile)
	fmt.Printf("You can now run:\n")
	fmt.Printf("hauler store sync -f %s\n", outputFile)
	fmt.Printf("hauler store save --filename enbuild-%s-haul.tar.zst\n", chartVersion)
	return nil
}
