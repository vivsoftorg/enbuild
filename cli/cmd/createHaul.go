package cmd

import (
	"bufio"
	"fmt"
	"log"
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
	log.Printf("Creating Haul Manifest file with BigBang version %s\n", bbVersion)

	targetDirectory := "target"
	if err := os.MkdirAll(targetDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDirectory, err)
	}

	bbImageListFile := fmt.Sprintf("%s/bigbang_images_list_%s.txt", targetDirectory, bbVersion)
	bbHelmListFile := fmt.Sprintf("%s/bigbang_helm_list_%s.txt", targetDirectory, bbVersion)
	bbInputList := fmt.Sprintf("%s/bigbang_merged_list_%s.txt", targetDirectory, bbVersion)
	bbHaulFile := fmt.Sprintf("%s/hauler_bb_images_%s.yaml", targetDirectory, bbVersion)

	_, err = os.Stat(bbImageListFile)
	if os.IsNotExist(err) {
		if err := downloadAndSaveFile(fmt.Sprintf("https://umbrella-bigbang-releases.s3-us-gov-west-1.amazonaws.com/umbrella/%s/images.txt", bbVersion), bbImageListFile); err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", bbImageListFile, err)
	} else {
		log.Printf("File %s already exists. Skipping download.\n", bbImageListFile)
	}

	_, err = os.Stat(bbHelmListFile) // this file is there only after BB version 2.25.0
	if os.IsNotExist(err) {
		if err := downloadAndSaveFile(fmt.Sprintf("https://umbrella-bigbang-releases.s3-us-gov-west-1.amazonaws.com/umbrella/%s/oci_package_list.txt", bbVersion), bbHelmListFile); err != nil {
			return err
		}
	} else if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", bbHelmListFile, err)
	} else {
		log.Printf("File %s already exists. Skipping download.\n", bbHelmListFile)
	}

	// merge the two files

	if err := mergeFiles(bbImageListFile, bbHelmListFile, bbInputList); err != nil {
		return err
	}

	return createHaulYaml(bbInputList, bbVersion, bbHaulFile)
}

func mergeFiles(bbImageListFile string, bbHelmListFile string, mergedFileName string) error {
	images, err := os.ReadFile(bbImageListFile)
	if err != nil {
		return fmt.Errorf("failed to read the images file: %w", err)
	}

	helms, err := os.ReadFile(bbHelmListFile)
	if err != nil {
		return fmt.Errorf("failed to read the helm file: %w", err)
	}

	imagesString := string(images)
	helmsString := string(helms)

	i_List := strings.Split(imagesString, "\n")
	// get only lines containing registry1.dso.mil from the images list
	var imagesList []string
	for _, image := range i_List {
		if strings.Contains(image, "registry1.dso.mil") {
			imagesList = append(imagesList, image)
		}
	}

	h_List := strings.Split(helmsString, "\n")

	// get only lines containing registry1.dso.mil from the helm list
	var helmsList []string
	for _, helm := range h_List {
		if strings.Contains(helm, "registry1.dso.mil") {
			helmsList = append(helmsList, helm)
		}
	}

	mergedList := append(imagesList, helmsList...)

	mergedFile := strings.Join(mergedList, "\n")

	if err := os.WriteFile(mergedFileName, []byte(mergedFile), 0644); err != nil {
		return fmt.Errorf("failed to write the merged file: %w", err)
	}
	log.Printf("Merged file created successfully: %s\n", mergedFileName)
	return nil
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
	fmt.Printf("You can now run \n")
	fmt.Printf("hauler login registry1.dso.mil -u <registry1_username> -p <registry1_password>\n")
	fmt.Printf("hauler store sync -f %s`\n", outpurFilePath)
	fmt.Printf("hauler store save --filename bb%s-haul.tar.zst\n", bbVersion)
	return nil
}

// func runHaulerCommands(haulFilePath string) error {
// 	log.Printf("Now You can run `hauler login registry1.dso.mil`\n")
// 	log.Printf("And then run `hauler store sync -f %s`\n", haulFilePath)

// 	// Run hauler login command
// 	loginCmd := exec.Command("hauler", "login", "registry1.dso.mil")
// 	loginCmd.Stdout = os.Stdout
// 	loginCmd.Stderr = os.Stderr
// 	if err := loginCmd.Run(); err != nil {
// 		return fmt.Errorf("failed to run hauler login command: %w", err)
// 	}

// 	// Run hauler store sync command
// 	syncCmd := exec.Command("hauler", "store", "sync", "-f", haulFilePath)
// 	syncCmd.Stdout = os.Stdout
// 	syncCmd.Stderr = os.Stderr
// 	if err := syncCmd.Run(); err != nil {
// 		return fmt.Errorf("failed to run hauler store sync command: %w", err)
// 	}

// 	return nil
// }
