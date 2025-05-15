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


var manifestCmd = &cobra.Command{
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
            fmt.Printf("Getting manifest with ID %s:\n", idFlag)
            manifest, err := client.Manifests.Get(idFlag, &manifests.ManifestListOptions{})
            if err != nil {
                log.Fatalf("Error getting manifest: %v", err)
            }
            printManifestTable([]*types.Manifest{manifest})
            return
        }

        opts := &manifests.ManifestListOptions{
            VCS:  vcsFlag,
            Type: typeFlag,
            Name: nameFlag,
        }
        fmt.Println("Listing manifests...")
        results, err := client.Manifests.List(opts)
        if err != nil {
            log.Fatalf("Error listing manifests: %v", err)
        }
        printManifestTable(results)
    },
}

func init() {
    rootCmd.AddCommand(manifestCmd)
    manifestCmd.Flags().StringVar(&token, "token", "", "API token (or set ENBUILD_API_TOKEN)")
    manifestCmd.Flags().StringVar(&baseURL, "base-url", "", "API base URL (or set ENBUILD_BASE_URL)")
    manifestCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
    manifestCmd.Flags().StringVar(&idFlag, "id", "", "Get manifest by ID")
    manifestCmd.Flags().StringVar(&vcsFlag, "vcs", "", "Filter manifests by VCS (e.g., github)")
    manifestCmd.Flags().StringVar(&typeFlag, "type", "", "Filter manifests by type (e.g., terraform)")
    manifestCmd.Flags().StringVar(&nameFlag, "name", "", "Search manifests by name")
}

func printManifestTable(manifests []*types.Manifest) {
    fmt.Printf("Found %d manifest(s):\n", len(manifests))
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
