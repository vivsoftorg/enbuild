package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	// fmt.Printf("Saved BigBang Images list file: %s\n", filepath)
	return nil
}

// Helper function to check if a slice contains a string
func contains(slice []string, str string) bool {
    for _, v := range slice {
        if v == str {
            return true
        }
    }
    return false
}
