package mojang

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/z3orc/minecraft-server-docker/internal/httpclient"
)

// Fetches user id and username from mojang api, using provided username.
func GetPlayerProfile(username string) (*Profile, error) {
	client := httpclient.New()
	url := fmt.Sprintf("%s/users/profiles/minecraft/%s", MOJANG_API_BASE_URL, username)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile from mojang api: %e", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got unexpected status code from mojang api: %d", resp.StatusCode)
	}

	profile := Profile{}
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, fmt.Errorf("failed decode user profile from mojang api: %e", err)
	}

	return &profile, nil
}
