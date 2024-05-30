package commands

import (
	botrcon "DiscordMinecraftHelper/server"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type (
	SlashCommand struct {
		Name        string
		Description string
		Options     []*discordgo.ApplicationCommandOption
	}

	HandleFunc func(s *discordgo.Session, i *discordgo.InteractionCreate, g botrcon.Server)
)

var (
	List = SlashCommand{
		Name:        "list",
		Description: "List the players currently active on the server.",
	}
	Start = SlashCommand{
		Name:        "start",
		Description: "Starts the Minecraft server. Will not restart it if already started.",
	}
	Restart = SlashCommand{
		Name:        "restart",
		Description: "Restarts the Minecraft server manually. This is done every night automatically.",
	}
	Address = SlashCommand{
		Name:        "address",
		Description: "Return the current server IP + port.",
	}
)

func AddCommandHandlers(s *discordgo.Session, server botrcon.Server, logger *slog.Logger) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandHandlers := GetCommandsHandlers()
		logger.Info("command received", "command", i.ApplicationCommandData().Name, "user", i.Member.User.Username)
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i, server)
		}
	})
}

func GetCommands() []SlashCommand {
	return []SlashCommand{
		List,
		Start,
		Restart,
		Address,
	}
}

// Command handlers must be present in the returned
// map along with the command itself or they will
// not be registered.
func GetCommandsHandlers() map[string]HandleFunc {
	return map[string]HandleFunc{
		"list":    PlayerListHandler,
		"restart": RestartServerHandler,
		"start":   StartServerHandler,
		"address": ServerAddressHandler,
	}
}

func RegisterCommands(s *discordgo.Session, guildID string, logger *slog.Logger) []*discordgo.ApplicationCommand {

	commands := GetCommands()

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, v := range commands {
		cmd := &discordgo.ApplicationCommand{
			Name:        v.Name,
			Description: v.Description,
			Options:     v.Options,
		}
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)

		if err != nil {
			logger.Error("error adding command", "command", v.Name, "error", err)
		}
		registeredCommands[i] = cmd
	}

	return registeredCommands
}
