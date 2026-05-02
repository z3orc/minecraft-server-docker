package manage

import (
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/z3orc/minecraft-server-docker/internal/data/mojang"
)

type WhitelistEntry struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	//IgnoresPlayerLimit bool   `json:"ignoresPlayerLimit"`
}

type Whitelist []WhitelistEntry

func AddPlayerToWhitelist(username string, list PlayerList, dir string) error {
	filename := string(list)
	playerList := Whitelist{}
	path := filepath.Join(dir, filename)

	err := readPlayerList(path, &playerList)
	if err != nil {
		return err
	}

	for _, player := range playerList {
		if strings.EqualFold(player.Name, username) {
			slog.Info("player already in whitelist", "username", username)
			return nil
		}
	}

	player, err := mojang.GetPlayerProfile(username)
	if err != nil {
		return err
	}

	playerList = append(playerList, WhitelistEntry{UUID: player.Id.String(), Name: player.Name})
	err = writePlayerList(path, &playerList)
	if err != nil {
		return err
	}

	return nil
}
