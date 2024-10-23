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

// Data holds the chart version and images list to be rendered in the template
type Data struct {
	ChartVersion string
	Images       []string
}

// createHaulCmd represents the createHaul command
var createENBUILDHaulCmd = &cobra.Command{
	Use:   `create-enbuild-haul`,
	Short: "Create a haul manifest file for ENBUILD Helm Chart",
	Long:  "Create a haul manifest.yaml file given the BigBang version.",
	RunE:  createENBUILDHaul,
}

func init() {
	rootCmd.AddCommand(createENBUILDHaulCmd)
}

func createENBUILDHaul(cmd *cobra.Command, args []string) error {
	// log.Println("Adding the Vivsoft Helm repo")
	_, err := runCommand("helm", "repo", "add", "vivsoft", "https://vivsoftorg.github.io/enbuild")
	if err != nil {
		return fmt.Errorf("failed to add Helm repo: %v", err)
	}
	log.Println("ENBUILD Helm repo added successfully")

	_, err = runCommand("helm", "repo", "update")
	if err != nil {
		return fmt.Errorf("failed to update Helm repo: %v", err)
	}
	log.Println("Helm repo updated successfully")

	versionJSON, err := runCommand("helm", "search", "repo", "vivsoft/enbuild", "-o", "json")
	if err != nil {
		return fmt.Errorf("failed to search Helm chart: %v", err)
	}

	log.Println("Extracting chart version from JSON using jq")
	chartVersion, err := runPipedCommand(fmt.Sprintf("echo '%s' | jq -r '.[0].version'", versionJSON))
	if err != nil {
		return fmt.Errorf("failed to extract chart version: %v", err)
	}
	log.Printf("Latest ENBUILD chart version is : %s", chartVersion)

	imagesOutput, err := runCommand("helm", "template", "vivsoft/enbuild", "--version", chartVersion)
	if err != nil {
		return fmt.Errorf("failed to template Helm chart: %v", err)
	}

	log.Println("Extracting images using yq and sorting them")
	images, err := runPipedCommand(fmt.Sprintf("echo '%s' | yq -N '..|.image? | select(.)' | sort -u", imagesOutput))
	if err != nil {
		return fmt.Errorf("failed to extract images: %v", err)
	}
	// log.Printf("Image list for ENBUILD Helm Chart: %s", images)

	log.Println("Preparing data for the HaulTemplate")
	data := Data{
		ChartVersion: chartVersion,
		Images:       strings.Split(images, "\n"),
	}
	// log.Println("Rendering the YAML template")
	tmpl, err := template.New("yaml").Parse(yamlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}
	log.Println("Template rendered successfully")

	// log.Println("Outputting the rendered YAML")
	// fmt.Println(rendered.String())

	targetDirectory := "target"
	if err := os.MkdirAll(targetDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDirectory, err)
	}

	enbuildHaulFile := fmt.Sprintf("%s/enbuild_%s_haul.yaml", targetDirectory, chartVersion)

	if err = os.WriteFile(enbuildHaulFile, rendered.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write the YAML file: %w", err)
	}

	fmt.Printf("Haul file created successfully: %s\n", enbuildHaulFile)
	fmt.Printf("You can now run \n")
	fmt.Printf("hauler store sync -f %s\n", enbuildHaulFile)
	fmt.Printf("hauler store save --filename enbuild-%s-haul.tar.zst\n", chartVersion)
	return nil
}
