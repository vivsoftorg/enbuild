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
	fmt.Printf("Downloaded file %s\n", filepath)
	return nil
}
