package bot

import (
	botrcon "DiscordMinecraftHelper/server"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func UpdateBotStatus(s *discordgo.Session, server *botrcon.Server) {
	playerCount, _ := server.GetPlayerCount()
	if server.ServerRunning() {
		activity := discordgo.Activity{
			Name:    fmt.Sprintf("Players: %v online", playerCount),
			Type:    discordgo.ActivityTypeWatching,
			State:   "Online",
			Details: fmt.Sprintf("%v player(s) online!", playerCount),
		}
		presence := discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				&activity,
			},
			Status: string(discordgo.StatusOnline),
			AFK:    false,
		}
		s.UpdateStatusComplex(presence)
	} else {
		activity := discordgo.Activity{
			Name:    "Server offline",
			Type:    discordgo.ActivityTypeWatching,
			State:   "Offline",
			Details: "Server offline",
		}
		presence := discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				&activity,
			},
			Status: string(discordgo.StatusOnline),
			AFK:    false,
		}
		s.UpdateStatusComplex(presence)
	}
}
