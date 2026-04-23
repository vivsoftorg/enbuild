package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/vivsoftorg/enbuild-sdk-go/pkg/enbuild"
)

var deleteStackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Delete a stack from Enbuild",
	Long:  `Delete a stack from the Enbuild platform by ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		options := prepareClientOptions(baseURL, username, password)
		fmt.Printf("Using base URL: %s\n", baseURL)
		fmt.Printf("Using username: %s\n", username)

		client, err := enbuild.NewClient(context.Background(), options...)
		if err != nil {
			log.Fatalf("Error creating client: %v", err)
		}

		deleteStack(client, idFlag)
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
		}

		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteStackCmd)
	deleteStackCmd.Flags().StringVar(&idFlag, "id", "", "Stack ID to delete (required)")

	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("password")
	rootCmd.MarkFlagRequired("base-url")
}

func deleteStack(client *enbuild.Client, id string) {
	if id == "" {
		log.Fatal("Stack ID is required. Use --id flag.")
	}

	fmt.Printf("Deleting stack with ID: %s\n", id)
	err := client.Stacks.DeleteStack(context.Background(), id)
	if err != nil {
		log.Fatalf("Error deleting stack: %v", err)
	}
	fmt.Printf("Stack %s deleted successfully\n", id)
}