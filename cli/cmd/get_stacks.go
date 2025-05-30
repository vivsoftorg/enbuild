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

// getStacksCmd represents the stacks command under get
var getStacksCmd = &cobra.Command{
	Use:   "stacks",
	Short: "Get stacks from Enbuild",
	Long:  `Retrieve and list stacks from the Enbuild platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		options := prepareClientOptions(baseURL, username, password)
		fmt.Printf("Using base URL: %s\n", baseURL)
		fmt.Printf("Using username: %s\n", username)

		client, err := enbuild.NewClient(context.Background(), options...)
		if err != nil {
			log.Fatalf("Error creating client: %v", err)
		}

		listStacks(client)
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
	getCmd.AddCommand(getStacksCmd)
	getStacksCmd.Flags().StringVar(&nameFlag, "name", "", "Search stacks by name")

	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("password")
	rootCmd.MarkFlagRequired("base-url")
}

func listStacks(client *enbuild.Client) {
	page := 0
	limit := 100
	fmt.Println("Listing stacks...")
	results, err := client.Stacks.ListStacks(context.Background(), page, limit, nameFlag)
	if err != nil {
		log.Fatalf("Error listing stacks: %v", err)
	}
	printStackTable(results)
}

func printStackTable(stacks []*enbuild.Stack) {
	fmt.Printf("Found %d stack(s):\n", len(stacks))
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Name", "Type", "Status"})

	for _, c := range stacks {
		id := fmt.Sprintf("%v", c.ID)
		table.Append([]string{id, c.Name, c.Type, c.Status})
	}
	table.Render()
}
