package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

func downloadAndSaveFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to save file %s: %w", filepath, err)
	}
	fmt.Printf("Downloaded file %s\n", filepath)
	return nil
}

func WriteInFile(fileName string, content []byte) string {
	fullPath := "/tmp/enbuild/"
	if runtime.GOOS == "windows" {
		fullPath = "C:\\Users\\Default\\AppData\\Local\\Temp\\enbuild\\"
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err := os.Mkdir(fullPath, 0777)
		if err != nil {
			log.Fatalf("Couldn't create folder : " + err.Error())
			os.Exit(1)
			panic("unreachable") // staticcheck false positive: https://staticcheck.io/docs/checks#SA5011
		}
	}

	err := os.WriteFile(fullPath+fileName, content, 0777)
	if err != nil {
		log.Fatalf("Couldn't write file : " + err.Error())
		return ""
	}

	return fullPath + fileName
}

func DeleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Error(err)
	}
}

func DeleteFolder(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Error(err)
	}
}

// runCommand runs a shell command and returns its output or an error
func runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

// runPipedCommand runs a shell command with piping, combining multiple commands
func runPipedCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
