package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
	"github.com/vivsoftorg/enbuild-sdk-go/pkg/manifests"
	"github.com/vivsoftorg/enbuild-sdk-go/pkg/types"
)

// getCatalogsCmd represents the catalogs command under get
var getCatalogsCmd = &cobra.Command{
	Use:   "catalogs",
	Short: "Get catalogs from Enbuild",
	Long:  `Retrieve and list catalogs from the Enbuild platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Fallback to ENV if flags not provided
		if token == "" {
			token = os.Getenv("ENBUILD_API_TOKEN")
		}
		if token == "" {
			log.Fatal("ENBUILD_API_TOKEN is required (via --token or env)")
		}
		if baseURL == "" {
			baseURL = os.Getenv("ENBUILD_BASE_URL")
		}

		options := []enbuild.ClientOption{
			enbuild.WithAuthToken(token),
		}
		if baseURL != "" {
			options = append(options, enbuild.WithBaseURL(baseURL))
		}
		if debug {
			options = append(options, enbuild.WithDebug(true))
		}

		client, err := enbuild.NewClient(options...)
		if err != nil {
			log.Fatalf("Error creating client: %v", err)
		}

		if idFlag != "" {
			fmt.Printf("Getting catalog with ID %s:\n", idFlag)
			manifest, err := client.Manifests.Get(idFlag, &manifests.ManifestListOptions{})
			if err != nil {
				log.Fatalf("Error getting catalog: %v", err)
			}
			printCatalogTable([]*types.Manifest{manifest})
			return
		}

		opts := &manifests.ManifestListOptions{
			VCS:  vcsFlag,
			Type: typeFlag,
			Name: nameFlag,
		}
		fmt.Println("Listing catalogs...")
		results, err := client.Manifests.List(opts)
		if err != nil {
			log.Fatalf("Error listing catalogs: %v", err)
		}
		printCatalogTable(results)
	},
}

func init() {
	getCmd.AddCommand(getCatalogsCmd)
	getCatalogsCmd.Flags().StringVar(&token, "token", "", "API token (or set ENBUILD_API_TOKEN)")
	getCatalogsCmd.Flags().StringVar(&baseURL, "base-url", "", "API base URL (or set ENBUILD_BASE_URL)")
	getCatalogsCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
	getCatalogsCmd.Flags().StringVar(&idFlag, "id", "", "Get catalog by ID")
	getCatalogsCmd.Flags().StringVar(&vcsFlag, "vcs", "", "Filter catalogs by VCS (e.g., github)")
	getCatalogsCmd.Flags().StringVar(&typeFlag, "type", "", "Filter catalogs by type (e.g., terraform)")
	getCatalogsCmd.Flags().StringVar(&nameFlag, "name", "", "Search catalogs by name")
}

func printCatalogTable(manifests []*types.Manifest) {
	fmt.Printf("Found %d catalog(s):\n", len(manifests))
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Name", "Type", "Slug", "VCS"})

	for _, m := range manifests {
		var id string
		if strID, ok := m.ID.(string); ok {
			id = strID
		} else {
			id = fmt.Sprintf("%v", m.ID)
		}
		table.Append([]string{id, m.Name, m.Type, m.Slug, m.VCS})
	}

	table.Render()
}
