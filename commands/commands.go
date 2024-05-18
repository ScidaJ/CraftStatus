package commands

import (
	"github.com/bwmarrin/discordgo"
)

type (
	SlashCommand struct {
		Name        string
		Description string
		Options     []*discordgo.ApplicationCommandOption
	}

	HandleFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)
)

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

// Command structs located in commands.go must
// be in the returned slice or they will not be applied
func GetCommands() []SlashCommand {
	return []SlashCommand{
		PlayerList,
		StartServer,
		RestartServer,
		ServerAddress,
	}
}

// Command handlers must be present in the returned
// map along with the command itself or they will
// not be registered.
func GetCommandsHandlers() map[string]HandleFunc {
	return map[string]HandleFunc{
		"player-list":    PlayerListHandler,
		"restart-server": RestartServerHandler,
		"start-server":   StartServerHandler,
		"server-address": ServerAddressHandler,
	}
}
