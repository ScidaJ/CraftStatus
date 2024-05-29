package botrcon

import (
	"os"
	"strings"
)

type ServerEnv struct {
	PLAYER_LIST       map[string]string
	RCON_ADDRESS      string
	RCON_PASSWORD     string
	START_SERVER_PATH string
	SERVER_ADDRESS    string
}

func NewServerEnv() ServerEnv {
	var env ServerEnv

	serverAddress, _ := os.LookupEnv("SERVER_ADDRESS")

	env = ServerEnv{
		PLAYER_LIST:       loadPlayerList(),
		RCON_ADDRESS:      os.Getenv("RCON_ADDRESS"),
		RCON_PASSWORD:     os.Getenv("RCON_PASSWORD"),
		START_SERVER_PATH: os.Getenv("START_SERVER_PATH"),
		SERVER_ADDRESS:    serverAddress,
	}

	return env
}

func loadPlayerList() map[string]string {
	playerList := map[string]string{}

	playersString, _ := os.LookupEnv("PLAYER_LIST")
	playersSlice := strings.Split(playersString, ",")

	for _, v := range playersSlice {
		player := strings.Split(v, ":")
		playerList[player[0]] = player[1]
	}

	return playerList
}
