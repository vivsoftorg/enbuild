package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Image struct {
	Name string `yaml:"name"`
}

type ImageSpec struct {
	Images []Image `yaml:"images"`
}

type ImagesYaml struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   Metadata  `yaml:"metadata"`
	Spec       ImageSpec `yaml:"spec"`
}

// Metadata structure if needed to be reused or extended.
type Metadata struct {
	Name string `yaml:"name"`
}

// createHaulCmd represents the createHaul command
var createHaulCmd = &cobra.Command{
	Use:   "create-haul",
	Short: "Create a haul manifest file",
	Long:  "Create a haul manifest.yaml file given the BigBang version.",
	RunE:  createHaul,
}

func init() {
	bigbangCmd.AddCommand(createHaulCmd)
	createHaulCmd.Flags().StringP("bb-version", "v", "", "Specify the BigBang version (required)")
	_ = createHaulCmd.MarkFlagRequired("bb-version")
}

func createHaul(cmd *cobra.Command, args []string) error {
	bbVersion, err := cmd.Flags().GetString("bb-version")
	if err != nil {
		return fmt.Errorf("failed to read 'bb-version' flag: %w", err)
	}
	fmt.Printf("Creating Haul Manifest file with BigBang version %s\n", bbVersion)

	targetDirectory := "target"
	if err := os.MkdirAll(targetDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDirectory, err)
	}

	bbImageListFile := fmt.Sprintf("%s/bigbang_images_list_%s.txt", targetDirectory, bbVersion)
	bbHaulFile := fmt.Sprintf("%s/hauler_bb_images_%s.yaml", targetDirectory, bbVersion)
	if err := downloadAndSaveFile(fmt.Sprintf("https://umbrella-bigbang-releases.s3-us-gov-west-1.amazonaws.com/umbrella/%s/images.txt", bbVersion), bbImageListFile); err != nil {
		return err
	}

	return createHaulYaml(bbImageListFile, bbVersion, bbHaulFile)
}


func createHaulYaml(inputFilePath string, bbVersion string, outpurFilePath string) error {
	file, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("failed to open the input file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var images []string
	for scanner.Scan() {
		if image := strings.TrimSpace(scanner.Text()); image != "" {
			images = append(images, image)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading the input file: %w", err)
	}

	imagesYaml := ImagesYaml{
		APIVersion: "content.hauler.cattle.io/v1alpha1",
		Kind:       "Images",
		Metadata:   Metadata{Name: "bigbang-images-" + bbVersion},
		Spec:       ImageSpec{Images: make([]Image, len(images))},
	}

	for i, image := range images {
		imagesYaml.Spec.Images[i].Name = image
	}

	yamlData, err := yaml.Marshal(&imagesYaml)
	if err != nil {
		return fmt.Errorf("failed to marshal images to YAML: %w", err)
	}

	if err = os.WriteFile(outpurFilePath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write the YAML file: %w", err)
	}

	fmt.Printf("Haul file created successfully: %s\n", outpurFilePath)
	fmt.Printf("Now You can run `hauler login registry1.dso.mil`\n")
	fmt.Printf("And then run `hauler store sync -f %s`\n", outpurFilePath)
	return nil
}
