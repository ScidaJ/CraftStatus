package botrcon

import (
	"os"
	"strings"
)

type ServerEnv struct {
	PLAYER_LIST    map[string]string
	RCON_ADDRESS   string
	RCON_PASSWORD  string
	SERVER_ADDRESS string
	SERVER_PORT    string
}

func NewServerEnv() ServerEnv {
	var env ServerEnv

	serverAddress, _ := os.LookupEnv("SERVER_ADDRESS")
	serverPort, _ := os.LookupEnv("SERVER_PORT")

	env = ServerEnv{
		PLAYER_LIST:    loadPlayerList(),
		RCON_ADDRESS:   os.Getenv("RCON_ADDRESS"),
		RCON_PASSWORD:  os.Getenv("RCON_PASSWORD"),
		SERVER_ADDRESS: serverAddress,
		SERVER_PORT:    serverPort,
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
