package cmd

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"runtime"
)

var (
	upClusterName string
	upDebug       bool
	valueFilePath string
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
		valueFilePath = "/tmp/enbuild/values.yaml"
		if runtime.GOOS == "windows" {
			valueFilePath = "C:\\Users\\Default\\AppData\\Local\\Temp\\enbuild\\values.yaml"
		}
		scriptPath := WriteInFile("create_enbuild_demo.sh", demoScriptsCreate)
		execCmd := exec.Command("sh", scriptPath, clusterNameArg, debugArg, valueFilePath)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		if err := execCmd.Run(); err != nil || !execCmd.ProcessState.Success() {
			fmt.Errorf("error executing the command %s", err)
			return
		}
		DeleteFile(scriptPath)
	},
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Remove k3s cluster with ENBUILD installed on your local machine",
	Run: func(cmd *cobra.Command, args []string) {
		scriptPath := WriteInFile("destroy_enbuild_demo.sh", demoScriptsDestroy)
		execCmd := exec.Command("sh", scriptPath)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		if err := execCmd.Run(); err != nil || !execCmd.ProcessState.Success() {
			fmt.Errorf("error executing the command %s", err)
			return
		}
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
