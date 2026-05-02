package management

type OpsListEntry struct {
	UUID                string `json:"uuid"`
	Name                string `json:"name"`
	Level               int    `json:"level"`
	BypassesPlayerLimit bool   `json:"bypassesPlayerLimit"`
}

type OpsList []OpsListEntry

func AddPlayerToOpsList(username string, list PlayerList, dir string) error {
	filename := string(list)
	playerList := OpsList{}
	path := filepath.Join(dir, filename)

	err := readPlayerList(path, &playerList)
	if err != nil {
		return err
	}

	for _, player := range playerList {
		if strings.EqualFold(player.Name, username) {
			slog.Info("player already in ops", "username", username)
			return nil
		}
	}

	player, err := mojang.GetPlayerProfile(username)
	if err != nil {
		return err
	}

	playerList = append(playerList,
		OpsListEntry{UUID: player.Id.String(), Name: player.Name, Level: 4, BypassesPlayerLimit: false})
	err = writePlayerList(path, &playerList)
	if err != nil {
		return err
	}

	return nil
}