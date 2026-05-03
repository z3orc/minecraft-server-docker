package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// ServerExecutable downloads the file found at the provided url. File is stored on disk at:
// destDir/jarName.
//
// Returns an error if the request fails, the server responds with a non-200 status or
// the file cannot be written,
func ServerExecutable(url string, destDir string, jarName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download server jar from %s: %s", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d downloading server jar from %s", resp.StatusCode, url)
	}

	dest := filepath.Join(destDir, jarName)
	file, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %s", dest, err)
	}
	defer file.Close()

	reader := io.Reader(resp.Body)
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to write server jar to %s: %s", dest, err)
	}

	return nil
}
