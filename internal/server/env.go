package botrcon

import (
	"fmt"
	"os"
	"strings"
)

type ServerEnv struct {
	PLAYER_LIST    map[string]string
	RCON_ADDRESS   string
	RCON_PASSWORD  string
	SERVER_ADDRESS string
}

func NewServerEnv() ServerEnv {
	var env ServerEnv

	env = ServerEnv{
		PLAYER_LIST:    loadPlayerList(),
		RCON_ADDRESS:   os.Getenv("RCON_ADDRESS"),
		RCON_PASSWORD:  os.Getenv("RCON_PASSWORD"),
		SERVER_ADDRESS: serverAddressParser(),
	}

	return env
}

func serverAddressParser() string {
	serverAddress, exists := os.LookupEnv("SERVER_ADDRESS")
	if !exists {
		return ""
	}
	serverPort, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		return serverAddress
	} else {
		return fmt.Sprintf("%v:%v", serverAddress, serverPort)
	}
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
