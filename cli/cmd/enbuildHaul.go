package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
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

// runCommand runs a shell command and returns its output or an error
func runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

// runPipedCommand runs a shell command with piping, combining multiple commands
func runPipedCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	// log.Println("Adding the Vivsoft Helm repo")
	_, err := runCommand("helm", "repo", "add", "vivsoft", "https://vivsoftorg.github.io/enbuild")
	if err != nil {
		log.Fatalf("Failed to add Helm repo: %v", err)
	}
	log.Println("ENBUILD Helm repo added successfully")

	_, err = runCommand("helm", "repo", "update")
	if err != nil {
		log.Fatalf("Failed to update Helm repo: %v", err)
	}
	log.Println("Helm repo updated successfully")

	// log.Println("Step 3: Getting the Helm chart version using Helm search and jq")
	versionJSON, err := runCommand("helm", "search", "repo", "vivsoft/enbuild", "-o", "json")
	if err != nil {
		log.Fatalf("Failed to search Helm chart: %v", err)
	}
	// log.Printf("Helm search output: %s", versionJSON)

	log.Println("Extracting chart version from JSON using jq")
	chartVersion, err := runPipedCommand(fmt.Sprintf("echo '%s' | jq -r '.[0].version'", versionJSON))
	if err != nil {
		log.Fatalf("Failed to extract chart version: %v", err)
	}
	log.Printf("Latest ENBUILD chart version is : %s", chartVersion)

	// log.Println("Step 4: Getting the list of images using Helm template and YQ")
	imagesOutput, err := runCommand("helm", "template", "vivsoft/enbuild", "--version", chartVersion)
	if err != nil {
		log.Fatalf("Failed to template Helm chart: %v", err)
	}

	log.Println("Extracting images using yq and sorting them")
	images, err := runPipedCommand(fmt.Sprintf("echo '%s' | yq -N '..|.image? | select(.)' | sort -u", imagesOutput))
	if err != nil {
		log.Fatalf("Failed to extract images: %v", err)
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
		log.Fatalf("Failed to parse template: %v", err)
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}
	log.Println("Template rendered successfully")

	// log.Println("Outputting the rendered YAML")
	// fmt.Println(rendered.String())

	outpurFilePath := "/tmp/enbuild-haul.yaml"

	if err = os.WriteFile(outpurFilePath, rendered.Bytes(), 0644); err != nil {
		log.Fatalf("Failed to write the YAML file: %v", err)
	}

	fmt.Printf("Haul file created successfully: %s\n", outpurFilePath)
	fmt.Printf("You can now run \n")
	fmt.Printf("hauler store sync -f %s\n", outpurFilePath)
	fmt.Printf("hauler store save --filename enbuild-%s-haul.tar.zst\n", chartVersion)
}
