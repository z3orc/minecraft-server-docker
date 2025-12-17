package jar

import (
	"fmt"
	"os/exec"
)

func DownloadServerJar(url string, destDir string) error {
	cmd := exec.Command("wget", url, "-O server.jar")
	cmd.Dir = destDir
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to download server jar from %s: %e", url, err)
	}

	return nil
}
