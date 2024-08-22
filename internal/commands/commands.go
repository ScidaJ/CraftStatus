package commands

import (
	botrcon "DiscordMinecraftHelper/internal/server"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type (
	SlashCommand struct {
		Name        string
		Description string
		Options     []*discordgo.ApplicationCommandOption
	}

	HandleFunc func(s *discordgo.Session, i *discordgo.InteractionCreate, g *botrcon.Server)
)

var (
	Address = SlashCommand{
		Name:        "address",
		Description: "Return the current server IP + port.",
	}
	List = SlashCommand{
		Name:        "list",
		Description: "List the players currently active on the server.",
	}
	Restart = SlashCommand{
		Name:        "restart",
		Description: "Restarts the Minecraft server manually. Admin user only.",
	}
)

func AddCommandHandlers(s *discordgo.Session, server *botrcon.Server, logger *slog.Logger) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandHandlers := GetCommandsHandlers()
		logger.Info("command received", "command", i.ApplicationCommandData().Name)
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i, server)
		}
	})
}

func GetCommands() []SlashCommand {
	return []SlashCommand{
		Address,
		List,
		Restart,
	}
}

// Command handlers must be present in the returned
// map along with the command itself or they will
// not be registered.
func GetCommandsHandlers() map[string]HandleFunc {
	return map[string]HandleFunc{
		"address": AddressHandler,
		"list":    ListHandler,
		"restart": RestartHandler,
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
