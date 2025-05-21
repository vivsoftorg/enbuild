package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
)

// getCatalogsCmd represents the catalogs command under get
var getCatalogsCmd = &cobra.Command{
	Use:   "catalogs",
	Short: "Get catalogs from Enbuild",
	Long:  `Retrieve and list catalogs from the Enbuild platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		username := getEnvOrFlag(username, "username", "ENBUILD_USERNAME")
		password := getEnvOrFlag(password, "password", "ENBUILD_PASSWORD")
		baseURL := getEnvOrFlag(baseURL, "base-url", "ENBUILD_BASE_URL")

		fmt.Printf("Using ENBUILD_BASE_URL: %s\n", baseURL)
		fmt.Printf("Using ENBUILD_USERNAME: %s\n", username)
		fmt.Printf("Using ENBUILD_PASSWOR: %s\n", password)

		options := []enbuild.ClientOption{
			enbuild.WithDebug(debug),
			enbuild.WithBaseURL(baseURL),
			enbuild.WithKeycloakAuth(username, password),
		}

		client, err := enbuild.NewClient(options...)
		if err != nil {
			log.Fatalf("Error creating client: %v", err)
		}

		if idFlag != "" {
			processSingleCatalog(client, idFlag)
			return
		}

		listCatalogs(client)
	},
}

func init() {
	getCmd.AddCommand(getCatalogsCmd)
	getCatalogsCmd.Flags().StringVar(&idFlag, "id", "", "Get catalog by ID")
	getCatalogsCmd.Flags().StringVar(&vcsFlag, "vcs", "", "Filter catalogs by VCS (e.g., github)")
	getCatalogsCmd.Flags().StringVar(&typeFlag, "type", "", "Filter catalogs by type (e.g., terraform)")
	getCatalogsCmd.Flags().StringVar(&nameFlag, "name", "", "Search catalogs by name")
}

func getEnvOrFlag(flag, flagname string, env string) string {
	if flag == "" {
		flag = os.Getenv(env)
	}
	if flag == "" {
		log.Fatalf("%s is required via flag or set via env variable %s", flagname, env)
	}
	return flag
}

func processSingleCatalog(client *enbuild.Client, id string) {
	log.Printf("ENBUILD_BASE_URL: %s\n", baseURL)
	fmt.Printf("Getting catalog with ID %s:\n", id)
	catalog, err := client.Catalogs.Get(id, &enbuild.CatalogListOptions{})
	if err != nil {
		log.Fatalf("Error getting catalog: %v", err)
	}
	printCatalogTable([]*enbuild.Catalog{catalog})
}

func listCatalogs(client *enbuild.Client) {
	opts := &enbuild.CatalogListOptions{
		VCS:  vcsFlag,
		Type: typeFlag,
		Name: nameFlag,
	}
	log.Printf("ENBUILD_BASE_URL: %s\n", baseURL)
	fmt.Println("Listing catalogs...")
	results, err := client.Catalogs.List(opts)
	if err != nil {
		log.Fatalf("Error listing catalogs: %v", err)
	}
	printCatalogTable(results)
}

func printCatalogTable(catalogs []*enbuild.Catalog) {
	fmt.Printf("Found %d catalog(s):\n", len(catalogs))
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Name", "Type", "Slug", "VCS"})

	for _, c := range catalogs {
		id := fmt.Sprintf("%v", c.ID) // Safe conversion; assumes no panics for valid IDs.
		table.Append([]string{id, c.Name, c.Type, c.Slug, c.VCS})
	}
	table.Render()
}
