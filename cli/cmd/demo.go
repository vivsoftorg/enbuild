package cmd

import (
	_ "embed"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	upClusterName string
	upDebug       bool
)

//go:embed demo_scripts/create_enbuild_demo.sh
var demoScriptsCreate []byte

//go:embed demo_scripts/destroy_enbuild_demo.sh
var demoScriptsDestroy []byte

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Try Enbuild on your local machine",
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create a k3s kubernetes cluster with ENBUILD installed on your local machine",
	Run: func(cmd *cobra.Command, args []string) {
		// Handle the 'up' action
		clusterNameArg := upClusterName
		debugArg := "false"
		if upDebug {
			debugArg = "true"
		}
		scriptPath := WriteInFile("create_enbuild_demo.sh", demoScriptsCreate)
		execCmd := exec.Command("sh", scriptPath, clusterNameArg, debugArg)
		output, err := execCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error executing create script: %v\n", err)
			fmt.Printf("Output: %s\n", string(output))
			return
		}

		fmt.Println(string(output))
		DeleteFile(scriptPath)
	},
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Remove k3s cluster with ENBUILD installed on your local machine",
	Run: func(cmd *cobra.Command, args []string) {
		scriptPath := WriteInFile("destroy_enbuild_demo.sh", demoScriptsDestroy)
		execCmd := exec.Command("sh", scriptPath)
		output, err := execCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error executing destroy script: %v\n", err)
			fmt.Printf("Output: %s\n", string(output))
			return
		}

		fmt.Println(string(output))
		DeleteFile(scriptPath)
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)

	// Define flags for the 'up' subcommand
	upCmd.Flags().StringVar(&upClusterName, "clusterName", "enbuild", "Name of the cluster")
	upCmd.Flags().BoolVar(&upDebug, "debug", false, "Enable debug mode")

	// Add the 'up' and 'destroy' subcommands to 'demo'
	demoCmd.AddCommand(upCmd)
	demoCmd.AddCommand(destroyCmd)
}
