package main

import "github.com/bwmarrin/discordgo"

type SlashCommand struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
}

var (
	PlayerList = SlashCommand{
		Name:        "player-list",
		Description: "List the players currently active on the server.",
	}
	StartServer = SlashCommand{
		Name:        "start-server",
		Description: "Starts the Minecraft server. Will not restart it if already started.",
	}
	RestartServer = SlashCommand{
		Name:        "restart-server",
		Description: "Restarts the Minecraft server manually. This is done every night automatically.",
	}
	ServerAddress = SlashCommand{
		Name:        "server-address",
		Description: "Return the current server IP + port.",
	}
)
