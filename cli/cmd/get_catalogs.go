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
        options := prepareClientOptions(baseURL, username, password)
		fmt.Printf("Using base URL: %s\n", baseURL)
		fmt.Printf("Using username: %s\n", username)

        client, err := enbuild.NewClient(options...)
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
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        baseURL = getEnvOrFlag(baseURL, "base-url", "ENBUILD_BASE_URL")
        username = getEnvOrFlag(username, "username", "ENBUILD_USERNAME")
        password = getEnvOrFlag(password, "password", "ENBUILD_PASSWORD")
        return nil // Return nil as it completes successfully
    },
}

func init() {
    getCmd.AddCommand(getCatalogsCmd) // Assuming getCmd is defined elsewhere
    getCatalogsCmd.Flags().StringVar(&idFlag, "id", "", "Get catalog by ID")
    getCatalogsCmd.Flags().StringVar(&vcsFlag, "vcs", "", "Filter catalogs by VCS (e.g., github)")
    getCatalogsCmd.Flags().StringVar(&typeFlag, "type", "", "Filter catalogs by type (e.g., terraform)")
    getCatalogsCmd.Flags().StringVar(&nameFlag, "name", "", "Search catalogs by name")
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
        id := fmt.Sprintf("%v", c.ID)
        table.Append([]string{id, c.Name, c.Type, c.Slug, c.VCS})
    }
    table.Render()
}
