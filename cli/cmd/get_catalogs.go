package cmd

import (
    "context"
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
        options := prepareClientOptions(baseURL, username, password)
		fmt.Printf("Using base URL: %s\n", baseURL)
		fmt.Printf("Using username: %s\n", username)

        client, err := enbuild.NewClient(context.Background(), options...)
        if err != nil {
            log.Fatalf("Error creating client: %v", err)
        }

        switch {
        case idFlag != "":
            processSingleCatalog(client, idFlag)
        default:
            listCatalogs(client)
        }
    },
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if username == "" {
			username = os.Getenv("ENBUILD_USERNAME")
		}
		if username == "" {
			return fmt.Errorf("username is required, set with --username or ENBUILD_USERNAME env var")
		}

		if password == "" {
			password = os.Getenv("ENBUILD_PASSWORD")
		}
		if password == "" {
			return fmt.Errorf("password is required, set with --password or ENBUILD_PASSWORD env var")
		}

		if baseURL == "" {
			baseURL = os.Getenv("ENBUILD_BASE_URL")
		}
		if baseURL == "" {
			baseURL = "https://enbuild.vivplatform.io"
			// return fmt.Errorf("base-url is required, set with --base-url or ENBUILD_BASE_URL env var")
		}

		return nil
	},
}

func init() {
    getCmd.AddCommand(getCatalogsCmd)
    getCatalogsCmd.Flags().StringVar(&idFlag, "id", "", "Get catalog by ID")
    getCatalogsCmd.Flags().StringVar(&vcsFlag, "vcs", "", "Filter catalogs by VCS (e.g., github)")
    getCatalogsCmd.Flags().StringVar(&typeFlag, "type", "", "Filter catalogs by type (e.g., terraform)")
    getCatalogsCmd.Flags().StringVar(&nameFlag, "name", "", "Search catalogs by name")

	rootCmd.MarkFlagRequired("username")
    rootCmd.MarkFlagRequired("password")
    rootCmd.MarkFlagRequired("base-url")
}

func getEnvOrFlag(flag, flagname string, env string) string {
    value := flag
    if value == "" {
        value = os.Getenv(env)
    }
    if value == "" {
        log.Fatalf("%s is required via flag or set via env variable %s", flagname, env)
    }
    return value
}

func prepareClientOptions(baseURL, username, password string) []enbuild.ClientOption {
    return []enbuild.ClientOption{
        enbuild.WithDebug(debug),
        enbuild.WithBaseURL(baseURL),
        enbuild.WithKeycloakAuth(username, password),
    }
}

func processSingleCatalog(client *enbuild.Client, id string) {
    fmt.Printf("Getting catalog with ID %s:\n", id)
    catalog, err := client.Catalogs.GetCatalog(context.Background(), id, &enbuild.CatalogListOptions{})
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
    fmt.Println("Listing catalogs...")
    results, err := client.Catalogs.ListCatalog(context.Background(), opts)
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
        id := fmt.Sprintf("%v", c.ID)
        table.Append([]string{id, c.Name, c.Type, c.Slug, c.VCS})
    }
    table.Render()
}
