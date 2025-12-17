package mojang

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Fetches user id and username from mojang api, using provided username.
func GetPlayerProfile(username string) (*Profile, error) {
	url := fmt.Sprintf("%s/users/profiles/minecraft/%s", MOJANG_API_BASE_URL, username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile from mojang api: %e", err)
	}
	defer resp.Body.Close()

	profile := Profile{}
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, fmt.Errorf("failed decode user profile from mojang api: %e", err)
	}

	return &profile, nil
}
