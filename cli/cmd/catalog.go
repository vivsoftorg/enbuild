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

var (
    baseURL  string
    debug    bool
    token    string
    idFlag   string
    vcsFlag  string
    typeFlag string
    nameFlag string
)


var catalogCmd = &cobra.Command{
    Use:   "catalog",
    Short: "CLI to interact with Enbuild catalogs",
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
            log.Printf("Getting catalog from ENBUILD %s", baseURL)
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
        log.Printf("Listing catalogs from ENBUILD %s", baseURL)
        results, err := client.Manifests.List(opts)
        if err != nil {
            log.Fatalf("Error listing catalogs: %v", err)
        }
        printCatalogTable(results)
    },
}

func init() {
    rootCmd.AddCommand(catalogCmd)
    catalogCmd.Flags().StringVar(&token, "token", "", "API token (or set ENBUILD_API_TOKEN)")
    catalogCmd.Flags().StringVar(&baseURL, "base-url", "", "API base URL (or set ENBUILD_BASE_URL)")
    catalogCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
    catalogCmd.Flags().StringVar(&idFlag, "id", "", "Get catalog by ID")
    catalogCmd.Flags().StringVar(&vcsFlag, "vcs", "", "Filter catalogs by VCS (e.g., github)")
    catalogCmd.Flags().StringVar(&typeFlag, "type", "", "Filter catalogs by type (e.g., terraform)")
    catalogCmd.Flags().StringVar(&nameFlag, "name", "", "Search catalogs by name")
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
