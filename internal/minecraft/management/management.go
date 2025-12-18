package management

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/z3orc/minecraft-server-docker/internal/data/mojang"
)

type PlayerList string

const (
	WHITELIST PlayerList = "whitelist.json"
	OPS_LIST  PlayerList = "ops.json"
)

// Reads the player list found at 'path' into struct 'target'. Returns error or nil.
func readPlayerList(path string, target interface{}) error {
	file, err := os.Open(path)
	if err != nil && os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to open player list file '%s': %s", path, err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(target)
	if err != nil && err == io.EOF {
		slog.Warn("reached EOF when decoding player list. assuming list stored on disk is empty")
		return nil
	} else if err != nil {
		return fmt.Errorf("failed decode player list: %s", err)
	}

	return nil
}

func writePlayerList(path string, source interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create/overwrite player list to file at '%s': %s", path, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(source)
	if err != nil {
		return fmt.Errorf("failed to write player list to file at '%s': %s", path, err)
	}

	return nil
}

func AddPlayerToList(username string, listType PlayerList, dir string) error {
	var playerList interface{}
	filename := string(listType)
	path := filepath.Join(dir, filename)

	switch listType {
	case WHITELIST:
		playerList = Whitelist{}
	case OPS_LIST:
		playerList = OpsList{}
	}

	err := readPlayerList(path, &playerList)
	if err != nil {
		return err
	}

	for _, player := range playerList.(Whitelist) {
		if strings.EqualFold(player.Name, username) {
			slog.Debug("player already in whitelist", "username", username)
			return nil
		}
	}

	player, err := mojang.GetPlayerProfile(username)
	if err != nil {
		return err
	}

	if listType == WHITELIST {
		playerList = append(playerList.(Whitelist), WhitelistEntry{UUID: player.Id.String(), Name: player.Name})
	} else if listType == OPS_LIST {
		playerList = append(playerList.(OpsList),
			OpsListEntry{UUID: player.Id.String(), Name: player.Name, Level: 4, BypassesPlayerLimit: false})
	}

	err = writePlayerList(path, &playerList)
	if err != nil {
		return err
	}

	return nil
}

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
			slog.Debug("player already in whitelist", "username", username)
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
			slog.Debug("player already in ops", "username", username)
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
