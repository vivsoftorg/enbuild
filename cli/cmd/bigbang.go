package cmd

import (
	// "fmt"

	"github.com/spf13/cobra"
)

// bigbangCmd represents the bigbang command
var bigbangCmd = &cobra.Command{
	Use:   "bigbang",
	Short: "bigbang",
	Long: `bigbang is a CLI library for enbuild that helps creating bigbang template.
This also helps to create the hauler yaml and haul for offline deployment of bigbang`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("Available subcommands for enbuild bigbang:\n  - create-haul")
	// },
}

func init() {
	rootCmd.AddCommand(bigbangCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bigbangCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bigbangCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
