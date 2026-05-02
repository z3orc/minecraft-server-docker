package download

import (
	"fmt"
	"os/exec"
)

func ServerExecutable(url string, destDir string, jarName string) error {
	cmd := exec.Command("wget", url, "-O", jarName)
	cmd.Dir = destDir
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to download server jar from %s: %s", url, err)
	}

	return nil
}
