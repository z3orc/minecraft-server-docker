package fabric

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Returns the version of the latest compatible fabric loader based on
// provided game version 'gameVersion'
func findLatestCompatibleLoader(gameVersion string) (string, error) {
	url := fmt.Sprintf("%s/v2/versions/loader/%s", FABRIC_API_BASE_URL, gameVersion)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get list of loaders from fabric api: %e", err)
	}
	defer resp.Body.Close()

	loaderResp := loaderResponse{}
	err = json.NewDecoder(resp.Body).Decode(&loaderResp)
	if err != nil {
		return "", fmt.Errorf("failed decode list of loaders fabric api: %e", err)
	}

	if len(loaderResp) <= 0 {
		return "", fmt.Errorf("list of loaders is empty. expected one or more")
	}

	return loaderResp[0].Loader.Version, nil
}

// Returns the version number of the latest installer version.
func findLatestInstaller() (string, error) {
	url := fmt.Sprintf("%s/v2/versions/installer", FABRIC_API_BASE_URL)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get list installer versions from fabric api: %e", err)
	}
	defer resp.Body.Close()

	installers := installerResponse{}
	err = json.NewDecoder(resp.Body).Decode(&installers)
	if err != nil {
		return "", fmt.Errorf("failed decode list of installer versions from fabric api: %e", err)
	}

	if len(installers) <= 0 {
		return "", fmt.Errorf("list of installers is empty. expected one or more")
	}

	return installers[0].Version, nil
}

func GetDownloadUrl(gameVersion string) (string, error) {
	loaderVersion, err := findLatestCompatibleLoader(gameVersion)
	if err != nil {
		return "", err
	}

	installerVersion, err := findLatestInstaller()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/v2/versions/loader/%s/%s/%s/server/jar",
		FABRIC_API_BASE_URL, gameVersion, loaderVersion, installerVersion), nil
}
